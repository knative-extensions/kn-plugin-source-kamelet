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
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"knative.dev/client-pkg/pkg/printers"
	"knative.dev/pkg/apis"

	"knative.dev/client-pkg/pkg/commands"
	knerrors "knative.dev/client-pkg/pkg/errors"

	"github.com/spf13/cobra"
)

var describeExample = `
  # Describe given Kamelets
  kn source kamelet describe NAME

  # Describe given Kamelets in YAML output format
  kn source kamelet describe NAME -o yaml`

// NewDescribeCommand implements 'kn-source-kamelet describe' command
func NewDescribeCommand(p *KameletPluginParams) *cobra.Command {
	printFlags := genericclioptions.NewPrintFlags("")

	cmd := &cobra.Command{
		Use:     "describe NAME",
		Short:   "Show details of given Kamelet source type",
		Example: describeExample,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) != 1 {
				return errors.New("'kn source kamelet describe' requires the Kamelet name given as single argument")
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

			kamelet, err := client.Kamelets(namespace).Get(p.Context, name, v1.GetOptions{})
			if err != nil {
				return knerrors.GetError(err)
			}
			updateKameletGvk(kamelet)

			out := cmd.OutOrStdout()

			if !isEventSourceType(kamelet) {
				return fmt.Errorf("kamelet %s is not an event source", name)
			}

			if printFlags.OutputFlagSpecified() {
				if strings.ToLower(*printFlags.OutputFormat) == "url" {
					fmt.Fprintf(out, "%s\n", kamelet.GetSelfLink())
					return nil
				}
				printer, err := printFlags.ToPrinter()
				if err != nil {
					return err
				}
				return printer.PrintObj(kamelet, out)
			}

			dw := printers.NewPrefixWriter(out)

			printDetails, err := cmd.Flags().GetBool("verbose")
			if err != nil {
				return err
			}

			writeKamelet(dw, kamelet, printDetails)
			dw.WriteLine()
			if err := dw.Flush(); err != nil {
				return err
			}

			// Condition info
			commands.WriteConditions(dw, asApiConditions(kamelet.Status.Conditions), printDetails)
			if err := dw.Flush(); err != nil {
				return err
			}

			return nil
		},
	}
	flags := cmd.Flags()
	commands.AddNamespaceFlags(flags, false)
	flags.BoolP("verbose", "v", false, "More output.")
	printFlags.AddFlags(cmd)
	cmd.Flag("output").Usage = fmt.Sprintf("Output format. One of: %s.", strings.Join(append(printFlags.AllowedFormats(), "url"), "|"))
	return cmd
}

func writeKamelet(dw printers.PrefixWriter, kamelet *v1alpha1.Kamelet, printDetails bool) {
	commands.WriteMetadata(dw, &kamelet.ObjectMeta, printDetails)
	if kamelet.Spec.Definition.Title != "" {
		dw.WriteAttribute("Description", fmt.Sprintf("%s - %s", kamelet.Spec.Definition.Title, kamelet.Spec.Definition.Description))
	} else {
		dw.WriteAttribute("Description", kamelet.Spec.Definition.Description)
	}

	dw.WriteAttribute("Provider", extractKameletProvider(kamelet))
	dw.WriteAttribute("Support Level", extractKameletSupportLevel(kamelet))
	dw.WriteAttribute("Phase", string(kamelet.Status.Phase))

	dw.WriteLine()
	writeKameletProperties(dw, kamelet)
}

func writeKameletProperties(dw printers.PrefixWriter, kamelet *v1alpha1.Kamelet) {
	section := dw.WriteAttribute("Properties", "")
	maxLen := getMaxPropertyNameLen(kamelet.Spec.Definition.Properties)
	format := "%-" + maxLen + "s %-4s %-8s %s\n"
	section.Writef(format, "Name", "Req", "Type", "Description")

	propertyNames := make([]string, 0, len(kamelet.Spec.Definition.Properties))
	for key := range kamelet.Spec.Definition.Properties {
		propertyNames = append(propertyNames, key)
	}
	sort.Strings(propertyNames)

	for _, propertyName := range propertyNames {
		property := kamelet.Spec.Definition.Properties[propertyName]
		section.Writef(format, propertyName, isRequired(propertyName, kamelet.Spec.Definition.Required), property.Type, property.Description)
	}
}

func isRequired(name string, required []string) string {
	for _, propertyName := range required {
		if propertyName == name {
			return "✓"
		}
	}

	return " "
}

func getMaxPropertyNameLen(properties map[string]v1alpha1.JSONSchemaProps) string {
	max := 0
	for name := range properties {
		if len(name) > max {
			max = len(name)
		}
	}
	return strconv.Itoa(max)
}

func asApiConditions(conditions []v1alpha1.KameletCondition) apis.Conditions {
	var aConditions apis.Conditions

	for _, condition := range conditions {
		aConditions = append(aConditions, apis.Condition{
			Type:   apis.ConditionReady,
			Status: condition.Status,
			LastTransitionTime: apis.VolatileTime{
				Inner: condition.LastTransitionTime,
			},
			Reason:  condition.Reason,
			Message: condition.Message,
		})
	}

	return aConditions
}
