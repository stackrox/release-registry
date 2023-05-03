#!/usr/bin/python3

import json
import os
import random
import requests
import yaml

BASE_URL="localhost:8443"

RELREG_TOKEN = os.getenv("RELREG_TOKEN")

HEADERS = {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'Authorization': f'Bearer {RELREG_TOKEN}'
    }

def create_new_release(release):
    url = f"https://{BASE_URL}/v1/release"

    payload = json.dumps(release)
    response = requests.request("POST", url, headers=HEADERS, data=payload, verify=False)

    print(response.text)

def create_new_quality_milestone_definition(qmd):
    url = f"https://{BASE_URL}/v1/qualitymilestonedefinition"
    payload = json.dumps(qmd)
    response = requests.request("POST", url, headers=HEADERS, data=payload, verify=False)
    print(response.text)

def approve_release(tag, qm):
    url = f"https://{BASE_URL}/v1/release/{tag}/approve"
    payload = json.dumps(qm)
    response = requests.request("POST", url, headers=HEADERS, data=payload, verify=False)
    print(response.text)


def read_demo_definition(path):
    with open(path) as f:
        data = yaml.load(f, yaml.SafeLoader)
    return data


def write_demo_definition(path, data):
    with open(path, "w") as f:
        yaml.dump(data, f)


def main():
    path = "demo.yaml"
    data = read_demo_definition(path)

    for release in data["releases"]:
        create_new_release(release)

    for qmd in data["qualityMilestoneDefinition"]:
        create_new_quality_milestone_definition(qmd)

    for release in data["releases"]:
        if "rc" in release["tag"]:
            if int(release["tag"][-1]) % 2 == 0:
                approve_release(release["tag"], {
                    "qualityMilestoneDefinitionName": "Nightly passed",
                    "metadata": [
                        {"key": "BuildUrl", "value": f"https://github.com/stackrox/stackrox/actions/runs/{random.randint(4787637767, 6787637767)}"},
                        {"key": "NumberOfCIFailures", "value": str(random.randint(0, 10))}
                    ]
                })
        if "nightly" in release["tag"]:
            if int(release["tag"][-1]) % 2 == 0:
                approve_release(release["tag"], {
                    "qualityMilestoneDefinitionName": "Nightly passed",
                    "metadata": [
                        {"key": "BuildUrl", "value": f"https://github.com/stackrox/stackrox/actions/runs/{random.randint(4787637767, 6787637767)}"},
                        {"key": "NumberOfCIFailures", "value": str(random.randint(0, 10))}
                    ]
                })
            if int(release["tag"][-1]) == 2:
                approve_release(release["tag"], {
                    "qualityMilestoneDefinitionName": "Cloud service validation passed",
                    "metadata": [
                        {"key": "BakeDuration", "value": str(random.randint(0, 10))},
                        {"key": "CanarySize",  "value": str(random.randint(0, 200))}
                    ]
                })


if __name__ == "__main__":
    main()
