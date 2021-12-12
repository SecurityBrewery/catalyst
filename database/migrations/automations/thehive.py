#!/usr/bin/env python

import subprocess
import sys
import json
from datetime import datetime
import io

subprocess.check_call(
    [sys.executable, "-m", "pip", "install", "thehive4py", "requests", "minio"],
    stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL,
)

defaultschema = {
  "definitions": {},
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "https://example.com/object1618746510.json",
  "title": "Default",
  "type": "object",
  "required": [
    "severity",
    "description",
    "summary",
    "tlp",
    "pap"
  ],
  "properties": {
    "severity": {
      "$id": "#root/severity",
      "title": "Severity",
      "type": "string",
      "default": "Medium",
      "x-cols": 6,
      "x-class": "pr-2",
      "x-display": "icon",
      "x-itemIcon": "icon",
      "oneOf": [
        {
          "const": "Unknown",
          "title": "Unknown",
          "icon": "mdi-help"
        },
        {
          "const": "Low",
          "title": "Low",
          "icon": "mdi-chevron-up"
        },
        {
          "const": "Medium",
          "title": "Medium",
          "icon": "mdi-chevron-double-up"
        },
        {
          "const": "High",
          "title": "High",
          "icon": "mdi-chevron-triple-up"
        },
        {
          "const": "Very High",
          "title": "Very High",
          "icon": "mdi-exclamation"
        }
      ]
    },
    "flag": {
      "title": "Flag",
      "type": "boolean",
      "x-cols": 6,
    },
    "tlp": {
      "$id": "#root/tlp",
      "title": "TLP",
      "type": "string",
      "x-cols": 6,
      "x-class": "pr-2",
      "x-display": "icon",
      "x-itemIcon": "icon",
      "oneOf": [
        {
          "const": "White",
          "title": "White",
          "icon": "mdi-alpha-w"
        },
        {
          "const": "Green",
          "title": "Green",
          "icon": "mdi-alpha-g"
        },
        {
          "const": "Amber",
          "title": "Amber",
          "icon": "mdi-alpha-a"
        },
        {
          "const": "Red",
          "title": "Red",
          "icon": "mdi-alpha-r"
        }
      ]
    },
    "pap": {
      "$id": "#root/pap",
      "title": "PAP",
      "type": "string",
      "x-cols": 6,
      "x-class": "pr-2",
      "x-display": "icon",
      "x-itemIcon": "icon",
      "oneOf": [
        {
          "const": "White",
          "title": "White",
          "icon": "mdi-alpha-w"
        },
        {
          "const": "Green",
          "title": "Green",
          "icon": "mdi-alpha-g"
        },
        {
          "const": "Amber",
          "title": "Amber",
          "icon": "mdi-alpha-a"
        },
        {
          "const": "Red",
          "title": "Red",
          "icon": "mdi-alpha-r"
        }
      ]
    },
    "tags": {
      "$id": "#root/tags",
      "title": "Tags",
      "type": "array",
      "items": {
          "type": "string"
      }
    },
    "description": {
      "$id": "#root/description",
      "title": "Description",
      "type": "string",
      "x-display": "textarea",
      "x-class": "pr-2"
    },
    "resolutionStatus": {
      "$id": "#root/resolutionStatus",
      "title": "Resolution Status",
      "type": "string",
      "x-cols": 6,
      "x-class": "pr-2",
    },
    "endDate": {
      "$id": "#root/endDate",
      "title": "End Data",
      "type": "string",
      "format": "date-time",
      "x-cols": 6,
      "x-class": "pr-2",
    },
    "summary": {
      "$id": "#root/summary",
      "title": "Summary",
      "type": "string",
      "x-display": "textarea",
      "x-class": "pr-2"
    }
  }
}

defaultalertschema = {
  "definitions": {},
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "https://example.com/object1618746510.json",
  "title": "Default",
  "type": "object",
  "required": [
    "severity",
    "description",
    "summary",
    "tlp",
    "pap"
  ],
  "properties": {
    "severity": {
      "$id": "#root/severity",
      "title": "Severity",
      "type": "string",
      "default": "Medium",
      "x-cols": 6,
      "x-class": "pr-2",
      "x-display": "icon",
      "x-itemIcon": "icon",
      "oneOf": [
        {
          "const": "Unknown",
          "title": "Unknown",
          "icon": "mdi-help"
        },
        {
          "const": "Low",
          "title": "Low",
          "icon": "mdi-chevron-up"
        },
        {
          "const": "Medium",
          "title": "Medium",
          "icon": "mdi-chevron-double-up"
        },
        {
          "const": "High",
          "title": "High",
          "icon": "mdi-chevron-triple-up"
        },
        {
          "const": "Very High",
          "title": "Very High",
          "icon": "mdi-exclamation"
        }
      ]
    },
    "tlp": {
      "$id": "#root/tlp",
      "title": "TLP",
      "type": "string",
      "x-cols": 6,
      "x-class": "pr-2",
      "x-display": "icon",
      "x-itemIcon": "icon",
      "oneOf": [
        {
          "const": "White",
          "title": "White",
          "icon": "mdi-alpha-w"
        },
        {
          "const": "Green",
          "title": "Green",
          "icon": "mdi-alpha-g"
        },
        {
          "const": "Amber",
          "title": "Amber",
          "icon": "mdi-alpha-a"
        },
        {
          "const": "Red",
          "title": "Red",
          "icon": "mdi-alpha-r"
        }
      ]
    },
    "source": {
      "$id": "#root/source",
      "title": "Source",
      "type": "string",
      "x-cols": 4,
      "x-class": "pr-2",
    },
    "sourceRef": {
      "$id": "#root/sourceRef",
      "title": "Source Ref",
      "type": "string",
      "x-cols": 4,
      "x-class": "pr-2",
    },
    "type": {
      "$id": "#root/type",
      "title": "Type",
      "type": "string",
      "x-cols": 4,
      "x-class": "pr-2",
    },
    "description": {
      "$id": "#root/description",
      "title": "Description",
      "type": "string",
      "x-display": "textarea",
      "x-class": "pr-2"
    }
  }
}


class schema:
    def __init__(self):
        self.schema = defaultschema

    def add_string(self, title):
        self.schema["properties"][title] = { "type": "string", "x-cols": 6, "x-class": "pr-2" }

    def add_boolean(self, title):
        self.schema["properties"][title] = { "type": "boolean", "x-cols": 6, "x-class": "pr-2" }

    def add_date(self, title):
        self.schema["properties"][title] = { "type": "string", "format": "date-time", "x-cols": 6, "x-class": "pr-2" }

    def add_integer(self, title):
        self.schema["properties"][title] = { "type": "integer", "x-cols": 6, "x-class": "pr-2" }

    def add_float(self, title):
        self.schema["properties"][title] = { "type": "number", "x-cols": 6, "x-class": "pr-2" }


class alertschema:
    def __init__(self):
        self.schema = defaultalertschema


def maptime(hivetime):
    if hivetime is None:
        return None
    return datetime.fromtimestamp(hivetime/1000).isoformat() + "Z"


def mapstatus(hivestatus):
    if hivestatus == "Open" or hivestatus == "New":
        return "open"
    return "closed"


def maptlp(hivetlp):
    if hivetlp == 0:
        return "White"
    if hivetlp == 1:
        return "Green"
    if hivetlp == 2:
        return "Amber"
    if hivetlp == 3:
        return "Red"
    return "White"


def mapseverity(hiveseverity):
    if hiveseverity == 1:
        return "Low"
    if hiveseverity == 2:
        return "Medium"
    if hiveseverity == 3:
        return "High"
    if hiveseverity == 4:
        return "Very High"
    return "Unknown"

# {
#   "_id": "~16416",
#   "id": "~16416",
#   "createdBy": "jonas@thehive.local",
#   "updatedBy": "jonas@thehive.local",
#   "createdAt": 1638704013583,
#   "updatedAt": 1638704061151,
#   "_type": "case",
#   "caseId": 1,
#   "title": "My Test 1",
#   "description": "My Testcase",
#   "severity": 2,
#   "startDate": 1638703980000,
#   "endDate": null,
#   "impactStatus": null,
#   "resolutionStatus": null,
#   "tags": [],
#   "flag": false,
#   "tlp": 2,
#   "pap": 2,
#   "status": "Open",
#   "summary": null,
#   "owner": "jonas@thehive.local",
#   "customFields": {},
#   "stats": {},
#   "permissions": [ "manageShare", "manageAnalyse", "manageTask", "manageCaseTemplate", "manageCase", "manageUser", "manageProcedure", "managePage", "manageObservable", "manageTag", "manageConfig", "manageAlert", "accessTheHiveFS", "manageAction" ]
# }
def mapcase(hivecase, url, keep_ids):

    s = schema()
    details = {}
    for name, data in hivecase["customFields"].items():
        if "string" in data and data["string"] is not None:
            s.add_string(name)
            details[name] = data["string"]
        if "boolean" in data and data["boolean"] is not None:
            s.add_boolean(name)
            details[name] = data["boolean"]
        if "date" in data and data["date"] is not None:
            s.add_date(name)
            details[name] = maptime(data["date"])
        if "integer" in data and data["integer"] is not None:
            s.add_integer(name)
            details[name] = data["integer"]
        if "float" in data and data["float"] is not None:
            s.add_float(name)
            details[name] = data["float"]

    case = {}
    if keep_ids:
        case["id"] = hivecase["caseId"]

    return {
        "name":       hivecase["title"],
        "type":       "incident",
        "status":     mapstatus(hivecase["status"]),

        "owner":      hivecase["owner"],
        # "write":      hivecase["write"],
        # "read":       hivecase["read"],

        "schema":     json.dumps(s.schema),
        "details":    {
            "tlp":              maptlp(hivecase["tlp"]),
            "pap":              maptlp(hivecase["pap"]),
            "severity":         mapseverity(hivecase["severity"]),
            "description":      hivecase["description"],
            "summary":          hivecase["summary"],
            "tags":             hivecase["tags"],
            "endDate":          maptime(hivecase["endDate"]),
            "resolutionStatus": hivecase["resolutionStatus"],
            "flag":             hivecase["flag"],
        } | details,
        "references": [
            { "name": "TheHive #%d" % hivecase["caseId"], "href": "%s/index.html#!/case/~%s/details" % (url, hivecase["id"]) }
        ],
        #
        # "playbooks":  hivecase["playbooks"],
        #
        "files":    [],
        "comments": [],
        # creator, created, message
        #
        "artifacts": [],
        # name, type, status, enrichment
        #                      name, data

        "created":    maptime(hivecase["createdAt"]),
        "modified":   maptime(hivecase["updatedAt"]),
    } | case

# {
#     "_id": "ce2c00f17132359cb3c50dfbb1901810",
#     "_type": "alert",
#     "artifacts": [],
#     "createdAt": 1495012062014,
#     "createdBy": "myuser",
#     "date": 1495012062016,
#     "description": "N/A",
#     "follow": true,
#     "id": "ce2c00f17132359cb3c50dfbb1901810",
#     "lastSyncDate": 1495012062016,
#     "severity": 2,
#     "source": "instance1",
#     "sourceRef": "alert-ref",
#     "status": "New",
#     "title": "New Alert",
#     "tlp": 2,
#     "type": "external",
#     "user": "myuser"
# }
def mapalert(hivealert, url):
    s = alertschema()
    details = {}

    return {
        "name":       hivealert["title"],
        "type":       "alert",
        "status":     mapstatus(hivealert["status"]),
        "owner":      hivealert["user"],
        "schema":     json.dumps(s.schema),
        "details":    {
            "tlp":              maptlp(hivealert["tlp"]),
            "severity":         mapseverity(hivealert["severity"]),
            "description":      hivealert["description"],
            "source":           hivealert["source"],
            "sourceRef":        hivealert["sourceRef"],
            "type":             hivealert["type"],
        } | details,
        "references": [
            { "name": "TheHive Alerts", "href": "%s/index.html#!/alert/list" % url }
        ],
        "files":    [],
        "comments": [],
        "artifacts": [],
        "created":    maptime(hivealert["createdAt"]),
        "modified":   maptime(hivealert["lastSyncDate"]),
    }

# {
#   "_id": "~41152",
#   "id": "~41152",
#   "createdBy": "jonas@thehive.local",
#   "createdAt": 1638723814523,
#   "_type": "case_artifact",
#   "dataType": "ip",
#   "data": "2.2.2.2",
#   "startDate": 1638723814523,
#   "tlp": 2,
#   "tags": [],
#   "ioc": false,
#   "sighted": false,
#   "message": ".",
#   "reports": {},
#   "stats": {},
#   "ignoreSimilarity": false
# }
def mapobservable(hiveobservable):
    status = "unknown"
    if hiveobservable["ioc"]:
        status = "malicious"
    return {
        "name": hiveobservable["data"],
        "type": hiveobservable["dataType"],
        "status": status,
    }

# {
#   "id": "~12296",
#   "_id": "~12296",
#   "createdBy": "jonas@thehive.local",
#   "createdAt": 1638704029800,
#   "_type": "case_task",
#   "title": "Start",
#   "group": "MyTaskGroup1",
#   "owner": "jonas@thehive.local",
#   "status": "InProgress",
#   "flag": false,
#   "startDate": 1638704115667,
#   "order": 0
# }
# {
#   "_id": "~24656",
#   "id": "~24656",
#   "createdBy": "jonas@thehive.local",
#   "createdAt": 1638729992590,
#   "_type": "case_task_log",
#   "message": "asd",
#   "startDate": 1638729992590,
#   "attachment": {
#     "name": "Chemistry Vector.eps",
#     "hashes": [
#       "adf2d4cd72f4141fe7f8eb4af035596415a29c048d3039be6449008f291258e9",
#       "180f66a6d22b1f09ed198afd814f701e42440e7c",
#       "b28ae347371df003b76cbb8c6199c97e"
#     ],
#     "size": 3421842,
#     "contentType": "application/postscript",
#     "id": "adf2d4cd72f4141fe7f8eb4af035596415a29c048d3039be6449008f291258e9"
#   },
#   "status": "Ok",
#   "owner": "jonas@thehive.local"
# }
def maptasklog(hivetask, hivetasklog):
    message = "**" + hivetask["group"] + ": " + hivetask["title"] + "** (" + hivetask["status"] + ")\n\n"
    message += hivetasklog["message"]
    if 'attachment' in hivetasklog:
        message += "\n\n*Attachment*: " + hivetasklog['attachment']["name"]
    return {
        "creator": hivetasklog["createdBy"],
        "created": maptime(hivetasklog["createdAt"]),
        "message": message,
    }


def run(msg):
    skip_files = msg["payload"]["skip_files"]
    keep_ids = msg["payload"]["keep_ids"]

    from thehive4py.api import TheHiveApi
    import requests
    from minio import Minio

    headers = {"PRIVATE-TOKEN": msg["secrets"]["catalyst_apikey"]}
    # minioclient = Minio("try.catalyst-soar.com:9000", access_key="minio", secret_key="password")
    if not skip_files:
        minioclient = Minio(
            msg["secrets"]["minio_host"],
            access_key=msg["secrets"]["minio_access_key"],
            secret_key=msg["secrets"]["minio_secret_key"])

    # url = "http://localhost:9000"
    url = msg["payload"]["thehiveurl"]
    # api = TheHiveApi(url, "dtUCnzY4h291GIFHJKW/Z2I2SgjTRQqo")
    api = TheHiveApi(url, msg["payload"]["thehivekey"])

    print("find alerts", file=sys.stderr)
    alerts = []
    resp = api.find_alerts(query={}, sort=['-createdAt'], range='all')
    resp.raise_for_status()
    for alert in resp.json():
        alerts.append(mapalert(alert, url))

    if alerts:
        print("create %s alerts" % len(alerts), file=sys.stderr)
        response = requests.post(msg["secrets"]["catalyst_apiurl"] + "/tickets/batch", json=alerts, headers=headers)
        response.raise_for_status()

    print("find incidents", file=sys.stderr)
    incidents = []
    resp = api.find_cases(query={}, sort=['-createdAt'], range='all')
    resp.raise_for_status()
    for case in resp.json():
        incident = mapcase(case, url, keep_ids)
        for observable in api.get_case_observables(case["id"]).json():
            incident["artifacts"].append(mapobservable(observable))
        for task in api.get_case_tasks(case["id"]).json():
            for log in api.get_task_logs(task["id"]).json():
                incident["comments"].append(maptasklog(task, log))
                if 'attachment' in log and not skip_files:
                    incident["files"].append({ "key": log['attachment']["id"], "name": log['attachment']["name"] })

                    bucket_name = "catalyst-%d" % incident["id"]
                    if not minioclient.bucket_exists(bucket_name):
                        minioclient.make_bucket(bucket_name)

                    response = api.download_attachment(log["attachment"]["id"])
                    data = io.BytesIO(response.content)

                    minioclient.put_object(bucket_name, log["attachment"]["id"], data, length=-1, part_size=10*1024*1024)
        incidents.append(incident)

    if incidents:
        if keep_ids:
            print("delete incidents", file=sys.stderr)
            for incident in incidents:
                requests.delete(msg["secrets"]["catalyst_apiurl"] + "/tickets/%d" % incident["id"], headers=headers)
        print("create %d incidents" % len(incidents), file=sys.stderr)
        response = requests.post(msg["secrets"]["catalyst_apiurl"] + "/tickets/batch", json=incidents, headers=headers)
        response.raise_for_status()

    return {"done": True}


print(json.dumps(run(json.loads(sys.argv[1]))))
