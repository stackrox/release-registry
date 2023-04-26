# Rotating secrets

## Authentication/OIDC secrets

OIDC client secrets expire after 180 days.
The occasion is marked by calendar invite 14 days before expiration.

1. Follow the OIDC provider instructions to create a new client secret with expiration.
1. Update the stored secret value with the new client secret.
1. Redeploy the application.
1. Verify login still works in an incognito window.
1. Create a new calendar invite 14 days before the next expiration.

## Certificates

The gRPC server serves a localhost certificate.
Besides local deployments, gRPC gateway expects this certificate to be valid as well.
The gRPC gateway also expects a valid certificate.
The expiration is marked by calendar invite 14 days before.

1. Run `make server-renew-cert`.
1. Update the stored secret value with the new certificate and private key.
1. Redeploy the application.
1. Verify everything works.
1. Create a new calendar invite 14 days before the next expiration.
