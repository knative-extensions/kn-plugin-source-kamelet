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
	knative.dev/client v0.27.1-0.20211213030018-8b8b56581c63
	knative.dev/eventing v0.28.0
	knative.dev/hack v0.0.0-20211203062838-e11ac125e707
	knative.dev/pkg v0.0.0-20211206113427-18589ac7627e
	knative.dev/serving v0.28.0
)
