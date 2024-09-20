local k8s = import 'k8s.libsonnet';

{
  createConfigMap(name, data, labels={})::
    k8s.addLabels(
      k8s.k8sResource('ConfigMap') +
      {
        metadata: {
          name: name,
        },
        data: data,
      },
      labels
    ),
}
