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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1alpha1 "saas-controller.cloud-native-coding.dev/pkg/apis/cndecontroller/v1alpha1"
)

// BuilderLister helps list Builders.
type BuilderLister interface {
	// List lists all Builders in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.Builder, err error)
	// Builders returns an object that can list and get Builders.
	Builders(namespace string) BuilderNamespaceLister
	BuilderListerExpansion
}

// builderLister implements the BuilderLister interface.
type builderLister struct {
	indexer cache.Indexer
}

// NewBuilderLister returns a new BuilderLister.
func NewBuilderLister(indexer cache.Indexer) BuilderLister {
	return &builderLister{indexer: indexer}
}

// List lists all Builders in the indexer.
func (s *builderLister) List(selector labels.Selector) (ret []*v1alpha1.Builder, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Builder))
	})
	return ret, err
}

// Builders returns an object that can list and get Builders.
func (s *builderLister) Builders(namespace string) BuilderNamespaceLister {
	return builderNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// BuilderNamespaceLister helps list and get Builders.
type BuilderNamespaceLister interface {
	// List lists all Builders in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.Builder, err error)
	// Get retrieves the Builder from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.Builder, error)
	BuilderNamespaceListerExpansion
}

// builderNamespaceLister implements the BuilderNamespaceLister
// interface.
type builderNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Builders in the indexer for a given namespace.
func (s builderNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Builder, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Builder))
	})
	return ret, err
}

// Get retrieves the Builder from the indexer for a given namespace and name.
func (s builderNamespaceLister) Get(name string) (*v1alpha1.Builder, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("builder"), name)
	}
	return obj.(*v1alpha1.Builder), nil
}
