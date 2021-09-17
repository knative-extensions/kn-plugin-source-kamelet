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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/client/pkg/kn/commands"

	camelkv1alpha1 "github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/spf13/cobra"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"knative.dev/client/pkg/kn/commands/flags"
	hprinters "knative.dev/client/pkg/printers"
)

var listExample = `
  # List available Kamelets
  kn-source-kamelet list

  # List available Kamelets in YAML output format
  kn-source-kamelet list -o yaml`

// NewListCommand implements 'kn-source-kamelet list' command
func NewListCommand(p *KameletPluginParams) *cobra.Command {
	kameletListFlags := flags.NewListPrintFlags(ListHandlers)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List available Kamelet source types",
		Aliases: []string{"ls"},
		Example: listExample,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			namespace, err := p.GetNamespace(cmd)
			if err != nil {
				return err
			}

			kameletClient, err := p.NewKameletClient()
			if err != nil {
				return err
			}

			filterCriteria := v1.ListOptions{
				LabelSelector: fmt.Sprintf("%s=%s", KameletTypeLabel, "source"),
			}

			kameletList, err := kameletClient.Kamelets(namespace).List(p.Context, filterCriteria)
			if err != nil {
				return err
			}
			if len(kameletList.Items) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No resources found.\n")
				return nil
			}

			// empty namespace indicates all-namespaces flag is specified
			if namespace == "" {
				kameletListFlags.EnsureWithNamespace()
			}

			err = kameletListFlags.Print(kameletList, cmd.OutOrStdout())
			if err != nil {
				return err
			}
			return nil
		},
	}
	commands.AddNamespaceFlags(cmd.Flags(), true)
	kameletListFlags.AddFlags(cmd)
	return cmd
}

// ListHandlers handles printing human readable table for `kn-source-kamelet list` command's output
func ListHandlers(h hprinters.PrintHandler) {
	kameletColumnDefinitions := []metav1beta1.TableColumnDefinition{
		{Name: "Namespace", Type: "string", Description: "Namespace of the Kamelet instance", Priority: 0},
		{Name: "Name", Type: "string", Description: "Name of the Kamelet instance", Priority: 1},
		{Name: "Phase", Type: "string", Description: "Phase of the Kamelet instance", Priority: 1},
		{Name: "Age", Type: "string", Description: "Age of the Kamelet instance", Priority: 1},
		{Name: "Conditions", Type: "string", Description: "Ready state conditions", Priority: 1},
		{Name: "Ready", Type: "string", Description: "Ready state of the Kamelet instance", Priority: 1},
		{Name: "Reason", Type: "string", Description: "Reason if state is not Ready", Priority: 1},
	}
	h.TableHandler(kameletColumnDefinitions, printKamelet)
	h.TableHandler(kameletColumnDefinitions, printKameletList)
}

// printKameletList populates the Kamelet list table rows
func printKameletList(kameletList *camelkv1alpha1.KameletList, options hprinters.PrintOptions) ([]metav1beta1.TableRow, error) {
	rows := make([]metav1beta1.TableRow, 0, len(kameletList.Items))

	for i := range kameletList.Items {
		ksvc := &kameletList.Items[i]
		r, err := printKamelet(ksvc, options)
		if err != nil {
			return nil, err
		}
		rows = append(rows, r...)
	}
	return rows, nil
}

// printKamelet populates the Kamelet table rows
func printKamelet(kamelet *camelkv1alpha1.Kamelet, options hprinters.PrintOptions) ([]metav1beta1.TableRow, error) {
	name := kamelet.Name
	phase := kamelet.Status.Phase
	age := commands.TranslateTimestampSince(kamelet.CreationTimestamp)
	conditions := conditionsValue(kamelet.Status.Conditions)
	ready := readyCondition(kamelet.Status.Conditions)
	reason := nonReadyConditionReason(kamelet.Status.Conditions)

	row := metav1beta1.TableRow{
		Object: runtime.RawExtension{Object: kamelet},
	}

	if options.AllNamespaces {
		row.Cells = append(row.Cells, kamelet.Namespace)
	}

	row.Cells = append(row.Cells,
		name,
		phase,
		age,
		conditions,
		ready,
		reason)
	return []metav1beta1.TableRow{row}, nil
}

// conditionsValue returns the True conditions count among total conditions
func conditionsValue(conditions []camelkv1alpha1.KameletCondition) string {
	var ok int
	for _, condition := range conditions {
		if condition.Status == "True" {
			ok++
		}
	}
	return fmt.Sprintf("%d OK / %d", ok, len(conditions))
}

// readyCondition returns status of resource's Ready type condition
func readyCondition(conditions []camelkv1alpha1.KameletCondition) string {
	for _, condition := range conditions {
		if condition.Type == camelkv1alpha1.KameletConditionReady {
			return string(condition.Status)
		}
	}
	return "<unknown>"
}

// NonReadyConditionReason returns formatted string of
// reason and message for non ready conditions
func nonReadyConditionReason(conditions []camelkv1alpha1.KameletCondition) string {
	for _, condition := range conditions {
		if condition.Type == camelkv1alpha1.KameletConditionReady {
			if condition.Status == corev1.ConditionTrue {
				return ""
			}
			if condition.Message != "" {
				return fmt.Sprintf("%s : %s", condition.Reason, condition.Message)
			}
			return condition.Reason
		}
	}
	return "<unknown>"
}
