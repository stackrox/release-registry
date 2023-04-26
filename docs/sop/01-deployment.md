# Deploying Release Registry

## Prerequisites

- GKE cluster
- gcloud authenticated for GKE cluster and GCP Secrets Manager
- kubectl and Helm setup to point to GKE cluster
- Identity provider with OIDC support

## Install operator

```bash
helm upgrade pgo \
    --install \
    --namespace postgres-operator \
    --create-namespace \
    oci://registry.developers.crunchydata.com/crunchydata/pgo \
    --version 5.3.1
```

## Configure secrets

For bootstrapping a new environment, configure the secrets for the deployment, e.g. image pull secrets, OIDC configuration, admin password, etc. in `deploy/chart/release-registry/values-<ENVIRONMENT>.yaml`.
The list of required values is kept in `deploy/chart/release-registry/secret-values.example.yaml`.
To find hints for OIDC configuration and localhost certificates for gRPC Gateway, consult the appropriate sections of this repository's README.

Then, upload the new configuration file to GCP Secrets Manager:

```bash
ENVIRONMENT=<new environment> make server-helm-upload-local-values`.
```

## Deploy

To deploy with stored secret configuration, use:

```bash
ENVIRONMENT=<development, production, ...> make server-helm-deploy
```

Alternatively, you can use download and modify local configuration.

```bash
# Download
ENVIRONMENT=<development, production, ...> make server-helm-download-local-values
# Do you modifications
# Deploy
ENVIRONMENT=<development, production, ...> make server-helm-deploy-local-values
```

This will use values from `deploy/chart/release-registry/configuration/values-<ENVIRONMENT>.yaml`.
After your modifications, upload the configuration with:

```bash
ENVIRONMENT=<development, production, ...> make server-helm-upload-local-values
```
