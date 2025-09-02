import sys
import json
import random
import os

import requests

# Parse the event from the webhook payload
event = json.loads(sys.argv[1])
body = json.loads(event["body"])

url = os.environ["CATALYST_APP_URL"]
header = {"Authorization": "Bearer " + os.environ["CATALYST_TOKEN"]}

# Create a new ticket
requests.post(url + "/api/tickets", headers=header, json={
    "name": body["name"],
    "type": "alert",
    "open": True,
})