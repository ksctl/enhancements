## Controllers & CRDs

### Application Stacks and Components

#### Application.Stacks()

apiVersion: `application.ksctl.com/v1alpha1`
kind: `Stack`

```yaml
spec:
  components:
  - name: !str ""
    version: !str ""
    type: !enum "app | cni"
```

### Storage Importer & Exporter

> [!NOTE]
> for now we are going to use this just for exporting the state files
> (given) the creation of the cluster took place from host local machine
> (constrains) it will not import when the storage falls under 
    **_external storage compatibility requirements_** (which are Mongodb)

#### StorageImporter
it will [Watch](#storageimport)
will create a ksctl agent rpc client to send the docuemnt to import in the kubernetes cluster

apiVersion: `storage.ksctl.com/v1alpha1`
kind: `ImportState`


```yaml
spec:
  rawData: !bytes ""
  Succeded: !bool false
```


### LoadBalancer Provisioning

> [!NOTE]
> Work in Progress

### Cluster Autoscaler

> [!NOTE]
> Work in Progress

