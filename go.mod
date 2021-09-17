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
	knative.dev/client v0.25.1-0.20210913155632-82a21a5773be
	knative.dev/eventing v0.25.1-0.20210916210831-9e66fc9bdb12
	knative.dev/hack v0.0.0-20210806075220-815cd312d65c
	knative.dev/pkg v0.0.0-20210914164111-4857ab6939e3
	knative.dev/serving v0.25.1-0.20210916175739-5792a5c1933d
)

replace github.com/go-openapi/spec => github.com/go-openapi/spec v0.19.3
