# Lawnmower
`Lawnmower` makes ClusterRole which explicity lists apiGroups and remove specified apiGroup.

Kubernetes RBAC can't use [Deny Rules](https://github.com/kubernetes/kubernetes/issues/85963).
This tool removes specified apiGroups insted of deny rules.

## Build

```
$ go build -o lawnmower .
```

## Usage

```
$ ./lawnmower --help
Usage of ./lawnmower:
  -clusterRoleName string
    	resource name of cluster role. (default "restricted-cluster-admin")
  -except string
    	except apiGroups sepalated by comma. EX) rbac.authorization.k8s.io,networking.k8s.io
  -kubeconfig string
    	(optional) absolute path to the kubeconfig file. (default "~/.kube/config")
  -outputFileName string
    	name of output cluster role file. (default "restricted-cluster-admin-cluster-role.yaml")
```

### Make ClusterRole

`lawnmower` read all apiGroups from kubernetes API and set a rule of ClusterRole.


```
$ ./lawnmower

$ cat restricted-cluster-admin-cluster-role.yaml 
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  creationTimestamp: null
  name: restricted-cluster-admin
rules:
- apiGroups:
  - apiregistration.k8s.io
  - extensions
  - apps
  - events.k8s.io
  - authentication.k8s.io
  - authorization.k8s.io
  - autoscaling
  - batch
  - certificates.k8s.io
  - networking.k8s.io
  - policy
  - rbac.authorization.k8s.io
  - storage.k8s.io
  - admissionregistration.k8s.io
  - apiextensions.k8s.io
  - scheduling.k8s.io
  - coordination.k8s.io
  - node.k8s.io
  - compose.docker.com
  resources:
  - '*'
  verbs:
  - '*'
- nonResourceURLs:
  - '*'
  verbs:
  - '*'

```

## Except rbac.authorization.k8s.io

```
$ ./lawnmower --except rbac.authorization.k8s.io 

$ cat restricted-cluster-admin-cluster-role.yaml 
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  creationTimestamp: null
  name: restricted-cluster-admin
rules:
- apiGroups:
  - ""
  - apiregistration.k8s.io
  - extensions
  - apps
  - events.k8s.io
  - authentication.k8s.io
  - authorization.k8s.io
  - autoscaling
  - batch
  - certificates.k8s.io
  - networking.k8s.io
  - policy
  - storage.k8s.io
  - admissionregistration.k8s.io
  - apiextensions.k8s.io
  - scheduling.k8s.io
  - coordination.k8s.io
  - node.k8s.io
  - compose.docker.com
  resources:
  - '*'
  verbs:
  - '*'
- nonResourceURLs:
  - '*'
  verbs:
  - '*'
```