module knative.dev/kn-plugin-source-kamelet

go 1.16

require (
	github.com/apache/camel-k/pkg/apis/camel v1.3.1
	github.com/apache/camel-k/pkg/client/camel v1.3.1
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.21.4
	k8s.io/cli-runtime v0.21.4
	k8s.io/client-go v0.21.4
	knative.dev/client v0.26.1-0.20211014105642-b1e4132be8b5
	knative.dev/eventing v0.26.1-0.20211014072442-a6a819dc71cf
	knative.dev/hack v0.0.0-20211015200324-86876688e735
	knative.dev/pkg v0.0.0-20211015194524-a5bb75923981
	knative.dev/serving v0.26.1-0.20211016013324-e5d8560f950c
)
