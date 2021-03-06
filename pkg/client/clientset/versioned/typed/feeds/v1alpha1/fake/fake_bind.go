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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/knative/eventing/pkg/apis/feeds/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeBinds implements BindInterface
type FakeBinds struct {
	Fake *FakeFeedsV1alpha1
	ns   string
}

var bindsResource = schema.GroupVersionResource{Group: "feeds.knative.dev", Version: "v1alpha1", Resource: "binds"}

var bindsKind = schema.GroupVersionKind{Group: "feeds.knative.dev", Version: "v1alpha1", Kind: "Bind"}

// Get takes name of the bind, and returns the corresponding bind object, and an error if there is any.
func (c *FakeBinds) Get(name string, options v1.GetOptions) (result *v1alpha1.Bind, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(bindsResource, c.ns, name), &v1alpha1.Bind{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Bind), err
}

// List takes label and field selectors, and returns the list of Binds that match those selectors.
func (c *FakeBinds) List(opts v1.ListOptions) (result *v1alpha1.BindList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(bindsResource, bindsKind, c.ns, opts), &v1alpha1.BindList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.BindList{}
	for _, item := range obj.(*v1alpha1.BindList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested binds.
func (c *FakeBinds) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(bindsResource, c.ns, opts))

}

// Create takes the representation of a bind and creates it.  Returns the server's representation of the bind, and an error, if there is any.
func (c *FakeBinds) Create(bind *v1alpha1.Bind) (result *v1alpha1.Bind, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(bindsResource, c.ns, bind), &v1alpha1.Bind{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Bind), err
}

// Update takes the representation of a bind and updates it. Returns the server's representation of the bind, and an error, if there is any.
func (c *FakeBinds) Update(bind *v1alpha1.Bind) (result *v1alpha1.Bind, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(bindsResource, c.ns, bind), &v1alpha1.Bind{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Bind), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeBinds) UpdateStatus(bind *v1alpha1.Bind) (*v1alpha1.Bind, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(bindsResource, "status", c.ns, bind), &v1alpha1.Bind{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Bind), err
}

// Delete takes name of the bind and deletes it. Returns an error if one occurs.
func (c *FakeBinds) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(bindsResource, c.ns, name), &v1alpha1.Bind{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBinds) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(bindsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.BindList{})
	return err
}

// Patch applies the patch and returns the patched bind.
func (c *FakeBinds) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Bind, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(bindsResource, c.ns, name, data, subresources...), &v1alpha1.Bind{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Bind), err
}
