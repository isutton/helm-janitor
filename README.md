# helm-janitor

`helm-janitor` is a Helm plugin that remove failed releases revisions from the cluster.

## Installing from sources

To compile and install `helm-janitor` from sources, perform the following commands:

```text
# creates the environment variables Helm provides to plugins to
# properly install in the host system
$ eval $(helm env)
$ echo $HELM_PLUGINS
/home/isuttonl/.local/share/helm/plugins

# compiles and install the plugin in HELM_PLUGINS directory
$ make install
mkdir -p ./dist
go build -o ./dist/helm-janitor main.go
mkdir -p ./dist/janitor
cp ./dist/helm-janitor ./dist/janitor
cp ./plugin.yaml ./dist/janitor
mkdir -p /home/isuttonl/.local/share/helm/plugins/janitor
install ./dist/janitor/* /home/isuttonl/.local/share/helm/plugins/janitor/
```

Once this is finished, the plugin should be available:

```text
$ helm janitor
A Helm plugin that remove failed releases revisions from the cluster

Usage:
  helm-janitor [command]

Available Commands:
  clean       remove unused artifacts of previous failed releases
  completion  generate the autocompletion script for the specified shell
  help        Help about any command

Flags:
      --config string   config file (default is $HOME/.helm_janitor.yaml)
  -h, --help            help for helm-janitor

Use "helm-janitor [command] --help" for more information about a command.
```

## Building `helm-janitor`

To build only the plugin program, use the `build` target:

```text
make build
```

# License

MIT
