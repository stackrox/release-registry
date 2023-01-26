# Supported use cases

This document describes the supported use cases for an MVP and lays out extensions and next steps for each use case.

## Table of contents

1. Support on-call engineer as the reviewer of the nightly marks a tag as successful
2. Cloud Service Upgrader finds the latest version

## 1. Support on-call engineer as the reviewer of the nightly marks a tag as successful

### Description / MVP

- Role: Support on-call engineer
- Actions:
  - They create a new tag in release registry when starting to review a nightly build, by running a CLI command.
  - Once they conclude that the nightly build was successful, e.g. all tests have passed, or all test failures are known/accepted/flakes, they approve the tag in the release registry by running a CLI command. That would be an approximation of reaching the Quality Milestone 1 of the redesigned release process.

### Extension / Next steps

- Nightly automation creates the tag in release registry through API, so engineer's task is reduced to the approval via CLI command.
- Approval is done automatically when all checks have passed through API.

## 2. Cloud Service Upgrader finds the latest version

### Description / MVP

- Role: Engineer who wants to upgrade a Cloud Service staging environment to a known good nightly version
- Actions:
  - They use CLI to search for latest approved (nightly) tag.

### Extension / Next steps

- Engineer can reject tags, i.e. remove approval
- Cloud Service upgrade process can use API to search for tags
