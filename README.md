# networksets-controller
With this controller, you can control access to/from custom resources by their names<br>
The addresses of these resources can be obtained from your own dedicated service<br>
For example, they can be salt hosts, net prefixes, dns names...<br>
Here is an example with dns names<br>

## Description
Controller watches by create/update/delete [Calico NetworkPolicy](https://docs.projectcalico.org/reference/resources/networkpolicy).<br>
If source/destination selector of NetworkPolicy have the specific label `DNS_RESOLVER=<domain>` then controller creates/updates [Calico NetworkSet](https://docs.projectcalico.org/reference/resources/networkset).<br>
IP networks/CIDRs for NetworkSet are requested from the http url. This url is customizable for specific label.<br>
Controller periodically updates the NetworkSet once per 5 seconds.<br>
Аlso works with GlobalNetworkPolicy/GlobalNetworkSet.

## Getting Started

### Prerequisites
- go version v1.20.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.
- Calico API v3

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/networksets-controller:tag
```

**NOTE:** This image ought to be published in the personal registry you specified. 
And it is required to have access to pull the image from the working environment. 
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/networksets-controller:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin 
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k examples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k examples/
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Metrics
```
# HELP networkset_controller_globalnetworkset_created Total number of successful created globalnetworksets.
# TYPE networkset_controller_globalnetworkset_created counter
networkset_controller_globalnetworkset_created 0
# HELP networkset_controller_globalnetworkset_deleted Total number of successful deleted globalnetworksets.
# TYPE networkset_controller_globalnetworkset_deleted counter
networkset_controller_globalnetworkset_deleted 0
# HELP networkset_controller_globalnetworkset_deletion_failed Total number of failed deletion globalnetworksets.
# TYPE networkset_controller_globalnetworkset_deletion_failed counter
networkset_controller_globalnetworkset_deletion_failed 0
# HELP networkset_controller_globalnetworkset_update_failed Total number of failed updating globalnetworksets.
# TYPE networkset_controller_globalnetworkset_update_failed counter
networkset_controller_globalnetworkset_update_failed 0
# HELP networkset_controller_globalnetworkset_updated Total number of successful updated globalnetworksets.
# TYPE networkset_controller_globalnetworkset_updated counter
networkset_controller_globalnetworkset_updated 0
# HELP networkset_controller_networkset_created Total number of successful created networksets.
# TYPE networkset_controller_networkset_created counter
networkset_controller_networkset_created 0
# HELP networkset_controller_networkset_deleted Total number of successful deleted networksets.
# TYPE networkset_controller_networkset_deleted counter
networkset_controller_networkset_deleted 0
# HELP networkset_controller_networkset_deletion_failed Total number of failed deletion networksets.
# TYPE networkset_controller_networkset_deletion_failed counter
networkset_controller_networkset_deletion_failed 0
# HELP networkset_controller_networkset_failed Total number of failed creating networksets.
# TYPE networkset_controller_networkset_failed counter
networkset_controller_networkset_failed 0
# HELP networkset_controller_networkset_update_failed Total number of failed updating networksets.
# TYPE networkset_controller_networkset_update_failed counter
networkset_controller_networkset_update_failed 0
# HELP networkset_controller_networkset_updated Total number of successful updated networksets.
# TYPE networkset_controller_networkset_updated counter
networkset_controller_networkset_updated 0
# HELP networkset_controller_resolve_failed Total number of failed resolve attempts.
# TYPE networkset_controller_resolve_failed counter
networkset_controller_resolve_failed 0
# HELP networkset_controller_resolve_succesful 
# TYPE networkset_controller_resolve_succesful counter
networkset_controller_resolve_succesful 0
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

