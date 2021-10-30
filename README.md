[![Go Reference](https://pkg.go.dev/badge/github.com/TibebeJs/go-alive.svg)](https://pkg.go.dev/github.com/TibebeJs/go-alive@v0.6.0) [![codecov](https://codecov.io/gh/TibebeJS/go-alive/branch/main/graph/badge.svg?token=k3AHKhTtqO)](https://codecov.io/gh/TibebeJS/go-alive) [![Go Report Card](https://goreportcard.com/badge/github.com/TibebeJS/go-alive)](https://goreportcard.com/report/github.com/TibebeJS/go-alive)

Robust services healthiness probing written in Go. (with notification support of telegram, slack, email and more)

> :warning: WARNING: Under heavy construction. API may have breaking changes frequently.

## Getting Started ##

Please follow the steps below to get started quick

## Using as a CLI tool (using the go binary directly) ##

If you are using Go 1.17+, run the following:
```bash
$ go install github.com/TibebeJs/go-alive@latest
```

If you are using an older version of golang,
```console
$ GO111MODULE=on go get github.com/TibebeJs/go-alive@latest
```
## Using Docker ##

### Alternative 1: Building the Image locally
Clone the source code repository first:

```console
$ git@github.com:TibebeJS/go-alive.git && cd go-alive
```

Then build the Docker Image (with [Docker BuildKit](https://docs.docker.com/develop/develop-images/build_enhancements/#to-enable-buildkit-builds) enabled):

```console
$ DOCKER_BUILDKIT=1 docker build . -t go-alive
```

Finally mount a folder where your `config.yml` file resides at as `/config` and run the image:
```console
$ docker run -v $(pwd):/config go-alive         # assuming config.yml is in the current working directory
```

### Alternative 2: Pull the image from the Docker Hub Registry and run it
To pull the docker image:

```console
$ docker pull tibebesjs/go-alive
```

Then simplyy mount a folder where you have your `config.yml` in as `/config` and run the image:
```console
$ docker run -v $(pwd):/config go-alive
```

## Configuration ##

Every operational aspects of go-alive is configured through the yaml file.

```yaml
targets:                                          # list of services to scan
  - name: "Test Server"  
    ip: "127.0.0.1"
    cron: "*/5 * * * *"                           # scan every 5 seconds
    strategy: status-code                         # can be "ping", "telnet" or "status-code"
    https: false
    ports:                                        # list of ports to scan for the specified host
      - port: 8000
        notify:                                   # notification channels for the result of the specific port scan
          - via: telegram
            chat: go-alive-test-group
            from: go-alive-test-bot               # defined in the notifications block
            template: ""
      - port: 8010
        notify:
    rules:                                        # conditional rules to check on host scan result
      - failures: ">0"                            # can be "<num", "num", ">num". eg. <4 (less than 4 failures)
        notify:
        - via: telegram
          chat: go-alive-test-group
          from: go-alive-test-bot
          template: >                             # template for telegram message (go template is supported)
            IP: {{.Host}}
            Scan Type: {{.Strategy}}
            Scan summary:
            {{range .Results}}
              port: {{.Port}}
              reachable: {{.IsReachable}}
                ------------{{end}}
        - via: email
          from: "tibebe"
          to: test@gmail.com
          subject: "Go-Alive Test Report"
          template: >
            IP: {{.Host}}
            Scan Type: {{.Strategy}}
            Scan summary:
            {{range .Results}}
              port: {{.Port}}
              reachable: {{.IsReachable}}
              {{ if .Error }}
                error:
                  {{ .Error }}
              {{end}}
                ------------
            {{end}}
        - via: slack
          channel: go-alive-test-group
          from: go-alive-test-bot
          template: >
            IP: {{.Host}}

            Scan Type: {{.Strategy}}

            Scan summary:
            {{range .Results}}
              port: {{.Port}}
              reachable: {{.IsReachable}}
                ------------{{end}}
notifications:                                      # notification channels configurations
  telegram:
    bots:                                           # telegram bots to send messages from
      - name: "go-alive-test-bot"
        token: "123456:bot-token"
    chats:                                          # list of telegram recipients
      - name: "go-alive-test-group"
        chatid: 1123232322
      - name: "tibebe"
        chatid: 12345678
  slack:
    apps:
      - name: "go-alive-test-bot"
        token: "bot user token"
    channels:
      - name: "go-alive-test-group"
        channelid: 'channel id'
  email:                                            # email configuration
    smtp:
      - name: "tibebe"
        sender: "test@gmail.com"
        auth:
          username: "test@gmail.com"
          password: "password"
        server: "smtp.gmail.com"
        port: 587
  webhook:
    - name: "webhook api"
      endpoint: "http://localhost:8000"
      auth:
        endpoint: "http://localhost:7000"
        email: "test@gmail.com"
        password: "password"
```

## Bugs ##

Bugs or suggestions? Visit the [issue tracker](https://github.com/TibebeJS/go-alive/issues) 

## Contribution

Feel free to fork, edit and send a PR.
