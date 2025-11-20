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

package command

import (
	"context"
	"errors"
	"testing"

	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	eventingv1 "knative.dev/eventing/pkg/apis/eventing/v1"
	messagingv1 "knative.dev/eventing/pkg/apis/messaging/v1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"

	camelkv1alpha1 "github.com/apache/camel-k/pkg/client/camel/clientset/versioned/typed/camel/v1alpha1"
	"knative.dev/client-pkg/pkg/commands"
	"knative.dev/kn-plugin-source-kamelet/internal/client"

	"gotest.tools/v3/assert"
)

func TestBindSetup(t *testing.T) {
	p := KameletPluginParams{
		Context: context.TODO(),
	}

	bindCmd := NewBindCommand(&p)
	assert.Equal(t, bindCmd.Use, "bind SOURCE")
	assert.Equal(t, bindCmd.Short, "Create Kamelet bindings and bind source to Knative broker, channel or service.")
	assert.Assert(t, bindCmd.RunE != nil)
}

func TestBindErrorCaseMissingArgument(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	err := runBindCmd(mockClient)
	assert.Error(t, err, "'kn-source-kamelet bind' requires the Kamelet source as argument")
	recorder.Validate()
}

func TestBindErrorCaseNotFound(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	kamelet := createKamelet("k1")
	recorder.Get(kamelet, errors.New("not found"))

	err := runBindCmd(mockClient, "k1", "--channel", "test")
	assert.Error(t, err, "not found")
	recorder.Validate()
}

func TestBindErrorCaseNoEventSource(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	kamelet := createKamelet("k1")
	kamelet.Labels = map[string]string{
		KameletTypeLabel: "sink",
	}
	recorder.Get(kamelet, nil)

	err := runBindCmd(mockClient, "k1", "--channel", "test")
	assert.Error(t, err, "kamelet k1 is not an event source")
	recorder.Validate()
}

func TestBindErrorCaseMissingRequiredProperty(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	kamelet := createKamelet("k1")
	recorder.Get(kamelet, nil)

	err := runBindCmd(mockClient, "k1", "--channel", "test")
	assert.Error(t, err, "binding is missing required property \"k1_prop\" for Kamelet \"k1\"")

	recorder.Validate()
}

func TestBindErrorCaseUnknownProperty(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	kamelet := createKamelet("k1")
	recorder.Get(kamelet, nil)

	err := runBindCmd(mockClient, "k1", "--channel", "test", "--property", "k1_prop=foo", "--property", "foo=unknown")
	assert.Error(t, err, "binding uses unknown property \"foo\" for Kamelet \"k1\"")

	recorder.Validate()
}

func TestBindErrorCaseUnsupportedSinkType(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	kamelet := createKamelet("k1")
	recorder.Get(kamelet, nil)

	err := runBindCmd(mockClient, "k1", "--sink", "foo:test", "--property", "k1_prop=foo")
	assert.Error(t, err, "unsupported sink type \"foo\"")

	recorder.Validate()
}

func TestBindErrorCaseUnsupportedSinkExpression(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	kamelet := createKamelet("k1")
	recorder.Get(kamelet, nil)

	err := runBindCmd(mockClient, "k1", "--sink", "foo", "--property", "k1_prop=foo")
	assert.Error(t, err, "unsupported sink expression \"foo\" - please use format <kind>:<name>")

	recorder.Validate()
}

func TestBindToChannel(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	namespace := "current"
	kamelet := createKameletInNamespace("k1", namespace)
	recorder.Get(kamelet, nil)

	recorder.CreateKameletBinding(&v1alpha1.KameletBinding{
		ObjectMeta: v1.ObjectMeta{
			Namespace: namespace,
			Name:      "k1-to-channel-test",
		},
		Spec: v1alpha1.KameletBindingSpec{
			Source: v1alpha1.Endpoint{
				Properties: &v1alpha1.EndpointProperties{
					RawMessage: []byte("{\"k1_prop\":\"foo\"}"),
				},
				Ref: &corev1.ObjectReference{
					Kind:       v1alpha1.KameletKind,
					APIVersion: v1alpha1.SchemeGroupVersion.String(),
					Namespace:  namespace,
					Name:       "k1",
				},
			},
			Sink: v1alpha1.Endpoint{
				Properties: &v1alpha1.EndpointProperties{},
				Ref: &corev1.ObjectReference{
					Kind:       "Channel",
					APIVersion: messagingv1.SchemeGroupVersion.String(),
					Namespace:  namespace,
					Name:       "test",
				},
			},
		},
	}, nil)
	err := runBindCmd(mockClient, "k1", "--channel", "test", "--property", "k1_prop=foo")
	assert.NilError(t, err)

	recorder.Validate()
}

func TestBindToBroker(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	namespace := "current"
	kamelet := createKameletInNamespace("k2", namespace)
	recorder.Get(kamelet, nil)

	recorder.CreateKameletBinding(&v1alpha1.KameletBinding{
		ObjectMeta: v1.ObjectMeta{
			Namespace: namespace,
			Name:      "k2-to-broker-test",
		},
		Spec: v1alpha1.KameletBindingSpec{
			Source: v1alpha1.Endpoint{
				Properties: &v1alpha1.EndpointProperties{
					RawMessage: []byte("{\"k2_optional\":\"bar\",\"k2_prop\":\"foo\"}"),
				},
				Ref: &corev1.ObjectReference{
					Kind:       v1alpha1.KameletKind,
					APIVersion: v1alpha1.SchemeGroupVersion.String(),
					Namespace:  namespace,
					Name:       "k2",
				},
			},
			Sink: v1alpha1.Endpoint{
				Properties: &v1alpha1.EndpointProperties{},
				Ref: &corev1.ObjectReference{
					Kind:       "Broker",
					APIVersion: eventingv1.SchemeGroupVersion.String(),
					Namespace:  namespace,
					Name:       "test",
				},
			},
		},
	}, nil)
	err := runBindCmd(mockClient, "k2", "--broker", "test", "--property", "k2_prop=foo", "--property", "k2_optional=bar")
	assert.NilError(t, err)

	recorder.Validate()
}

func TestBindToService(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	namespace := "current"
	kamelet := createKameletInNamespace("k3", namespace)
	recorder.Get(kamelet, nil)

	recorder.CreateKameletBinding(&v1alpha1.KameletBinding{
		ObjectMeta: v1.ObjectMeta{
			Namespace: namespace,
			Name:      "k3-to-service-test",
		},
		Spec: v1alpha1.KameletBindingSpec{
			Source: v1alpha1.Endpoint{
				Properties: &v1alpha1.EndpointProperties{
					RawMessage: []byte("{\"k3_prop\":\"foo\"}"),
				},
				Ref: &corev1.ObjectReference{
					Kind:       v1alpha1.KameletKind,
					APIVersion: v1alpha1.SchemeGroupVersion.String(),
					Namespace:  namespace,
					Name:       "k3",
				},
			},
			Sink: v1alpha1.Endpoint{
				Properties: &v1alpha1.EndpointProperties{},
				Ref: &corev1.ObjectReference{
					Kind:       "Service",
					APIVersion: servingv1.SchemeGroupVersion.String(),
					Namespace:  namespace,
					Name:       "test",
				},
			},
		},
	}, nil)
	err := runBindCmd(mockClient, "k3", "--service", "test", "--property", "k3_prop=foo")
	assert.NilError(t, err)

	recorder.Validate()
}

func TestBindAutoUpdate(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	namespace := "current"
	kamelet := createKameletInNamespace("k1", namespace)
	recorder.Get(kamelet, nil)

	binding := &v1alpha1.KameletBinding{
		ObjectMeta: v1.ObjectMeta{
			Namespace: namespace,
			Name:      "k1-to-channel-test",
		},
		Spec: v1alpha1.KameletBindingSpec{
			Source: v1alpha1.Endpoint{
				Properties: &v1alpha1.EndpointProperties{
					RawMessage: []byte("{\"k1_prop\":\"foo\"}"),
				},
				Ref: &corev1.ObjectReference{
					Kind:       v1alpha1.KameletKind,
					APIVersion: v1alpha1.SchemeGroupVersion.String(),
					Namespace:  namespace,
					Name:       "k1",
				},
			},
			Sink: v1alpha1.Endpoint{
				Properties: &v1alpha1.EndpointProperties{},
				Ref: &corev1.ObjectReference{
					Kind:       "Channel",
					APIVersion: messagingv1.SchemeGroupVersion.String(),
					Namespace:  namespace,
					Name:       "test",
				},
			},
		},
	}

	recorder.CreateKameletBinding(binding, k8serrors.NewAlreadyExists(v1alpha1.Resource("bindings"), "k1-to-channel-test"))

	recorder.GetKameletBinding(binding, nil)

	recorder.UpdateKameletBinding(binding, nil)

	err := runBindCmd(mockClient, "k1", "--channel", "test", "--property", "k1_prop=foo")
	assert.NilError(t, err)

	recorder.Validate()
}

func TestBindWithCustomCloudEventSettings(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	namespace := "current"
	kamelet := createKameletInNamespace("k4", namespace)
	recorder.Get(kamelet, nil)

	recorder.CreateKameletBinding(&v1alpha1.KameletBinding{
		ObjectMeta: v1.ObjectMeta{
			Namespace: namespace,
			Name:      "k4-ce-settings-test",
		},
		Spec: v1alpha1.KameletBindingSpec{
			Source: v1alpha1.Endpoint{
				Properties: &v1alpha1.EndpointProperties{
					RawMessage: []byte("{\"k4_prop\":\"foo\"}"),
				},
				Ref: &corev1.ObjectReference{
					Kind:       v1alpha1.KameletKind,
					APIVersion: v1alpha1.SchemeGroupVersion.String(),
					Namespace:  namespace,
					Name:       "k4",
				},
			},
			Sink: v1alpha1.Endpoint{
				Properties: &v1alpha1.EndpointProperties{
					RawMessage: []byte("{\"ce.override.subject\":\"custom\",\"cloudEventsSpecVersion\":\"1.0.1\",\"cloudEventsType\":\"custom-type\"}"),
				},
				Ref: &corev1.ObjectReference{
					Kind:       "Channel",
					APIVersion: messagingv1.SchemeGroupVersion.String(),
					Namespace:  namespace,
					Name:       "test",
				},
			},
		},
	}, nil)
	err := runBindCmd(mockClient, "k4", "--channel", "test", "--name", "k4-ce-settings-test", "--property", "k4_prop=foo", "--ce-spec", "1.0.1", "--ce-type", "custom-type", "--ce-override", "subject=custom")
	assert.NilError(t, err)

	recorder.Validate()
}

func runBindCmd(c *client.MockClient, options ...string) error {
	p := KameletPluginParams{
		KnParams: &commands.KnParams{},
		Context:  context.TODO(),
		NewKameletClient: func() (camelkv1alpha1.CamelV1alpha1Interface, error) {
			return c, nil
		},
	}

	bindCmd, _, _ := commands.CreateSourcesTestKnCommand(NewBindCommand(&p), p.KnParams)

	args := []string{"bind"}
	args = append(args, options...)
	bindCmd.SetArgs(args)
	err := bindCmd.Execute()

	return err
}
