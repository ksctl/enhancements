local k8s = import 'k8s.libsonnet';

{
  createService(name, selector, ports)::
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
}
