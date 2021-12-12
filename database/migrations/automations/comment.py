#!/usr/bin/env python

import subprocess
import sys

subprocess.call(
    [sys.executable, "-m", "pip", "install", "requests"],
    stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL,
)

import json
import requests


def run(msg):
    if "ticket" in msg["context"]:
        headers = {"PRIVATE-TOKEN": msg["secrets"]["catalyst_apikey"]}
        url = "%s/tickets/%d/comments" % (msg["secrets"]["catalyst_apiurl"], msg["context"]["ticket"]["id"])
        data = {'message': msg["payload"]["default"], 'creator': 'automation'}
        requests.post(url, json=data, headers=headers).json()

    return {"done": True}


print(json.dumps(run(json.loads(sys.argv[1]))))
