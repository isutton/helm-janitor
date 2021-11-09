all: build

build:
	mkdir -p $(PWD)/dist
	go build -o $(PWD)/dist/helm-janitor main.go

clean:
	rm -fr $(PWD)/dist
