# Document on Ksctl agent related things

It gets directly installed on the cluster created by ksctl (Both Managed and Self-managed)

It is a GRPC Based app which is a ksctl core with all of its benefits

It will listen on Port 8080 `agent.ksctl.svc.cluster.local:8080`

It contains ClusterLevel permissions with all access

> [!CAUTION]
> need to add netpol for ksctl agent and who can call

Using the clusterPolicy it recieves state from the ksctl
stoage controller to store the ksctl `struct StorageDocument` and `struct CredentialDocument`

## State Transfer

### If Cluster was created and the state where it stored was local (host)

For that we need to deploy the storageImporter which can call storage.Export()
and storage.Import() inside ksctl agent to successfully transfer the sate from 
the local machine to inside kubernetes cluster

### If the cluster was created and the state where it stored used a external store (MongoDB)

for this we only nneed to transfer the credentials to the kubernetes cluster as **Secrets**
and the rest of the state management by the ksctl will be taken care of using
Database URI with credentials


