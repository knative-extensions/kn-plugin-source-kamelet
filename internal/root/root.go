// Copyright © 2021 The Knative Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package root

import (
	"context"
	"github.com/spf13/cobra"
	"knative.dev/kn-plugin-source-kamelet/internal/command"
)

// NewSourceKameletCommand represents the plugin's entrypoint
func NewSourceKameletCommand() *cobra.Command {

	var rootCmd = &cobra.Command{
		Use:   "kn-source-kamelet",
		Short: "Knative eventing Kamelet source plugin",
		Long:  `Plugin manages Kamelets and KameletBindings as Knative eventing sources.`,
	}

	ctx, cancel := context.WithCancel(context.Background())

	p := &command.KameletPluginParams{
		Context:       ctx,
		ContextCancel: cancel,
	}
	p.Initialize()

	rootCmd.AddCommand(command.NewListCommand(p))
	rootCmd.AddCommand(command.NewVersionCommand())

	return rootCmd
}
