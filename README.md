[![Go Reference](https://pkg.go.dev/badge/github.com/TibebeJs/go-alive.svg)](https://pkg.go.dev/github.com/TibebeJs/go-alive) [![codecov](https://codecov.io/gh/TibebeJS/go-alive/branch/main/graph/badge.svg?token=k3AHKhTtqO)](https://codecov.io/gh/TibebeJS/go-alive) ![Go Report Card](https://goreportcard.com/badge/github.com/tibebejs/go-alive)

Robust services healthiness probing written in Go. (with notification support of webhook, telegram and more)

> :warning: WARNING: Under heavy construction. API may have breaking changes frequently.

## Getting Started ##

Please follow the steps below to get started quick

## Using as a CLI tool (using the go binary directly) ##

If you are using Go 1.17+, run the following:
```
$ go install github.com/TibebeJs/go-alive@latest
```

If you are using an older version of golang,
```
$ GO111MODULE=on go get github.com/TibebeJs/go-alive@latest
```
## Using Docker ##

### Alternative 1: Building the Image locally
Clone the source code repository first:

```
$ git@github.com:TibebeJS/go-alive.git && cd go-alive
```

Then build the Docker Image (with [Docker BuildKit](https://docs.docker.com/develop/develop-images/build_enhancements/#to-enable-buildkit-builds) enabled):

```
$ DOCKER_BUILDKIT=1 docker build . -t go-alive
```

Finally mount a folder where your `config.yml` file resides at as `/config` and run the image:
```
$ docker run -v $(pwd):/config go-alive
```

## Bugs ##

Bugs or suggestions? Visit the [issue tracker](https://github.com/TibebeJS/go-alive/issues) 

## Contribution

Feel free to fork, edit and send a PR.
