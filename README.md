# kubeswipe

Have you found your cluster filled with redundant resources over time? Spending hours, days, or even months debugging and manually removing them can be a hassle. But fear not, there's something called kubeswipe that can efficiently clean up your resources and save you time.

Regular swipe operations should clean terminating namespaces as well, but for that, you need cluster scope permissions.

kubeswipe not only identifies and removes idle resources but also intelligently detects resources that may not appear idle at first glance but aren't serving any value. For example, if you deployed two apps with the intention to use only one, kubeswipe can identify and delete the unnecessary pods based on the resource history of each application. It can even detect cases where your main application isn't receiving traffic and remove unused pods. For your convenience, you can set an expiration time, or use the default.

To enable cleanup based on resource consumption, set `swipePolicy: moderate`.

While some resources might not be redundant but kept idle for future use, there's a high chance you'll forget about them, especially if they contain sensitive data. kubeswipe ensures these idle resources are identified and removed, preventing unauthorized access.

In cases where your pods are pending and you have no plans to scale, kubeswipe can clean them up as well by setting `swipePolicy: moderate`.

## Swipe Policy:

- **Low** involves:
  - Cleaning idle resources.
  - Cleaning terminating namespaces.

- **Moderate** involves:
  - Cleaning resources that aren't serving your users anymore.

Sometimes, resources go missing and never return. But with kubeswipe, you can ensure they come back. Setting `backup: true` will record the YAML of the resource and leave it in the `/kubeswipe/` directory.

**Todo:** 
- Add cloud backup support.
- Implement an easy apply option in the UI.

## Reasons to use kubeswipe:

- You're in a production cluster and want to avoid unnecessary costs.
- You're in a dev cluster and don't want to repeatedly delete and reset dev dependencies.
- You prioritize resource creation over deletion.

### Usuage 

```yaml
apiVersion: kubeswipe.kubefit.com/v1
kind: ResourceCleaner
metadata:
  name: resourcecleaner-sample
spec:
  resources:
    include:
      - name: Service
        namespace: default
        backup: false
      - name: Pod
        namespace: default
        backup: false
    exclude:
      - name: Namespace
        namespace: kube-system
        backup: false
  schedule: "@every 1m"
  operation: CLEANUP
```  

If you want to hand pick the resources to include and exclude use the include and exclude fields under the resources or else if you leave them empty it becomes cluster wide for all supported resources by kubeswipe . 

set schedule based on the time you want to schedule the reconcillation of the cleanup process 

operation you can set CLEANUP or SERVE . CLEANUP finds used resources and cleans them automatically serve helps to just retrieve and delete it by clicking the button on the UI


Future support features:
- to track the deleted resources 
- backup them via cloudprovider
- reapply deleted resources
- advance features with usuage of swipePolicy


### Running on the cluster


1. Build

```sh
make build
```

2. Install

```sh
make install
```

## Contributing

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

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

