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
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	eventingv1 "knative.dev/eventing/pkg/apis/eventing/v1"
	messagingv1 "knative.dev/eventing/pkg/apis/messaging/v1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"

	camelkv1alpha1 "github.com/apache/camel-k/pkg/client/camel/clientset/versioned/typed/camel/v1alpha1"
	"knative.dev/client/pkg/kn/commands"
	"knative.dev/kn-plugin-source-kamelet/internal/client"

	"gotest.tools/v3/assert"
)

func TestBindingCreateSetup(t *testing.T) {
	p := KameletPluginParams{
		Context: context.TODO(),
	}

	command := newBindingCreateCommand(&p)
	assert.Equal(t, command.Use, "create")
	assert.Equal(t, command.Short, "Create Kamelet bindings and bind source to Knative broker, channel or service.")
}
func TestBindingCreateErrorCaseMissingArgument(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	err := runBindingCreateCmd(mockClient)
	assert.Error(t, err, "'kn-source-kamelet binding create' requires the binding name as argument")
	recorder.Validate()
}

func TestBindingCreateErrorCaseNotFound(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	kamelet := createKamelet("k1")
	recorder.Get(kamelet, errors.New("not found"))

	err := runBindingCreateCmd(mockClient, "k1-to-sink", "--kamelet", "k1", "--channel", "test")
	assert.Error(t, err, "not found")
	recorder.Validate()
}

func TestBindingCreateErrorCaseNoEventSource(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	kamelet := createKamelet("k1")
	kamelet.Labels = map[string]string{
		KameletTypeLabel: "sink",
	}
	recorder.Get(kamelet, nil)

	err := runBindingCreateCmd(mockClient, "k1-to-sink", "--kamelet", "k1", "--channel", "test")
	assert.Error(t, err, "Kamelet k1 is not an event source")
	recorder.Validate()
}

func TestBindingCreateErrorCaseMissingRequiredProperty(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	kamelet := createKamelet("k1")
	recorder.Get(kamelet, nil)

	err := runBindingCreateCmd(mockClient, "k1-to-sink", "--kamelet", "k1", "--channel", "test")
	assert.Error(t, err, "binding is missing required property \"k1_prop\" for Kamelet \"k1\"")

	recorder.Validate()
}

func TestBindingCreateErrorCaseUnsupportedSinkType(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	kamelet := createKamelet("k1")
	recorder.Get(kamelet, nil)

	err := runBindingCreateCmd(mockClient, "k1-to-foo", "--kamelet", "k1", "--sink", "foo:test", "--source-property", "k1_prop=foo")
	assert.Error(t, err, "unsupported sink type \"foo\"")

	recorder.Validate()
}

func TestBindingCreateErrorCaseUnsupportedSinkExpression(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	kamelet := createKamelet("k1")
	recorder.Get(kamelet, nil)

	err := runBindingCreateCmd(mockClient, "k1-to-foo", "--kamelet", "k1", "--sink", "foo", "--source-property", "k1_prop=foo")
	assert.Error(t, err, "unsupported sink expression \"foo\" - please use format <kind>:<name>")

	recorder.Validate()
}

func TestBindingCreateToChannel(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	namespace := "current"
	kamelet := createKameletInNamespace("k1", namespace)
	recorder.Get(kamelet, nil)

	recorder.CreateKameletBinding(&v1alpha1.KameletBinding{
		ObjectMeta: v1.ObjectMeta{
			Namespace: namespace,
			Name:      "k1-to-channel",
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
				Ref: &corev1.ObjectReference{
					Kind:       "Channel",
					APIVersion: messagingv1.SchemeGroupVersion.String(),
					Namespace:  namespace,
					Name:       "test",
				},
			},
		},
	}, nil)
	err := runBindingCreateCmd(mockClient, "k1-to-channel", "--kamelet", "k1", "--channel", "test", "--source-property", "k1_prop=foo")
	assert.NilError(t, err)

	recorder.Validate()
}

func TestBindingCreateToBroker(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	namespace := "current"
	kamelet := createKameletInNamespace("k2", namespace)
	recorder.Get(kamelet, nil)

	recorder.CreateKameletBinding(&v1alpha1.KameletBinding{
		ObjectMeta: v1.ObjectMeta{
			Namespace: namespace,
			Name:      "k2-to-broker",
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
				Ref: &corev1.ObjectReference{
					Kind:       "Broker",
					APIVersion: eventingv1.SchemeGroupVersion.String(),
					Namespace:  namespace,
					Name:       "test",
				},
			},
		},
	}, nil)
	err := runBindingCreateCmd(mockClient, "k2-to-broker", "--kamelet", "k2", "--broker", "test", "--source-property", "k2_prop=foo", "--source-property", "k2_optional=bar")
	assert.NilError(t, err)

	recorder.Validate()
}

func TestBindingCreateToService(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	namespace := "current"
	kamelet := createKameletInNamespace("k3", namespace)
	recorder.Get(kamelet, nil)

	recorder.CreateKameletBinding(&v1alpha1.KameletBinding{
		ObjectMeta: v1.ObjectMeta{
			Namespace: namespace,
			Name:      "k3-to-service",
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
				Ref: &corev1.ObjectReference{
					Kind:       "Service",
					APIVersion: servingv1.SchemeGroupVersion.String(),
					Namespace:  namespace,
					Name:       "test",
				},
			},
		},
	}, nil)
	err := runBindingCreateCmd(mockClient, "k3-to-service", "--kamelet", "k3", "--service", "test", "--source-property", "k3_prop=foo")
	assert.NilError(t, err)

	recorder.Validate()
}

func runBindingCreateCmd(c *client.MockClient, options ...string) error {
	p := KameletPluginParams{
		KnParams: &commands.KnParams{},
		Context:  context.TODO(),
		NewKameletClient: func() (camelkv1alpha1.CamelV1alpha1Interface, error) {
			return c, nil
		},
	}

	command, _, _ := commands.CreateSourcesTestKnCommand(newBindingCreateCommand(&p), p.KnParams)

	args := []string{"create"}
	args = append(args, options...)
	command.SetArgs(args)
	err := command.Execute()

	return err
}
