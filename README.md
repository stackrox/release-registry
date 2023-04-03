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
