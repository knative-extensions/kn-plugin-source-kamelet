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
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	camelv1 "github.com/apache/camel-k/pkg/apis/camel/v1"
	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	knerrors "knative.dev/client/pkg/errors"
	"knative.dev/client/pkg/kn/commands"
)

var bindExample = `
  # Bind Kamelets in an integration flow
  kn-source-kamelet bind SOURCE

  # Add a binding properties
  kn-source-kamelet bind SOURCE --sink|broker|channel|service=<name> --source-property=<key>=<value> --sink-property=<key>=<value>`

// NewBindCommand implements 'kn-source-kamelet bind' command
func NewBindCommand(p *KameletPluginParams) *cobra.Command {
	printFlags := genericclioptions.NewPrintFlags("")

	var sourceProperties []string
	var sinkProperties []string
	var sink string
	var broker string
	var channel string
	var service string
	cmd := &cobra.Command{
		Use:     "bind",
		Short:   "Create KameletBindings and bind Kamelet source to Knative broker, channel or service.",
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

			kamelet, err := client.Kamelets(namespace).Get(p.Context, source, v1.GetOptions{})
			if err != nil {
				return knerrors.GetError(err)
			}

			if !isEventSourceType(kamelet) {
				return fmt.Errorf("Kamelet %s is not an event source", source)
			}

			sourceProps, err := parseProperties(sourceProperties)
			if err != nil {
				return knerrors.GetError(err)
			}
			sourceEndpoint := v1alpha1.Endpoint{
				Properties: &sourceProps,
				Ref: &corev1.ObjectReference{
					Kind:       v1alpha1.KameletKind,
					APIVersion: v1alpha1.SchemeGroupVersion.String(),
					Name:       kamelet.Name,
					Namespace:  kamelet.Namespace,
				},
			}

			if err := verifyProperties(kamelet, sourceEndpoint); err != nil {
				return knerrors.GetError(err)
			}

			var sinkRef corev1.ObjectReference
			if sink != "" {
				sinkRef, err = decodeSink(sink)
			} else if broker != "" {
				sinkRef, err = decodeSink("broker:" + broker)
			} else if channel != "" {
				sinkRef, err = decodeSink("channel:" + channel)
			} else if service != "" {
				sinkRef, err = decodeSink("service:" + service)
			} else {
				err = fmt.Errorf("missing sink for binding - please use one of --sink, --broker, --channel, --service")
			}

			if err != nil {
				return knerrors.GetError(err)
			}

			if sinkRef.Namespace == "" {
				sinkRef.Namespace = namespace
			}

			sinkProps, err := parseProperties(sinkProperties)
			if err != nil {
				return knerrors.GetError(err)
			}
			sinkEndpoint := v1alpha1.Endpoint{
				Properties: &sinkProps,
				Ref:        &sinkRef,
			}

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return knerrors.GetError(err)
			}

			name = nameFor(name, source, sinkRef)

			binding := v1alpha1.KameletBinding{
				ObjectMeta: v1.ObjectMeta{
					Namespace: namespace,
					Name:      name,
				},
				Spec: v1alpha1.KameletBindingSpec{
					Source: sourceEndpoint,
					Sink:   sinkEndpoint,
				},
			}

			existed := false
			_, err = client.KameletBindings(namespace).Create(p.Context, &binding, v1.CreateOptions{})
			if err != nil && k8serrors.IsAlreadyExists(err) {
				existed = true

				existing, err := client.KameletBindings(namespace).Get(p.Context, binding.Name, v1.GetOptions{})
				if err != nil {
					return knerrors.GetError(err)
				}
				// Update the custom resource
				binding.ResourceVersion = existing.ResourceVersion
				_, err = client.KameletBindings(namespace).Update(p.Context, &binding, v1.UpdateOptions{})
				if err != nil {
					return knerrors.GetError(err)
				}
			}

			if !existed {
				fmt.Printf("kamelet binding \"%s\" created\n", name)
			} else {
				fmt.Printf("kamelet binding \"%s\" updated\n", name)
			}

			out := cmd.OutOrStdout()
			if printFlags.OutputFlagSpecified() {
				printer, err := printFlags.ToPrinter()
				if err != nil {
					return err
				}
				return printer.PrintObj(&binding, out)
			}

			return nil
		},
	}
	flags := cmd.Flags()
	commands.AddNamespaceFlags(flags, false)

	flags.String("name", "", "Binding name.")
	flags.StringVar(&sink, "sink", "", "Sink expression to define the binding sink.")
	flags.StringVar(&broker, "broker", "", "Uses a broker as binding sink.")
	flags.StringVar(&channel, "channel", "", "Uses a channel as binding sink.")
	flags.StringVar(&service, "service", "", "Uses a Knative service as binding sink.")
	flags.StringArrayVar(&sourceProperties, "source-property", nil, `Add a source property in the form of "<key>=<value>"`)
	flags.StringArrayVar(&sinkProperties, "sink-property", nil, `Add a sink property in the form of "<key>=<value>"`)

	printFlags.AddFlags(cmd)
	cmd.Flag("output").Usage = fmt.Sprintf("Output format. One of: %s.", strings.Join(append(printFlags.AllowedFormats(), "url"), "|"))
	return cmd
}

func nameFor(name, source string, sinkRef corev1.ObjectReference) string {
	if name != "" {
		return name
	}

	generated := fmt.Sprintf("%s-to-%s-%s", source, sinkRef.Kind, sinkRef.Name)

	generated = filepath.Base(generated)
	generated = strings.Split(generated, ".")[0]
	generated = strings.ToLower(generated)
	generated = disallowedChars.ReplaceAllString(generated, "")
	generated = strings.TrimFunc(generated, isDisallowedStartEndChar)

	return generated
}

func decodeSink(sink string) (corev1.ObjectReference, error) {
	ref := corev1.ObjectReference{}

	if sinkExpression.MatchString(sink) {
		groupNames := sinkExpression.SubexpNames()
		for _, match := range sinkExpression.FindAllStringSubmatch(sink, -1) {
			for idx, text := range match {
				groupName := groupNames[idx]
				switch groupName {
				case "apiVersion":
					ref.APIVersion = text
				case "namespace":
					ref.Namespace = text
				case "kind":
					ref.Kind = text
				case "name":
					ref.Name = text
				}
			}
		}

		if sinkType, ok := sinkTypes[ref.Kind]; ok {
			if sinkType.Kind != "" {
				ref.Kind = sinkType.Kind
			}
			if ref.APIVersion == "" && sinkType.APIVersion != "" {
				ref.APIVersion = sinkType.APIVersion
			}
		} else {
			return ref, fmt.Errorf("unsupported sink type %q", ref.Kind)
		}
	} else {
		return ref, fmt.Errorf("unsupported sink expression %q - please use format <kind>:<name>", sink)
	}

	return ref, nil
}

func verifyProperties(kamelet *v1alpha1.Kamelet, endpoint v1alpha1.Endpoint) error {
	if kamelet.Spec.Definition != nil && len(kamelet.Spec.Definition.Required) > 0 {
		pMap, err := endpoint.Properties.GetPropertyMap()
		if err != nil {
			return err
		}
		for _, reqProp := range kamelet.Spec.Definition.Required {
			found := false
			if endpoint.Properties != nil {
				if _, contains := pMap[reqProp]; contains {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("binding is missing required property %q for Kamelet %q", reqProp, kamelet.Name)
			}
		}
	}

	return nil
}

func parseProperties(properties []string) (v1alpha1.EndpointProperties, error) {
	props := make(map[string]string)
	for _, p := range properties {
		key, value, err := parseProperty(p)
		if err != nil {
			continue
		}
		props[key] = value
	}
	return asEndpointProperties(props)
}

func asEndpointProperties(props map[string]string) (v1alpha1.EndpointProperties, error) {
	if len(props) == 0 {
		return v1alpha1.EndpointProperties{}, nil
	}
	data, err := json.Marshal(props)
	if err != nil {
		return v1alpha1.EndpointProperties{}, err
	}
	return v1alpha1.EndpointProperties{
		RawMessage: camelv1.RawMessage(data),
	}, nil
}

func parseProperty(prop string) (string, string, error) {
	parts := strings.SplitN(prop, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf(`property %q does not follow format "<key>=<value>"`, prop)
	}
	return parts[0], parts[1], nil
}
