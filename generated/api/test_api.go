package api

import "time"

type Args struct {
	Method string
	URL    string
	Data   interface{}
}
type Want struct {
	Status int
	Body   interface{}
}

var Tests = []struct {
	Name string
	Args Args
	Want Want
}{

	{
		Name: "ListAutomations",
		Args: Args{Method: "Get", URL: "/automations"},
		Want: Want{
			Status: 200,
			Body:   []interface{}{map[string]interface{}{"id": "comment", "image": "docker.io/python:3", "script": "", "type": []interface{}{"playbook"}}, map[string]interface{}{"id": "hash.sha1", "image": "docker.io/python:3", "schema": "{\"title\":\"Input\",\"type\":\"object\",\"properties\":{\"default\":{\"type\":\"string\",\"title\":\"Value\"}},\"required\":[\"default\"]}", "script": "", "type": []interface{}{"global", "artifact", "playbook"}}, map[string]interface{}{"id": "vt.hash", "image": "docker.io/python:3", "schema": "{\"title\":\"Input\",\"type\":\"object\",\"properties\":{\"default\":{\"type\":\"string\",\"title\":\"Value\"}},\"required\":[\"default\"]}", "script": "", "type": []interface{}{"global", "artifact", "playbook"}}},
		},
	},

	{
		Name: "CreateAutomation",
		Args: Args{Method: "Post", URL: "/automations", Data: map[string]interface{}{"id": "hash-sha-256", "image": "docker.io/python:3", "script": "import sys\nimport json\nimport hashlib\n\n\ndef run(msg):\n    sha256 = hashlib.sha256(msg['payload']['default'].encode('utf-8'))\n    return {'hash': sha256.hexdigest()}\n\n\nprint(json.dumps(run(json.loads(sys.argv[1]))))\n", "type": []interface{}{"global"}}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"id": "hash-sha-256", "image": "docker.io/python:3", "script": "import sys\nimport json\nimport hashlib\n\n\ndef run(msg):\n    sha256 = hashlib.sha256(msg['payload']['default'].encode('utf-8'))\n    return {'hash': sha256.hexdigest()}\n\n\nprint(json.dumps(run(json.loads(sys.argv[1]))))\n", "type": []interface{}{"global"}},
		},
	},

	{
		Name: "GetAutomation",
		Args: Args{Method: "Get", URL: "/automations/hash.sha1"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"id": "hash.sha1", "image": "docker.io/python:3", "schema": "{\"title\":\"Input\",\"type\":\"object\",\"properties\":{\"default\":{\"type\":\"string\",\"title\":\"Value\"}},\"required\":[\"default\"]}", "script": "#!/usr/bin/env python\n\nimport sys\nimport json\nimport hashlib\n\n\ndef run(msg):\n    sha1 = hashlib.sha1(msg['payload']['default'].encode('utf-8'))\n    return {\"hash\": sha1.hexdigest()}\n\n\nprint(json.dumps(run(json.loads(sys.argv[1]))))\n", "type": []interface{}{"global", "artifact", "playbook"}},
		},
	},

	{
		Name: "UpdateAutomation",
		Args: Args{Method: "Put", URL: "/automations/hash.sha1", Data: map[string]interface{}{"id": "hash.sha1", "image": "docker.io/python:3", "script": "import sys\nimport json\nimport hashlib\n\n\ndef run(msg):\n    sha1 = hashlib.sha1(msg['payload'].encode('utf-8'))\n    return {'hash': sha1.hexdigest()}\n\n\nprint(json.dumps(run(json.loads(sys.argv[1]))))\n", "type": []interface{}{"global", "artifact", "playbook"}}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"id": "hash.sha1", "image": "docker.io/python:3", "script": "import sys\nimport json\nimport hashlib\n\n\ndef run(msg):\n    sha1 = hashlib.sha1(msg['payload'].encode('utf-8'))\n    return {'hash': sha1.hexdigest()}\n\n\nprint(json.dumps(run(json.loads(sys.argv[1]))))\n", "type": []interface{}{"global", "artifact", "playbook"}},
		},
	},

	{
		Name: "DeleteAutomation",
		Args: Args{Method: "Delete", URL: "/automations/hash.sha1"},
		Want: Want{
			Status: 204,
			Body:   nil,
		},
	},

	{
		Name: "CurrentUser",
		Args: Args{Method: "Get", URL: "/currentuser"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"apikey": false, "blocked": false, "id": "bob", "roles": []interface{}{"admin:backup:read", "admin:backup:restore", "admin:dashboard:write", "admin:group:write", "admin:job:read", "admin:job:write", "admin:log:read", "admin:settings:write", "admin:ticket:delete", "admin:user:write", "admin:userdata:read", "admin:userdata:write", "analyst:automation:read", "analyst:currentsettings:write", "analyst:currentuser:read", "analyst:currentuserdata:read", "analyst:dashboard:read", "analyst:file", "analyst:group:read", "analyst:playbook:read", "analyst:rule:read", "analyst:settings:read", "analyst:template:read", "analyst:ticket:read", "analyst:ticket:write", "analyst:tickettype:read", "analyst:user:read", "engineer:automation:write", "engineer:playbook:write", "engineer:rule:write", "engineer:template:write", "engineer:tickettype:write"}},
		},
	},

	{
		Name: "CurrentUserData",
		Args: Args{Method: "Get", URL: "/currentuserdata"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"email": "bob@example.org", "id": "bob", "name": "Bob Bad"},
		},
	},

	{
		Name: "UpdateCurrentUserData",
		Args: Args{Method: "Put", URL: "/currentuserdata", Data: map[string]interface{}{"email": "bob@example.org", "name": "Bob Bad"}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"email": "bob@example.org", "id": "bob", "name": "Bob Bad"},
		},
	},

	{
		Name: "DashboardData",
		Args: Args{Method: "Get", URL: "/dashboard/data?aggregation=type&filter=status+%3D%3D+%22closed%22"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"alert": 2, "incident": 1},
		},
	},

	{
		Name: "ListDashboards",
		Args: Args{Method: "Get", URL: "/dashboards"},
		Want: Want{
			Status: 200,
			Body:   []interface{}{map[string]interface{}{"id": "simple", "name": "Simple", "widgets": []interface{}{map[string]interface{}{"aggregation": "owner", "filter": "status == \"open\"", "name": "open_tickets_per_user", "type": "bar", "width": 4}, map[string]interface{}{"aggregation": "CONCAT(DATE_YEAR(created), \"-\", DATE_ISOWEEK(created) < 10 ? \"0\" : \"\", DATE_ISOWEEK(created))", "name": "tickets_per_week", "type": "line", "width": 8}}}},
		},
	},

	{
		Name: "CreateDashboard",
		Args: Args{Method: "Post", URL: "/dashboards", Data: map[string]interface{}{"name": "My Dashboard", "widgets": []interface{}{}}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"id": "my-dashboard", "name": "My Dashboard", "widgets": []interface{}{}},
		},
	},

	{
		Name: "GetDashboard",
		Args: Args{Method: "Get", URL: "/dashboards/simple"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"id": "simple", "name": "Simple", "widgets": []interface{}{map[string]interface{}{"aggregation": "owner", "filter": "status == \"open\"", "name": "open_tickets_per_user", "type": "bar", "width": 4}, map[string]interface{}{"aggregation": "CONCAT(DATE_YEAR(created), \"-\", DATE_ISOWEEK(created) < 10 ? \"0\" : \"\", DATE_ISOWEEK(created))", "name": "tickets_per_week", "type": "line", "width": 8}}},
		},
	},

	{
		Name: "UpdateDashboard",
		Args: Args{Method: "Put", URL: "/dashboards/simple", Data: map[string]interface{}{"name": "Simple", "widgets": []interface{}{}}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"id": "simple", "name": "Simple", "widgets": []interface{}{}},
		},
	},

	{
		Name: "DeleteDashboard",
		Args: Args{Method: "Delete", URL: "/dashboards/simple"},
		Want: Want{
			Status: 204,
			Body:   nil,
		},
	},

	{
		Name: "ListJobs",
		Args: Args{Method: "Get", URL: "/jobs"},
		Want: Want{
			Status: 200,
			Body:   []interface{}{map[string]interface{}{"automation": "hash.sha1", "id": "b81c2366-ea37-43d2-b61b-03afdc21d985", "payload": "test", "status": "created"}},
		},
	},

	{
		Name: "RunJob",
		Args: Args{Method: "Post", URL: "/jobs", Data: map[string]interface{}{"automation": "hash.sha1", "payload": "test"}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"automation": "hash.sha1", "id": "87390749-2125-4a87-91c5-da7e3f9bebf1", "payload": "test", "status": "created"},
		},
	},

	{
		Name: "GetJob",
		Args: Args{Method: "Get", URL: "/jobs/b81c2366-ea37-43d2-b61b-03afdc21d985"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"automation": "hash.sha1", "id": "b81c2366-ea37-43d2-b61b-03afdc21d985", "payload": "test", "status": "created"},
		},
	},

	{
		Name: "UpdateJob",
		Args: Args{Method: "Put", URL: "/jobs/b81c2366-ea37-43d2-b61b-03afdc21d985", Data: map[string]interface{}{"running": false, "status": "failed"}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"automation": "hash.sha1", "id": "b81c2366-ea37-43d2-b61b-03afdc21d985", "payload": "test", "status": "failed"},
		},
	},

	{
		Name: "GetLogs",
		Args: Args{Method: "Get", URL: "/logs/tickets%252F294511"},
		Want: Want{
			Status: 200,
			Body:   []interface{}{map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "creator": "bob", "message": "Fail run account resist lend solve incident centre priority temperature. Cause change distribution examine location technique shape partner milk customer. Rail tea plate soil report cook railway interpretation breath action. Exercise dream accept park conclusion addition shoot assistance may answer. Gold writer link stop combine hear power name commitment operation. Determine lifespan support grow degree henry exclude detail set religion. Direct library policy convention chain retain discover ride walk student. Gather proposal select march aspect play noise avoid encourage employ. Assessment preserve transport combine wish influence income guess run stand. Charge limit crime ignore statement foundation study issue stop claim.", "reference": "tickets/294511", "type": "manual"}},
		},
	},

	{
		Name: "ListPlaybooks",
		Args: Args{Method: "Get", URL: "/playbooks"},
		Want: Want{
			Status: 200,
			Body:   []interface{}{map[string]interface{}{"id": "malware", "name": "Malware", "yaml": "name: Malware\ntasks:\n  file-or-hash:\n    name: Do you have the file or the hash?\n    type: input\n    schema:\n      title: Malware\n      type: object\n      properties:\n        file:\n          type: string\n          title: \"I have the\"\n          enum: [ \"File\", \"Hash\" ]\n    next:\n      enter-hash: \"file == 'Hash'\"\n      upload: \"file == 'File'\"\n\n  enter-hash:\n    name: Please enter the hash\n    type: input\n    schema:\n      title: Malware\n      type: object\n      properties:\n        hash:\n          type: string\n          title: Please enter the hash value\n          minlength: 32\n    next:\n      virustotal: \"hash != ''\"\n\n  upload:\n    name: Upload the malware\n    type: input\n    schema:\n      title: Malware\n      type: object\n      properties:\n        malware:\n          type: object\n          x-display: file\n          title: Please upload the malware\n    next:\n      hash: \"malware\"\n\n  hash:\n    name: Hash the malware\n    type: automation\n    automation: hash.sha1\n    payload:\n      default: \"playbook.tasks['upload'].data['malware']\"\n    next:\n      virustotal:\n\n  virustotal:\n    name: Send hash to VirusTotal\n    type: automation\n    automation: vt.hash\n    args:\n      hash: \"playbook.tasks['enter-hash'].data['hash'] || playbook.tasks['hash'].data['hash']\"\n    # next:\n    #   known-malware: \"score > 5\"\n    #   sandbox: \"score < 6\" # unknown-malware\n"}, map[string]interface{}{"id": "phishing", "name": "Phishing", "yaml": "name: Phishing\ntasks:\n  board:\n    name: Board Involvement?\n    description: Is a board member involved?\n    type: input\n    schema:\n      properties:\n        boardInvolved:\n          default: false\n          title: A board member is involved.\n          type: boolean\n      required:\n        - boardInvolved\n      title: Board Involvement?\n      type: object\n    next:\n      escalate: \"boardInvolved == true\"\n      mail-available: \"boardInvolved == false\"\n\n  escalate:\n    name: Escalate to CISO\n    description: Please escalate the task to the CISO\n    type: task\n\n  mail-available:\n    name: Mail available\n    type: input\n    schema:\n      oneOf:\n        - properties:\n            mail:\n              title: Mail\n              type: string\n              x-display: textarea\n            schemaKey:\n              const: 'yes'\n              type: string\n          required:\n            - mail\n          title: 'Yes'\n        - properties:\n            schemaKey:\n              const: 'no'\n              type: string\n          title: 'No'\n      title: Mail available\n      type: object\n    next:\n      block-sender: \"schemaKey == 'yes'\"\n      extract-iocs: \"schemaKey == 'yes'\"\n      search-email-gateway: \"schemaKey == 'no'\"\n\n  search-email-gateway:\n    name: Search email gateway\n    description: Please search email-gateway for the phishing mail.\n    type: task\n    next:\n      extract-iocs:\n\n  block-sender:\n    name: Block sender\n    type: task\n    next:\n      extract-iocs:\n\n  extract-iocs:\n    name: Extract IOCs\n    description: Please insert the IOCs\n    type: input\n    schema:\n      properties:\n        iocs:\n          items:\n            type: string\n          title: IOCs\n          type: array\n      title: Extract IOCs\n      type: object\n    next:\n      block-iocs:\n\n  block-iocs:\n    name: Block IOCs\n    type: task\n"}, map[string]interface{}{"id": "simple", "name": "Simple", "yaml": "name: Simple\ntasks:\n  input:\n    name: Enter something to hash\n    type: input\n    schema:\n      title: Something\n      type: object\n      properties:\n        something:\n          type: string\n          title: Something\n          default: \"\"\n    next:\n      hash: \"something != ''\"\n\n  hash:\n    name: Hash the something\n    type: automation\n    automation: hash.sha1\n    payload:\n      default: \"playbook.tasks['input'].data['something']\"\n    next:\n      comment: \"hash != ''\"\n\n  comment:\n    name: Comment the hash\n    type: automation\n    automation: comment\n    payload:\n      default: \"playbook.tasks['hash'].data['hash']\"\n    next:\n      done: \"done\"\n\n  done:\n    name: You can close this case now\n    type: task\n"}},
		},
	},

	{
		Name: "CreatePlaybook",
		Args: Args{Method: "Post", URL: "/playbooks", Data: map[string]interface{}{"yaml": "name: Simple2\ntasks:\n  input:\n    name: Upload malware if possible\n    type: input\n    schema:\n      title: Malware\n      type: object\n      properties:\n        malware:\n          type: string\n          title: Select malware\n          default: \"\"\n    next:\n      hash: \"malware != ''\"\n\n  hash:\n    name: Hash the malware\n    type: automation\n    automation: hash.sha1\n    payload:\n      default: \"playbook.tasks['input'].data['malware']\"\n    next:\n      escalate:\n\n  escalate:\n    name: Escalate to malware team\n    type: task\n"}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"id": "simple-2", "name": "Simple2", "yaml": "name: Simple2\ntasks:\n  input:\n    name: Upload malware if possible\n    type: input\n    schema:\n      title: Malware\n      type: object\n      properties:\n        malware:\n          type: string\n          title: Select malware\n          default: \"\"\n    next:\n      hash: \"malware != ''\"\n\n  hash:\n    name: Hash the malware\n    type: automation\n    automation: hash.sha1\n    payload:\n      default: \"playbook.tasks['input'].data['malware']\"\n    next:\n      escalate:\n\n  escalate:\n    name: Escalate to malware team\n    type: task\n"},
		},
	},

	{
		Name: "GetPlaybook",
		Args: Args{Method: "Get", URL: "/playbooks/simple"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"id": "simple", "name": "Simple", "yaml": "name: Simple\ntasks:\n  input:\n    name: Enter something to hash\n    type: input\n    schema:\n      title: Something\n      type: object\n      properties:\n        something:\n          type: string\n          title: Something\n          default: \"\"\n    next:\n      hash: \"something != ''\"\n\n  hash:\n    name: Hash the something\n    type: automation\n    automation: hash.sha1\n    payload:\n      default: \"playbook.tasks['input'].data['something']\"\n    next:\n      comment: \"hash != ''\"\n\n  comment:\n    name: Comment the hash\n    type: automation\n    automation: comment\n    payload:\n      default: \"playbook.tasks['hash'].data['hash']\"\n    next:\n      done: \"done\"\n\n  done:\n    name: You can close this case now\n    type: task\n"},
		},
	},

	{
		Name: "UpdatePlaybook",
		Args: Args{Method: "Put", URL: "/playbooks/simple", Data: map[string]interface{}{"yaml": "name: Simple\ntasks:\n  input:\n    name: Upload malware if possible\n    type: input\n    schema:\n      title: Malware\n      type: object\n      properties:\n        malware:\n          type: string\n          title: Select malware\n          default: \"\"\n    next:\n      hash: \"malware != ''\"\n\n  hash:\n    name: Hash the malware\n    type: automation\n    automation: hash.sha1\n    payload:\n      default: \"playbook.tasks['input'].data['malware']\"\n    next:\n      escalate:\n\n  escalate:\n    name: Escalate to malware team\n    type: task\n"}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"id": "simple", "name": "Simple", "yaml": "name: Simple\ntasks:\n  input:\n    name: Upload malware if possible\n    type: input\n    schema:\n      title: Malware\n      type: object\n      properties:\n        malware:\n          type: string\n          title: Select malware\n          default: \"\"\n    next:\n      hash: \"malware != ''\"\n\n  hash:\n    name: Hash the malware\n    type: automation\n    automation: hash.sha1\n    payload:\n      default: \"playbook.tasks['input'].data['malware']\"\n    next:\n      escalate:\n\n  escalate:\n    name: Escalate to malware team\n    type: task\n"},
		},
	},

	{
		Name: "DeletePlaybook",
		Args: Args{Method: "Delete", URL: "/playbooks/simple"},
		Want: Want{
			Status: 204,
			Body:   nil,
		},
	},

	{
		Name: "GetSettings",
		Args: Args{Method: "Get", URL: "/settings"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"artifactKinds": []interface{}{map[string]interface{}{"icon": "mdi-server", "id": "asset", "name": "Asset"}, map[string]interface{}{"icon": "mdi-bullseye", "id": "ioc", "name": "IOC"}}, "artifactStates": []interface{}{map[string]interface{}{"color": "info", "icon": "mdi-help-circle-outline", "id": "unknown", "name": "Unknown"}, map[string]interface{}{"color": "error", "icon": "mdi-skull", "id": "malicious", "name": "Malicious"}, map[string]interface{}{"color": "success", "icon": "mdi-check", "id": "clean", "name": "Clean"}}, "roles": []interface{}{"admin:backup:read", "admin:backup:restore", "admin:dashboard:write", "admin:group:write", "admin:job:read", "admin:job:write", "admin:log:read", "admin:settings:write", "admin:ticket:delete", "admin:user:write", "admin:userdata:read", "admin:userdata:write", "analyst:automation:read", "analyst:currentsettings:write", "analyst:currentuser:read", "analyst:currentuserdata:read", "analyst:dashboard:read", "analyst:file", "analyst:group:read", "analyst:playbook:read", "analyst:rule:read", "analyst:settings:read", "analyst:template:read", "analyst:ticket:read", "analyst:ticket:write", "analyst:tickettype:read", "analyst:user:read", "engineer:automation:write", "engineer:playbook:write", "engineer:rule:write", "engineer:template:write", "engineer:tickettype:write"}, "ticketTypes": []interface{}{map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-alert", "id": "alert", "name": "Alerts"}, map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-radioactive", "id": "incident", "name": "Incidents"}, map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-fingerprint", "id": "investigation", "name": "Forensic Investigations"}, map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-target", "id": "hunt", "name": "Threat Hunting"}}, "tier": "community", "timeformat": "YYYY-MM-DDThh:mm:ss", "version": "0.0.0-test"},
		},
	},

	{
		Name: "SaveSettings",
		Args: Args{Method: "Post", URL: "/settings", Data: map[string]interface{}{"artifactKinds": []interface{}{map[string]interface{}{"icon": "mdi-server", "id": "asset", "name": "Asset"}, map[string]interface{}{"icon": "mdi-bullseye", "id": "ioc", "name": "IOC"}}, "artifactStates": []interface{}{map[string]interface{}{"color": "info", "icon": "mdi-help-circle-outline", "id": "unknown", "name": "Unknown"}, map[string]interface{}{"color": "error", "icon": "mdi-skull", "id": "malicious", "name": "Malicious"}, map[string]interface{}{"color": "success", "icon": "mdi-check", "id": "clean", "name": "Clean"}}, "timeformat": "YYYY-MM-DDThh:mm:ss"}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"artifactKinds": []interface{}{map[string]interface{}{"icon": "mdi-server", "id": "asset", "name": "Asset"}, map[string]interface{}{"icon": "mdi-bullseye", "id": "ioc", "name": "IOC"}}, "artifactStates": []interface{}{map[string]interface{}{"color": "info", "icon": "mdi-help-circle-outline", "id": "unknown", "name": "Unknown"}, map[string]interface{}{"color": "error", "icon": "mdi-skull", "id": "malicious", "name": "Malicious"}, map[string]interface{}{"color": "success", "icon": "mdi-check", "id": "clean", "name": "Clean"}}, "roles": []interface{}{"admin:backup:read", "admin:backup:restore", "admin:dashboard:write", "admin:group:write", "admin:job:read", "admin:job:write", "admin:log:read", "admin:settings:write", "admin:ticket:delete", "admin:user:write", "admin:userdata:read", "admin:userdata:write", "analyst:automation:read", "analyst:currentsettings:write", "analyst:currentuser:read", "analyst:currentuserdata:read", "analyst:dashboard:read", "analyst:file", "analyst:group:read", "analyst:playbook:read", "analyst:rule:read", "analyst:settings:read", "analyst:template:read", "analyst:ticket:read", "analyst:ticket:write", "analyst:tickettype:read", "analyst:user:read", "engineer:automation:write", "engineer:playbook:write", "engineer:rule:write", "engineer:template:write", "engineer:tickettype:write"}, "ticketTypes": []interface{}{map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-alert", "id": "alert", "name": "Alerts"}, map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-radioactive", "id": "incident", "name": "Incidents"}, map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-fingerprint", "id": "investigation", "name": "Forensic Investigations"}, map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-target", "id": "hunt", "name": "Threat Hunting"}}, "tier": "community", "timeformat": "YYYY-MM-DDThh:mm:ss", "version": "0.0.0-test"},
		},
	},

	{
		Name: "GetStatistics",
		Args: Args{Method: "Get", URL: "/statistics"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"open_tickets_per_user": map[string]interface{}{}, "tickets_per_type": map[string]interface{}{"alert": 2, "incident": 1}, "tickets_per_week": map[string]interface{}{"2021-39": 3}, "unassigned": 0},
		},
	},

	{
		Name: "ListTasks",
		Args: Args{Method: "Get", URL: "/tasks"},
		Want: Want{
			Status: 200,
			Body:   nil,
		},
	},

	{
		Name: "ListTemplates",
		Args: Args{Method: "Get", URL: "/templates"},
		Want: Want{
			Status: 200,
			Body:   []interface{}{map[string]interface{}{"id": "default", "name": "Default", "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Default\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n"}},
		},
	},

	{
		Name: "CreateTemplate",
		Args: Args{Method: "Post", URL: "/templates", Data: map[string]interface{}{"name": "My Template", "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n"}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"id": "my-template", "name": "My Template", "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n"},
		},
	},

	{
		Name: "GetTemplate",
		Args: Args{Method: "Get", URL: "/templates/default"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"id": "default", "name": "Default", "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Default\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n"},
		},
	},

	{
		Name: "UpdateTemplate",
		Args: Args{Method: "Put", URL: "/templates/default", Data: map[string]interface{}{"name": "My Template", "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n"}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"id": "default", "name": "My Template", "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n"},
		},
	},

	{
		Name: "DeleteTemplate",
		Args: Args{Method: "Delete", URL: "/templates/default"},
		Want: Want{
			Status: 204,
			Body:   nil,
		},
	},

	{
		Name: "ListTickets",
		Args: Args{Method: "Get", URL: "/tickets"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"count": 3, "tickets": []interface{}{map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}, map[string]interface{}{"name": "leadreintermediate.io", "status": "malicious"}}, "created": time.Date(2021, time.October, 2, 16, 4, 59, 78206000, time.UTC), "id": 8123, "modified": time.Date(2021, time.October, 2, 16, 4, 59, 78206000, time.UTC), "name": "live zebra", "owner": "demo", "playbooks": map[string]interface{}{"phishing": map[string]interface{}{"name": "Phishing", "tasks": map[string]interface{}{"block-iocs": map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block IOCs", "type": "task"}, "block-sender": map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block sender", "next": map[string]interface{}{"extract-iocs": ""}, "type": "task"}, "board": map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}, "escalate": map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Escalate to CISO", "type": "task"}, "extract-iocs": map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Extract IOCs", "next": map[string]interface{}{"block-iocs": ""}, "schema": map[string]interface{}{"properties": map[string]interface{}{"iocs": map[string]interface{}{"items": map[string]interface{}{"type": "string"}, "title": "IOCs", "type": "array"}}, "title": "Extract IOCs", "type": "object"}, "type": "input"}, "mail-available": map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Mail available", "next": map[string]interface{}{"block-sender": "schemaKey == 'yes'", "extract-iocs": "schemaKey == 'yes'", "search-email-gateway": "schemaKey == 'no'"}, "schema": map[string]interface{}{"oneOf": []interface{}{map[string]interface{}{"properties": map[string]interface{}{"mail": map[string]interface{}{"title": "Mail", "type": "string", "x-display": "textarea"}, "schemaKey": map[string]interface{}{"const": "yes", "type": "string"}}, "required": []interface{}{"mail"}, "title": "Yes"}, map[string]interface{}{"properties": map[string]interface{}{"schemaKey": map[string]interface{}{"const": "no", "type": "string"}}, "title": "No"}}, "title": "Mail available", "type": "object"}, "type": "input"}, "search-email-gateway": map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Search email gateway", "next": map[string]interface{}{"extract-iocs": ""}, "type": "task"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"}, map[string]interface{}{"created": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "id": 8125, "modified": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "name": "phishing from selenafadel@von.com detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "https://www.seniorleading-edge.name/users/efficient", "name": "recovery"}, map[string]interface{}{"href": "http://www.dynamicseamless.com/clicks-and-mortar", "name": "force"}, map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}, "schema": "{}", "status": "closed", "type": "alert"}, map[string]interface{}{"created": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "id": 8126, "modified": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "name": "Surfaceintroduce virus detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.centralworld-class.io/synthesize", "name": "university"}, map[string]interface{}{"href": "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", "name": "goal"}, map[string]interface{}{"href": "http://www.chiefsyndicate.io/action-items", "name": "unemployment"}}, "schema": "{}", "status": "closed", "type": "alert"}}},
		},
	},

	{
		Name: "CreateTicket",
		Args: Args{Method: "Post", URL: "/tickets", Data: map[string]interface{}{"id": 123, "name": "Wannacry infection", "owner": "bob", "status": "open", "type": "incident"}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "id": 123, "modified": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "name": "Wannacry infection", "owner": "bob", "schema": "{}", "status": "open", "type": "incident"},
		},
	},

	{
		Name: "CreateTicketBatch",
		Args: Args{Method: "Post", URL: "/tickets/batch", Data: []interface{}{map[string]interface{}{"id": 123, "name": "Wannacry infection", "owner": "bob", "status": "open", "type": "incident"}}},
		Want: Want{
			Status: 204,
			Body:   nil,
		},
	},

	{
		Name: "GetTicket",
		Args: Args{Method: "Get", URL: "/tickets/8125"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"created": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "id": 8125, "modified": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "name": "phishing from selenafadel@von.com detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "https://www.seniorleading-edge.name/users/efficient", "name": "recovery"}, map[string]interface{}{"href": "http://www.dynamicseamless.com/clicks-and-mortar", "name": "force"}, map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}, "schema": "{}", "status": "closed", "tickets": []interface{}{map[string]interface{}{"created": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "id": 8126, "modified": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "name": "Surfaceintroduce virus detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.centralworld-class.io/synthesize", "name": "university"}, map[string]interface{}{"href": "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", "name": "goal"}, map[string]interface{}{"href": "http://www.chiefsyndicate.io/action-items", "name": "unemployment"}}, "schema": "{}", "status": "closed", "type": "alert"}}, "type": "alert"},
		},
	},

	{
		Name: "UpdateTicket",
		Args: Args{Method: "Put", URL: "/tickets/8125", Data: map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "modified": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "name": "phishing from selenafadel@von.org detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "https://www.seniorleading-edge.name/users/efficient", "name": "recovery"}, map[string]interface{}{"href": "http://www.dynamicseamless.com/clicks-and-mortar", "name": "force"}, map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}, "schema": "{}", "status": "closed", "type": "alert"}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "id": 8125, "modified": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "name": "phishing from selenafadel@von.org detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "https://www.seniorleading-edge.name/users/efficient", "name": "recovery"}, map[string]interface{}{"href": "http://www.dynamicseamless.com/clicks-and-mortar", "name": "force"}, map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}, "schema": "{}", "status": "closed", "tickets": []interface{}{map[string]interface{}{"created": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "id": 8126, "modified": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "name": "Surfaceintroduce virus detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.centralworld-class.io/synthesize", "name": "university"}, map[string]interface{}{"href": "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", "name": "goal"}, map[string]interface{}{"href": "http://www.chiefsyndicate.io/action-items", "name": "unemployment"}}, "schema": "{}", "status": "closed", "type": "alert"}}, "type": "alert"},
		},
	},

	{
		Name: "DeleteTicket",
		Args: Args{Method: "Delete", URL: "/tickets/8125"},
		Want: Want{
			Status: 204,
			Body:   nil,
		},
	},

	{
		Name: "AddArtifact",
		Args: Args{Method: "Post", URL: "/tickets/8123/artifacts", Data: map[string]interface{}{"name": "2.2.2.2"}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}, map[string]interface{}{"name": "leadreintermediate.io", "status": "malicious"}, map[string]interface{}{"name": "2.2.2.2", "status": "unknown", "type": "ip"}}, "created": time.Date(2021, time.October, 2, 16, 4, 59, 78206000, time.UTC), "id": 8123, "modified": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "name": "live zebra", "owner": "demo", "playbooks": map[string]interface{}{"phishing": map[string]interface{}{"name": "Phishing", "tasks": map[string]interface{}{"block-iocs": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block IOCs", "order": 6, "type": "task"}, "block-sender": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block sender", "next": map[string]interface{}{"extract-iocs": ""}, "order": 3, "type": "task"}, "board": map[string]interface{}{"active": true, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "order": 0, "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}, "escalate": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Escalate to CISO", "order": 1, "type": "task"}, "extract-iocs": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Extract IOCs", "next": map[string]interface{}{"block-iocs": ""}, "order": 5, "schema": map[string]interface{}{"properties": map[string]interface{}{"iocs": map[string]interface{}{"items": map[string]interface{}{"type": "string"}, "title": "IOCs", "type": "array"}}, "title": "Extract IOCs", "type": "object"}, "type": "input"}, "mail-available": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Mail available", "next": map[string]interface{}{"block-sender": "schemaKey == 'yes'", "extract-iocs": "schemaKey == 'yes'", "search-email-gateway": "schemaKey == 'no'"}, "order": 2, "schema": map[string]interface{}{"oneOf": []interface{}{map[string]interface{}{"properties": map[string]interface{}{"mail": map[string]interface{}{"title": "Mail", "type": "string", "x-display": "textarea"}, "schemaKey": map[string]interface{}{"const": "yes", "type": "string"}}, "required": []interface{}{"mail"}, "title": "Yes"}, map[string]interface{}{"properties": map[string]interface{}{"schemaKey": map[string]interface{}{"const": "no", "type": "string"}}, "title": "No"}}, "title": "Mail available", "type": "object"}, "type": "input"}, "search-email-gateway": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Search email gateway", "next": map[string]interface{}{"extract-iocs": ""}, "order": 4, "type": "task"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"},
		},
	},

	{
		Name: "GetArtifact",
		Args: Args{Method: "Get", URL: "/tickets/8123/artifacts/leadreintermediate.io"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"name": "leadreintermediate.io", "status": "malicious"},
		},
	},

	{
		Name: "SetArtifact",
		Args: Args{Method: "Put", URL: "/tickets/8123/artifacts/leadreintermediate.io", Data: map[string]interface{}{"name": "leadreintermediate.io", "status": "clean"}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}, map[string]interface{}{"name": "leadreintermediate.io", "status": "clean"}}, "created": time.Date(2021, time.October, 2, 16, 4, 59, 78206000, time.UTC), "id": 8123, "modified": time.Date(2021, time.October, 2, 16, 4, 59, 78206000, time.UTC), "name": "live zebra", "owner": "demo", "playbooks": map[string]interface{}{"phishing": map[string]interface{}{"name": "Phishing", "tasks": map[string]interface{}{"block-iocs": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block IOCs", "order": 6, "type": "task"}, "block-sender": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block sender", "next": map[string]interface{}{"extract-iocs": ""}, "order": 3, "type": "task"}, "board": map[string]interface{}{"active": true, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "order": 0, "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}, "escalate": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Escalate to CISO", "order": 1, "type": "task"}, "extract-iocs": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Extract IOCs", "next": map[string]interface{}{"block-iocs": ""}, "order": 5, "schema": map[string]interface{}{"properties": map[string]interface{}{"iocs": map[string]interface{}{"items": map[string]interface{}{"type": "string"}, "title": "IOCs", "type": "array"}}, "title": "Extract IOCs", "type": "object"}, "type": "input"}, "mail-available": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Mail available", "next": map[string]interface{}{"block-sender": "schemaKey == 'yes'", "extract-iocs": "schemaKey == 'yes'", "search-email-gateway": "schemaKey == 'no'"}, "order": 2, "schema": map[string]interface{}{"oneOf": []interface{}{map[string]interface{}{"properties": map[string]interface{}{"mail": map[string]interface{}{"title": "Mail", "type": "string", "x-display": "textarea"}, "schemaKey": map[string]interface{}{"const": "yes", "type": "string"}}, "required": []interface{}{"mail"}, "title": "Yes"}, map[string]interface{}{"properties": map[string]interface{}{"schemaKey": map[string]interface{}{"const": "no", "type": "string"}}, "title": "No"}}, "title": "Mail available", "type": "object"}, "type": "input"}, "search-email-gateway": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Search email gateway", "next": map[string]interface{}{"extract-iocs": ""}, "order": 4, "type": "task"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"},
		},
	},

	{
		Name: "RemoveArtifact",
		Args: Args{Method: "Delete", URL: "/tickets/8123/artifacts/leadreintermediate.io"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}}, "created": time.Date(2021, time.October, 2, 16, 4, 59, 78206000, time.UTC), "id": 8123, "modified": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "name": "live zebra", "owner": "demo", "playbooks": map[string]interface{}{"phishing": map[string]interface{}{"name": "Phishing", "tasks": map[string]interface{}{"block-iocs": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block IOCs", "order": 6, "type": "task"}, "block-sender": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block sender", "next": map[string]interface{}{"extract-iocs": ""}, "order": 3, "type": "task"}, "board": map[string]interface{}{"active": true, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "order": 0, "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}, "escalate": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Escalate to CISO", "order": 1, "type": "task"}, "extract-iocs": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Extract IOCs", "next": map[string]interface{}{"block-iocs": ""}, "order": 5, "schema": map[string]interface{}{"properties": map[string]interface{}{"iocs": map[string]interface{}{"items": map[string]interface{}{"type": "string"}, "title": "IOCs", "type": "array"}}, "title": "Extract IOCs", "type": "object"}, "type": "input"}, "mail-available": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Mail available", "next": map[string]interface{}{"block-sender": "schemaKey == 'yes'", "extract-iocs": "schemaKey == 'yes'", "search-email-gateway": "schemaKey == 'no'"}, "order": 2, "schema": map[string]interface{}{"oneOf": []interface{}{map[string]interface{}{"properties": map[string]interface{}{"mail": map[string]interface{}{"title": "Mail", "type": "string", "x-display": "textarea"}, "schemaKey": map[string]interface{}{"const": "yes", "type": "string"}}, "required": []interface{}{"mail"}, "title": "Yes"}, map[string]interface{}{"properties": map[string]interface{}{"schemaKey": map[string]interface{}{"const": "no", "type": "string"}}, "title": "No"}}, "title": "Mail available", "type": "object"}, "type": "input"}, "search-email-gateway": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Search email gateway", "next": map[string]interface{}{"extract-iocs": ""}, "order": 4, "type": "task"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"},
		},
	},

	{
		Name: "EnrichArtifact",
		Args: Args{Method: "Post", URL: "/tickets/8123/artifacts/leadreintermediate.io/enrich", Data: map[string]interface{}{"data": map[string]interface{}{"hash": "b7a067a742c20d07a7456646de89bc2d408a1153"}, "name": "hash.sha1"}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}, map[string]interface{}{"enrichments": map[string]interface{}{"hash.sha1": map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "data": map[string]interface{}{"hash": "b7a067a742c20d07a7456646de89bc2d408a1153"}, "name": "hash.sha1"}}, "name": "leadreintermediate.io", "status": "malicious"}}, "created": time.Date(2021, time.October, 2, 16, 4, 59, 78206000, time.UTC), "id": 8123, "modified": time.Date(2021, time.October, 2, 16, 4, 59, 78206000, time.UTC), "name": "live zebra", "owner": "demo", "playbooks": map[string]interface{}{"phishing": map[string]interface{}{"name": "Phishing", "tasks": map[string]interface{}{"block-iocs": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block IOCs", "order": 6, "type": "task"}, "block-sender": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block sender", "next": map[string]interface{}{"extract-iocs": ""}, "order": 3, "type": "task"}, "board": map[string]interface{}{"active": true, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "order": 0, "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}, "escalate": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Escalate to CISO", "order": 1, "type": "task"}, "extract-iocs": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Extract IOCs", "next": map[string]interface{}{"block-iocs": ""}, "order": 5, "schema": map[string]interface{}{"properties": map[string]interface{}{"iocs": map[string]interface{}{"items": map[string]interface{}{"type": "string"}, "title": "IOCs", "type": "array"}}, "title": "Extract IOCs", "type": "object"}, "type": "input"}, "mail-available": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Mail available", "next": map[string]interface{}{"block-sender": "schemaKey == 'yes'", "extract-iocs": "schemaKey == 'yes'", "search-email-gateway": "schemaKey == 'no'"}, "order": 2, "schema": map[string]interface{}{"oneOf": []interface{}{map[string]interface{}{"properties": map[string]interface{}{"mail": map[string]interface{}{"title": "Mail", "type": "string", "x-display": "textarea"}, "schemaKey": map[string]interface{}{"const": "yes", "type": "string"}}, "required": []interface{}{"mail"}, "title": "Yes"}, map[string]interface{}{"properties": map[string]interface{}{"schemaKey": map[string]interface{}{"const": "no", "type": "string"}}, "title": "No"}}, "title": "Mail available", "type": "object"}, "type": "input"}, "search-email-gateway": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Search email gateway", "next": map[string]interface{}{"extract-iocs": ""}, "order": 4, "type": "task"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"},
		},
	},

	{
		Name: "RunArtifact",
		Args: Args{Method: "Post", URL: "/tickets/8123/artifacts/leadreintermediate.io/run/hash.sha1"},
		Want: Want{
			Status: 204,
			Body:   nil,
		},
	},

	{
		Name: "AddComment",
		Args: Args{Method: "Post", URL: "/tickets/8125/comments", Data: map[string]interface{}{"message": "My first comment"}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"comments": []interface{}{map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "creator": "bob", "message": "My first comment"}}, "created": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "id": 8125, "modified": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "name": "phishing from selenafadel@von.com detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "https://www.seniorleading-edge.name/users/efficient", "name": "recovery"}, map[string]interface{}{"href": "http://www.dynamicseamless.com/clicks-and-mortar", "name": "force"}, map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}, "schema": "{}", "status": "closed", "tickets": []interface{}{map[string]interface{}{"created": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "id": 8126, "modified": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "name": "Surfaceintroduce virus detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.centralworld-class.io/synthesize", "name": "university"}, map[string]interface{}{"href": "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", "name": "goal"}, map[string]interface{}{"href": "http://www.chiefsyndicate.io/action-items", "name": "unemployment"}}, "schema": "{}", "status": "closed", "type": "alert"}}, "type": "alert"},
		},
	},

	{
		Name: "RemoveComment",
		Args: Args{Method: "Delete", URL: "/tickets/8123/comments/0"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}, map[string]interface{}{"name": "leadreintermediate.io", "status": "malicious"}}, "created": time.Date(2021, time.October, 2, 16, 4, 59, 78206000, time.UTC), "id": 8123, "modified": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "name": "live zebra", "owner": "demo", "playbooks": map[string]interface{}{"phishing": map[string]interface{}{"name": "Phishing", "tasks": map[string]interface{}{"block-iocs": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block IOCs", "order": 6, "type": "task"}, "block-sender": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block sender", "next": map[string]interface{}{"extract-iocs": ""}, "order": 3, "type": "task"}, "board": map[string]interface{}{"active": true, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "order": 0, "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}, "escalate": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Escalate to CISO", "order": 1, "type": "task"}, "extract-iocs": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Extract IOCs", "next": map[string]interface{}{"block-iocs": ""}, "order": 5, "schema": map[string]interface{}{"properties": map[string]interface{}{"iocs": map[string]interface{}{"items": map[string]interface{}{"type": "string"}, "title": "IOCs", "type": "array"}}, "title": "Extract IOCs", "type": "object"}, "type": "input"}, "mail-available": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Mail available", "next": map[string]interface{}{"block-sender": "schemaKey == 'yes'", "extract-iocs": "schemaKey == 'yes'", "search-email-gateway": "schemaKey == 'no'"}, "order": 2, "schema": map[string]interface{}{"oneOf": []interface{}{map[string]interface{}{"properties": map[string]interface{}{"mail": map[string]interface{}{"title": "Mail", "type": "string", "x-display": "textarea"}, "schemaKey": map[string]interface{}{"const": "yes", "type": "string"}}, "required": []interface{}{"mail"}, "title": "Yes"}, map[string]interface{}{"properties": map[string]interface{}{"schemaKey": map[string]interface{}{"const": "no", "type": "string"}}, "title": "No"}}, "title": "Mail available", "type": "object"}, "type": "input"}, "search-email-gateway": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Search email gateway", "next": map[string]interface{}{"extract-iocs": ""}, "order": 4, "type": "task"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"},
		},
	},

	{
		Name: "AddTicketPlaybook",
		Args: Args{Method: "Post", URL: "/tickets/8125/playbooks", Data: map[string]interface{}{"yaml": "name: Simple\ntasks:\n  input:\n    name: Upload malware if possible\n    type: input\n    schema:\n      title: Malware\n      type: object\n      properties:\n        malware:\n          type: string\n          title: Select malware\n          default: \"\"\n    next:\n      hash: \"malware != ''\"\n\n  hash:\n    name: Hash the malware\n    type: automation\n    automation: hash.sha1\n    payload:\n      default: \"playbook.tasks['input'].data['malware']\"\n    next:\n      escalate:\n\n  escalate:\n    name: Escalate to malware team\n    type: task\n"}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"created": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "id": 8125, "modified": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "name": "phishing from selenafadel@von.com detected", "owner": "demo", "playbooks": map[string]interface{}{"simple": map[string]interface{}{"name": "Simple", "tasks": map[string]interface{}{"escalate": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Escalate to malware team", "order": 2, "type": "task"}, "hash": map[string]interface{}{"active": false, "automation": "hash.sha1", "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Hash the malware", "next": map[string]interface{}{"escalate": ""}, "order": 1, "payload": map[string]interface{}{"default": "playbook.tasks['input'].data['malware']"}, "type": "automation"}, "input": map[string]interface{}{"active": true, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Upload malware if possible", "next": map[string]interface{}{"hash": "malware != ''"}, "order": 0, "schema": map[string]interface{}{"properties": map[string]interface{}{"malware": map[string]interface{}{"default": "", "title": "Select malware", "type": "string"}}, "title": "Malware", "type": "object"}, "type": "input"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.seniorleading-edge.name/users/efficient", "name": "recovery"}, map[string]interface{}{"href": "http://www.dynamicseamless.com/clicks-and-mortar", "name": "force"}, map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}, "schema": "{}", "status": "closed", "tickets": []interface{}{map[string]interface{}{"created": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "id": 8126, "modified": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "name": "Surfaceintroduce virus detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.centralworld-class.io/synthesize", "name": "university"}, map[string]interface{}{"href": "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", "name": "goal"}, map[string]interface{}{"href": "http://www.chiefsyndicate.io/action-items", "name": "unemployment"}}, "schema": "{}", "status": "closed", "type": "alert"}}, "type": "alert"},
		},
	},

	{
		Name: "RemoveTicketPlaybook",
		Args: Args{Method: "Delete", URL: "/tickets/8123/playbooks/phishing"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}, map[string]interface{}{"name": "leadreintermediate.io", "status": "malicious"}}, "created": time.Date(2021, time.October, 2, 16, 4, 59, 78206000, time.UTC), "id": 8123, "modified": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "name": "live zebra", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"},
		},
	},

	{
		Name: "SetTaskData",
		Args: Args{Method: "Put", URL: "/tickets/8123/playbooks/phishing/task/board", Data: map[string]interface{}{"boardInvolved": true}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}, map[string]interface{}{"name": "leadreintermediate.io", "status": "malicious"}}, "created": time.Date(2021, time.October, 2, 16, 4, 59, 78206000, time.UTC), "id": 8123, "modified": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "name": "live zebra", "owner": "demo", "playbooks": map[string]interface{}{"phishing": map[string]interface{}{"name": "Phishing", "tasks": map[string]interface{}{"block-iocs": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block IOCs", "order": 6, "type": "task"}, "block-sender": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block sender", "next": map[string]interface{}{"extract-iocs": ""}, "order": 3, "type": "task"}, "board": map[string]interface{}{"active": true, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "data": map[string]interface{}{"boardInvolved": true}, "done": false, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "order": 0, "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}, "escalate": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Escalate to CISO", "order": 1, "type": "task"}, "extract-iocs": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Extract IOCs", "next": map[string]interface{}{"block-iocs": ""}, "order": 5, "schema": map[string]interface{}{"properties": map[string]interface{}{"iocs": map[string]interface{}{"items": map[string]interface{}{"type": "string"}, "title": "IOCs", "type": "array"}}, "title": "Extract IOCs", "type": "object"}, "type": "input"}, "mail-available": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Mail available", "next": map[string]interface{}{"block-sender": "schemaKey == 'yes'", "extract-iocs": "schemaKey == 'yes'", "search-email-gateway": "schemaKey == 'no'"}, "order": 2, "schema": map[string]interface{}{"oneOf": []interface{}{map[string]interface{}{"properties": map[string]interface{}{"mail": map[string]interface{}{"title": "Mail", "type": "string", "x-display": "textarea"}, "schemaKey": map[string]interface{}{"const": "yes", "type": "string"}}, "required": []interface{}{"mail"}, "title": "Yes"}, map[string]interface{}{"properties": map[string]interface{}{"schemaKey": map[string]interface{}{"const": "no", "type": "string"}}, "title": "No"}}, "title": "Mail available", "type": "object"}, "type": "input"}, "search-email-gateway": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Search email gateway", "next": map[string]interface{}{"extract-iocs": ""}, "order": 4, "type": "task"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"},
		},
	},

	{
		Name: "CompleteTask",
		Args: Args{Method: "Put", URL: "/tickets/8123/playbooks/phishing/task/board/complete", Data: map[string]interface{}{"boardInvolved": true}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}, map[string]interface{}{"name": "leadreintermediate.io", "status": "malicious"}}, "created": time.Date(2021, time.October, 2, 16, 4, 59, 78206000, time.UTC), "id": 8123, "modified": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "name": "live zebra", "owner": "demo", "playbooks": map[string]interface{}{"phishing": map[string]interface{}{"name": "Phishing", "tasks": map[string]interface{}{"block-iocs": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block IOCs", "order": 6, "type": "task"}, "block-sender": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block sender", "next": map[string]interface{}{"extract-iocs": ""}, "order": 3, "type": "task"}, "board": map[string]interface{}{"active": false, "closed": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "data": map[string]interface{}{"boardInvolved": true}, "done": true, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "order": 0, "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}, "escalate": map[string]interface{}{"active": true, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Escalate to CISO", "order": 1, "type": "task"}, "extract-iocs": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Extract IOCs", "next": map[string]interface{}{"block-iocs": ""}, "order": 5, "schema": map[string]interface{}{"properties": map[string]interface{}{"iocs": map[string]interface{}{"items": map[string]interface{}{"type": "string"}, "title": "IOCs", "type": "array"}}, "title": "Extract IOCs", "type": "object"}, "type": "input"}, "mail-available": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Mail available", "next": map[string]interface{}{"block-sender": "schemaKey == 'yes'", "extract-iocs": "schemaKey == 'yes'", "search-email-gateway": "schemaKey == 'no'"}, "order": 2, "schema": map[string]interface{}{"oneOf": []interface{}{map[string]interface{}{"properties": map[string]interface{}{"mail": map[string]interface{}{"title": "Mail", "type": "string", "x-display": "textarea"}, "schemaKey": map[string]interface{}{"const": "yes", "type": "string"}}, "required": []interface{}{"mail"}, "title": "Yes"}, map[string]interface{}{"properties": map[string]interface{}{"schemaKey": map[string]interface{}{"const": "no", "type": "string"}}, "title": "No"}}, "title": "Mail available", "type": "object"}, "type": "input"}, "search-email-gateway": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Search email gateway", "next": map[string]interface{}{"extract-iocs": ""}, "order": 4, "type": "task"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"},
		},
	},

	{
		Name: "SetTaskOwner",
		Args: Args{Method: "Put", URL: "/tickets/8123/playbooks/phishing/task/board/owner", Data: "eve"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}, map[string]interface{}{"name": "leadreintermediate.io", "status": "malicious"}}, "created": time.Date(2021, time.October, 2, 16, 4, 59, 78206000, time.UTC), "id": 8123, "modified": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "name": "live zebra", "owner": "demo", "playbooks": map[string]interface{}{"phishing": map[string]interface{}{"name": "Phishing", "tasks": map[string]interface{}{"block-iocs": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block IOCs", "order": 6, "type": "task"}, "block-sender": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block sender", "next": map[string]interface{}{"extract-iocs": ""}, "order": 3, "type": "task"}, "board": map[string]interface{}{"active": true, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "order": 0, "owner": "eve", "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}, "escalate": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Escalate to CISO", "order": 1, "type": "task"}, "extract-iocs": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Extract IOCs", "next": map[string]interface{}{"block-iocs": ""}, "order": 5, "schema": map[string]interface{}{"properties": map[string]interface{}{"iocs": map[string]interface{}{"items": map[string]interface{}{"type": "string"}, "title": "IOCs", "type": "array"}}, "title": "Extract IOCs", "type": "object"}, "type": "input"}, "mail-available": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Mail available", "next": map[string]interface{}{"block-sender": "schemaKey == 'yes'", "extract-iocs": "schemaKey == 'yes'", "search-email-gateway": "schemaKey == 'no'"}, "order": 2, "schema": map[string]interface{}{"oneOf": []interface{}{map[string]interface{}{"properties": map[string]interface{}{"mail": map[string]interface{}{"title": "Mail", "type": "string", "x-display": "textarea"}, "schemaKey": map[string]interface{}{"const": "yes", "type": "string"}}, "required": []interface{}{"mail"}, "title": "Yes"}, map[string]interface{}{"properties": map[string]interface{}{"schemaKey": map[string]interface{}{"const": "no", "type": "string"}}, "title": "No"}}, "title": "Mail available", "type": "object"}, "type": "input"}, "search-email-gateway": map[string]interface{}{"active": false, "created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Search email gateway", "next": map[string]interface{}{"extract-iocs": ""}, "order": 4, "type": "task"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"},
		},
	},

	{
		Name: "RunTask",
		Args: Args{Method: "Post", URL: "/tickets/8123/playbooks/phishing/task/board/run"},
		Want: Want{
			Status: 204,
			Body:   nil,
		},
	},

	{
		Name: "SetReferences",
		Args: Args{Method: "Put", URL: "/tickets/8125/references", Data: []interface{}{map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"created": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "id": 8125, "modified": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "name": "phishing from selenafadel@von.com detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}, "schema": "{}", "status": "closed", "tickets": []interface{}{map[string]interface{}{"created": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "id": 8126, "modified": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "name": "Surfaceintroduce virus detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.centralworld-class.io/synthesize", "name": "university"}, map[string]interface{}{"href": "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", "name": "goal"}, map[string]interface{}{"href": "http://www.chiefsyndicate.io/action-items", "name": "unemployment"}}, "schema": "{}", "status": "closed", "type": "alert"}}, "type": "alert"},
		},
	},

	{
		Name: "SetSchema",
		Args: Args{Method: "Put", URL: "/tickets/8125/schema", Data: "{}"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"created": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "id": 8125, "modified": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "name": "phishing from selenafadel@von.com detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "https://www.seniorleading-edge.name/users/efficient", "name": "recovery"}, map[string]interface{}{"href": "http://www.dynamicseamless.com/clicks-and-mortar", "name": "force"}, map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}, "schema": "{}", "status": "closed", "tickets": []interface{}{map[string]interface{}{"created": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "id": 8126, "modified": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "name": "Surfaceintroduce virus detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.centralworld-class.io/synthesize", "name": "university"}, map[string]interface{}{"href": "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", "name": "goal"}, map[string]interface{}{"href": "http://www.chiefsyndicate.io/action-items", "name": "unemployment"}}, "schema": "{}", "status": "closed", "type": "alert"}}, "type": "alert"},
		},
	},

	{
		Name: "LinkTicket",
		Args: Args{Method: "Patch", URL: "/tickets/8126/tickets", Data: 8123},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"created": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "id": 8126, "modified": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "name": "Surfaceintroduce virus detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.centralworld-class.io/synthesize", "name": "university"}, map[string]interface{}{"href": "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", "name": "goal"}, map[string]interface{}{"href": "http://www.chiefsyndicate.io/action-items", "name": "unemployment"}}, "schema": "{}", "status": "closed", "tickets": []interface{}{map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}, map[string]interface{}{"name": "leadreintermediate.io", "status": "malicious"}}, "created": time.Date(2021, time.October, 2, 16, 4, 59, 78206000, time.UTC), "id": 8123, "modified": time.Date(2021, time.October, 2, 16, 4, 59, 78206000, time.UTC), "name": "live zebra", "owner": "demo", "playbooks": map[string]interface{}{"phishing": map[string]interface{}{"name": "Phishing", "tasks": map[string]interface{}{"block-iocs": map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block IOCs", "type": "task"}, "block-sender": map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Block sender", "next": map[string]interface{}{"extract-iocs": ""}, "type": "task"}, "board": map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}, "escalate": map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Escalate to CISO", "type": "task"}, "extract-iocs": map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Extract IOCs", "next": map[string]interface{}{"block-iocs": ""}, "schema": map[string]interface{}{"properties": map[string]interface{}{"iocs": map[string]interface{}{"items": map[string]interface{}{"type": "string"}, "title": "IOCs", "type": "array"}}, "title": "Extract IOCs", "type": "object"}, "type": "input"}, "mail-available": map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Mail available", "next": map[string]interface{}{"block-sender": "schemaKey == 'yes'", "extract-iocs": "schemaKey == 'yes'", "search-email-gateway": "schemaKey == 'no'"}, "schema": map[string]interface{}{"oneOf": []interface{}{map[string]interface{}{"properties": map[string]interface{}{"mail": map[string]interface{}{"title": "Mail", "type": "string", "x-display": "textarea"}, "schemaKey": map[string]interface{}{"const": "yes", "type": "string"}}, "required": []interface{}{"mail"}, "title": "Yes"}, map[string]interface{}{"properties": map[string]interface{}{"schemaKey": map[string]interface{}{"const": "no", "type": "string"}}, "title": "No"}}, "title": "Mail available", "type": "object"}, "type": "input"}, "search-email-gateway": map[string]interface{}{"created": time.Date(2021, time.December, 12, 12, 12, 12, 12, time.UTC), "done": false, "name": "Search email gateway", "next": map[string]interface{}{"extract-iocs": ""}, "type": "task"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"}, map[string]interface{}{"created": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "id": 8125, "modified": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "name": "phishing from selenafadel@von.com detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "https://www.seniorleading-edge.name/users/efficient", "name": "recovery"}, map[string]interface{}{"href": "http://www.dynamicseamless.com/clicks-and-mortar", "name": "force"}, map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}, "schema": "{}", "status": "closed", "type": "alert"}}, "type": "alert"},
		},
	},

	{
		Name: "UnlinkTicket",
		Args: Args{Method: "Delete", URL: "/tickets/8126/tickets", Data: 8125},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"created": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "id": 8126, "modified": time.Date(2021, time.October, 2, 16, 4, 59, 78186000, time.UTC), "name": "Surfaceintroduce virus detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.centralworld-class.io/synthesize", "name": "university"}, map[string]interface{}{"href": "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", "name": "goal"}, map[string]interface{}{"href": "http://www.chiefsyndicate.io/action-items", "name": "unemployment"}}, "schema": "{}", "status": "closed", "type": "alert"},
		},
	},

	{
		Name: "ListTicketTypes",
		Args: Args{Method: "Get", URL: "/tickettypes"},
		Want: Want{
			Status: 200,
			Body:   []interface{}{map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-alert", "id": "alert", "name": "Alerts"}, map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-radioactive", "id": "incident", "name": "Incidents"}, map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-fingerprint", "id": "investigation", "name": "Forensic Investigations"}, map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-target", "id": "hunt", "name": "Threat Hunting"}},
		},
	},

	{
		Name: "CreateTicketType",
		Args: Args{Method: "Post", URL: "/tickettypes", Data: map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-newspaper-variant-outline", "name": "TI Tickets"}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-newspaper-variant-outline", "id": "ti-tickets", "name": "TI Tickets"},
		},
	},

	{
		Name: "GetTicketType",
		Args: Args{Method: "Get", URL: "/tickettypes/alert"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-alert", "id": "alert", "name": "Alerts"},
		},
	},

	{
		Name: "UpdateTicketType",
		Args: Args{Method: "Put", URL: "/tickettypes/alert", Data: map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-bell", "id": "alert", "name": "Alerts"}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-bell", "id": "alert", "name": "Alerts"},
		},
	},

	{
		Name: "DeleteTicketType",
		Args: Args{Method: "Delete", URL: "/tickettypes/alert"},
		Want: Want{
			Status: 204,
			Body:   nil,
		},
	},

	{
		Name: "ListUserData",
		Args: Args{Method: "Get", URL: "/userdata"},
		Want: Want{
			Status: 200,
			Body:   []interface{}{map[string]interface{}{"email": "bob@example.org", "id": "bob", "name": "Bob Bad"}},
		},
	},

	{
		Name: "GetUserData",
		Args: Args{Method: "Get", URL: "/userdata/bob"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"email": "bob@example.org", "id": "bob", "name": "Bob Bad"},
		},
	},

	{
		Name: "UpdateUserData",
		Args: Args{Method: "Put", URL: "/userdata/bob", Data: map[string]interface{}{"blocked": false, "email": "bob@example.org", "name": "Bob Bad"}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"email": "bob@example.org", "id": "bob", "name": "Bob Bad"},
		},
	},

	{
		Name: "ListUsers",
		Args: Args{Method: "Get", URL: "/users"},
		Want: Want{
			Status: 200,
			Body:   []interface{}{map[string]interface{}{"apikey": false, "blocked": false, "id": "bob", "roles": []interface{}{"admin:backup:read", "admin:backup:restore", "admin:dashboard:write", "admin:group:write", "admin:job:read", "admin:job:write", "admin:log:read", "admin:settings:write", "admin:ticket:delete", "admin:user:write", "admin:userdata:read", "admin:userdata:write", "analyst:automation:read", "analyst:currentsettings:write", "analyst:currentuser:read", "analyst:currentuserdata:read", "analyst:dashboard:read", "analyst:file", "analyst:group:read", "analyst:playbook:read", "analyst:rule:read", "analyst:settings:read", "analyst:template:read", "analyst:ticket:read", "analyst:ticket:write", "analyst:tickettype:read", "analyst:user:read", "engineer:automation:write", "engineer:playbook:write", "engineer:rule:write", "engineer:template:write", "engineer:tickettype:write"}}, map[string]interface{}{"apikey": true, "blocked": false, "id": "script", "roles": []interface{}{"analyst:automation:read", "analyst:currentsettings:write", "analyst:currentuser:read", "analyst:currentuserdata:read", "analyst:dashboard:read", "analyst:file", "analyst:group:read", "analyst:playbook:read", "analyst:rule:read", "analyst:settings:read", "analyst:template:read", "analyst:ticket:read", "analyst:ticket:write", "analyst:tickettype:read", "analyst:user:read", "engineer:automation:write", "engineer:playbook:write", "engineer:rule:write", "engineer:template:write", "engineer:tickettype:write"}}},
		},
	},

	{
		Name: "CreateUser",
		Args: Args{Method: "Post", URL: "/users", Data: map[string]interface{}{"apikey": true, "blocked": false, "id": "syncscript", "roles": []interface{}{"analyst"}}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"blocked": false, "id": "syncscript", "roles": []interface{}{"analyst:automation:read", "analyst:currentsettings:write", "analyst:currentuser:read", "analyst:currentuserdata:read", "analyst:dashboard:read", "analyst:file", "analyst:group:read", "analyst:playbook:read", "analyst:rule:read", "analyst:settings:read", "analyst:template:read", "analyst:ticket:read", "analyst:ticket:write", "analyst:tickettype:read", "analyst:user:read"}, "secret": "v39bOuobnlEljfWzjAgoKzhmnh1xSMxH"},
		},
	},

	{
		Name: "GetUser",
		Args: Args{Method: "Get", URL: "/users/script"},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"apikey": true, "blocked": false, "id": "script", "roles": []interface{}{"analyst:automation:read", "analyst:currentsettings:write", "analyst:currentuser:read", "analyst:currentuserdata:read", "analyst:dashboard:read", "analyst:file", "analyst:group:read", "analyst:playbook:read", "analyst:rule:read", "analyst:settings:read", "analyst:template:read", "analyst:ticket:read", "analyst:ticket:write", "analyst:tickettype:read", "analyst:user:read", "engineer:automation:write", "engineer:playbook:write", "engineer:rule:write", "engineer:template:write", "engineer:tickettype:write"}},
		},
	},

	{
		Name: "UpdateUser",
		Args: Args{Method: "Put", URL: "/users/bob", Data: map[string]interface{}{"apikey": false, "blocked": false, "id": "syncscript", "roles": []interface{}{"analyst", "admin"}}},
		Want: Want{
			Status: 200,
			Body:   map[string]interface{}{"apikey": false, "blocked": false, "id": "bob", "roles": []interface{}{"admin:backup:read", "admin:backup:restore", "admin:dashboard:write", "admin:group:write", "admin:job:read", "admin:job:write", "admin:log:read", "admin:settings:write", "admin:ticket:delete", "admin:user:write", "admin:userdata:read", "admin:userdata:write", "analyst:automation:read", "analyst:currentsettings:write", "analyst:currentuser:read", "analyst:currentuserdata:read", "analyst:dashboard:read", "analyst:file", "analyst:group:read", "analyst:playbook:read", "analyst:rule:read", "analyst:settings:read", "analyst:template:read", "analyst:ticket:read", "analyst:ticket:write", "analyst:tickettype:read", "analyst:user:read", "engineer:automation:write", "engineer:playbook:write", "engineer:rule:write", "engineer:template:write", "engineer:tickettype:write"}},
		},
	},

	{
		Name: "DeleteUser",
		Args: Args{Method: "Delete", URL: "/users/script"},
		Want: Want{
			Status: 204,
			Body:   nil,
		},
	},
}
