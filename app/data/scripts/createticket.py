import sys
import json
import random
import os

import requests

url = os.environ["CATALYST_APP_URL"]
header = {"Authorization": "Bearer " + os.environ["CATALYST_TOKEN"]}

newtickets = requests.get(url + "/api/tickets", headers=header).json()
for ticket in newtickets.items:
    requests.delete(url + "/api/tickets/" + ticket.id, headers=header)

# Create a new ticket
requests.post(url + "/api/tickets", headers=header, json={
    "name": "New Ticket",
    "type": "alert",
    "open": True,
})