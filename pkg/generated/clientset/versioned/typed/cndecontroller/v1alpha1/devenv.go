/*
Copyright The Kubernetes Authors.

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

package v1alpha1

import (
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1alpha1 "saas-controller.cloud-native-coding.dev/pkg/apis/cndecontroller/v1alpha1"
	scheme "saas-controller.cloud-native-coding.dev/pkg/generated/clientset/versioned/scheme"
)

// DevEnvsGetter has a method to return a DevEnvInterface.
// A group's client should implement this interface.
type DevEnvsGetter interface {
	DevEnvs() DevEnvInterface
}

// DevEnvInterface has methods to work with DevEnv resources.
type DevEnvInterface interface {
	Create(*v1alpha1.DevEnv) (*v1alpha1.DevEnv, error)
	Update(*v1alpha1.DevEnv) (*v1alpha1.DevEnv, error)
	UpdateStatus(*v1alpha1.DevEnv) (*v1alpha1.DevEnv, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.DevEnv, error)
	List(opts v1.ListOptions) (*v1alpha1.DevEnvList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.DevEnv, err error)
	DevEnvExpansion
}

// devEnvs implements DevEnvInterface
type devEnvs struct {
	client rest.Interface
}

// newDevEnvs returns a DevEnvs
func newDevEnvs(c *CndecontrollerV1alpha1Client) *devEnvs {
	return &devEnvs{
		client: c.RESTClient(),
	}
}

// Get takes name of the devEnv, and returns the corresponding devEnv object, and an error if there is any.
func (c *devEnvs) Get(name string, options v1.GetOptions) (result *v1alpha1.DevEnv, err error) {
	result = &v1alpha1.DevEnv{}
	err = c.client.Get().
		Resource("devenvs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of DevEnvs that match those selectors.
func (c *devEnvs) List(opts v1.ListOptions) (result *v1alpha1.DevEnvList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.DevEnvList{}
	err = c.client.Get().
		Resource("devenvs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested devEnvs.
func (c *devEnvs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("devenvs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a devEnv and creates it.  Returns the server's representation of the devEnv, and an error, if there is any.
func (c *devEnvs) Create(devEnv *v1alpha1.DevEnv) (result *v1alpha1.DevEnv, err error) {
	result = &v1alpha1.DevEnv{}
	err = c.client.Post().
		Resource("devenvs").
		Body(devEnv).
		Do().
		Into(result)
	return
}

// Update takes the representation of a devEnv and updates it. Returns the server's representation of the devEnv, and an error, if there is any.
func (c *devEnvs) Update(devEnv *v1alpha1.DevEnv) (result *v1alpha1.DevEnv, err error) {
	result = &v1alpha1.DevEnv{}
	err = c.client.Put().
		Resource("devenvs").
		Name(devEnv.Name).
		Body(devEnv).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *devEnvs) UpdateStatus(devEnv *v1alpha1.DevEnv) (result *v1alpha1.DevEnv, err error) {
	result = &v1alpha1.DevEnv{}
	err = c.client.Put().
		Resource("devenvs").
		Name(devEnv.Name).
		SubResource("status").
		Body(devEnv).
		Do().
		Into(result)
	return
}

// Delete takes name of the devEnv and deletes it. Returns an error if one occurs.
func (c *devEnvs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("devenvs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *devEnvs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("devenvs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched devEnv.
func (c *devEnvs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.DevEnv, err error) {
	result = &v1alpha1.DevEnv{}
	err = c.client.Patch(pt).
		Resource("devenvs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
