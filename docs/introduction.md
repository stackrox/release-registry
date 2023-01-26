# Introduction

## Definitions

- *Quality Milestones* (QM) are intermediate goals that release process instances go through.
  They are aimed at simplifying the communication about the progress of the release process instance.
  Each QM acts as a gate to enable further tests or publishing of release artifacts.

- *Release artifacts* describe the artifacts published by the release process for consumers, for example images published to Quay.io, Helm charts, GitHub release, etc.

## Requirements

### Functional Requirements

- Create a release from commit (or tags)
- Mark a release as "has passed Quality Milestone X"
- Get the latest release for a specific QM
  - overall
  - for a specific channel, e.g. 3.73.x
- Trigger additional step when a release has been marked
- Keep (searchable) history of releases
- Get status for a release
  - Which Quality Milestones have passed?
- Links to release artifacts
  - Images
  - Workflows

### Non-functional Requirements

- Low runtime costs
- Reduce development costs and dependencies on third parties
- Extensibility towards a service for additionally tracking
  - quality milestones
  - deployment state
  - deployment history
- Testability
