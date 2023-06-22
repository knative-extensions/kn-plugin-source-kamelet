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
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
	eventingv1 "knative.dev/eventing/pkg/apis/eventing/v1"
	messagingv1 "knative.dev/eventing/pkg/apis/messaging/v1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"

	camelkapis "github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	camelkv1alpha1 "github.com/apache/camel-k/pkg/client/camel/clientset/versioned/typed/camel/v1alpha1"
	"knative.dev/client-pkg/pkg/kn/commands"
	"knative.dev/client-pkg/pkg/util"
	"knative.dev/kn-plugin-source-kamelet/internal/client"

	"gotest.tools/v3/assert"
)

func TestBindingListSetup(t *testing.T) {
	p := KameletPluginParams{
		Context: context.TODO(),
	}

	listCmd := newBindingListCommand(&p)
	assert.Equal(t, listCmd.Use, "list")
	assert.Equal(t, listCmd.Short, "List Kamelet bindings.")
	assert.Assert(t, listCmd.RunE != nil)
}

func TestBindingListOutput(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	binding1 := createKameletBinding("k1-to-broker", "k1", &corev1.ObjectReference{
		Kind:       "Broker",
		APIVersion: eventingv1.SchemeGroupVersion.String(),
		Namespace:  "default",
		Name:       "b1",
	})
	binding1.Status = statusReady()
	binding2 := createKameletBinding("k2-to-channel", "k2", &corev1.ObjectReference{
		Kind:       "Channel",
		APIVersion: eventingv1.SchemeGroupVersion.String(),
		Namespace:  "default",
		Name:       "c1",
	})
	binding2.Status = statusReady()
	binding3 := createKameletBinding("k3-to-service", "k3", &corev1.ObjectReference{
		Kind:       "Service",
		APIVersion: servingv1.SchemeGroupVersion.String(),
		Namespace:  "default",
		Name:       "s1",
	})
	binding3.Status = statusReady()
	bindingList := &camelkapis.KameletBindingList{Items: []camelkapis.KameletBinding{*binding1, *binding2, *binding3}}
	recorder.ListBindings(bindingList, nil)

	output, err := runBindingListCmd(mockClient)
	assert.NilError(t, err)

	outputLines := strings.Split(output, "\n")

	assert.Check(t, util.ContainsAll(outputLines[0], "NAME", "PHASE", "AGE", "CONDITIONS", "READY", "REASON"))
	assert.Check(t, util.ContainsAll(outputLines[1], "k1-to-broker", "Ready", "1 OK / 1", "True"))
	assert.Check(t, util.ContainsAll(outputLines[2], "k2-to-channel", "Ready", "1 OK / 1", "True"))
	assert.Check(t, util.ContainsAll(outputLines[3], "k3-to-service", "Ready", "1 OK / 1", "True"))

	recorder.Validate()
}

func TestBindingListEmpty(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	recorder.ListBindings(&camelkapis.KameletBindingList{}, nil)
	output, err := runBindingListCmd(mockClient)
	assert.NilError(t, err)

	assert.Assert(t, util.ContainsAll(output, "No", "resources", "found"))

	recorder.Validate()
}

func TestBindingListNoReadyReasonOutput(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	binding1 := createKameletBinding("k1-to-broker", "k1", &corev1.ObjectReference{
		Kind:       "Broker",
		APIVersion: eventingv1.SchemeGroupVersion.String(),
		Namespace:  "default",
		Name:       "b1",
	})
	binding1.Status = statusReady()
	binding2 := createKameletBinding("k2-to-channel", "k2", &corev1.ObjectReference{
		Kind:       "Channel",
		APIVersion: eventingv1.SchemeGroupVersion.String(),
		Namespace:  "default",
		Name:       "c1",
	})
	binding2.Status = camelkapis.KameletBindingStatus{
		Phase: camelkapis.KameletBindingPhaseError,
		Conditions: []camelkapis.KameletBindingCondition{
			{
				Type:    camelkapis.KameletBindingConditionReady,
				Status:  corev1.ConditionFalse,
				Reason:  "Internal",
				Message: "Something went wrong",
			},
		},
	}
	binding3 := createKameletBinding("k3-to-channel", "k3", &corev1.ObjectReference{
		Kind:       "Channel",
		APIVersion: messagingv1.SchemeGroupVersion.String(),
		Namespace:  "default",
		Name:       "c2",
	})
	binding3.Status = statusReady()

	bindingList := &camelkapis.KameletBindingList{Items: []camelkapis.KameletBinding{*binding1, *binding2, *binding3}}
	recorder.ListBindings(bindingList, nil)

	output, err := runBindingListCmd(mockClient)
	assert.NilError(t, err)

	outputLines := strings.Split(output, "\n")

	assert.Check(t, util.ContainsAll(outputLines[0], "NAME", "PHASE", "AGE", "CONDITIONS", "READY", "REASON"))
	assert.Check(t, util.ContainsAll(outputLines[1], "k1-to-broker", "Ready", "1 OK / 1", "True"))
	assert.Check(t, util.ContainsAll(outputLines[2], "k2-to-channel", "Error", "0 OK / 1", "False", "Internal : Something went wrong"))
	assert.Check(t, util.ContainsAll(outputLines[3], "k3-to-channel", "Ready", "1 OK / 1", "True"))

	recorder.Validate()
}

func TestBindingListAllNamespace(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	binding1 := createKameletBindingInNamespace("k1-to-broker", "k1", "default1", &corev1.ObjectReference{
		Kind:       "Broker",
		APIVersion: eventingv1.SchemeGroupVersion.String(),
		Namespace:  "default1",
		Name:       "b1",
	})
	binding1.Status = statusReady()
	binding2 := createKameletBindingInNamespace("k2-to-channel", "k2", "default2", &corev1.ObjectReference{
		Kind:       "Channel",
		APIVersion: eventingv1.SchemeGroupVersion.String(),
		Namespace:  "default2",
		Name:       "c1",
	})
	binding2.Status = statusReady()
	binding3 := createKameletBindingInNamespace("k3-to-channel", "k3", "default3", &corev1.ObjectReference{
		Kind:       "Channel",
		APIVersion: messagingv1.SchemeGroupVersion.String(),
		Namespace:  "default3",
		Name:       "c2",
	})
	binding3.Status = statusReady()
	bindingList := &camelkapis.KameletBindingList{Items: []camelkapis.KameletBinding{*binding1, *binding2, *binding3}}
	recorder.ListBindings(bindingList, nil)

	output, err := runBindingListCmd(mockClient, "--all-namespaces")
	assert.NilError(t, err)

	outputLines := strings.Split(output, "\n")
	assert.Check(t, util.ContainsAll(outputLines[0], "NAMESPACE", "NAME", "PHASE", "AGE", "CONDITIONS", "READY", "REASON"))
	assert.Check(t, util.ContainsAll(outputLines[1], "default1", "k1-to-broker"))
	assert.Check(t, util.ContainsAll(outputLines[2], "default2", "k2-to-channel"))
	assert.Check(t, util.ContainsAll(outputLines[3], "default3", "k3-to-channel"))

	recorder.Validate()
}

func runBindingListCmd(c *client.MockClient, options ...string) (string, error) {
	p := KameletPluginParams{
		KnParams: &commands.KnParams{},
		Context:  context.TODO(),
		NewKameletClient: func() (camelkv1alpha1.CamelV1alpha1Interface, error) {
			return c, nil
		},
	}

	command, _, output := commands.CreateSourcesTestKnCommand(newBindingListCommand(&p), p.KnParams)

	args := []string{"list"}
	args = append(args, options...)
	command.SetArgs(args)
	err := command.Execute()

	return output.String(), err
}
