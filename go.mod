module knative.dev/kn-plugin-source-kamelet

go 1.15

require (
	github.com/apache/camel-k/pkg/apis/camel v1.3.1
	github.com/apache/camel-k/pkg/client/camel v1.3.1
	github.com/spf13/cobra v1.1.3
	gotest.tools/v3 v3.0.3
	k8s.io/api v0.19.7
	k8s.io/apimachinery v0.19.7
	k8s.io/cli-runtime v0.19.7
	k8s.io/client-go v0.19.7
	knative.dev/client v0.22.1-0.20210428162854-dccf3e30fa14
	knative.dev/hack v0.0.0-20210428122153-93ad9129c268
	knative.dev/pkg v0.0.0-20210428141353-878c85083565
)

replace github.com/go-openapi/spec => github.com/go-openapi/spec v0.19.3
