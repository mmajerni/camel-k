/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

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
	"context"
	time "time"

	camelv1alpha1 "github.com/apache/camel-k/v2/pkg/apis/camel/v1alpha1"
	versioned "github.com/apache/camel-k/v2/pkg/client/camel/clientset/versioned"
	internalinterfaces "github.com/apache/camel-k/v2/pkg/client/camel/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/apache/camel-k/v2/pkg/client/camel/listers/camel/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// KameletBindingInformer provides access to a shared informer and lister for
// KameletBindings.
type KameletBindingInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.KameletBindingLister
}

type kameletBindingInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewKameletBindingInformer constructs a new informer for KameletBinding type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewKameletBindingInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredKameletBindingInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredKameletBindingInformer constructs a new informer for KameletBinding type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredKameletBindingInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CamelV1alpha1().KameletBindings(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CamelV1alpha1().KameletBindings(namespace).Watch(context.TODO(), options)
			},
		},
		&camelv1alpha1.KameletBinding{},
		resyncPeriod,
		indexers,
	)
}

func (f *kameletBindingInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredKameletBindingInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *kameletBindingInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&camelv1alpha1.KameletBinding{}, f.defaultInformer)
}

func (f *kameletBindingInformer) Lister() v1alpha1.KameletBindingLister {
	return v1alpha1.NewKameletBindingLister(f.Informer().GetIndexer())
}