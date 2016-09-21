# tugbot-kube
[![CircleCI](https://circleci.com/gh/gaia-docker/tugbot-kubernetes.svg?style=shield)](https://circleci.com/gh/gaia-docker/tugbot-kubernetes)
[![codecov](https://codecov.io/gh/gaia-docker/tugbot-kubernetes/branch/master/graph/badge.svg)](https://codecov.io/gh/gaia-docker/tugbot-kubernetes)
[![Go Report Card](https://goreportcard.com/badge/github.com/gaia-docker/tugbot-kubernetes)](https://goreportcard.com/report/github.com/gaia-docker/tugbot-kubernetes)
[![Docker](https://img.shields.io/docker/pulls/gaiadocker/tugbot-kube.svg)](https://hub.docker.com/r/gaiadocker/tugbot-kube/)
[![Docker Image Layers](https://imagelayers.io/badge/gaiadocker/tugbot-kube:latest.svg)](https://imagelayers.io/?images=gaiadocker/tugbot-kube:latest 'Get your own badge on imagelayers.io')

**Tugbot Kube** is a Continuous Testing Framework for Kubernetes based production/staging/testing environment. **Tugbot** executes *Kubernetes Test Jobs* upon some *event*, like Kubernetes node registration, deployment.

## Kubernetes Test Job

*Kubernetes Test Job* is a regular Kubernetes job. Docker & Kubernetes `LABEL` is used to discover *Kubernetes test job* and **Tugbot** related test metadata.
**Tugbot Kube** will trigger a sequential *test job* execution upon *event* (see `tugbot.kubernetes.events` label).
Running a *test job* should result running a *docker test container* (1 or more). In order to collect test results and depoly those results to elasticsearch, for example, you should [tugbot-collect](https://github.com/gaia-docker/tugbot-collect) related labels.
### Tugbot labels

All **Tugbot** labels must be prefixed with `tugbot.` to avoid potential conflict with other labels.
**Tugbot** labels are divided into:

1) Container labels:
- `tugbot.results.dir` - directory, where *test container* reports test results; default to `/var/tests/results`

2) Kubernetess Service labels:

- `tugbot.kubernetes.events` - list of comma separated Kubernetes events format: Kind.Reason. For example: Node.Starting,ReplicaSet.SuccessfulCreate

#####Example Kubernetes Test Job (adding tugbot events' label to [Kubernetes example](http://kubernetes.io/docs/user-guide/jobs/#running-an-example-job/)), *test-job.yaml*:
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
  labels:
    tugbot.kubernetes.events: Node.Starting,ReplicaSet.SuccessfulCreate
spec:
  template:
    metadata:
      name: pi
    spec:
      containers:
      - name: pi
        image: perl
        command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never
```
> It is highly recommended to determain `restartPolicy: Never` when creating a *Kubernetes Test Job*, so that it will rerun upon defined events using tugbot.

> Use `tugbot.kubernetes.events: Node.Starting,ReplicaSet.SuccessfulCreate` *Kubernetes label* to tell tugbot framework that this is a test job that should be updated each a replica-set successfuly created ot a node starting.

## Running Tugbot Kubernetes inside a Docker container
```bash
docker run -d -e KUBERNETES_HOST=<Address> -e KUBERNETES_CERT_PATH=<Kubernetes Certificate Path> --name tugbot-kube gaiadocker/tugbot-kube
```
- `KUBERNETES_HOST` - Kubernetes host address (default: https://192.168.99.100:8443).
- `KUBERNETES_CERT_PATH` - dirctory should contain: ca.pem, apiserver.crt, apiserver.key.
- `TUGBOT_KUBERNETES_NAMESPACE` - Namespace where jobs created by *Tugbot* should run.
- `TUGBOT_KUBERNETES_LOG_LEVEL` - Enable debug mode. When this option set to `debug` you'll see more verbose logging in the Tugbot Kubernetes log file.
