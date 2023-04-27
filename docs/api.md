# API

## Quick links

- [Latest protocol documentation](https://release-registry.rox.systems/docs/)
- [OpenAPI spec](https://release-registry.rox.systems/docs/proto/)

## Protocol Documentation

The procotol documentation with rendered comments is generated at [`gen/openapiv2/index.html`](../gen/openapiv2/index.html).
It is served under the root of the `/docs/` endpoint of their release registry.

## OpenAPI (HTTP) spec

The OpenAPI spec is generated at [`gen/openapiv2/proto`](../gen/openapiv2/proto/).
The subdirectory structure of the `proto` directory is mirrored here, i.e. `proto/api/v1/release.proto` can be found at `gen/openapiv2/proto/api/v1/release.swagger.json`.
Users can explore the spec files at the `/docs/proto/` endpoint of their release registry.

For developing & testing the API, users can [import these files into Postman as Collections directly from the URL](https://learning.postman.com/docs/getting-started/importing-and-exporting-data/#importing-postman-data).
Then, set the baseUrl variable value to the deployed release registry and provide the service account token as Bearer Token.
Users can also [generate client code from Postman](https://learning.postman.com/docs/sending-requests/generate-code-snippets/), e.g. for cURL, Python and native Go.
