# drlm-core v3

drlm-Core v3

## Install protoc plugin for GO

`$ go get -u google.golang.org/grpc`

## Install cobra 

`$ go get -u github.com/spf13/cobra/cobra`

```
$ cat $HOME/.cobra.yaml
author: Pau Roura <pau@brainupdaters.net>
license: AGPL
```
`$ cobra init github.com/brainupdaters/drlm-cli`
`$ cd github.com/brainupdaters/drlm-cli`
`$ cobra add newCommand`


https://github.com/spf13/cobra/blob/master/cobra/README.md

## Install logrus

Logrus is a structured logger for Go (golang), completely API compatible with the standard library logger.

`$ go get "github.com/Sirupsen/logrus"`


### install Protocol Buffers V3

The simplest way to do this is to download pre-compiled binaries for your platform(protoc-<version>-<platform>.zip) from here: https://github.com/google/protobuf/releases
*Unzip this file.
*Update the environment variable PATH to include the path to the protoc binary file.

`$ go get -u github.com/golang/protobuf/protoc-gen-go`

`$ protoc -I drlm-comm/ drlm-comm/drlm-comm.proto --go_out=plugins=grpc:drlm-comm/`



## etcd
## mgmtconfig
## Sentry
