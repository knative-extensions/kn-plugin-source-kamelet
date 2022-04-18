module knative.dev/kn-plugin-source-kamelet

go 1.16

require (
	github.com/apache/camel-k/pkg/apis/camel v1.3.1
	github.com/apache/camel-k/pkg/client/camel v1.3.1
	github.com/spf13/cobra v1.3.0
	github.com/stretchr/testify v1.7.0
	gotest.tools/v3 v3.1.0
	k8s.io/api v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/cli-runtime v0.23.4
	k8s.io/client-go v0.23.5
	knative.dev/client v0.30.2-0.20220414141510-76f17f686f4a
	knative.dev/eventing v0.30.1-0.20220415141711-ff55a456c3f9
	knative.dev/hack v0.0.0-20220411131823-6ffd8417de7c
	knative.dev/pkg v0.0.0-20220412134708-e325df66cb51
	knative.dev/serving v0.30.1-0.20220416140111-2e5ca679a71e
)
