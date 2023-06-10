# helm-clear

A Helm plugin to clear old release's resource.
It search revisioned configmap, and delete old configmaps.

## Requirement

- helm >= 3.x

## Install

```bash
helm plugin install https://github.com/mocyuto/helm-clear --version v0.x.x
```

## Usage

Show what configmap would be deleted.

```bash
helm clear configmap test-chart --dry-run --namespace test
```

### Example

```bash
$ cd example
$ helm package configmap-sample

## create two revisions
$ helm install configmap-sample ./configmap-sample-0.1.0.tgz
$ helm upgrade configmap-sample ./configmap-sample-0.1.0.tgz

## show what configmap would be deleted
$ helm clear configmap configmap-sample --dry-run --history 1
dry-run mode: configmaps to be removed
sample-txt-1
```
