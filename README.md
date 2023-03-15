# Release Registry

A mechanism to mark, identify and search release artifacts using Quality Milestones.

## Documentation

- [Introduction and Requirements](./docs/introduction.md)
- [Use Cases](./docs/use-cases.md)
- [Data Model](./docs/data-model.md)
- [Architecture](./docs/architecture.md)

## (Re-)generate certificates

To (re-)generate the certificates used by the server, run the following on openSSL >= 1.1.1:

```bash
openssl req -x509 -newkey rsa:4096 -sha256 -days 365 -nodes \
  -keyout example/server.key -out example/server.crt -subj "/CN=localhost" \
  -addext "subjectAltName=DNS:localhost"
```
