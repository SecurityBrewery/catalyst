#!/usr/bin/env python

import sys
import json
import hashlib


def run(msg):
    sha1 = hashlib.sha1(msg['payload']['default'].encode('utf-8'))
    return {"hash": sha1.hexdigest()}


print(json.dumps(run(json.loads(sys.argv[1]))))
