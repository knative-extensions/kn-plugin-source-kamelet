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
	knative.dev/client v0.27.1-0.20211110021545-e3022e662f29
	knative.dev/eventing v0.27.1-0.20211109172546-31e4e2eb26d1
	knative.dev/hack v0.0.0-20211108170701-96aac1c30be3
	knative.dev/pkg v0.0.0-20211109100843-91d1932616a7
	knative.dev/serving v0.27.1-0.20211110154546-6c3394bb4d5a
)
