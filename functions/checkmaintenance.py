import requests

url = "https://raw.githubusercontent.com/heckelmann/kyverno-keptn-workshop/main/functions/maintenance.json"
expected_json = {"maintenance": False}

response = requests.get(url)
if response.status_code == 200:
    json_data = response.json()
    if json_data == expected_json:
        print("Maintenance check passed!")
        print(json_data)
        exit(0)
    else:
        print("Maintenance check failed!")
        print(json_data)
        exit(1)
else:
    print("Failed to retrieve JSON data from the URL.")
    exit(1)