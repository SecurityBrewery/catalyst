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
    api_key = msg['secrets']['vt_api_key'].encode('utf-8')
    resource = msg['payload']['default'].encode('utf-8')
    params = {'apikey': api_key, 'resource': resource}
    return requests.get("https://www.virustotal.com/vtapi/v2/file/report", params=params).json()


print(json.dumps(run(json.loads(sys.argv[1]))))
