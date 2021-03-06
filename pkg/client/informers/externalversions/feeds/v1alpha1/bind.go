/*
Copyright 2018 Google, Inc. All rights reserved.

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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	feeds_v1alpha1 "github.com/knative/eventing/pkg/apis/feeds/v1alpha1"
	versioned "github.com/knative/eventing/pkg/client/clientset/versioned"
	internalinterfaces "github.com/knative/eventing/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/knative/eventing/pkg/client/listers/feeds/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// BindInformer provides access to a shared informer and lister for
// Binds.
type BindInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.BindLister
}

type bindInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewBindInformer constructs a new informer for Bind type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewBindInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredBindInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredBindInformer constructs a new informer for Bind type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredBindInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.FeedsV1alpha1().Binds(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.FeedsV1alpha1().Binds(namespace).Watch(options)
			},
		},
		&feeds_v1alpha1.Bind{},
		resyncPeriod,
		indexers,
	)
}

func (f *bindInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredBindInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *bindInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&feeds_v1alpha1.Bind{}, f.defaultInformer)
}

func (f *bindInformer) Lister() v1alpha1.BindLister {
	return v1alpha1.NewBindLister(f.Informer().GetIndexer())
}
