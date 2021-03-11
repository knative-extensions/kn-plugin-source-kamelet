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

	camelkapis "github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	camelkv1alpha1 "github.com/apache/camel-k/pkg/client/camel/clientset/versioned/typed/camel/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/client/pkg/kn/commands"
	"knative.dev/client/pkg/util"
	"knative.dev/kn-plugin-source-kamelet/internal/client"

	"gotest.tools/v3/assert"
)

func TestListSetup(t *testing.T) {
	p := KameletPluginParams{
		Context: context.TODO(),
	}

	listCmd := NewListTypesCommand(&p)
	assert.Equal(t, listCmd.Use, "list-types")
	assert.Equal(t, listCmd.Short, "List available Kamelet source types")
	assert.Assert(t, listCmd.RunE != nil)
}

func TestListOutput(t *testing.T) {
	mockClient := client.NewMockKameletClient(t)
	recorder := mockClient.Recorder()

	kamelet1 := createKamelet("k1")
	kamelet2 := createKamelet("k2")
	kamelet3 := createKamelet("k3")
	kameletList := &camelkapis.KameletList{Items: []camelkapis.Kamelet{*kamelet1, *kamelet2, *kamelet3}}
	recorder.List(kameletList, nil)

	output, err := runListCmd(mockClient)
	assert.NilError(t, err)

	outputLines := strings.Split(output, "\n")

	assert.Check(t, util.ContainsAll(outputLines[0], "NAME", "PHASE", "AGE", "CONDITIONS", "READY", "REASON"))
	assert.Check(t, util.ContainsAll(outputLines[1], "k1", "Ready", "1 OK / 1", "True"))
	assert.Check(t, util.ContainsAll(outputLines[2], "k2", "Ready", "1 OK / 1", "True"))
	assert.Check(t, util.ContainsAll(outputLines[3], "k3", "Ready", "1 OK / 1", "True"))

	recorder.Validate()
}

func TestListEmpty(t *testing.T) {
	mockClient := client.NewMockKameletClient(t)
	recorder := mockClient.Recorder()

	recorder.List(&camelkapis.KameletList{}, nil)
	output, err := runListCmd(mockClient)
	assert.NilError(t, err)

	assert.Assert(t, util.ContainsAll(output, "No", "resources", "found"))

	recorder.Validate()
}

func TestListNoReadyReasonOutput(t *testing.T) {
	mockClient := client.NewMockKameletClient(t)
	recorder := mockClient.Recorder()

	kamelet1 := createKamelet("k1")
	kamelet2 := createKamelet("k2")
	kamelet3 := createKamelet("k3")

	kamelet2.Status.Phase = camelkapis.KameletPhaseError
	kamelet2.Status.Conditions[0].Status = "False"
	kamelet2.Status.Conditions[0].Reason = "Internal"
	kamelet2.Status.Conditions[0].Message = "Something went wrong"

	kameletList := &camelkapis.KameletList{Items: []camelkapis.Kamelet{*kamelet1, *kamelet2, *kamelet3}}
	recorder.List(kameletList, nil)

	output, err := runListCmd(mockClient)
	assert.NilError(t, err)

	outputLines := strings.Split(output, "\n")

	assert.Check(t, util.ContainsAll(outputLines[0], "NAME", "PHASE", "AGE", "CONDITIONS", "READY", "REASON"))
	assert.Check(t, util.ContainsAll(outputLines[1], "k1", "Ready", "1 OK / 1", "True"))
	assert.Check(t, util.ContainsAll(outputLines[2], "k2", "Error", "0 OK / 1", "False", "Internal : Something went wrong"))
	assert.Check(t, util.ContainsAll(outputLines[3], "k3", "Ready", "1 OK / 1", "True"))

	recorder.Validate()
}

func TestListAllNamespace(t *testing.T) {
	mockClient := client.NewMockKameletClient(t)
	recorder := mockClient.Recorder()

	kamelet1 := createKameletInNamespace("k1", "default1")
	kamelet2 := createKameletInNamespace("k2", "default2")
	kamelet3 := createKameletInNamespace("k3", "default3")
	kameletList := &camelkapis.KameletList{Items: []camelkapis.Kamelet{*kamelet1, *kamelet2, *kamelet3}}
	recorder.List(kameletList, nil)

	output, err := runListCmd(mockClient, "--all-namespaces")
	assert.NilError(t, err)

	outputLines := strings.Split(output, "\n")
	assert.Check(t, util.ContainsAll(outputLines[0], "NAMESPACE", "NAME", "PHASE", "AGE", "CONDITIONS", "READY", "REASON"))
	assert.Check(t, util.ContainsAll(outputLines[1], "default1", "k1"))
	assert.Check(t, util.ContainsAll(outputLines[2], "default2", "k2"))
	assert.Check(t, util.ContainsAll(outputLines[3], "default3", "k3"))

	recorder.Validate()
}

func runListCmd(c *client.MockKameletClient, options ...string) (string, error) {
	p := KameletPluginParams{
		KnParams: &commands.KnParams{},
		Context:  context.TODO(),
		NewKameletClient: func() (camelkv1alpha1.CamelV1alpha1Interface, error) {
			return c, nil
		},
	}

	listCmd, _, output := commands.CreateSourcesTestKnCommand(NewListTypesCommand(&p), p.KnParams)

	args := []string{"list-types"}
	args = append(args, options...)
	listCmd.SetArgs(args)
	err := listCmd.Execute()

	return output.String(), err
}

func createKamelet(kameletName string) *camelkapis.Kamelet {
	return createKameletInNamespace(kameletName, "default")
}

func createKameletInNamespace(kameletName string, namespace string) *camelkapis.Kamelet {
	return &camelkapis.Kamelet{
		TypeMeta: v1.TypeMeta{
			APIVersion: camelkapis.SchemeGroupVersion.String(),
			Kind:       camelkapis.KameletKind,
		},
		ObjectMeta: v1.ObjectMeta{
			Namespace:         namespace,
			Name:              kameletName,
			CreationTimestamp: v1.Now(),
		},
		Spec: camelkapis.KameletSpec{},
		Status: camelkapis.KameletStatus{
			Phase: camelkapis.KameletPhaseReady,
			Conditions: []camelkapis.KameletCondition{
				{
					Type:   "Ready",
					Status: "True",
				},
			},
		},
	}
}
