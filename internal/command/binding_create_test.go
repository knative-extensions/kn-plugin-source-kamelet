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

	corev1 "k8s.io/api/core/v1"
	eventingv1 "knative.dev/eventing/pkg/apis/eventing/v1"
	messagingv1 "knative.dev/eventing/pkg/apis/messaging/v1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"

	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	camelkv1alpha1 "github.com/apache/camel-k/pkg/client/camel/clientset/versioned/typed/camel/v1alpha1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"knative.dev/client/pkg/kn/commands"
	"knative.dev/kn-plugin-source-kamelet/internal/client"

	"gotest.tools/v3/assert"
)

func TestBindingCreateSetup(t *testing.T) {
	p := KameletPluginParams{
		Context: context.TODO(),
	}

	command := newBindingCreateCommand(&p)
	assert.Equal(t, command.Use, "create NAME")
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

func TestBindingCreateErrorCaseUnknownProperty(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	kamelet := createKamelet("k1")
	recorder.Get(kamelet, nil)

	err := runBindingCreateCmd(mockClient, "k1-to-sink", "--kamelet", "k1", "--channel", "test", "--property", "k1_prop=foo", "--property", "foo=unknown")
	assert.Error(t, err, "binding uses unknown property \"foo\" for Kamelet \"k1\"")

	recorder.Validate()
}

func TestBindingCreateErrorCaseUnsupportedSinkType(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	kamelet := createKamelet("k1")
	recorder.Get(kamelet, nil)

	err := runBindingCreateCmd(mockClient, "k1-to-foo", "--kamelet", "k1", "--sink", "foo:test", "--property", "k1_prop=foo")
	assert.Error(t, err, "unsupported sink type \"foo\"")

	recorder.Validate()
}

func TestBindingCreateErrorCaseUnsupportedSinkExpression(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	kamelet := createKamelet("k1")
	recorder.Get(kamelet, nil)

	err := runBindingCreateCmd(mockClient, "k1-to-foo", "--kamelet", "k1", "--sink", "foo", "--property", "k1_prop=foo")
	assert.Error(t, err, "unsupported sink expression \"foo\" - please use format <kind>:<name>")

	recorder.Validate()
}

func TestBindingCreateErrorCaseAlreadyExists(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	namespace := "current"
	kamelet := createKameletInNamespace("k1", namespace)
	recorder.Get(kamelet, nil)

	recorder.CreateKameletBinding(createKameletBindingInNamespace("k1-to-channel", "k1", namespace, &corev1.ObjectReference{
		Kind:       "Channel",
		APIVersion: messagingv1.SchemeGroupVersion.String(),
		Namespace:  namespace,
		Name:       "test",
	}), k8serrors.NewAlreadyExists(v1alpha1.Resource("bindings"), "k1-to-channel-test"))

	err := runBindingCreateCmd(mockClient, "k1-to-channel", "--kamelet", "k1", "--channel", "test", "--property", "k1_prop=foo")
	assert.Error(t, err, "kamelet binding with name \"k1-to-channel\" already exists. Use --force to recreate the binding")

	recorder.Validate()
}

func TestBindingCreateToChannel(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	namespace := "current"
	kamelet := createKameletInNamespace("k1", namespace)
	recorder.Get(kamelet, nil)

	recorder.CreateKameletBinding(createKameletBindingInNamespace("k1-to-channel", "k1", namespace, &corev1.ObjectReference{
		Kind:       "Channel",
		APIVersion: messagingv1.SchemeGroupVersion.String(),
		Namespace:  namespace,
		Name:       "test",
	}), nil)
	err := runBindingCreateCmd(mockClient, "k1-to-channel", "--kamelet", "k1", "--channel", "test", "--property", "k1_prop=foo")
	assert.NilError(t, err)

	recorder.Validate()
}

func TestBindingCreateToBroker(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	namespace := "current"
	kamelet := createKameletInNamespace("k2", namespace)
	recorder.Get(kamelet, nil)

	binding := createKameletBindingInNamespace("k2-to-broker", "k2", namespace, &corev1.ObjectReference{
		Kind:       "Broker",
		APIVersion: eventingv1.SchemeGroupVersion.String(),
		Namespace:  namespace,
		Name:       "test",
	})

	binding.Spec.Source.Properties.RawMessage = []byte("{\"k2_optional\":\"bar\",\"k2_prop\":\"foo\"}")

	recorder.CreateKameletBinding(binding, nil)
	err := runBindingCreateCmd(mockClient, "k2-to-broker", "--kamelet", "k2", "--broker", "test", "--property", "k2_prop=foo", "--property", "k2_optional=bar")
	assert.NilError(t, err)

	recorder.Validate()
}

func TestBindingCreateToService(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	namespace := "current"
	kamelet := createKameletInNamespace("k3", namespace)
	recorder.Get(kamelet, nil)

	recorder.CreateKameletBinding(createKameletBindingInNamespace("k3-to-service", "k3", namespace, &corev1.ObjectReference{
		Kind:       "Service",
		APIVersion: servingv1.SchemeGroupVersion.String(),
		Namespace:  namespace,
		Name:       "test",
	}), nil)
	err := runBindingCreateCmd(mockClient, "k3-to-service", "--kamelet", "k3", "--service", "test", "--property", "k3_prop=foo")
	assert.NilError(t, err)

	recorder.Validate()
}

func TestBindingCreateWithCloudEventsSettings(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	namespace := "current"
	kamelet := createKameletInNamespace("k4", namespace)
	recorder.Get(kamelet, nil)

	binding := createKameletBindingInNamespace("k4-to-channel", "k4", namespace, &corev1.ObjectReference{
		Kind:       "Channel",
		APIVersion: messagingv1.SchemeGroupVersion.String(),
		Namespace:  namespace,
		Name:       "test",
	})

	binding.Spec.Sink.Properties.RawMessage = []byte("{\"ce.override.subject\":\"custom\",\"cloudEventsSpecVersion\":\"1.0.1\",\"cloudEventsType\":\"custom-type\"}")

	recorder.CreateKameletBinding(binding, nil)
	err := runBindingCreateCmd(mockClient, "k4-to-channel", "--kamelet", "k4", "--channel", "test", "--property", "k4_prop=foo", "--ce-spec", "1.0.1", "--ce-type", "custom-type", "--ce-override", "subject=custom")
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
