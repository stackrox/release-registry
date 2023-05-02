# Creating new release-registry release

On a branch based on the main branch:

1. Bump version in
    - `ui/package.json`
    - `deploy/chart/release-registry/Chart.yaml`

1. Update CHANGELOG

1. Merge changes to main

1. `git tag <tag> && git push <remote> <tag>`
