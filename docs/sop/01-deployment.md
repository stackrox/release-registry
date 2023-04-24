# Deploying Release Registry

## Prerequisites

- GKE cluster
- kubectl and Helm setup to point to GKE cluster
- Idendity provider with OIDC support

## Install operator

```bash
helm install pgo \
    --namespace postgres-operator \
    --create-namespace \
    oci://registry.developers.crunchydata.com/crunchydata/pgo \
    --version 5.3.1
```

## Configure secrets

Configure the secrets for the deployment, e.g. image pull secrets, OIDC configuration, admin password, etc.

The list of required values is kept in `deploy/chart/release-registry/secret-values.example.yaml`.

To find hints for OIDC configuration and localhost certificates for gRPC Gateway, consult the appropriate sections of this repository's README.

## Deploy

```bash
ENVIRONMENT=<development, production, ...> make server-helm-deploy
```

This will take values from `deploy/configuration/values-<ENVIRONMENT>.yaml`.
