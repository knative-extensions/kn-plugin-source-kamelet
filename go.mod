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
	k8s.io/cli-runtime v0.22.5
	k8s.io/client-go v0.22.5
	knative.dev/client v0.29.1-0.20220308195105-c2fe56c83038
	knative.dev/eventing v0.30.0
	knative.dev/hack v0.0.0-20220224013837-e1785985d364
	knative.dev/pkg v0.0.0-20220301181942-2fdd5f232e77
	knative.dev/serving v0.30.0
)
