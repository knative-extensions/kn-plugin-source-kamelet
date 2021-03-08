/*
 * Copyright © 2021 The Knative Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package client

import (
	"context"
	camelkapis "github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	camelkv1 "github.com/apache/camel-k/pkg/client/camel/clientset/versioned/typed/camel/v1"
	camelkv1alpha1 "github.com/apache/camel-k/pkg/client/camel/clientset/versioned/typed/camel/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"knative.dev/client/pkg/util/mock"
	"testing"
)

// MockKameletClient is a combine of test object and recorder
type MockKameletClient struct {
	t        *testing.T
	recorder *KameletRecorder
}

func (c *MockKameletClient) RESTClient() rest.Interface {
	panic("should not be called")
}

// NewMockKameletClient returns a new mock instance which you need to record for
func NewMockKameletClient(t *testing.T, ns ...string) *MockKameletClient {
	namespace := "default"
	if len(ns) > 0 {
		namespace = ns[0]
	}
	return &MockKameletClient{
		t:        t,
		recorder: &KameletRecorder{mock.NewRecorder(t, namespace)},
	}
}

// Ensure that the interface is implemented
var _ camelkv1alpha1.CamelV1alpha1Interface = &MockKameletClient{}
var _ camelkv1alpha1.KameletInterface = &MockKameletClient{}

// KameletRecorder is recorder for eventing objects
type KameletRecorder struct {
	r *mock.Recorder
}

func (c *MockKameletClient) CamelV1() camelkv1.CamelV1Interface {
	panic("implement me")
}

func (c *MockKameletClient) CamelV1alpha1() *camelkv1alpha1.CamelV1alpha1Interface {
	var i camelkv1alpha1.CamelV1alpha1Interface = c
	return &i
}

func (c *MockKameletClient) GetScheme() *runtime.Scheme {
	panic("implement me")
}

func (c *MockKameletClient) GetConfig() *rest.Config {
	panic("implement me")
}

func (c *MockKameletClient) GetCurrentNamespace(kubeConfig string) (string, error) {
	panic("implement me")
}

func (c *MockKameletClient) Kamelets(namespace string) camelkv1alpha1.KameletInterface {
	var i camelkv1alpha1.KameletInterface = c
	return i
}

func (c *MockKameletClient) KameletBindings(namespace string) camelkv1alpha1.KameletBindingInterface {
	panic("implement me")
}

// Recorder returns the recorder for registering API calls
func (c *MockKameletClient) Recorder() *KameletRecorder {
	return c.recorder
}

// List records a call for ListKamelets with the expected result and error (nil if none)
func (sr *KameletRecorder) List(kameletList *camelkapis.KameletList, err error) {
	sr.r.Add("List", nil, []interface{}{kameletList, err})
}

// List performs a previously recorded action
func (c *MockKameletClient) List(ctx context.Context, opts v1.ListOptions) (*camelkapis.KameletList, error) {
	call := c.recorder.r.VerifyCall("List")
	return call.Result[0].(*camelkapis.KameletList), mock.ErrorOrNil(call.Result[1])
}

func (c *MockKameletClient) Create(ctx context.Context, kamelet *camelkapis.Kamelet, opts v1.CreateOptions) (*camelkapis.Kamelet, error) {
	panic("implement me")
}

func (c *MockKameletClient) Update(ctx context.Context, kamelet *camelkapis.Kamelet, opts v1.UpdateOptions) (*camelkapis.Kamelet, error) {
	panic("implement me")
}

func (c *MockKameletClient) UpdateStatus(ctx context.Context, kamelet *camelkapis.Kamelet, opts v1.UpdateOptions) (*camelkapis.Kamelet, error) {
	panic("implement me")
}

func (c *MockKameletClient) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	panic("implement me")
}

func (c *MockKameletClient) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	panic("implement me")
}

func (c *MockKameletClient) Get(ctx context.Context, name string, opts v1.GetOptions) (*camelkapis.Kamelet, error) {
	panic("implement me")
}

func (c *MockKameletClient) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	panic("implement me")
}

func (c *MockKameletClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *camelkapis.Kamelet, err error) {
	panic("implement me")
}

// Validate validates whether every recorded action has been called
func (sr *KameletRecorder) Validate() {
	sr.r.CheckThatAllRecordedMethodsHaveBeenCalled()
}
