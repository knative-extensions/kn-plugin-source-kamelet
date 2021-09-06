module knative.dev/kn-plugin-source-kamelet

go 1.16

require (
	github.com/apache/camel-k/pkg/apis/camel v1.3.1
	github.com/apache/camel-k/pkg/client/camel v1.3.1
	github.com/spf13/cobra v1.2.1
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.21.4
	k8s.io/cli-runtime v0.21.4
	k8s.io/client-go v0.21.4
	knative.dev/client v0.25.1-0.20210903081157-817068bd7f98
	knative.dev/hack v0.0.0-20210806075220-815cd312d65c
	knative.dev/pkg v0.0.0-20210906105443-bb433c98147d
)

replace github.com/go-openapi/spec => github.com/go-openapi/spec v0.19.3
