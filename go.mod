module knative.dev/kn-plugin-source-kamelet

go 1.16

require (
	github.com/apache/camel-k/pkg/apis/camel v1.3.1
	github.com/apache/camel-k/pkg/client/camel v1.3.1
	github.com/spf13/cobra v1.3.0
	github.com/stretchr/testify v1.7.0
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.22.5
	k8s.io/apimachinery v0.22.5
	k8s.io/cli-runtime v0.21.4
	k8s.io/client-go v0.22.5
	knative.dev/client v0.28.1-0.20220121192504-ed096b9669cb
	knative.dev/eventing v0.28.1-0.20220124132429-9e0e90c9d603
	knative.dev/hack v0.0.0-20220118141833-9b2ed8471e30
	knative.dev/pkg v0.0.0-20220118160532-77555ea48cd4
	knative.dev/serving v0.28.1-0.20220121202304-ff545005beab
)
