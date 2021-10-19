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
	knative.dev/client v0.26.1-0.20211018155636-a736af7129a1
	knative.dev/eventing v0.26.1-0.20211018174236-a34aaa09f7d2
	knative.dev/hack v0.0.0-20211018110626-47ac3b032e60
	knative.dev/pkg v0.0.0-20211018141937-a34efd6b409d
	knative.dev/serving v0.26.1-0.20211018142437-2dc26f102ade
)
