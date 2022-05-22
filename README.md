# silencer
Create new silences in the Prometheus alertmanager through a REST API.

## Usage
Send an HTTP POST request to `$silencerUri/silence/$service/$tenant` with the variable `duration=$silenceDuration` in seconds in the HTTP body.

For example, to create 5 minute silence for the service `my-service` with the tenant `tenant1` and the stage `prod`, you can issue the following command:

```shell script
curl -X POST -F 'duration=300' -F 'comment=Jenkins Job XYZ' http://localhost:8000/silence/my-service/tenant1-prod
```

During the creation, it is automatically checked if there are other active silences which covers the received silence request. If there is already a silence, a new one will __not__ be created.

## Configuration

```yaml
log_level: info                                 # log level
port: 8000                                      # listening port of the silencer application
alertmanager_host: alertmanager.prometheus.svc  # AlertManager's hostname
alertmanager_scheme: http                       # AlertManager's protocol scheme
alertmanager_port: 80                           # AlertManager's port
known_services:
  - "*"                                         # list of known services or "*" as wildcard to accept *any* service
```