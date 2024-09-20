local deployment = import 'lib/deployment.libsonnet';
local service = import 'lib/service.libsonnet';
local configmap = import 'lib/configmap.libsonnet';
local secret = import 'lib/secret.libsonnet';

local appName = 'ksctl-poc-nginx-jsonnet';
local appPort = 80;

local commonLabels = {
  'app.kubernetes.io/name': appName,
  'app.kubernetes.io/instance': appName + '-' + std.extVar('environment'),
};

local myDeployment = deployment.createDeployment(
  appName,
  3,
  [{
    name: appName,
    image: 'nginx:latest',
    ports: [
      {
        name: 'http',
        containerPort: appPort,
      },
    ],
  }],
  commonLabels + { 'app.kubernetes.io/component': 'web' }
);

local myService = service.createService(
  appName,
  { app: appName },
  [{
    name: 'http',
    port: appPort,
    targetPort: appPort,
  }],
  commonLabels + { 'app.kubernetes.io/component': 'web' }
);

local myConfigMap = configmap.createConfigMap(
  appName + '-config',
  {
    'app.properties': 'key1=value1\nkey2=value2',
  },
  commonLabels + { 'app.kubernetes.io/component': 'web' }
);

local mySecret = secret.createSecret(
  appName + '-secret',
  'Opaque',
  {
    username: std.base64('admin'),
    password: std.base64('password123'),
  },
  commonLabels + { 'app.kubernetes.io/component': 'web' }
);

{
  apiVersion: 'v1',
  kind: 'List',
  items: [
    myDeployment,
    myService,
    myConfigMap,
    mySecret,
  ],
}
