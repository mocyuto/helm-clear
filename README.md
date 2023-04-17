# helm-clear

A Helm plugin to clear old release's resource.
It search revisioned configmap, and delete old configmaps.

## Install

```bash
helm plugin install https://github.com/mocyuto/helm-clear --version v0.x.x
```

## Usage

Show what configmap would be deleted.

```bash
$ helm clear configmap test-chart --dry-run --namespace test
```

