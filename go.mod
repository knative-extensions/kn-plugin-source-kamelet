module knative.dev/kn-plugin-source-kamelet

go 1.15

require (
	github.com/apache/camel-k/pkg/apis/camel v1.3.1
	github.com/apache/camel-k/pkg/client/camel v1.3.1
	github.com/spf13/cobra v1.1.3
	gotest.tools/v3 v3.0.3
	k8s.io/apimachinery v0.19.7
	k8s.io/client-go v0.19.7
	knative.dev/client v0.21.1-0.20210319143953-43b669fe9fb6
	knative.dev/hack v0.0.0-20210325223819-b6ab329907d3
)

replace github.com/go-openapi/spec => github.com/go-openapi/spec v0.19.3
