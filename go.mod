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
	knative.dev/client v0.26.1-0.20211019150534-534d91319f7d
	knative.dev/eventing v0.26.1-0.20211019092333-7af98bbb4491
	knative.dev/hack v0.0.0-20211019034732-ced8ce706528
	knative.dev/pkg v0.0.0-20211019132235-ba2b2b1bf268
	knative.dev/serving v0.26.1-0.20211019142434-ab2fcc63551e
)
