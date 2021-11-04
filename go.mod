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
	knative.dev/client v0.27.1-0.20211104101401-4fb6bdb95a9c
	knative.dev/eventing v0.27.1-0.20211104083501-4cc5ecf9635e
	knative.dev/hack v0.0.0-20211104075903-0f69979bbb7d
	knative.dev/pkg v0.0.0-20211104101302-51b9e7f161b4
	knative.dev/serving v0.27.1-0.20211104132102-416ded2f9a62
)
