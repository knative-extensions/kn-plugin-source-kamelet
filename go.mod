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
	knative.dev/client v0.28.1-0.20220114130644-a96ecf55f6de
	knative.dev/eventing v0.28.1-0.20220116182528-68d441013ca8
	knative.dev/hack v0.0.0-20220111151514-59b0cf17578e
	knative.dev/pkg v0.0.0-20220114141842-0a429cba1c73
	knative.dev/serving v0.28.1-0.20220113203312-9073261f9b89
)
