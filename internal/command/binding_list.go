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

	camelkv1alpha1 "github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1beta1 "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/client/pkg/kn/commands"
	"knative.dev/client/pkg/kn/commands/flags"
	hprinters "knative.dev/client/pkg/printers"
)

var bindingListExample = `
  # List Kamelet bindings.
  kn source kamelet binding list

  # List available Kamelet bindings in YAML output format
  kn source kamelet binding list -o yaml`

// newBindingCreateCommand implements 'kn-source-kamelet binding list' command
func newBindingListCommand(p *KameletPluginParams) *cobra.Command {
	listFlags := flags.NewListPrintFlags(ListBindingHandlers)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List Kamelet bindings.",
		Aliases: []string{"ls"},
		Example: bindingListExample,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			namespace, err := p.GetNamespace(cmd)
			if err != nil {
				return err
			}

			kameletClient, err := p.NewKameletClient()
			if err != nil {
				return err
			}

			bindingList, err := kameletClient.KameletBindings(namespace).List(p.Context, v1.ListOptions{})
			if err != nil {
				return err
			}
			if len(bindingList.Items) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "No resources found.\n")
				return nil
			}
			updateKameletBindingListGvk(bindingList)

			// empty namespace indicates all-namespaces flag is specified
			if namespace == "" {
				listFlags.EnsureWithNamespace()
			}

			err = listFlags.Print(bindingList, cmd.OutOrStdout())
			if err != nil {
				return err
			}

			return nil
		},
	}
	commands.AddNamespaceFlags(cmd.Flags(), true)
	listFlags.AddFlags(cmd)
	return cmd
}

// ListBindingHandlers handles printing human readable table for `kn-source-kamelet binding list` command's output
func ListBindingHandlers(h hprinters.PrintHandler) {
	columnDefinitions := []metav1beta1.TableColumnDefinition{
		{Name: "Namespace", Type: "string", Description: "Namespace of the Kamelet binding", Priority: 0},
		{Name: "Name", Type: "string", Description: "Name of the Kamelet binding ", Priority: 1},
		{Name: "Phase", Type: "string", Description: "Phase of the Kamelet binding ", Priority: 1},
		{Name: "Age", Type: "string", Description: "Age of the Kamelet binding ", Priority: 1},
		{Name: "Conditions", Type: "string", Description: "Ready state conditions", Priority: 1},
		{Name: "Ready", Type: "string", Description: "Ready state of the Kamelet binding ", Priority: 1},
		{Name: "Reason", Type: "string", Description: "Reason if state is not Ready", Priority: 1},
	}
	h.TableHandler(columnDefinitions, printBinding)
	h.TableHandler(columnDefinitions, printBindingList)
}

// printKameletList populates the Kamelet list table rows
func printBindingList(bindingList *camelkv1alpha1.KameletBindingList, options hprinters.PrintOptions) ([]metav1beta1.TableRow, error) {
	rows := make([]metav1beta1.TableRow, 0, len(bindingList.Items))

	for i := range bindingList.Items {
		binding := &bindingList.Items[i]
		r, err := printBinding(binding, options)
		if err != nil {
			return nil, err
		}
		rows = append(rows, r...)
	}
	return rows, nil
}

// printBinding populates the Kamelet binding table rows
func printBinding(binding *camelkv1alpha1.KameletBinding, options hprinters.PrintOptions) ([]metav1beta1.TableRow, error) {
	name := binding.Name
	phase := binding.Status.Phase
	age := commands.TranslateTimestampSince(binding.CreationTimestamp)
	conditions := bindingConditionsValue(binding.Status.Conditions)
	ready := bindingReadyCondition(binding.Status.Conditions)
	reason := bindingNonReadyConditionReason(binding.Status.Conditions)

	row := metav1beta1.TableRow{
		Object: runtime.RawExtension{Object: binding},
	}

	if options.AllNamespaces {
		row.Cells = append(row.Cells, binding.Namespace)
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

// bindingConditionsValue returns the True conditions count among total conditions
func bindingConditionsValue(conditions []camelkv1alpha1.KameletBindingCondition) string {
	var ok int
	for _, condition := range conditions {
		if condition.Status == "True" {
			ok++
		}
	}
	return fmt.Sprintf("%d OK / %d", ok, len(conditions))
}

// bindingReadyCondition returns status of resource's Ready type condition
func bindingReadyCondition(conditions []camelkv1alpha1.KameletBindingCondition) string {
	for _, condition := range conditions {
		if condition.Type == camelkv1alpha1.KameletBindingConditionReady {
			return string(condition.Status)
		}
	}
	return "<unknown>"
}

// bindingNonReadyConditionReason returns formatted string of
// reason and message for non ready conditions
func bindingNonReadyConditionReason(conditions []camelkv1alpha1.KameletBindingCondition) string {
	for _, condition := range conditions {
		if condition.Type == camelkv1alpha1.KameletBindingConditionReady {
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
