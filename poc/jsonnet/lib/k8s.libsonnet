{
  local apiVersions = {
    Pod: 'v1',
    Service: 'v1',
    Deployment: 'apps/v1',
    ConfigMap: 'v1',
    Secret: 'v1',
  },

  k8sResource(kind):: {
    apiVersion: if kind in apiVersions then apiVersions[kind] else error 'Unknown resource kind: ' + kind,
    kind: kind,
  },
}
