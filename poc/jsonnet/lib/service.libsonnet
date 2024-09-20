local k8s = import 'k8s.libsonnet';

{
  createService(name, selector, ports, labels={})::
    k8s.addLabels(
      k8s.k8sResource('Service') +
      {
        metadata: {
          name: name,
        },
        spec: {
          selector: selector,
          ports: ports,
          type: 'ClusterIP',
        },
      },
      labels
    ),
}
