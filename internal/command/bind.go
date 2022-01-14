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
	"errors"

	"github.com/spf13/cobra"
	knerrors "knative.dev/client/pkg/errors"
	"knative.dev/client/pkg/kn/commands"
)

var bindExample = `
  # Bind Kamelets to a Knative sink
  kn source kamelet bind SOURCE

  # Add a binding properties
  kn source kamelet bind SOURCE --sink|broker|channel|service=<name> --property=<key>=<value>`

// NewBindCommand implements 'kn-source-kamelet bind' command
func NewBindCommand(p *KameletPluginParams) *cobra.Command {
	var properties []string
	var sink string
	var broker string
	var channel string
	var service string
	cmd := &cobra.Command{
		Use:     "bind SOURCE",
		Short:   "Create Kamelet bindings and bind source to Knative broker, channel or service.",
		Example: bindExample,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) != 1 {
				return errors.New("'kn-source-kamelet bind' requires the Kamelet source as argument")
			}
			source := args[0]

			namespace, err := p.GetNamespace(cmd)
			if err != nil {
				return err
			}

			client, err := p.NewKameletClient()
			if err != nil {
				return err
			}

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return knerrors.GetError(err)
			}

			options := CreateBindingOptions{
				Name:             name,
				Source:           source,
				SourceProperties: properties,
				Sink:             sink,
				Broker:           broker,
				Channel:          channel,
				Service:          service,
				Force:            true,
				CmdOut:           cmd.OutOrStdout(),
			}

			err = createBinding(client, p.Context, namespace, options)
			if err != nil {
				return err
			}

			return nil
		},
	}
	flags := cmd.Flags()
	commands.AddNamespaceFlags(flags, false)

	flags.String("name", "", "Binding name.")
	flags.StringVarP(&sink, "sink", "s", "", "Sink expression to define the binding sink.")
	flags.StringVar(&broker, "broker", "", "Uses a broker as binding sink.")
	flags.StringVar(&channel, "channel", "", "Uses a channel as binding sink.")
	flags.StringVar(&service, "service", "", "Uses a Knative service as binding sink.")
	flags.StringArrayVar(&properties, "property", nil, `Add a source property in the form of "<key>=<value>"`)
	return cmd
}
