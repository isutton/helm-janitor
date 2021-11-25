HELM_PLUGINS ?= $(HOME)/.local/share/helm/plugins

all: build

build: ./dist/helm-janitor

./dist/helm-janitor:
	mkdir -p ./dist
	go build -o ./dist/helm-janitor main.go

./dist/janitor/plugin.yaml: ./dist/helm-janitor
	mkdir -p ./dist/janitor
	cp ./dist/helm-janitor ./dist/janitor
	cp ./plugin.yaml ./dist/janitor

plugin: ./dist/janitor/plugin.yaml

install: plugin
	mkdir -p $(HELM_PLUGINS)/janitor
	install ./dist/janitor/* $(HELM_PLUGINS)/janitor/

uninstall:
	rm -fr $(HELM_PLUGINS)/janitor

.PHONY: clean
clean:
	rm -fr ./dist

.PHONY: tidy
tidy:
	go mod tidy
