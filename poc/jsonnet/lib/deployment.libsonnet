local k8s = import 'k8s.libsonnet';

{
  createDeployment(name, replicas, containers)::
    k8s.k8sResource('Deployment') +
    {
      metadata: {
        name: name,
      },
      spec: {
        replicas: replicas,
        selector: {
          matchLabels: {
            app: name,
          },
        },
        template: {
          metadata: {
            labels: {
              app: name,
            },
          },
          spec: {
            containers: containers,
          },
        },
      },
    },
}