# GelOp
K8s Operator for Grafana Enterprise Logs plugin

## Description
Operator for tenant management and LBAC creation.

For more information about Grafana Enterprise Logs visit [GEL](https://grafana.com/docs/enterprise-logs/latest/)

**Note:** This Operator uses GEL v3 Admin API requests to create configurations

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).
 
1. Create a Tenant
```
apiVersion: loki.hamravesh.com/v1alpha1
kind: GrafanaEnterpriseLogsTenant
metadata:
  name: grafanaenterpriselogstenant-sample
spec:
  tenantInfo:
    name: "mytenant"
    displayName: "mytenant"
    clusterName: "loki"
```
2. Create a LBAC
```
apiVersion: loki.hamravesh.com/v1alpha1
kind: GrafanaEnterpriseLogsAccessPolicy
metadata:
  name: grafanaenterpriselogsaccesspolicy-sample
spec:
  tenantInfoRef: 
    tenantName: "mytenant"
    clusterName: "loki"
    accessPolicies:
      -  "logs:read"
      -  "logs:write"
    labelSelectors:
      - '{foo="bar"}'
      - '{name="hello"}'
```


### Running on the cluster
1. Checkout the deployment env
```
export Loki_Endpoint_Address="http://loki.com"
export Loki_Admin_Api_Token="X"
```

2. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

3. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/gel:tag
```

4. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/gel:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

## Contributing

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out
1. Add Loki Endpoint credentials as an ENV
```
export Loki_Endpoint_Address="http://loki.com"
export Loki_Admin_Api_Token="X"
```

2. Install the CRDs into the cluster:

```sh
make install
```

3. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

