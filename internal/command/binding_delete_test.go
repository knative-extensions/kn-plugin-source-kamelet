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

	camelkv1alpha1 "github.com/apache/camel-k/pkg/client/camel/clientset/versioned/typed/camel/v1alpha1"
	"knative.dev/client-pkg/pkg/commands"
	"knative.dev/kn-plugin-source-kamelet/internal/client"

	"gotest.tools/v3/assert"
)

func TestBindingDeleteSetup(t *testing.T) {
	p := KameletPluginParams{
		Context: context.TODO(),
	}

	command := newBindingDeleteCommand(&p)
	assert.Equal(t, command.Use, "delete NAME")
	assert.Equal(t, command.Short, "Delete Kamelet binding by its name.")
}

func TestBindingDeleteErrorCaseMissingArgument(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	err := runBindingDeleteCmd(mockClient)
	assert.Error(t, err, "'kn-source-kamelet binding delete' requires the binding name as argument")
	recorder.Validate()
}

func TestBindingDeleteErrorCaseNotExists(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	recorder.DeleteKameletBinding("k1-to-x", errors.New("not found"))
	err := runBindingDeleteCmd(mockClient, "k1-to-x")
	assert.Error(t, err, "not found")

	recorder.Validate()
}

func TestBindingDelete(t *testing.T) {
	mockClient := client.NewMockClient(t)
	recorder := mockClient.Recorder()

	recorder.DeleteKameletBinding("k1-to-foo", nil)
	err := runBindingDeleteCmd(mockClient, "k1-to-foo")
	assert.NilError(t, err)

	recorder.Validate()
}

func runBindingDeleteCmd(c *client.MockClient, options ...string) error {
	p := KameletPluginParams{
		KnParams: &commands.KnParams{},
		Context:  context.TODO(),
		NewKameletClient: func() (camelkv1alpha1.CamelV1alpha1Interface, error) {
			return c, nil
		},
	}

	command, _, _ := commands.CreateSourcesTestKnCommand(newBindingDeleteCommand(&p), p.KnParams)

	args := []string{"delete"}
	args = append(args, options...)
	command.SetArgs(args)
	err := command.Execute()

	return err
}
