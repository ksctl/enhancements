{
  local apiVersions = {
    Pod: 'v1',
    Service: 'v1',
    Deployment: 'apps/v1',
    ConfigMap: 'v1',
    Secret: 'v1',
  },

  defaultLabels: {
    'managed-by': 'jsonnet',
    'environment': std.extVar('environment'),  // This will be set when running jsonnet
  },

  k8sResource(kind):: {
    apiVersion: if kind in apiVersions then apiVersions[kind] else error 'Unknown resource kind: ' + kind,
    kind: kind,
  },

  addLabels(resource, labels)::
    resource + {
      metadata+: {
        labels: $.defaultLabels + (if 'labels' in resource.metadata then resource.metadata.labels else {}) + labels,
      },
    },
}
