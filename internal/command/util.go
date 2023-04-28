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
	"log"
	"regexp"
	"unicode"

	"github.com/apache/camel-k/pkg/client/camel/clientset/versioned/scheme"
	"knative.dev/client-pkg/pkg/util"

	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
)

var (
	disallowedChars = regexp.MustCompile(`[^a-z0-9-]`)
	sinkExpression  = regexp.MustCompile(`^(?:(?P<apiVersion>(?:[a-z0-9-.]+/)?[a-z0-9-.]+):)?(?P<kind>[A-Za-z0-9-.]+):(?:(?P<namespace>[a-z0-9-.]+)/)?(?P<name>[a-z0-9-.]+)(?:$|[?].*$)`)
)

func isEventSourceType(kamelet *v1alpha1.Kamelet) bool {
	return kamelet.Labels[KameletTypeLabel] == "source"
}

func extractKameletProvider(kamelet *v1alpha1.Kamelet) string {
	return kamelet.Annotations[KameletProviderAnnotation]
}

func extractKameletSupportLevel(kamelet *v1alpha1.Kamelet) string {
	return kamelet.Annotations[KameletSupportLevelAnnotation]
}

func isDisallowedStartEndChar(rune rune) bool {
	return !unicode.IsLetter(rune) && !unicode.IsNumber(rune)
}

func updateKameletListGvk(list *v1alpha1.KameletList) {
	err := util.UpdateGroupVersionKindWithScheme(list, v1alpha1.SchemeGroupVersion, scheme.Scheme)
	if err != nil {
		log.Fatalf("Internal error: %v", err)
	}

	for idx := range list.Items {
		updateKameletGvk(&list.Items[idx])
	}
}

func updateKameletGvk(kamelet *v1alpha1.Kamelet) {
	_ = util.UpdateGroupVersionKindWithScheme(kamelet, v1alpha1.SchemeGroupVersion, scheme.Scheme)
}

func updateKameletBindingListGvk(list *v1alpha1.KameletBindingList) {
	_ = util.UpdateGroupVersionKindWithScheme(list, v1alpha1.SchemeGroupVersion, scheme.Scheme)

	for i := range list.Items {
		updateKameletBindingGvk(&list.Items[i])
	}
}

func updateKameletBindingGvk(kb *v1alpha1.KameletBinding) {
	_ = util.UpdateGroupVersionKindWithScheme(kb, v1alpha1.SchemeGroupVersion, scheme.Scheme)
}
