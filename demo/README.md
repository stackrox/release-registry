# Demo script

## Demo environment description

SQLite database file `pre-demo.sqlite` with

- multiple RC releases
  - even RCs have passed "nightly" quality milestone
- multiple nightly releases
  - nightlies from even days have passed "nightly" quality milestone
  - nightly from days ending in "2" have also passed "cloud service validation" quality milestone

This was produced by `demo.py` based on the information in `demo.yaml`.
To reproduce, set the `RELREG_TEST_TOKEN` environment variable and run the script against a development server.

## Runbook

- Show all releases and quality milestone definitions
  - https://localhost:8443/v1/release?preload=true
  - https://localhost:8443/v1/qualitymilestonedefinition

- Find latest release for:
  - Nightly: https://localhost:8443/v1/find?preload=true&qualityMilestoneName=Nightly%20passed
  - Cloud Service validation: https://localhost:8443/v1/find?preload=true&qualityMilestoneName=Cloud%20service%20validation%20passed

- Approve 4.1.0 release for Nightly and show that it is now latest by revisiting https://localhost:8443/v1/find?preload=true&qualityMilestoneName=Nightly%20passed:

```bash
curl \
  --location "https://localhost:8443/v1/release/4.1.0/approve" \
  --insecure \
  --header "Content-Type: application/json" \
  --header "Accept: application/json" \
  --header "Authorization: Bearer ${RELREG_TEST_TOKEN}" \
  --data '{
    "qualityMilestoneDefinitionName": "Nightly passed",
    "metadata": [
      {
        "key": "BuildUrl",
        "value": "https://github.com/stackrox/stackrox/actions/runs/4787637767"
      },
      {
        "key": "NumberOfCIFailures",
        "value": "4"
      }
    ]
  }'
```
