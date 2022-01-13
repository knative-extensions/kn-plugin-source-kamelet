module knative.dev/kn-plugin-source-kamelet

go 1.16

require (
	github.com/apache/camel-k/pkg/apis/camel v1.3.1
	github.com/apache/camel-k/pkg/client/camel v1.3.1
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.22.5
	k8s.io/apimachinery v0.22.5
	k8s.io/cli-runtime v0.21.4
	k8s.io/client-go v0.22.5
	knative.dev/client v0.28.1-0.20220112141951-d8670f576217
	knative.dev/eventing v0.28.1-0.20220112214912-de8918823b0e
	knative.dev/hack v0.0.0-20220111151514-59b0cf17578e
	knative.dev/pkg v0.0.0-20220112181951-2b23ad111bc2
	knative.dev/serving v0.28.1-0.20220112163651-8beafe7f1683
)
