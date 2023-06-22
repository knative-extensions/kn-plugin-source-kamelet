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
	"fmt"
	"io"

	knerrors "knative.dev/client-pkg/pkg/errors"

	camelkv1alpha1 "github.com/apache/camel-k/pkg/client/camel/clientset/versioned/typed/camel/v1alpha1"

	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/client-pkg/pkg/kn/commands"
)

var bindingDeleteExample = `
  # Delete Kamelet binding.
  kn source kamelet binding delete NAME`

// newBindingDeleteCommand implements 'kn-source-kamelet binding delete' command
func newBindingDeleteCommand(p *KameletPluginParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete NAME",
		Short:   "Delete Kamelet binding by its name.",
		Example: bindingDeleteExample,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) != 1 {
				return errors.New("'kn-source-kamelet binding delete' requires the binding name as argument")
			}
			name := args[0]

			namespace, err := p.GetNamespace(cmd)
			if err != nil {
				return err
			}

			client, err := p.NewKameletClient()
			if err != nil {
				return err
			}

			err = deleteBinding(client, p.Context, name, namespace, cmd.OutOrStdout())
			if err != nil {
				return err
			}

			return nil
		},
	}
	flags := cmd.Flags()
	commands.AddNamespaceFlags(flags, false)

	return cmd
}

func deleteBinding(client camelkv1alpha1.CamelV1alpha1Interface, ctx context.Context, name string, namespace string, cmdOut io.Writer) error {
	err := client.KameletBindings(namespace).Delete(ctx, name, v1.DeleteOptions{})
	if err != nil {
		return knerrors.GetError(err)
	}

	_, _ = fmt.Fprintf(cmdOut, "kamelet binding %q deleted\n", name)

	return nil
}
