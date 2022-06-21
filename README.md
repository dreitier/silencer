# silencer
Create new silences in the Prometheus Alertmanager through a REST API.

## Description
During the deployment of a new release, you might want to silence alerts configured in the Prometheus Alertmanager. *silencer* provides a simple REST API, which then can be called by your build agents before the deployments starts.

## Getting started

### Usage
Send an HTTP POST request to `$silencerUri/silence/$service/$tenant` with the variable `duration=$silenceDuration` in seconds in the HTTP body.

For example, to create 5 minute silence for the service `my-service` with the tenant `tenant1` and the stage `prod`, you can issue the following command:

```bash
curl -X POST -F 'duration=300' -F 'comment=Jenkins Job XYZ' http://localhost:8000/silence/my-service/tenant1-prod
```

During the creation, it is automatically checked if there are other active silences which covers the received silence request. If there is already a silence, a new one will __not__ be created.

### Configuration

```yaml
log_level: info                                 # log level
port: 8000                                      # listening port of the silencer application
alertmanager_host: alertmanager.prometheus.svc  # AlertManager's hostname
alertmanager_scheme: http                       # AlertManager's protocol scheme
alertmanager_port: 80                           # AlertManager's port
known_services:
  - "*"                                         # list of known services or "*" as wildcard to accept *any* service
```

### Docker container
You can find ready-to-run Docker containers at [dreitier/silencer](https://hub.docker.com/repository/docker/dreitier/silencer).

## Development
### Creating new releases
A new release (artifact & Docker container) is automatically created when a new Git tag is pushed:

```bash
git tag x.y.z
git push origin x.y.z
```

## Support
This software is provided as-is. You can open an issue in GitHub's issue tracker at any time. But we can't promise to get it fixed in the near future.
If you need professionally support, consulting or a dedicated feature, please get in contact with us through our [website](https://dreitier.com).

## Contribution
Feel free to provide a pull request.

## License
This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.