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
	knative.dev/client v0.26.1-0.20211028082427-f027b38e200a
	knative.dev/eventing v0.26.1-0.20211028192027-b498c7fd6eb7
	knative.dev/hack v0.0.0-20211028194650-b96d65a5ff5e
	knative.dev/pkg v0.0.0-20211027105800-3b33e02e5b9c
	knative.dev/serving v0.26.1-0.20211028155847-785c55ae7c0d
)
