/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/golang/glog"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	elaclientset "github.com/elafros/elafros/pkg/client/clientset/versioned"
	elainformers "github.com/elafros/elafros/pkg/client/informers/externalversions"

	clientset "github.com/elafros/eventing/pkg/client/clientset/versioned"
	informers "github.com/elafros/eventing/pkg/client/informers/externalversions"
	"github.com/elafros/eventing/pkg/controller"
	"github.com/elafros/eventing/pkg/controller/bind"
	"github.com/elafros/eventing/pkg/signals"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	threadsPerController = 2
	metricsScrapeAddr    = ":9090"
	metricsScrapePath    = "/metrics"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	client, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building clientset: %s", err.Error())
	}

	elaClient, err := elaclientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatalf("Error building ela clientset: %s", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	informerFactory := informers.NewSharedInformerFactory(client, time.Second*30)
	elaInformerFactory := elainformers.NewSharedInformerFactory(elaClient, time.Second*30)

	// Add new controllers here.
	ctors := []controller.Constructor{
		bind.NewController,
	}

	// Build all of our controllers, with the clients constructed above.
	controllers := make([]controller.Interface, 0, len(ctors))
	for _, ctor := range ctors {
		controllers = append(controllers,
			ctor(kubeClient, client, kubeInformerFactory, informerFactory, elaInformerFactory))
	}

	go kubeInformerFactory.Start(stopCh)
	go informerFactory.Start(stopCh)
	go elaInformerFactory.Start(stopCh)

	// Start all of the controllers.
	for _, ctrlr := range controllers {
		go func(ctrlr controller.Interface) {
			// We don't expect this to return until stop is called,
			// but if it does, propagate it back.
			if err := ctrlr.Run(threadsPerController, stopCh); err != nil {
				glog.Fatalf("Error running controller: %s", err.Error())
			}
		}(ctrlr)
	}

	// Start the endpoint that Prometheus scraper talks to
	srv := &http.Server{Addr: metricsScrapeAddr}
	http.Handle(metricsScrapePath, promhttp.Handler())
	go func() {
		glog.Info("Starting metrics listener at %s", metricsScrapeAddr)
		if err := srv.ListenAndServe(); err != nil {
			glog.Infof("Httpserver: ListenAndServe() finished with error: %s", err)
		}
	}()

	<-stopCh

	// Close the http server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

	glog.Flush()
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}