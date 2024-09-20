local k8s = import 'k8s.libsonnet';

{
  createSecret(name, type, data)::
    k8s.k8sResource('Secret') +
    {
      metadata: {
        name: name,
      },
      type: type,
      data: data,
    },
}
