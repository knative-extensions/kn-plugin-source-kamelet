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
	"testing"

	"gotest.tools/v3/assert"

	camelkapis "github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	camelkv1 "github.com/apache/camel-k/pkg/client/camel/clientset/versioned/typed/camel/v1"
	camelkv1alpha1 "github.com/apache/camel-k/pkg/client/camel/clientset/versioned/typed/camel/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"knative.dev/client/pkg/util/mock"
)

// MockClient is a combine of test object and recorder
type MockClient struct {
	t        *testing.T
	recorder *KameletRecorder
}

func (c *MockClient) RESTClient() rest.Interface {
	panic("should not be called")
}

// NewMockClient returns a new mock instance which you need to record for
func NewMockClient(t *testing.T, ns ...string) *MockClient {
	namespace := "default"
	if len(ns) > 0 {
		namespace = ns[0]
	}
	return &MockClient{
		t:        t,
		recorder: &KameletRecorder{mock.NewRecorder(t, namespace)},
	}
}

// MockKameletClient is a combine of test object and recorder
type MockKameletClient struct {
	t        *testing.T
	recorder *KameletRecorder
}

// NewMockKameletClient returns a new mock instance which you need to record for
func newMockKameletClient(c *MockClient) *MockKameletClient {
	return &MockKameletClient{
		t:        c.t,
		recorder: c.recorder,
	}
}

// MockKameletBindingsClient is a combine of test object and recorder
type MockKameletBindingsClient struct {
	t        *testing.T
	recorder *KameletRecorder
}

// NewMockKameletBindingsClient returns a new mock instance which you need to record for
func newMockKameletBindingsClient(c *MockClient) *MockKameletBindingsClient {
	return &MockKameletBindingsClient{
		t:        c.t,
		recorder: c.recorder,
	}
}

// Ensure that the interface is implemented
var _ camelkv1alpha1.CamelV1alpha1Interface = &MockClient{}
var _ camelkv1alpha1.KameletInterface = &MockKameletClient{}
var _ camelkv1alpha1.KameletBindingInterface = &MockKameletBindingsClient{}

// KameletRecorder is recorder for eventing objects
type KameletRecorder struct {
	r *mock.Recorder
}

func (c *MockClient) CamelV1() camelkv1.CamelV1Interface {
	panic("implement me")
}

func (c *MockClient) CamelV1alpha1() *camelkv1alpha1.CamelV1alpha1Interface {
	var i camelkv1alpha1.CamelV1alpha1Interface = c
	return &i
}

func (c *MockClient) GetScheme() *runtime.Scheme {
	panic("implement me")
}

func (c *MockClient) GetConfig() *rest.Config {
	panic("implement me")
}

func (c *MockClient) GetCurrentNamespace(kubeConfig string) (string, error) {
	panic("implement me")
}

func (c *MockClient) Kamelets(namespace string) camelkv1alpha1.KameletInterface {
	var i camelkv1alpha1.KameletInterface = newMockKameletClient(c)
	return i
}

func (c *MockClient) KameletBindings(namespace string) camelkv1alpha1.KameletBindingInterface {
	var i camelkv1alpha1.KameletBindingInterface = newMockKameletBindingsClient(c)
	return i
}

// Recorder returns the recorder for registering API calls
func (c *MockClient) Recorder() *KameletRecorder {
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

// Get records a call for GetKamelet with the expected result and error (nil if none)
func (sr *KameletRecorder) Get(kamelet *camelkapis.Kamelet, err error) {
	sr.r.Add("Get", nil, []interface{}{kamelet, err})
}

// Get performs a previously recorded action
func (c *MockKameletClient) Get(ctx context.Context, name string, opts v1.GetOptions) (*camelkapis.Kamelet, error) {
	call := c.recorder.r.VerifyCall("Get")
	return call.Result[0].(*camelkapis.Kamelet), mock.ErrorOrNil(call.Result[1])
}

func (c *MockKameletClient) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	panic("implement me")
}

func (c *MockKameletClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *camelkapis.Kamelet, err error) {
	panic("implement me")
}

// CreateKameletBinding records a call for CreateKameletBinding with the expected result and error (nil if none)
func (sr *KameletRecorder) CreateKameletBinding(binding *camelkapis.KameletBinding, err error) {
	sr.r.Add("Create", nil, []interface{}{binding, err})
}

// Create performs a previously recorded action
func (c *MockKameletBindingsClient) Create(ctx context.Context, binding *camelkapis.KameletBinding, opts v1.CreateOptions) (*camelkapis.KameletBinding, error) {
	call := c.recorder.r.VerifyCall("Create")
	assert.DeepEqual(c.t, call.Result[0].(*camelkapis.KameletBinding), binding)
	return call.Result[0].(*camelkapis.KameletBinding), mock.ErrorOrNil(call.Result[1])
}

func (c *MockKameletBindingsClient) Update(ctx context.Context, binding *camelkapis.KameletBinding, opts v1.UpdateOptions) (*camelkapis.KameletBinding, error) {
	panic("implement me")
}

func (c *MockKameletBindingsClient) UpdateStatus(ctx context.Context, binding *camelkapis.KameletBinding, opts v1.UpdateOptions) (*camelkapis.KameletBinding, error) {
	panic("implement me")
}

func (c *MockKameletBindingsClient) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	panic("implement me")
}

func (c *MockKameletBindingsClient) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	panic("implement me")
}
func (c *MockKameletBindingsClient) List(ctx context.Context, opts v1.ListOptions) (*camelkapis.KameletBindingList, error) {
	panic("implement me")
}

// GetKameletBinding records a call for Get with the expected result and error (nil if none)
func (sr *KameletRecorder) GetKameletBinding(binding *camelkapis.KameletBinding, err error) {
	sr.r.Add("Get", nil, []interface{}{binding, err})
}

// Get performs a previously recorded action
func (c *MockKameletBindingsClient) Get(ctx context.Context, name string, opts v1.GetOptions) (*camelkapis.KameletBinding, error) {
	call := c.recorder.r.VerifyCall("Get")
	return call.Result[0].(*camelkapis.KameletBinding), mock.ErrorOrNil(call.Result[1])
}

func (c *MockKameletBindingsClient) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	panic("implement me")
}

func (c *MockKameletBindingsClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *camelkapis.KameletBinding, err error) {
	panic("implement me")
}

// Validate validates whether every recorded action has been called
func (sr *KameletRecorder) Validate() {
	sr.r.CheckThatAllRecordedMethodsHaveBeenCalled()
}
