import sys
import json
import random
import os

from pocketbase import PocketBase

# Connect to the PocketBase server
client = PocketBase(os.environ["CATALYST_APP_URL"])
client.auth_store.save(token=os.environ["CATALYST_TOKEN"])

newtickets = client.collection("tickets").get_list(1, 200, {"filter": 'name = "New Ticket"'})
for ticket in newtickets.items:
    client.collection("tickets").delete(ticket.id)

# Create a new ticket
client.collection("tickets").create({
    "name": "New Ticket",
    "type": "alert",
    "open": True,
})