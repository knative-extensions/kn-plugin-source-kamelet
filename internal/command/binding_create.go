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
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	camelv1 "github.com/apache/camel-k/pkg/apis/camel/v1"
	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	camelkv1alpha1 "github.com/apache/camel-k/pkg/client/camel/clientset/versioned/typed/camel/v1alpha1"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	knerrors "knative.dev/client/pkg/errors"
	"knative.dev/client/pkg/kn/commands"
)

var bindingCreateExample = `
  # Create Kamelet binding with source and sink.
  kn-source-kamelet binding create NAME

  # Add a binding properties
  kn-source-kamelet binding create NAME --kamelet=name --sink|broker|channel|service=<name> --property=<key>=<value>`

// newBindingCreateCommand implements 'kn-source-kamelet bind' command
func newBindingCreateCommand(p *KameletPluginParams) *cobra.Command {
	var properties []string
	var source string
	var sink string
	var broker string
	var channel string
	var service string
	var force bool
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create Kamelet bindings and bind source to Knative broker, channel or service.",
		Example: bindingCreateExample,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if len(args) != 1 {
				return errors.New("'kn-source-kamelet binding create' requires the binding name as argument")
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

			options := CreateBindingOptions{
				Name:             name,
				Source:           source,
				SourceProperties: properties,
				Sink:             sink,
				Broker:           broker,
				Channel:          channel,
				Service:          service,
				Force:            force,
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

	flags.StringVar(&source, "kamelet", "", "Kamelet source.")
	flags.StringVarP(&sink, "sink", "s", "", "Sink expression to define the binding sink.")
	flags.StringVar(&broker, "broker", "", "Uses a broker as binding sink.")
	flags.StringVar(&channel, "channel", "", "Uses a channel as binding sink.")
	flags.StringVar(&service, "service", "", "Uses a Knative service as binding sink.")
	flags.BoolVar(&force, "force", false, "Apply the changes even if the binding already exists.")
	flags.StringArrayVar(&properties, "property", nil, `Add a source property in the form of "<key>=<value>"`)
	return cmd
}

func createBinding(client camelkv1alpha1.CamelV1alpha1Interface, ctx context.Context, namespace string, options CreateBindingOptions) error {
	kamelet, err := client.Kamelets(namespace).Get(ctx, options.Source, v1.GetOptions{})
	if err != nil {
		return knerrors.GetError(err)
	}

	if !isEventSourceType(kamelet) {
		return fmt.Errorf("Kamelet %s is not an event source", options.Source)
	}

	sourceProps, err := parseProperties(options.SourceProperties)
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
	if options.Sink != "" {
		sinkRef, err = decodeSink(options.Sink)
	} else if options.Broker != "" {
		sinkRef, err = decodeSink("broker:" + options.Broker)
	} else if options.Channel != "" {
		sinkRef, err = decodeSink("channel:" + options.Channel)
	} else if options.Service != "" {
		sinkRef, err = decodeSink("ksvc:" + options.Service)
	} else {
		err = fmt.Errorf("missing sink for binding - please use one of --sink, --broker, --channel, --service")
	}

	if err != nil {
		return knerrors.GetError(err)
	}

	if sinkRef.Namespace == "" {
		sinkRef.Namespace = namespace
	}

	sinkEndpoint := v1alpha1.Endpoint{
		Ref: &sinkRef,
	}

	name := nameFor(options.Name, options.Source, sinkRef)

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
	_, err = client.KameletBindings(namespace).Create(ctx, &binding, v1.CreateOptions{})
	if err != nil && k8serrors.IsAlreadyExists(err) {
		if options.Force {
			existed = true

			existing, err := client.KameletBindings(namespace).Get(ctx, binding.Name, v1.GetOptions{})
			if err != nil {
				return knerrors.GetError(err)
			}
			// Update the custom resource
			binding.ResourceVersion = existing.ResourceVersion
			_, err = client.KameletBindings(namespace).Update(ctx, &binding, v1.UpdateOptions{})
			if err != nil {
				return knerrors.GetError(err)
			}
		} else {
			return fmt.Errorf("kamelet binding with name %q already exists. Use --force to recreate the binding", binding.Name)
		}
	}

	if existed {
		fmt.Printf("kamelet binding %q updated\n", name)
	} else {
		fmt.Printf("kamelet binding %q created\n", name)
	}

	return nil
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
	pMap, err := endpoint.Properties.GetPropertyMap()

	if kamelet.Spec.Definition != nil && len(kamelet.Spec.Definition.Required) > 0 {
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

	for propName := range pMap {
		if _, ok := kamelet.Spec.Definition.Properties[propName]; !ok {
			return fmt.Errorf("binding uses unknown property %q for Kamelet %q", propName, kamelet.Name)
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
