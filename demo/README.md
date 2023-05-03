# Demo script

## Demo environment description

SQLite database file `pre-demo.sqlite` with

- multiple RC releases
  - even RCs have passed "nightly" quality milestone
- multiple nightly releases
  - nightlies from even days have passed "nightly" quality milestone
  - nightly from days ending in "2" have also passed "cloud service validation" quality milestone

This was produced by `demo.py` based on the information in `demo.yaml`.
To reproduce, set the `RELREG_TOKEN` environment variable and run the script against a development server.

## Runbook

- Show all releases and quality milestone definitions
  - https://localhost:8443/v1/release
  - https://localhost:8443/v1/qualitymilestonedefinition

- Find latest release for:
  - Nightly: https://localhost:8443/v1/find?preload=true&qualityMilestoneName=Nightly%20passed
  - Cloud Service validation: https://localhost:8443/v1/find?preload=true&qualityMilestoneName=Cloud%20service%20validation%20passed

- Approve 4.0.0 release for Nightly and show that it is now latest:
  - Approval needs to be done in Postman
  - Nightly: https://localhost:8443/v1/find?preload=true&qualityMilestoneName=Nightly%20passed
