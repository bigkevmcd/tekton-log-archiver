# tekton-log-archiver Go Operator

## Overview

This operator tracks completed [Tekton](https://github.com/tektoncd/pipeline) [PipelineRuns](https://github.com/tektoncd/pipeline/blob/master/docs/pipelineruns.md) and [TaskRuns](https://github.com/tektoncd/pipeline/blob/master/docs/taskruns.md) and attempts to archive the logged output of all containers to a storage location.

## Why?

You want to remove the runs that are cluttering your UI, but want to keep the logged output

## Prerequisites

- [go][go_tool] version v1.12+.
- [docker][docker_tool] version 17.03+
- [kubectl][kubectl_tool] v1.11.3+
- [operator-sdk][operator_install]
- Access to a Kubernetes v1.11.3+ cluster

## Getting Started

### Cloning the repository

Checkout the Operator repository

```
$ git clone https://github.com/bigkevmcd/tekton-log-archiver.git
$ cd tekton-log-archiver
```
### Pulling the dependencies

Run the following command

```
$ go mod tidy
```

### Building the operator

Build the operator image and push it to a public registry, such as quay.io:

```
$ export IMAGE=quay.io/example-inc/tekton-log-archiver:v0.0.1
$ operator-sdk build $IMAGE
$ docker push $IMAGE
```

### Using the image

```shell
# Update the operator manifest to use the built image name (if you are performing these steps on OSX, see note below)
$ sed -i 's|REPLACE_IMAGE|quay.io/example-inc/tekton-log-archiver:v0.0.1|g' deploy/operator.yaml
# On OSX use:
$ sed -i "" 's|REPLACE_IMAGE|quay.io/example-inc/tekton-log-archiver:v0.0.1|g' deploy/operator.yaml
```

**NOTE** The `quay.io/example-inc/tekton-log-archiver:v0.0.1` is an example. You should build and push the image for your repository.

### Installing

You *must* have Tekton [Pipeline](https://github.com/tektoncd/pipeline/) installed before installing this operator:

```shell
$ kubectl apply -f https://github.com/tektoncd/pipeline/releases/download/v0.9.1/release.yaml
```

And then you can install the statuses operator with:

```shell
$ kubectl apply -f https://github.com/bigkevmcd/operator-statuses/releases/download/v0.0.1/release.yaml
```

### Uninstalling

```shell
$ kubectl delete -f https://github.com/bigkevmcd/operator-statuses/releases/download/v0.0.1/release.yaml
```

### Troubleshooting

Use the following command to check the operator logs.

```shell
$ kubectl logs tekton-log-archiver
```

### Running Tests

```shell
$ go test -v ./...
```

[dep_tool]: https://golang.github.io/dep/docs/installation.html
[go_tool]: https://golang.org/dl/
[kubectl_tool]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
[docker_tool]: https://docs.docker.com/install/
[operator_sdk]: https://github.com/operator-framework/operator-sdk
[operator_install]: https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md
