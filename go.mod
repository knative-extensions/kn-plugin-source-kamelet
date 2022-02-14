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
	knative.dev/client v0.29.1-0.20220204171521-6690a20e8f56
	knative.dev/eventing v0.29.1-0.20220209143041-f13248e5a7de
	knative.dev/hack v0.0.0-20220209225905-7331bb16ba00
	knative.dev/pkg v0.0.0-20220210201907-fc93ac76d0b6
	knative.dev/serving v0.29.1-0.20220212025718-98f70b35a8ff
)
