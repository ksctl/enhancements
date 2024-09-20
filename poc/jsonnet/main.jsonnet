local deployment = import 'lib/deployment.libsonnet';
local service = import 'lib/service.libsonnet';
local configmap = import 'lib/configmap.libsonnet';
local secret = import 'lib/secret.libsonnet';

local appName = 'my-web-app';
local appPort = 8080;

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
  }]
);

local myService = service.createService(
  appName,
  { app: appName },
  [{
    name: 'http',
    port: appPort,
    targetPort: appPort,
  }]
);

local myConfigMap = configmap.createConfigMap(
  appName + '-config',
  {
    'app.properties': 'key1=value1\nkey2=value2',
  }
);

local mySecret = secret.createSecret(
  appName + '-secret',
  'Opaque',
  {
    username: std.base64('admin'),
    password: std.base64('password123'),
  }
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
