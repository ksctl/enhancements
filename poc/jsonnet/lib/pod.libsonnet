local k8s = import 'k8s.libsonnet';

{
  createPod(name, containers)::
    k8s.k8sResource('Pod') +
    {
      metadata: {
        name: name,
        labels: {
          app: name,
        },
      },
      spec: {
        containers: containers,
      },
    },
}
