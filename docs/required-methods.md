# List of required methods to support the use cases

This document describes all methods that the (data model for the) release registry should support.
It is the basis for evaluating the data model and creating end-to-end tests for a TDD approach.

## Defining Quality Milestones

### Create a Quality Milestone Definition

### Get a specific Quality Milestone Definition

### List Quality Milestone Definitions

### Edit a Quality Milestone Definition (out of MVP)

## Releases

### Create a release

### Get specific release

* and which Quality Milestones it has reached

### Update a field on a specific release (out of MVP)

### List all releases

* all releases overall
* all at a specific Quality Milestone overall
* all releases within a specific release channel (e.g. 3.73.x)
* all at a specific Quality Milestone within a specific release channel (e.g. 3.73.x)

### Release Lifecycle

### Approve a Quality Milestone for a release

* and trigger a webhook with a payload (out of MVP)

### Find the latest release

* latest created release overall
* latest release at a specific Quality Milestone overall
* latest created release within a specific release channel (e.g. 3.73.x)
* latest release at a specific Quality Milestone within a specific release channel (e.g. 3.73.x)

### Reject a release

* and you can still get this specific release by tag
* and it does not show up when listing all releases
  * can be shown based on flag
* and it does not show up when listing all releases for a specific Quality Milestone
  * can be shown based on flag

### Retry webhooks for Quality Milestone (out of MVP)

* single webhook
* all webhooks
