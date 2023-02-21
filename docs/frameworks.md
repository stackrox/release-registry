# Frameworks and Libraries

This document describes the requirements and choices for Go frameworks and libraries used in this project.

**Disclaimer: Without stating explicitly in each category, choices are also informed by our past experience in the Stackrox and Infra projects.**

| | Requirements           | Choice of technology and framework |
|--------------------------|--------------|------------------------------------|
| Server | <ul><li>lightweight</li><li>testable</li><li>RPC + HTTP API support</li><li>self-documenting</li></ul> | <ul><li>gRPC generated from ProtoBuf with JSON Gateway + Swagger plugins</li><li>Testing [buf CLI](https://github.com/bufbuild/buf) as manager</li></ul> |
| ORM Framework | <ul><li>ORM support incl. relations</li><li>Transactions, prepared statements</li><li>Migration support</li><li>Support for multiple databases</li></ul> | [GORM](https://gorm.io/) |
| Database | *see below* | Postgres + local SQLite |
| AuthnZ | <ul><li>Integration with Red Hat SSO (Keycloak)</li><li>Create expiring, signed tokens</li></ul> | [jwt-go](https://github.com/golang-jwt/jwt) |
| Logging | <ul><li>leveled</li><li>structured</li><li>easy to use</li><li>fast</li></ul> | [zap](https://github.com/uber-go/zap) |
| Build and Release | *Decision based on previous experience* | CI & CD environment GHA, distribution: Container image, Quay registry |
| Deployment | *Decision based on previous experience* | Helm + GKE, details TBD |
| Configuration Management | <ul><li>Load from files, env vars</li><li>Set defaults and overrides</li><li>Hot reload</li><li>Ease of use</li></ul>| [Viper](https://github.com/spf13/viper) |
| Command Line tools | <ul><li>Nested subcommands</li><li>Flags</li><li>Integration with configuration management</li></ul>| [Cobra](https://github.com/spf13/cobra) |
| Testing | <ul><li>Native Go</li><li>Fast</li><li>Readable test outputs</li><li>Deterministic</li></ul> | <ul><li>Default Go testing for speed and overall framework<ul><li>Utilize [table driven tests](https://github.com/golang/go/wiki/TableDrivenTests)</li></ul></li><li><a href="https://github.com/stretchr/testify">Testify</a> for its mock, setup and teardown functionalities</li></ul> |
| Developer Workflow | | The following aids shall be implemented as precommit hooks and Makefile targets: <ul><li>Linting (check for todos, check for unused code, other static checks)</li><li>Code formatting</li></ul> |

## Choice of Database

The GORM framework allows to plug and play different databases.
Therefore, we can pick the best database for two use cases:

We choose:

* Postgres for deployments, scale tests, ...
  * Industry standard, reliable and stable database
  * Already used in StackRox project
* SQLite for CI tests, local development
  * Lightweight, ease of use
  * Possibility to create empty database in memory and throw away after test run
  * Possibility to check SQLite database files as fixtures into the repository to setup (integration) tests
