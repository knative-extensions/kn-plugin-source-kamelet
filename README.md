# kn-plugin-source-kamelet

**[This component is BETA](<https://github.com/knative/community/tree/main/mechanics/MATURITY-LEVELS.md>)**

`kn-source-kamelet` Knative eventing sources plugin manages Kamelet
event sources on your cluster.

# Description

With this plugin, you can list available [Kamelets](<https://github.com/apache/camel-kamelets>) and bindings on
your cluster. Kamelets can act as Knative eventing sources where each
binding connects a Kamelet source to a Knative sink (broker, channel,
service).

# Usage

    Plugin manages Kamelets and KameletBindings as Knative eventing sources.

    Usage:
      kn-source-kamelet [command]

    Available Commands:
      bind          Create Kamelet bindings and bind source to Knative broker, channel or service.
      binding       Configure and manage a Kamelet binding.
      completion    generate the autocompletion script for the specified shell
      describe      Show details of given Kamelet source type
      help          Help about any command
      list          List available Kamelet source types
      version       Prints the plugin version

    Flags:
      -h, --help   help for kn-source-kamelet

    Use "kn-source-kamelet [command] --help" for more information about a command.

# Commands

## `list`

    List available Kamelet source types

    Usage:
      kn-source-kamelet list [flags]

    Aliases:
      list, ls

    Examples:

      # List available Kamelets
      kn-source-kamelet list

      # List available Kamelets in YAML output format
      kn-source-kamelet list -o yaml

    Flags:
      -A, --all-namespaces                If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.
          --allow-missing-template-keys   If true, ignore any errors in templates when a field or map key is missing in the template. Only applies to golang and jsonpath output formats. (default true)
      -h, --help                          help for list
      -n, --namespace string              Specify the namespace to operate in.
          --no-headers                    When using the default output format, don't print headers (default: print headers).
      -o, --output string                 Output format. One of: json|yaml|name|go-template|go-template-file|template|templatefile|jsonpath|jsonpath-as-json|jsonpath-file.
          --show-managed-fields           If true, keep the managedFields when printing objects in JSON or YAML format.
          --template string               Template string or path to template file to use when -o=go-template, -o=go-template-file. The template format is golang templates [http://golang.org/pkg/text/template/#pkg-overview].

## `describe`

    Show details of given Kamelet source type

    Usage:
      kn-source-kamelet describe NAME [flags]

    Aliases:
      describe

    Examples:

      # Describe given Kamelets
      kn-source-kamelet describe NAME

      # Describe given Kamelets in YAML output format
      kn-source-kamelet describe NAME -o yaml

    Flags:
      -h, --help                          help for describe
      -n, --namespace string              Specify the namespace to operate in.
      -o, --output string                 Output format. One of: json|yaml|name|url.
      -v, --verbose                       More output.

## `binding`

    Configure and manage a Kamelet binding.

    Usage:
      kn-source-kamelet binding [command]

    Examples:

      # Configure and manage a Kamelet binding.
      kn-source-kamelet binding create|update|delete

    Available Commands:
      create      Create Kamelet bindings and bind source to Knative broker, channel or service.
      delete      Delete Kamelet binding by its name.
      list        List Kamelet bindings.

    Flags:
      -h, --help   help for binding

    Use "kn-source-kamelet binding [command] --help" for more information about a command.

### `binding create`

    Create Kamelet bindings and bind source to Knative broker, channel or service.

    Usage:
      kn-source-kamelet binding create NAME [flags]

    Examples:

      # Create Kamelet binding with source and sink.
      kn-source-kamelet binding create NAME

      # Add a binding properties
      kn-source-kamelet binding create NAME --kamelet=name --sink|broker|channel|service=<name> --property=<key>=<value>

    Flags:
          --broker string                 Uses a broker as binding sink.
          --channel string                Uses a channel as binding sink.
      -h, --help                          help for create
          --force bool                    Apply the changes even if the binding already exists.
          --kamelet string                Kamelet source.
      -n, --namespace string              Specify the namespace to operate in.
          --service string                Uses a Knative service as binding sink.
      -s  --sink string                   Sink expression to define the binding sink.
          --property stringArray          Add a source property in the form of "<key>=<value>"
          --ce-override stringArray       Customize cloud events property in the form of "<key>=<value>"
          --ce-spec string                Customize cloud events spec version provided to the binding sink.
          --ce-type string                Customize cloud events type provided to the binding sink.

### `binding delete`

    Delete Kamelet binding by its name.

    Usage:
      kn-source-kamelet binding delete NAME [flags]

    Examples:

      # Delete Kamelet binding by its name.
      kn-source-kamelet binding delete NAME

    Flags:
      -h, --help                          help for create
      -n, --namespace string              Specify the namespace to operate in.

### `binding list`

    List Kamelet bindings.

    Usage:
      kn-source-kamelet binding list [flags]

    Aliases:
      list, ls

    Examples:

      # List Kamelet bindings.
      kn source kamelet binding list

      # List available Kamelet bindings in YAML output format
      kn source kamelet binding list -o yaml

    Flags:
      -A, --all-namespaces                If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.
          --allow-missing-template-keys   If true, ignore any errors in templates when a field or map key is missing in the template. Only applies to golang and jsonpath output formats. (default true)
      -h, --help                          help for list
      -n, --namespace string              Specify the namespace to operate in.
          --no-headers                    When using the default output format, don't print headers (default: print headers).
      -o, --output string                 Output format. One of: json|yaml|name|go-template|go-template-file|template|templatefile|jsonpath|jsonpath-as-json|jsonpath-file.
          --show-managed-fields           If true, keep the managedFields when printing objects in JSON or YAML format.
          --template string               Template string or path to template file to use when -o=go-template, -o=go-template-file. The template format is golang templates [http://golang.org/pkg/text/template/#pkg-overview].

## `bind`

Shortcut version of `kn-source-kamelet binding create` with Kamelet
source as positional argument. The shortcut command auto generates a
binding name in case no explicit name is given as command option
`--name`.

    Create Kamelet bindings and bind source to Knative broker, channel or service.

    Usage:
      kn-source-kamelet bind SOURCE [flags]

    Examples:

      # Bind Kamelets to a Knative sink
      kn-source-kamelet bind SOURCE

      # Add a binding properties
      kn-source-kamelet bind SOURCE --sink|broker|channel|service=<name> --property=<key>=<value>

    Flags:
          --broker string                 Uses a broker as binding sink.
          --channel string                Uses a channel as binding sink.
      -h, --help                          help for bind
          --force bool                    Apply the changes even if the binding already exists.
          --name string                   Binding name.
      -n, --namespace string              Specify the namespace to operate in.
          --service string                Uses a Knative service as binding sink.
      -s  --sink string                   Sink expression to define the binding sink.
          --property stringArray          Add a source property in the form of "<key>=<value>"
          --ce-override stringArray       Customize cloud events property in the form of "<key>=<value>"
          --ce-spec string                Customize cloud events spec version provided to the binding sink.
          --ce-type string                Customize cloud events type provided to the binding sink.

## `version`

This command prints out the version of this plugin and all extra
information which might help, for example when creating bug reports.

    Prints the plugin version

    Usage:
      kn-source-kamelet version [flags]

    Flags:
      -h, --help   help for version

# Examples

## List available Kamelet sources

You want to list all available Kamelets on your cluster. In this case,
you can use the `kn-source-kamelet list` command.

    $ kn-source-kamelet list

    Kamelet_1
    Kamelet_2
    Kamelet_3

## Print out the version of this plugin

The `kn-source-kamelet version` command helps you to identify the
version of this plugin.

    $ kn-source-kamelet version

    Version:      v20200402-local-a099aaf-dirty
    Build Date:   2020-04-02 18:16:20
    Git Revision: a099aaf

As you can see it prints out the version, (or a generated timestamp when
this plugin is built from a non-released commit) the date when the
plugin has been built and the actual Git revision.
