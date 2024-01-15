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
	"testing"

	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestUpdateKameletGvk(t *testing.T) {
	k := v1alpha1.Kamelet{}
	assert.True(t, k.GroupVersionKind().Empty())
	updateKameletGvk(&k)
	verifyGvk(t, "Kamelet", &k)
}

func TestUpdateKameletListGvk(t *testing.T) {
	kl := v1alpha1.KameletList{}
	kl.Items = []v1alpha1.Kamelet{
		{},
	}
	assert.True(t, kl.GroupVersionKind().Empty())
	assert.True(t, kl.Items[0].GroupVersionKind().Empty())
	updateKameletListGvk(&kl)
	verifyGvk(t, "KameletList", &kl)
	verifyGvk(t, "Kamelet", &kl.Items[0])
}

func TestUpdateKameletBindingGvk(t *testing.T) {
	kb := v1alpha1.KameletBinding{}
	assert.True(t, kb.GroupVersionKind().Empty())
	updatePipeGvk(&kb)
	verifyGvk(t, "KameletBinding", &kb)
}

func TestUpdateKameletBindingListGvk(t *testing.T) {
	kl := v1alpha1.KameletBindingList{}
	kl.Items = []v1alpha1.KameletBinding{
		{},
	}
	assert.True(t, kl.GroupVersionKind().Empty())
	assert.True(t, kl.Items[0].GroupVersionKind().Empty())
	updatePipeListGvk(&kl)
	verifyGvk(t, "KameletBindingList", &kl)
	verifyGvk(t, "KameletBinding", &kl.Items[0])
}

func verifyGvk(t *testing.T, kind string, o runtime.Object) {
	assert.Equal(t, v1alpha1.SchemeGroupVersion, o.GetObjectKind().GroupVersionKind().GroupVersion())
	assert.Equal(t, kind, o.GetObjectKind().GroupVersionKind().Kind)
}
