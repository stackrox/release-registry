# Release Registry

A mechanism to mark, identify and search release artifacts using Quality Milestones.

## Documentation

- [Introduction and Requirements](./docs/introduction.md)
- [Use Cases](./docs/use-cases.md)
- [Data Model](./docs/data-model.md)
- [Architecture](./docs/architecture.md)

## Bootstrapping

1. Update `example/config.yaml` to reflect your environment. All options can be found in the configuration package.
1. Generate localhost certificates for the gRPC gateway: `make server-renew-cert`. They will be placed in the `example` directory.

If you want to deploy to GKE, create a new global IP:

```bash
gcloud compute addresses create --project <PROJECT> <NAME OF ADDRESS> --global --ip-version IPV4
```

The name of the address should be recorded in `.Values.reservedAddressName`.

## Developing

Tests are available in the `tests` directory and the Go packages.
They contain unit, integration and end-to-end tests.

- *Unit tests* assert function output on individual package level.
- *Integration tests* test interplay between components with mocked dependencies.
- *End-to-end tests* emulate user behaviour and test the external interfaces of the service.

### Requirements

For the end-to-end tests, must obtain a test token or create your own for the robot user `roxbot+release-registry-e2e@redhat.com`.
This token must be available in the environment variable `RELREG_TEST_TOKEN` when running end-to-end-tests.
