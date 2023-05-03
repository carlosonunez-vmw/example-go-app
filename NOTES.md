## Namespace Ops

- Do these things whenever a new dev namespace is created:
  - Create a registry secret so that TBS can push images to your reg of choice,
    and
  - Create a ServiceAccount for doing things within the namespace and a
    RoleBinding to attach it to the cluster.

### Creating the Registry Secret

```sh
tanzu secret registry add registry-credentials \
  --server https://$REGISTRY_HOSTNAME \
  --username $REGISTRY_USERNAME \
  --password $REGISTRY_PASSWORD \
  --namespace $NAMESPACE
```
