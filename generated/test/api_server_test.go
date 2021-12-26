package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/test"
	ctime "github.com/SecurityBrewery/catalyst/time"
)

type testClock struct{}

func (testClock) Now() time.Time {
	return time.Date(2021, 12, 12, 12, 12, 12, 12, time.UTC)
}

func TestService(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctime.DefaultClock = testClock{}

	type args struct {
		method string
		url    string
		data   interface{}
	}
	type want struct {
		status int
		body   interface{}
	}
	tests := []struct {
		name string
		args args
		want want
	}{

		{
			name: "AddArtifact",
			args: args{method: "POST", url: "/api/tickets/8123/artifacts", data: map[string]interface{}{"name": "2.2.2.2"}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}, map[string]interface{}{"name": "leadreintermediate.io", "status": "malicious"}, map[string]interface{}{"name": "2.2.2.2", "status": "unknown", "type": "ip"}}, "created": "2021-12-12T12:12:12Z", "id": 8123, "modified": "2021-12-12T12:12:12Z", "name": "live zebra", "owner": "demo", "playbooks": map[string]interface{}{"phishing": map[string]interface{}{"name": "Phishing", "tasks": map[string]interface{}{"block-iocs": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Block IOCs", "order": 6, "type": "task"}, "block-sender": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Block sender", "next": map[string]interface{}{"extract-iocs": ""}, "order": 3, "type": "task"}, "board": map[string]interface{}{"active": true, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "order": 0, "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}, "escalate": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Escalate to CISO", "order": 1, "type": "task"}, "extract-iocs": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Extract IOCs", "next": map[string]interface{}{"block-iocs": ""}, "order": 5, "schema": map[string]interface{}{"properties": map[string]interface{}{"iocs": map[string]interface{}{"items": map[string]interface{}{"type": "string"}, "title": "IOCs", "type": "array"}}, "title": "Extract IOCs", "type": "object"}, "type": "input"}, "mail-available": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Mail available", "next": map[string]interface{}{"block-sender": "schemaKey == 'yes'", "extract-iocs": "schemaKey == 'yes'", "search-email-gateway": "schemaKey == 'no'"}, "order": 2, "schema": map[string]interface{}{"oneOf": []interface{}{map[string]interface{}{"properties": map[string]interface{}{"mail": map[string]interface{}{"title": "Mail", "type": "string", "x-display": "textarea"}, "schemaKey": map[string]interface{}{"const": "yes", "type": "string"}}, "required": []interface{}{"mail"}, "title": "Yes"}, map[string]interface{}{"properties": map[string]interface{}{"schemaKey": map[string]interface{}{"const": "no", "type": "string"}}, "title": "No"}}, "title": "Mail available", "type": "object"}, "type": "input"}, "search-email-gateway": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Search email gateway", "next": map[string]interface{}{"extract-iocs": ""}, "order": 4, "type": "task"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"},
			},
		},
		{
			name: "AddComment",
			args: args{method: "POST", url: "/api/tickets/8125/comments", data: map[string]interface{}{"message": "My first comment"}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"comments": []interface{}{map[string]interface{}{"created": "2021-12-12T12:12:12Z", "creator": "bob", "message": "My first comment"}}, "created": "2021-12-12T12:12:12Z", "id": 8125, "modified": "2021-12-12T12:12:12Z", "name": "phishing from selenafadel@von.com detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "https://www.seniorleading-edge.name/users/efficient", "name": "recovery"}, map[string]interface{}{"href": "http://www.dynamicseamless.com/clicks-and-mortar", "name": "force"}, map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}, "schema": "{}", "status": "closed", "tickets": []interface{}{map[string]interface{}{"created": "2021-12-12T12:12:12Z", "id": 8126, "modified": "2021-12-12T12:12:12Z", "name": "Surfaceintroduce virus detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.centralworld-class.io/synthesize", "name": "university"}, map[string]interface{}{"href": "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", "name": "goal"}, map[string]interface{}{"href": "http://www.chiefsyndicate.io/action-items", "name": "unemployment"}}, "schema": "{}", "status": "closed", "type": "alert"}}, "type": "alert"},
			},
		},
		{
			name: "AddTicketPlaybook",
			args: args{method: "POST", url: "/api/tickets/8125/playbooks", data: map[string]interface{}{"yaml": "name: Simple\ntasks:\n  input:\n    name: Upload malware if possible\n    type: input\n    schema:\n      title: Malware\n      type: object\n      properties:\n        malware:\n          type: string\n          title: Select malware\n          default: \"\"\n    next:\n      hash: \"malware != ''\"\n\n  hash:\n    name: Hash the malware\n    type: automation\n    automation: hash.sha1\n    payload:\n      default: \"playbook.tasks['input'].data['malware']\"\n    next:\n      escalate:\n\n  escalate:\n    name: Escalate to malware team\n    type: task\n"}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"created": "1985-04-12T23:20:50.52Z", "id": 8125, "modified": "1985-04-12T23:20:50.52Z", "name": "phishing from selenafadel@von.com detected", "owner": "demo", "playbooks": map[string]interface{}{"simple": map[string]interface{}{"name": "Simple", "tasks": map[string]interface{}{"escalate": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Escalate to malware team", "order": 2, "type": "task"}, "hash": map[string]interface{}{"active": false, "automation": "hash.sha1", "created": "2021-12-12T12:12:12Z", "done": false, "name": "Hash the malware", "next": map[string]interface{}{"escalate": ""}, "order": 1, "payload": map[string]interface{}{"default": "playbook.tasks['input'].data['malware']"}, "type": "automation"}, "input": map[string]interface{}{"active": true, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Upload malware if possible", "next": map[string]interface{}{"hash": "malware != ''"}, "order": 0, "schema": map[string]interface{}{"properties": map[string]interface{}{"malware": map[string]interface{}{"default": "", "title": "Select malware", "type": "string"}}, "title": "Malware", "type": "object"}, "type": "input"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.seniorleading-edge.name/users/efficient", "name": "recovery"}, map[string]interface{}{"href": "http://www.dynamicseamless.com/clicks-and-mortar", "name": "force"}, map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}, "schema": "{}", "status": "closed", "tickets": []interface{}{map[string]interface{}{"created": "2021-12-12T12:12:12Z", "id": 8126, "modified": "2021-12-12T12:12:12Z", "name": "Surfaceintroduce virus detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.centralworld-class.io/synthesize", "name": "university"}, map[string]interface{}{"href": "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", "name": "goal"}, map[string]interface{}{"href": "http://www.chiefsyndicate.io/action-items", "name": "unemployment"}}, "schema": "{}", "status": "closed", "type": "alert"}}, "type": "alert"},
			},
		},
		{
			name: "CompleteTask",
			args: args{method: "PUT", url: "/api/tickets/8123/playbooks/phishing/task/board/complete", data: map[string]interface{}{"boardInvolved": true}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}, map[string]interface{}{"name": "leadreintermediate.io", "status": "malicious"}}, "created": "2021-12-12T12:12:12Z", "id": 8123, "modified": "2021-12-12T12:12:12Z", "name": "live zebra", "owner": "demo", "playbooks": map[string]interface{}{"phishing": map[string]interface{}{"name": "Phishing", "tasks": map[string]interface{}{"block-iocs": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Block IOCs", "order": 6, "type": "task"}, "block-sender": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Block sender", "next": map[string]interface{}{"extract-iocs": ""}, "order": 3, "type": "task"}, "board": map[string]interface{}{"active": false, "closed": "2021-12-12T12:12:12Z", "created": "2021-12-12T12:12:12Z", "data": map[string]interface{}{"boardInvolved": true}, "done": true, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "order": 0, "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}, "escalate": map[string]interface{}{"active": true, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Escalate to CISO", "order": 1, "type": "task"}, "extract-iocs": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Extract IOCs", "next": map[string]interface{}{"block-iocs": ""}, "order": 5, "schema": map[string]interface{}{"properties": map[string]interface{}{"iocs": map[string]interface{}{"items": map[string]interface{}{"type": "string"}, "title": "IOCs", "type": "array"}}, "title": "Extract IOCs", "type": "object"}, "type": "input"}, "mail-available": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Mail available", "next": map[string]interface{}{"block-sender": "schemaKey == 'yes'", "extract-iocs": "schemaKey == 'yes'", "search-email-gateway": "schemaKey == 'no'"}, "order": 2, "schema": map[string]interface{}{"oneOf": []interface{}{map[string]interface{}{"properties": map[string]interface{}{"mail": map[string]interface{}{"title": "Mail", "type": "string", "x-display": "textarea"}, "schemaKey": map[string]interface{}{"const": "yes", "type": "string"}}, "required": []interface{}{"mail"}, "title": "Yes"}, map[string]interface{}{"properties": map[string]interface{}{"schemaKey": map[string]interface{}{"const": "no", "type": "string"}}, "title": "No"}}, "title": "Mail available", "type": "object"}, "type": "input"}, "search-email-gateway": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Search email gateway", "next": map[string]interface{}{"extract-iocs": ""}, "order": 4, "type": "task"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"},
			},
		},
		{
			name: "CreateAutomation",
			args: args{method: "POST", url: "/api/automations", data: map[string]interface{}{"id": "hash-sha-256", "image": "docker.io/python:3", "script": "import sys\nimport json\nimport hashlib\n\n\ndef run(msg):\n    sha256 = hashlib.sha256(msg['payload']['default'].encode('utf-8'))\n    return {'hash': sha256.hexdigest()}\n\n\nprint(json.dumps(run(json.loads(sys.argv[1]))))\n", "type": []interface{}{"global"}}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"id": "hash-sha-256", "image": "docker.io/python:3", "script": "import sys\nimport json\nimport hashlib\n\n\ndef run(msg):\n    sha256 = hashlib.sha256(msg['payload']['default'].encode('utf-8'))\n    return {'hash': sha256.hexdigest()}\n\n\nprint(json.dumps(run(json.loads(sys.argv[1]))))\n", "type": []interface{}{"global"}},
			},
		},
		{
			name: "CreatePlaybook",
			args: args{method: "POST", url: "/api/playbooks", data: map[string]interface{}{"yaml": "name: Simple2\ntasks:\n  input:\n    name: Upload malware if possible\n    type: input\n    schema:\n      title: Malware\n      type: object\n      properties:\n        malware:\n          type: string\n          title: Select malware\n          default: \"\"\n    next:\n      hash: \"malware != ''\"\n\n  hash:\n    name: Hash the malware\n    type: automation\n    automation: hash.sha1\n    payload:\n      default: \"playbook.tasks['input'].data['malware']\"\n    next:\n      escalate:\n\n  escalate:\n    name: Escalate to malware team\n    type: task\n"}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"id": "simple-2", "name": "Simple2", "yaml": "name: Simple2\ntasks:\n  input:\n    name: Upload malware if possible\n    type: input\n    schema:\n      title: Malware\n      type: object\n      properties:\n        malware:\n          type: string\n          title: Select malware\n          default: \"\"\n    next:\n      hash: \"malware != ''\"\n\n  hash:\n    name: Hash the malware\n    type: automation\n    automation: hash.sha1\n    payload:\n      default: \"playbook.tasks['input'].data['malware']\"\n    next:\n      escalate:\n\n  escalate:\n    name: Escalate to malware team\n    type: task\n"},
			},
		},
		{
			name: "CreateTemplate",
			args: args{method: "POST", url: "/api/templates", data: map[string]interface{}{"name": "My Template", "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n"}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"id": "my-template", "name": "My Template", "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n"},
			},
		},
		{
			name: "CreateTicket",
			args: args{method: "POST", url: "/api/tickets", data: map[string]interface{}{"id": 123, "name": "Wannacry infection", "owner": "bob", "status": "open", "type": "incident"}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"created": "1985-04-12T23:20:50.52Z", "id": 123, "modified": "1985-04-12T23:20:50.52Z", "name": "Wannacry infection", "owner": "bob", "schema": "{}", "status": "open", "type": "incident"},
			},
		},
		{
			name: "CreateTicketBatch",
			args: args{method: "POST", url: "/api/tickets/batch", data: []interface{}{map[string]interface{}{"id": 123, "name": "Wannacry infection", "owner": "bob", "status": "open", "type": "incident"}}},
			want: want{
				status: 204,
				body:   nil,
			},
		},
		{
			name: "CreateTicketType",
			args: args{method: "POST", url: "/api/tickettypes", data: map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-newspaper-variant-outline", "name": "TI Tickets"}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-newspaper-variant-outline", "id": "ti-tickets", "name": "TI Tickets"},
			},
		},
		{
			name: "CreateUser",
			args: args{method: "POST", url: "/api/users", data: map[string]interface{}{"id": "syncscript", "roles": []interface{}{"analyst"}}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"blocked": false, "id": "syncscript", "roles": []interface{}{"analyst:automation:read", "analyst:currentsettings:write", "analyst:currentuser:read", "analyst:currentuserdata:read", "analyst:file", "analyst:group:read", "analyst:playbook:read", "analyst:rule:read", "analyst:settings:read", "analyst:template:read", "analyst:ticket:read", "analyst:ticket:write", "analyst:tickettype:read", "analyst:user:read"}, "secret": "v39bOuobnlEljfWzjAgoKzhmnh1xSMxH"},
			},
		},
		{
			name: "CurrentUser",
			args: args{method: "GET", url: "/api/currentuser"},
			want: want{
				status: 200,
				body:   map[string]interface{}{"apikey": false, "blocked": false, "id": "bob", "roles": []interface{}{"admin:backup:read", "admin:backup:restore", "admin:group:write", "admin:job:read", "admin:job:write", "admin:log:read", "admin:ticket:delete", "admin:user:write", "admin:userdata:read", "admin:userdata:write", "analyst:automation:read", "analyst:currentsettings:write", "analyst:currentuser:read", "analyst:currentuserdata:read", "analyst:file", "analyst:group:read", "analyst:playbook:read", "analyst:rule:read", "analyst:settings:read", "analyst:template:read", "analyst:ticket:read", "analyst:ticket:write", "analyst:tickettype:read", "analyst:user:read", "engineer:automation:write", "engineer:playbook:write", "engineer:rule:write", "engineer:template:write", "engineer:tickettype:write"}},
			},
		},
		{
			name: "CurrentUserData",
			args: args{method: "GET", url: "/api/currentuserdata"},
			want: want{
				status: 200,
				body:   map[string]interface{}{"email": "bob@example.org", "id": "bob", "name": "Bob Bad"},
			},
		},
		{
			name: "DeleteAutomation",
			args: args{method: "DELETE", url: "/api/automations/hash.sha1"},
			want: want{
				status: 204,
				body:   nil,
			},
		},
		{
			name: "DeletePlaybook",
			args: args{method: "DELETE", url: "/api/playbooks/simple"},
			want: want{
				status: 204,
				body:   nil,
			},
		},
		{
			name: "DeleteTemplate",
			args: args{method: "DELETE", url: "/api/templates/default"},
			want: want{
				status: 204,
				body:   nil,
			},
		},
		{
			name: "DeleteTicket",
			args: args{method: "DELETE", url: "/api/tickets/8125"},
			want: want{
				status: 204,
				body:   nil,
			},
		},
		{
			name: "DeleteTicketType",
			args: args{method: "DELETE", url: "/api/tickettypes/alert"},
			want: want{
				status: 204,
				body:   nil,
			},
		},
		{
			name: "DeleteUser",
			args: args{method: "DELETE", url: "/api/users/script"},
			want: want{
				status: 204,
				body:   nil,
			},
		},
		{
			name: "EnrichArtifact",
			args: args{method: "POST", url: "/api/tickets/8123/artifacts/leadreintermediate.io/enrich", data: map[string]interface{}{"data": map[string]interface{}{"hash": "b7a067a742c20d07a7456646de89bc2d408a1153"}, "name": "hash.sha1"}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}, map[string]interface{}{"enrichments": map[string]interface{}{"hash.sha1": map[string]interface{}{"created": "2021-12-12T12:12:12Z", "data": map[string]interface{}{"hash": "b7a067a742c20d07a7456646de89bc2d408a1153"}, "name": "hash.sha1"}}, "name": "leadreintermediate.io", "status": "malicious"}}, "created": "2021-12-12T12:12:12Z", "id": 8123, "modified": "2021-12-12T12:12:12Z", "name": "live zebra", "owner": "demo", "playbooks": map[string]interface{}{"phishing": map[string]interface{}{"name": "Phishing", "tasks": map[string]interface{}{"block-iocs": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Block IOCs", "order": 6, "type": "task"}, "block-sender": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Block sender", "next": map[string]interface{}{"extract-iocs": ""}, "order": 3, "type": "task"}, "board": map[string]interface{}{"active": true, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "order": 0, "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}, "escalate": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Escalate to CISO", "order": 1, "type": "task"}, "extract-iocs": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Extract IOCs", "next": map[string]interface{}{"block-iocs": ""}, "order": 5, "schema": map[string]interface{}{"properties": map[string]interface{}{"iocs": map[string]interface{}{"items": map[string]interface{}{"type": "string"}, "title": "IOCs", "type": "array"}}, "title": "Extract IOCs", "type": "object"}, "type": "input"}, "mail-available": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Mail available", "next": map[string]interface{}{"block-sender": "schemaKey == 'yes'", "extract-iocs": "schemaKey == 'yes'", "search-email-gateway": "schemaKey == 'no'"}, "order": 2, "schema": map[string]interface{}{"oneOf": []interface{}{map[string]interface{}{"properties": map[string]interface{}{"mail": map[string]interface{}{"title": "Mail", "type": "string", "x-display": "textarea"}, "schemaKey": map[string]interface{}{"const": "yes", "type": "string"}}, "required": []interface{}{"mail"}, "title": "Yes"}, map[string]interface{}{"properties": map[string]interface{}{"schemaKey": map[string]interface{}{"const": "no", "type": "string"}}, "title": "No"}}, "title": "Mail available", "type": "object"}, "type": "input"}, "search-email-gateway": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Search email gateway", "next": map[string]interface{}{"extract-iocs": ""}, "order": 4, "type": "task"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"},
			},
		},
		{
			name: "GetArtifact",
			args: args{method: "GET", url: "/api/tickets/8123/artifacts/leadreintermediate.io"},
			want: want{
				status: 200,
				body:   map[string]interface{}{"name": "leadreintermediate.io", "status": "malicious"},
			},
		},
		{
			name: "GetAutomation",
			args: args{method: "GET", url: "/api/automations/hash.sha1"},
			want: want{
				status: 200,
				body:   map[string]interface{}{"id": "hash.sha1", "image": "docker.io/python:3", "schema": "{\"title\":\"Input\",\"type\":\"object\",\"properties\":{\"default\":{\"type\":\"string\",\"title\":\"Value\"}},\"required\":[\"default\"]}", "script": "#!/usr/bin/env python\n\nimport sys\nimport json\nimport hashlib\n\n\ndef run(msg):\n    sha1 = hashlib.sha1(msg['payload']['default'].encode('utf-8'))\n    return {\"hash\": sha1.hexdigest()}\n\n\nprint(json.dumps(run(json.loads(sys.argv[1]))))\n", "type": []interface{}{"global", "artifact", "playbook"}},
			},
		},
		{
			name: "GetJob",
			args: args{method: "GET", url: "/api/jobs/99cd67131b48"},
			want: want{
				status: 200,
				body:   map[string]interface{}{"automation": "hash.sha1", "id": "99cd67131b48", "payload": "test", "status": "created"},
			},
		},
		{
			name: "GetLogs",
			args: args{method: "GET", url: "/api/logs/tickets%252F294511"},
			want: want{
				status: 200,
				body:   []interface{}{map[string]interface{}{"created": "2021-12-12T12:12:12Z", "creator": "bob", "message": "Fail run account resist lend solve incident centre priority temperature. Cause change distribution examine location technique shape partner milk customer. Rail tea plate soil report cook railway interpretation breath action. Exercise dream accept park conclusion addition shoot assistance may answer. Gold writer link stop combine hear power name commitment operation. Determine lifespan support grow degree henry exclude detail set religion. Direct library policy convention chain retain discover ride walk student. Gather proposal select march aspect play noise avoid encourage employ. Assessment preserve transport combine wish influence income guess run stand. Charge limit crime ignore statement foundation study issue stop claim.", "reference": "tickets/294511"}},
			},
		},
		{
			name: "GetPlaybook",
			args: args{method: "GET", url: "/api/playbooks/simple"},
			want: want{
				status: 200,
				body:   map[string]interface{}{"id": "simple", "name": "Simple", "yaml": "name: Simple\ntasks:\n  input:\n    name: Enter something to hash\n    type: input\n    schema:\n      title: Something\n      type: object\n      properties:\n        something:\n          type: string\n          title: Something\n          default: \"\"\n    next:\n      hash: \"something != ''\"\n\n  hash:\n    name: Hash the something\n    type: automation\n    automation: hash.sha1\n    payload:\n      default: \"playbook.tasks['input'].data['something']\"\n    next:\n      comment: \"hash != ''\"\n\n  comment:\n    name: Comment the hash\n    type: automation\n    automation: comment\n    payload:\n      default: \"playbook.tasks['hash'].data['hash']\"\n    next:\n      done: \"done\"\n\n  done:\n    name: You can close this case now\n    type: task\n"},
			},
		},
		{
			name: "GetSettings",
			args: args{method: "GET", url: "/api/settings"},
			want: want{
				status: 200,
				body:   map[string]interface{}{"artifactStates": []interface{}{map[string]interface{}{"color": "info", "icon": "mdi-help-circle-outline", "id": "unknown", "name": "Unknown"}, map[string]interface{}{"color": "error", "icon": "mdi-skull", "id": "malicious", "name": "Malicious"}, map[string]interface{}{"color": "success", "icon": "mdi-check", "id": "clean", "name": "Clean"}}, "roles": []interface{}{"admin:backup:read", "admin:backup:restore", "admin:group:write", "admin:job:read", "admin:job:write", "admin:log:read", "admin:ticket:delete", "admin:user:write", "admin:userdata:read", "admin:userdata:write", "analyst:automation:read", "analyst:currentsettings:write", "analyst:currentuser:read", "analyst:currentuserdata:read", "analyst:file", "analyst:group:read", "analyst:playbook:read", "analyst:rule:read", "analyst:settings:read", "analyst:template:read", "analyst:ticket:read", "analyst:ticket:write", "analyst:tickettype:read", "analyst:user:read", "engineer:automation:write", "engineer:playbook:write", "engineer:rule:write", "engineer:template:write", "engineer:tickettype:write"}, "ticketTypes": []interface{}{map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-alert", "id": "alert", "name": "Alerts"}, map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-radioactive", "id": "incident", "name": "Incidents"}, map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-fingerprint", "id": "investigation", "name": "Forensic Investigations"}, map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-target", "id": "hunt", "name": "Threat Hunting"}}, "tier": "community", "timeformat": "YYYY-MM-DDThh:mm:ss", "version": "0.0.0-test"},
			},
		},
		{
			name: "GetStatistics",
			args: args{method: "GET", url: "/api/statistics"},
			want: want{
				status: 200,
				body:   map[string]interface{}{"open_tickets_per_user": map[string]interface{}{}, "tickets_per_type": map[string]interface{}{"alert": 2, "incident": 1}, "tickets_per_week": map[string]interface{}{"2021-39": 3}, "unassigned": 0},
			},
		},
		{
			name: "GetTemplate",
			args: args{method: "GET", url: "/api/templates/default"},
			want: want{
				status: 200,
				body:   map[string]interface{}{"id": "default", "name": "Default", "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Default\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n"},
			},
		},
		{
			name: "GetTicket",
			args: args{method: "GET", url: "/api/tickets/8125"},
			want: want{
				status: 200,
				body:   map[string]interface{}{"created": "2021-12-12T12:12:12Z", "id": 8125, "modified": "2021-12-12T12:12:12Z", "name": "phishing from selenafadel@von.com detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "https://www.seniorleading-edge.name/users/efficient", "name": "recovery"}, map[string]interface{}{"href": "http://www.dynamicseamless.com/clicks-and-mortar", "name": "force"}, map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}, "schema": "{}", "status": "closed", "tickets": []interface{}{map[string]interface{}{"created": "2021-12-12T12:12:12Z", "id": 8126, "modified": "2021-12-12T12:12:12Z", "name": "Surfaceintroduce virus detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.centralworld-class.io/synthesize", "name": "university"}, map[string]interface{}{"href": "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", "name": "goal"}, map[string]interface{}{"href": "http://www.chiefsyndicate.io/action-items", "name": "unemployment"}}, "schema": "{}", "status": "closed", "type": "alert"}}, "type": "alert"},
			},
		},
		{
			name: "GetTicketType",
			args: args{method: "GET", url: "/api/tickettypes/alert"},
			want: want{
				status: 200,
				body:   map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-alert", "id": "alert", "name": "Alerts"},
			},
		},
		{
			name: "GetUser",
			args: args{method: "GET", url: "/api/users/script"},
			want: want{
				status: 200,
				body:   map[string]interface{}{"apikey": true, "blocked": false, "id": "script", "roles": []interface{}{"analyst:automation:read", "analyst:currentsettings:write", "analyst:currentuser:read", "analyst:currentuserdata:read", "analyst:file", "analyst:group:read", "analyst:playbook:read", "analyst:rule:read", "analyst:settings:read", "analyst:template:read", "analyst:ticket:read", "analyst:ticket:write", "analyst:tickettype:read", "analyst:user:read", "engineer:automation:write", "engineer:playbook:write", "engineer:rule:write", "engineer:template:write", "engineer:tickettype:write"}},
			},
		},
		{
			name: "GetUserData",
			args: args{method: "GET", url: "/api/userdata/bob"},
			want: want{
				status: 200,
				body:   map[string]interface{}{"email": "bob@example.org", "id": "bob", "name": "Bob Bad"},
			},
		},
		{
			name: "LinkFiles",
			args: args{method: "PUT", url: "/api/tickets/8125/files", data: []interface{}{map[string]interface{}{"key": "myfile", "name": "document.doc"}}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"created": "2021-12-12T12:12:12Z", "files": []interface{}{map[string]interface{}{"key": "myfile", "name": "document.doc"}}, "id": 8125, "modified": "2021-12-12T12:12:12Z", "name": "phishing from selenafadel@von.com detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "https://www.seniorleading-edge.name/users/efficient", "name": "recovery"}, map[string]interface{}{"href": "http://www.dynamicseamless.com/clicks-and-mortar", "name": "force"}, map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}, "schema": "{}", "status": "closed", "tickets": []interface{}{map[string]interface{}{"created": "2021-12-12T12:12:12Z", "id": 8126, "modified": "2021-12-12T12:12:12Z", "name": "Surfaceintroduce virus detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.centralworld-class.io/synthesize", "name": "university"}, map[string]interface{}{"href": "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", "name": "goal"}, map[string]interface{}{"href": "http://www.chiefsyndicate.io/action-items", "name": "unemployment"}}, "schema": "{}", "status": "closed", "type": "alert"}}, "type": "alert"},
			},
		},
		{
			name: "LinkTicket",
			args: args{method: "PATCH", url: "/api/tickets/8126/tickets", data: 8123},
			want: want{
				status: 200,
				body:   map[string]interface{}{"created": "2021-12-12T12:12:12Z", "id": 8126, "modified": "2021-12-12T12:12:12Z", "name": "Surfaceintroduce virus detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.centralworld-class.io/synthesize", "name": "university"}, map[string]interface{}{"href": "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", "name": "goal"}, map[string]interface{}{"href": "http://www.chiefsyndicate.io/action-items", "name": "unemployment"}}, "schema": "{}", "status": "closed", "tickets": []interface{}{map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}, map[string]interface{}{"name": "leadreintermediate.io", "status": "malicious"}}, "created": "2021-12-12T12:12:12Z", "id": 8123, "modified": "2021-12-12T12:12:12Z", "name": "live zebra", "owner": "demo", "playbooks": map[string]interface{}{"phishing": map[string]interface{}{"name": "Phishing", "tasks": map[string]interface{}{"block-iocs": map[string]interface{}{"created": "2021-12-12T12:12:12Z", "done": false, "name": "Block IOCs", "type": "task"}, "block-sender": map[string]interface{}{"created": "2021-12-12T12:12:12Z", "done": false, "name": "Block sender", "next": map[string]interface{}{"extract-iocs": ""}, "type": "task"}, "board": map[string]interface{}{"created": "2021-12-12T12:12:12Z", "done": false, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}, "escalate": map[string]interface{}{"created": "2021-12-12T12:12:12Z", "done": false, "name": "Escalate to CISO", "type": "task"}, "extract-iocs": map[string]interface{}{"created": "2021-12-12T12:12:12Z", "done": false, "name": "Extract IOCs", "next": map[string]interface{}{"block-iocs": ""}, "schema": map[string]interface{}{"properties": map[string]interface{}{"iocs": map[string]interface{}{"items": map[string]interface{}{"type": "string"}, "title": "IOCs", "type": "array"}}, "title": "Extract IOCs", "type": "object"}, "type": "input"}, "mail-available": map[string]interface{}{"created": "2021-12-12T12:12:12Z", "done": false, "name": "Mail available", "next": map[string]interface{}{"block-sender": "schemaKey == 'yes'", "extract-iocs": "schemaKey == 'yes'", "search-email-gateway": "schemaKey == 'no'"}, "schema": map[string]interface{}{"oneOf": []interface{}{map[string]interface{}{"properties": map[string]interface{}{"mail": map[string]interface{}{"title": "Mail", "type": "string", "x-display": "textarea"}, "schemaKey": map[string]interface{}{"const": "yes", "type": "string"}}, "required": []interface{}{"mail"}, "title": "Yes"}, map[string]interface{}{"properties": map[string]interface{}{"schemaKey": map[string]interface{}{"const": "no", "type": "string"}}, "title": "No"}}, "title": "Mail available", "type": "object"}, "type": "input"}, "search-email-gateway": map[string]interface{}{"created": "2021-12-12T12:12:12Z", "done": false, "name": "Search email gateway", "next": map[string]interface{}{"extract-iocs": ""}, "type": "task"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"}, map[string]interface{}{"created": "2021-12-12T12:12:12Z", "id": 8125, "modified": "2021-12-12T12:12:12Z", "name": "phishing from selenafadel@von.com detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "https://www.seniorleading-edge.name/users/efficient", "name": "recovery"}, map[string]interface{}{"href": "http://www.dynamicseamless.com/clicks-and-mortar", "name": "force"}, map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}, "schema": "{}", "status": "closed", "type": "alert"}}, "type": "alert"},
			},
		},
		{
			name: "ListAutomations",
			args: args{method: "GET", url: "/api/automations"},
			want: want{
				status: 200,
				body:   []interface{}{map[string]interface{}{"id": "comment", "image": "docker.io/python:3", "script": "", "type": []interface{}{"playbook"}}, map[string]interface{}{"id": "hash.sha1", "image": "docker.io/python:3", "schema": "{\"title\":\"Input\",\"type\":\"object\",\"properties\":{\"default\":{\"type\":\"string\",\"title\":\"Value\"}},\"required\":[\"default\"]}", "script": "", "type": []interface{}{"global", "artifact", "playbook"}}, map[string]interface{}{"id": "thehive", "image": "docker.io/python:3", "schema": "{\"title\":\"TheHive credentials\",\"type\":\"object\",\"properties\":{\"thehiveurl\":{\"type\":\"string\",\"title\":\"TheHive URL (e.g. 'https://thehive.example.org')\"},\"thehivekey\":{\"type\":\"string\",\"title\":\"TheHive API Key\"},\"skip_files\":{\"type\":\"boolean\", \"default\": true, \"title\":\"Skip Files (much faster)\"},\"keep_ids\":{\"type\":\"boolean\", \"default\": true, \"title\":\"Keep IDs and overwrite existing IDs\"}},\"required\":[\"thehiveurl\", \"thehivekey\", \"skip_files\", \"keep_ids\"]}", "script": "", "type": []interface{}{"global"}}, map[string]interface{}{"id": "vt.hash", "image": "docker.io/python:3", "schema": "{\"title\":\"Input\",\"type\":\"object\",\"properties\":{\"default\":{\"type\":\"string\",\"title\":\"Value\"}},\"required\":[\"default\"]}", "script": "", "type": []interface{}{"global", "artifact", "playbook"}}},
			},
		},
		{
			name: "ListJobs",
			args: args{method: "GET", url: "/api/jobs"},
			want: want{
				status: 200,
				body:   []interface{}{map[string]interface{}{"automation": "hash.sha1", "id": "99cd67131b48", "payload": "test", "status": "created"}},
			},
		},
		{
			name: "ListPlaybooks",
			args: args{method: "GET", url: "/api/playbooks"},
			want: want{
				status: 200,
				body:   []interface{}{map[string]interface{}{"id": "malware", "name": "Malware", "yaml": "name: Malware\ntasks:\n  file-or-hash:\n    name: Do you have the file or the hash?\n    type: input\n    schema:\n      title: Malware\n      type: object\n      properties:\n        file:\n          type: string\n          title: \"I have the\"\n          enum: [ \"File\", \"Hash\" ]\n    next:\n      enter-hash: \"file == 'Hash'\"\n      upload: \"file == 'File'\"\n\n  enter-hash:\n    name: Please enter the hash\n    type: input\n    schema:\n      title: Malware\n      type: object\n      properties:\n        hash:\n          type: string\n          title: Please enter the hash value\n          minlength: 32\n    next:\n      virustotal: \"hash != ''\"\n\n  upload:\n    name: Upload the malware\n    type: input\n    schema:\n      title: Malware\n      type: object\n      properties:\n        malware:\n          type: object\n          x-display: file\n          title: Please upload the malware\n    next:\n      hash: \"malware\"\n\n  hash:\n    name: Hash the malware\n    type: automation\n    automation: hash.sha1\n    payload:\n      default: \"playbook.tasks['upload'].data['malware']\"\n    next:\n      virustotal:\n\n  virustotal:\n    name: Send hash to VirusTotal\n    type: automation\n    automation: vt.hash\n    args:\n      hash: \"playbook.tasks['enter-hash'].data['hash'] || playbook.tasks['hash'].data['hash']\"\n    # next:\n    #   known-malware: \"score > 5\"\n    #   sandbox: \"score < 6\" # unknown-malware\n"}, map[string]interface{}{"id": "phishing", "name": "Phishing", "yaml": "name: Phishing\ntasks:\n  board:\n    name: Board Involvement?\n    description: Is a board member involved?\n    type: input\n    schema:\n      properties:\n        boardInvolved:\n          default: false\n          title: A board member is involved.\n          type: boolean\n      required:\n        - boardInvolved\n      title: Board Involvement?\n      type: object\n    next:\n      escalate: \"boardInvolved == true\"\n      mail-available: \"boardInvolved == false\"\n\n  escalate:\n    name: Escalate to CISO\n    description: Please escalate the task to the CISO\n    type: task\n\n  mail-available:\n    name: Mail available\n    type: input\n    schema:\n      oneOf:\n        - properties:\n            mail:\n              title: Mail\n              type: string\n              x-display: textarea\n            schemaKey:\n              const: 'yes'\n              type: string\n          required:\n            - mail\n          title: 'Yes'\n        - properties:\n            schemaKey:\n              const: 'no'\n              type: string\n          title: 'No'\n      title: Mail available\n      type: object\n    next:\n      block-sender: \"schemaKey == 'yes'\"\n      extract-iocs: \"schemaKey == 'yes'\"\n      search-email-gateway: \"schemaKey == 'no'\"\n\n  search-email-gateway:\n    name: Search email gateway\n    description: Please search email-gateway for the phishing mail.\n    type: task\n    next:\n      extract-iocs:\n\n  block-sender:\n    name: Block sender\n    type: task\n    next:\n      extract-iocs:\n\n  extract-iocs:\n    name: Extract IOCs\n    description: Please insert the IOCs\n    type: input\n    schema:\n      properties:\n        iocs:\n          items:\n            type: string\n          title: IOCs\n          type: array\n      title: Extract IOCs\n      type: object\n    next:\n      block-iocs:\n\n  block-iocs:\n    name: Block IOCs\n    type: task\n"}, map[string]interface{}{"id": "simple", "name": "Simple", "yaml": "name: Simple\ntasks:\n  input:\n    name: Enter something to hash\n    type: input\n    schema:\n      title: Something\n      type: object\n      properties:\n        something:\n          type: string\n          title: Something\n          default: \"\"\n    next:\n      hash: \"something != ''\"\n\n  hash:\n    name: Hash the something\n    type: automation\n    automation: hash.sha1\n    payload:\n      default: \"playbook.tasks['input'].data['something']\"\n    next:\n      comment: \"hash != ''\"\n\n  comment:\n    name: Comment the hash\n    type: automation\n    automation: comment\n    payload:\n      default: \"playbook.tasks['hash'].data['hash']\"\n    next:\n      done: \"done\"\n\n  done:\n    name: You can close this case now\n    type: task\n"}},
			},
		},
		{
			name: "ListTasks",
			args: args{method: "GET", url: "/api/tasks"},
			want: want{
				status: 200,
				body:   []interface{}{},
			},
		},
		{
			name: "ListTemplates",
			args: args{method: "GET", url: "/api/templates"},
			want: want{
				status: 200,
				body:   []interface{}{map[string]interface{}{"id": "default", "name": "Default", "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Default\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n"}},
			},
		},
		{
			name: "ListTicketTypes",
			args: args{method: "GET", url: "/api/tickettypes"},
			want: want{
				status: 200,
				body:   []interface{}{map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-alert", "id": "alert", "name": "Alerts"}, map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-radioactive", "id": "incident", "name": "Incidents"}, map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-fingerprint", "id": "investigation", "name": "Forensic Investigations"}, map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-target", "id": "hunt", "name": "Threat Hunting"}},
			},
		},
		{
			name: "ListTickets",
			args: args{method: "GET", url: "/api/tickets"},
			want: want{
				status: 200,
				body:   map[string]interface{}{"count": 3, "tickets": []interface{}{map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}, map[string]interface{}{"name": "leadreintermediate.io", "status": "malicious"}}, "created": "2021-12-12T12:12:12Z", "id": 8123, "modified": "2021-12-12T12:12:12Z", "name": "live zebra", "owner": "demo", "playbooks": map[string]interface{}{"phishing": map[string]interface{}{"name": "Phishing", "tasks": map[string]interface{}{"block-iocs": map[string]interface{}{"created": "2021-12-12T12:12:12Z", "done": false, "name": "Block IOCs", "type": "task"}, "block-sender": map[string]interface{}{"created": "2021-12-12T12:12:12Z", "done": false, "name": "Block sender", "next": map[string]interface{}{"extract-iocs": ""}, "type": "task"}, "board": map[string]interface{}{"created": "2021-12-12T12:12:12Z", "done": false, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}, "escalate": map[string]interface{}{"created": "2021-12-12T12:12:12Z", "done": false, "name": "Escalate to CISO", "type": "task"}, "extract-iocs": map[string]interface{}{"created": "2021-12-12T12:12:12Z", "done": false, "name": "Extract IOCs", "next": map[string]interface{}{"block-iocs": ""}, "schema": map[string]interface{}{"properties": map[string]interface{}{"iocs": map[string]interface{}{"items": map[string]interface{}{"type": "string"}, "title": "IOCs", "type": "array"}}, "title": "Extract IOCs", "type": "object"}, "type": "input"}, "mail-available": map[string]interface{}{"created": "2021-12-12T12:12:12Z", "done": false, "name": "Mail available", "next": map[string]interface{}{"block-sender": "schemaKey == 'yes'", "extract-iocs": "schemaKey == 'yes'", "search-email-gateway": "schemaKey == 'no'"}, "schema": map[string]interface{}{"oneOf": []interface{}{map[string]interface{}{"properties": map[string]interface{}{"mail": map[string]interface{}{"title": "Mail", "type": "string", "x-display": "textarea"}, "schemaKey": map[string]interface{}{"const": "yes", "type": "string"}}, "required": []interface{}{"mail"}, "title": "Yes"}, map[string]interface{}{"properties": map[string]interface{}{"schemaKey": map[string]interface{}{"const": "no", "type": "string"}}, "title": "No"}}, "title": "Mail available", "type": "object"}, "type": "input"}, "search-email-gateway": map[string]interface{}{"created": "2021-12-12T12:12:12Z", "done": false, "name": "Search email gateway", "next": map[string]interface{}{"extract-iocs": ""}, "type": "task"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"}, map[string]interface{}{"created": "2021-12-12T12:12:12Z", "id": 8125, "modified": "2021-12-12T12:12:12Z", "name": "phishing from selenafadel@von.com detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "https://www.seniorleading-edge.name/users/efficient", "name": "recovery"}, map[string]interface{}{"href": "http://www.dynamicseamless.com/clicks-and-mortar", "name": "force"}, map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}, "schema": "{}", "status": "closed", "type": "alert"}, map[string]interface{}{"created": "2021-12-12T12:12:12Z", "id": 8126, "modified": "2021-12-12T12:12:12Z", "name": "Surfaceintroduce virus detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.centralworld-class.io/synthesize", "name": "university"}, map[string]interface{}{"href": "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", "name": "goal"}, map[string]interface{}{"href": "http://www.chiefsyndicate.io/action-items", "name": "unemployment"}}, "schema": "{}", "status": "closed", "type": "alert"}}},
			},
		},
		{
			name: "ListUserData",
			args: args{method: "GET", url: "/api/userdata"},
			want: want{
				status: 200,
				body:   []interface{}{map[string]interface{}{"email": "bob@example.org", "id": "bob", "name": "Bob Bad"}},
			},
		},
		{
			name: "ListUsers",
			args: args{method: "GET", url: "/api/users"},
			want: want{
				status: 200,
				body:   []interface{}{map[string]interface{}{"apikey": false, "blocked": false, "id": "bob", "roles": []interface{}{"admin:backup:read", "admin:backup:restore", "admin:group:write", "admin:job:read", "admin:job:write", "admin:log:read", "admin:ticket:delete", "admin:user:write", "admin:userdata:read", "admin:userdata:write", "analyst:automation:read", "analyst:currentsettings:write", "analyst:currentuser:read", "analyst:currentuserdata:read", "analyst:file", "analyst:group:read", "analyst:playbook:read", "analyst:rule:read", "analyst:settings:read", "analyst:template:read", "analyst:ticket:read", "analyst:ticket:write", "analyst:tickettype:read", "analyst:user:read", "engineer:automation:write", "engineer:playbook:write", "engineer:rule:write", "engineer:template:write", "engineer:tickettype:write"}}, map[string]interface{}{"apikey": true, "blocked": false, "id": "script", "roles": []interface{}{"analyst:automation:read", "analyst:currentsettings:write", "analyst:currentuser:read", "analyst:currentuserdata:read", "analyst:file", "analyst:group:read", "analyst:playbook:read", "analyst:rule:read", "analyst:settings:read", "analyst:template:read", "analyst:ticket:read", "analyst:ticket:write", "analyst:tickettype:read", "analyst:user:read", "engineer:automation:write", "engineer:playbook:write", "engineer:rule:write", "engineer:template:write", "engineer:tickettype:write"}}},
			},
		},
		{
			name: "RemoveArtifact",
			args: args{method: "DELETE", url: "/api/tickets/8123/artifacts/leadreintermediate.io"},
			want: want{
				status: 200,
				body:   map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}}, "created": "2021-12-12T12:12:12Z", "id": 8123, "modified": "2021-12-12T12:12:12Z", "name": "live zebra", "owner": "demo", "playbooks": map[string]interface{}{"phishing": map[string]interface{}{"name": "Phishing", "tasks": map[string]interface{}{"block-iocs": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Block IOCs", "order": 6, "type": "task"}, "block-sender": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Block sender", "next": map[string]interface{}{"extract-iocs": ""}, "order": 3, "type": "task"}, "board": map[string]interface{}{"active": true, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "order": 0, "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}, "escalate": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Escalate to CISO", "order": 1, "type": "task"}, "extract-iocs": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Extract IOCs", "next": map[string]interface{}{"block-iocs": ""}, "order": 5, "schema": map[string]interface{}{"properties": map[string]interface{}{"iocs": map[string]interface{}{"items": map[string]interface{}{"type": "string"}, "title": "IOCs", "type": "array"}}, "title": "Extract IOCs", "type": "object"}, "type": "input"}, "mail-available": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Mail available", "next": map[string]interface{}{"block-sender": "schemaKey == 'yes'", "extract-iocs": "schemaKey == 'yes'", "search-email-gateway": "schemaKey == 'no'"}, "order": 2, "schema": map[string]interface{}{"oneOf": []interface{}{map[string]interface{}{"properties": map[string]interface{}{"mail": map[string]interface{}{"title": "Mail", "type": "string", "x-display": "textarea"}, "schemaKey": map[string]interface{}{"const": "yes", "type": "string"}}, "required": []interface{}{"mail"}, "title": "Yes"}, map[string]interface{}{"properties": map[string]interface{}{"schemaKey": map[string]interface{}{"const": "no", "type": "string"}}, "title": "No"}}, "title": "Mail available", "type": "object"}, "type": "input"}, "search-email-gateway": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Search email gateway", "next": map[string]interface{}{"extract-iocs": ""}, "order": 4, "type": "task"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"},
			},
		},
		{
			name: "RemoveComment",
			args: args{method: "DELETE", url: "/api/tickets/8123/comments/0"},
			want: want{
				status: 200,
				body:   map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}, map[string]interface{}{"name": "leadreintermediate.io", "status": "malicious"}}, "created": "2021-12-12T12:12:12Z", "id": 8123, "modified": "2021-12-12T12:12:12Z", "name": "live zebra", "owner": "demo", "playbooks": map[string]interface{}{"phishing": map[string]interface{}{"name": "Phishing", "tasks": map[string]interface{}{"block-iocs": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Block IOCs", "order": 6, "type": "task"}, "block-sender": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Block sender", "next": map[string]interface{}{"extract-iocs": ""}, "order": 3, "type": "task"}, "board": map[string]interface{}{"active": true, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "order": 0, "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}, "escalate": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Escalate to CISO", "order": 1, "type": "task"}, "extract-iocs": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Extract IOCs", "next": map[string]interface{}{"block-iocs": ""}, "order": 5, "schema": map[string]interface{}{"properties": map[string]interface{}{"iocs": map[string]interface{}{"items": map[string]interface{}{"type": "string"}, "title": "IOCs", "type": "array"}}, "title": "Extract IOCs", "type": "object"}, "type": "input"}, "mail-available": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Mail available", "next": map[string]interface{}{"block-sender": "schemaKey == 'yes'", "extract-iocs": "schemaKey == 'yes'", "search-email-gateway": "schemaKey == 'no'"}, "order": 2, "schema": map[string]interface{}{"oneOf": []interface{}{map[string]interface{}{"properties": map[string]interface{}{"mail": map[string]interface{}{"title": "Mail", "type": "string", "x-display": "textarea"}, "schemaKey": map[string]interface{}{"const": "yes", "type": "string"}}, "required": []interface{}{"mail"}, "title": "Yes"}, map[string]interface{}{"properties": map[string]interface{}{"schemaKey": map[string]interface{}{"const": "no", "type": "string"}}, "title": "No"}}, "title": "Mail available", "type": "object"}, "type": "input"}, "search-email-gateway": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Search email gateway", "next": map[string]interface{}{"extract-iocs": ""}, "order": 4, "type": "task"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"},
			},
		},
		{
			name: "RemoveTicketPlaybook",
			args: args{method: "DELETE", url: "/api/tickets/8123/playbooks/phishing"},
			want: want{
				status: 200,
				body:   map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}, map[string]interface{}{"name": "leadreintermediate.io", "status": "malicious"}}, "created": "1985-04-12T23:20:50.52Z", "id": 8123, "modified": "1985-04-12T23:20:50.52Z", "name": "live zebra", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"},
			},
		},
		{
			name: "RunArtifact",
			args: args{method: "POST", url: "/api/tickets/8123/artifacts/leadreintermediate.io/run/hash.sha1"},
			want: want{
				status: 204,
				body:   nil,
			},
		},
		{
			name: "RunJob",
			args: args{method: "POST", url: "/api/jobs", data: map[string]interface{}{"automation": "hash.sha1", "message": map[string]interface{}{"payload": "test"}}},
			want: want{
				status: 204,
				body:   nil,
			},
		},
		{
			name: "RunTask",
			args: args{method: "POST", url: "/api/tickets/8123/playbooks/phishing/task/board/run"},
			want: want{
				status: 204,
				body:   nil,
			},
		},
		{
			name: "SetArtifact",
			args: args{method: "PUT", url: "/api/tickets/8123/artifacts/leadreintermediate.io", data: map[string]interface{}{"name": "leadreintermediate.io", "status": "clean"}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}, map[string]interface{}{"name": "leadreintermediate.io", "status": "clean"}}, "created": "2021-12-12T12:12:12Z", "id": 8123, "modified": "2021-12-12T12:12:12Z", "name": "live zebra", "owner": "demo", "playbooks": map[string]interface{}{"phishing": map[string]interface{}{"name": "Phishing", "tasks": map[string]interface{}{"block-iocs": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Block IOCs", "order": 6, "type": "task"}, "block-sender": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Block sender", "next": map[string]interface{}{"extract-iocs": ""}, "order": 3, "type": "task"}, "board": map[string]interface{}{"active": true, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "order": 0, "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}, "escalate": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Escalate to CISO", "order": 1, "type": "task"}, "extract-iocs": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Extract IOCs", "next": map[string]interface{}{"block-iocs": ""}, "order": 5, "schema": map[string]interface{}{"properties": map[string]interface{}{"iocs": map[string]interface{}{"items": map[string]interface{}{"type": "string"}, "title": "IOCs", "type": "array"}}, "title": "Extract IOCs", "type": "object"}, "type": "input"}, "mail-available": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Mail available", "next": map[string]interface{}{"block-sender": "schemaKey == 'yes'", "extract-iocs": "schemaKey == 'yes'", "search-email-gateway": "schemaKey == 'no'"}, "order": 2, "schema": map[string]interface{}{"oneOf": []interface{}{map[string]interface{}{"properties": map[string]interface{}{"mail": map[string]interface{}{"title": "Mail", "type": "string", "x-display": "textarea"}, "schemaKey": map[string]interface{}{"const": "yes", "type": "string"}}, "required": []interface{}{"mail"}, "title": "Yes"}, map[string]interface{}{"properties": map[string]interface{}{"schemaKey": map[string]interface{}{"const": "no", "type": "string"}}, "title": "No"}}, "title": "Mail available", "type": "object"}, "type": "input"}, "search-email-gateway": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Search email gateway", "next": map[string]interface{}{"extract-iocs": ""}, "order": 4, "type": "task"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"},
			},
		},
		{
			name: "SetReferences",
			args: args{method: "PUT", url: "/api/tickets/8125/references", data: []interface{}{map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"created": "2021-12-12T12:12:12Z", "id": 8125, "modified": "2021-12-12T12:12:12Z", "name": "phishing from selenafadel@von.com detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}, "schema": "{}", "status": "closed", "tickets": []interface{}{map[string]interface{}{"created": "2021-12-12T12:12:12Z", "id": 8126, "modified": "2021-12-12T12:12:12Z", "name": "Surfaceintroduce virus detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.centralworld-class.io/synthesize", "name": "university"}, map[string]interface{}{"href": "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", "name": "goal"}, map[string]interface{}{"href": "http://www.chiefsyndicate.io/action-items", "name": "unemployment"}}, "schema": "{}", "status": "closed", "type": "alert"}}, "type": "alert"},
			},
		},
		{
			name: "SetSchema",
			args: args{method: "PUT", url: "/api/tickets/8125/schema", data: "{}"},
			want: want{
				status: 200,
				body:   map[string]interface{}{"created": "2021-12-12T12:12:12Z", "id": 8125, "modified": "2021-12-12T12:12:12Z", "name": "phishing from selenafadel@von.com detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "https://www.seniorleading-edge.name/users/efficient", "name": "recovery"}, map[string]interface{}{"href": "http://www.dynamicseamless.com/clicks-and-mortar", "name": "force"}, map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}, "schema": "{}", "status": "closed", "tickets": []interface{}{map[string]interface{}{"created": "2021-12-12T12:12:12Z", "id": 8126, "modified": "2021-12-12T12:12:12Z", "name": "Surfaceintroduce virus detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.centralworld-class.io/synthesize", "name": "university"}, map[string]interface{}{"href": "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", "name": "goal"}, map[string]interface{}{"href": "http://www.chiefsyndicate.io/action-items", "name": "unemployment"}}, "schema": "{}", "status": "closed", "type": "alert"}}, "type": "alert"},
			},
		},
		{
			name: "SetTask",
			args: args{method: "PUT", url: "/api/tickets/8123/playbooks/phishing/task/board", data: map[string]interface{}{"active": true, "data": map[string]interface{}{"boardInvolved": true}, "done": false, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "order": 0, "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"artifacts": []interface{}{map[string]interface{}{"name": "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", "status": "unknown"}, map[string]interface{}{"name": "http://www.customerviral.io/scalable/vertical/killer", "status": "clean"}, map[string]interface{}{"name": "leadreintermediate.io", "status": "malicious"}}, "created": "2021-12-12T12:12:12Z", "id": 8123, "modified": "2021-12-12T12:12:12Z", "name": "live zebra", "owner": "demo", "playbooks": map[string]interface{}{"phishing": map[string]interface{}{"name": "Phishing", "tasks": map[string]interface{}{"block-iocs": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Block IOCs", "order": 6, "type": "task"}, "block-sender": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Block sender", "next": map[string]interface{}{"extract-iocs": ""}, "order": 3, "type": "task"}, "board": map[string]interface{}{"active": true, "created": "2021-12-12T12:12:12Z", "data": map[string]interface{}{"boardInvolved": true}, "done": false, "name": "Board Involvement?", "next": map[string]interface{}{"escalate": "boardInvolved == true", "mail-available": "boardInvolved == false"}, "order": 0, "schema": map[string]interface{}{"properties": map[string]interface{}{"boardInvolved": map[string]interface{}{"default": false, "title": "A board member is involved.", "type": "boolean"}}, "required": []interface{}{"boardInvolved"}, "title": "Board Involvement?", "type": "object"}, "type": "input"}, "escalate": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Escalate to CISO", "order": 1, "type": "task"}, "extract-iocs": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Extract IOCs", "next": map[string]interface{}{"block-iocs": ""}, "order": 5, "schema": map[string]interface{}{"properties": map[string]interface{}{"iocs": map[string]interface{}{"items": map[string]interface{}{"type": "string"}, "title": "IOCs", "type": "array"}}, "title": "Extract IOCs", "type": "object"}, "type": "input"}, "mail-available": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Mail available", "next": map[string]interface{}{"block-sender": "schemaKey == 'yes'", "extract-iocs": "schemaKey == 'yes'", "search-email-gateway": "schemaKey == 'no'"}, "order": 2, "schema": map[string]interface{}{"oneOf": []interface{}{map[string]interface{}{"properties": map[string]interface{}{"mail": map[string]interface{}{"title": "Mail", "type": "string", "x-display": "textarea"}, "schemaKey": map[string]interface{}{"const": "yes", "type": "string"}}, "required": []interface{}{"mail"}, "title": "Yes"}, map[string]interface{}{"properties": map[string]interface{}{"schemaKey": map[string]interface{}{"const": "no", "type": "string"}}, "title": "No"}}, "title": "Mail available", "type": "object"}, "type": "input"}, "search-email-gateway": map[string]interface{}{"active": false, "created": "2021-12-12T12:12:12Z", "done": false, "name": "Search email gateway", "next": map[string]interface{}{"extract-iocs": ""}, "order": 4, "type": "task"}}}}, "references": []interface{}{map[string]interface{}{"href": "https://www.leadmaximize.net/e-services/back-end", "name": "performance"}, map[string]interface{}{"href": "http://www.corporateinteractive.name/rich", "name": "autumn"}, map[string]interface{}{"href": "https://www.corporateintuitive.org/intuitive/platforms/integrate", "name": "suggest"}}, "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n", "status": "closed", "type": "incident"},
			},
		},
		{
			name: "UnlinkTicket",
			args: args{method: "DELETE", url: "/api/tickets/8126/tickets", data: 8125},
			want: want{
				status: 200,
				body:   map[string]interface{}{"created": "2021-12-12T12:12:12Z", "id": 8126, "modified": "2021-12-12T12:12:12Z", "name": "Surfaceintroduce virus detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.centralworld-class.io/synthesize", "name": "university"}, map[string]interface{}{"href": "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", "name": "goal"}, map[string]interface{}{"href": "http://www.chiefsyndicate.io/action-items", "name": "unemployment"}}, "schema": "{}", "status": "closed", "type": "alert"},
			},
		},
		{
			name: "UpdateAutomation",
			args: args{method: "PUT", url: "/api/automations/hash.sha1", data: map[string]interface{}{"id": "hash.sha1", "image": "docker.io/python:3", "script": "import sys\nimport json\nimport hashlib\n\n\ndef run(msg):\n    sha1 = hashlib.sha1(msg['payload'].encode('utf-8'))\n    return {'hash': sha1.hexdigest()}\n\n\nprint(json.dumps(run(json.loads(sys.argv[1]))))\n", "type": []interface{}{"global", "artifact", "playbook"}}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"id": "hash.sha1", "image": "docker.io/python:3", "script": "import sys\nimport json\nimport hashlib\n\n\ndef run(msg):\n    sha1 = hashlib.sha1(msg['payload'].encode('utf-8'))\n    return {'hash': sha1.hexdigest()}\n\n\nprint(json.dumps(run(json.loads(sys.argv[1]))))\n", "type": []interface{}{"global", "artifact", "playbook"}},
			},
		},
		{
			name: "UpdateCurrentUserData",
			args: args{method: "PUT", url: "/api/currentuserdata", data: map[string]interface{}{"email": "bob@example.org", "name": "Bob Bad"}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"email": "bob@example.org", "id": "bob", "name": "Bob Bad"},
			},
		},
		{
			name: "UpdateJob",
			args: args{method: "PUT", url: "/api/jobs/99cd67131b48", data: map[string]interface{}{"automation": "hash.sha1", "id": "99cd67131b48", "payload": "test", "status": "failed"}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"automation": "hash.sha1", "id": "99cd67131b48", "payload": "test", "status": "failed"},
			},
		},
		{
			name: "UpdatePlaybook",
			args: args{method: "PUT", url: "/api/playbooks/simple", data: map[string]interface{}{"yaml": "name: Simple\ntasks:\n  input:\n    name: Upload malware if possible\n    type: input\n    schema:\n      title: Malware\n      type: object\n      properties:\n        malware:\n          type: string\n          title: Select malware\n          default: \"\"\n    next:\n      hash: \"malware != ''\"\n\n  hash:\n    name: Hash the malware\n    type: automation\n    automation: hash.sha1\n    payload:\n      default: \"playbook.tasks['input'].data['malware']\"\n    next:\n      escalate:\n\n  escalate:\n    name: Escalate to malware team\n    type: task\n"}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"id": "simple", "name": "Simple", "yaml": "name: Simple\ntasks:\n  input:\n    name: Upload malware if possible\n    type: input\n    schema:\n      title: Malware\n      type: object\n      properties:\n        malware:\n          type: string\n          title: Select malware\n          default: \"\"\n    next:\n      hash: \"malware != ''\"\n\n  hash:\n    name: Hash the malware\n    type: automation\n    automation: hash.sha1\n    payload:\n      default: \"playbook.tasks['input'].data['malware']\"\n    next:\n      escalate:\n\n  escalate:\n    name: Escalate to malware team\n    type: task\n"},
			},
		},
		{
			name: "UpdateTemplate",
			args: args{method: "PUT", url: "/api/templates/default", data: map[string]interface{}{"name": "My Template", "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n"}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"id": "default", "name": "My Template", "schema": "{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n"},
			},
		},
		{
			name: "UpdateTicket",
			args: args{method: "PUT", url: "/api/tickets/8125", data: map[string]interface{}{"created": "2021-12-12T12:12:12Z", "modified": "2021-12-12T12:12:12Z", "name": "phishing from selenafadel@von.org detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "https://www.seniorleading-edge.name/users/efficient", "name": "recovery"}, map[string]interface{}{"href": "http://www.dynamicseamless.com/clicks-and-mortar", "name": "force"}, map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}, "schema": "{}", "status": "closed", "type": "alert"}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"created": "2021-12-12T12:12:12Z", "id": 8125, "modified": "2021-12-12T12:12:12Z", "name": "phishing from selenafadel@von.org detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "https://www.seniorleading-edge.name/users/efficient", "name": "recovery"}, map[string]interface{}{"href": "http://www.dynamicseamless.com/clicks-and-mortar", "name": "force"}, map[string]interface{}{"href": "http://www.leadscalable.biz/envisioneer", "name": "fund"}}, "schema": "{}", "status": "closed", "tickets": []interface{}{map[string]interface{}{"created": "2021-12-12T12:12:12Z", "id": 8126, "modified": "2021-12-12T12:12:12Z", "name": "Surfaceintroduce virus detected", "owner": "demo", "references": []interface{}{map[string]interface{}{"href": "http://www.centralworld-class.io/synthesize", "name": "university"}, map[string]interface{}{"href": "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", "name": "goal"}, map[string]interface{}{"href": "http://www.chiefsyndicate.io/action-items", "name": "unemployment"}}, "schema": "{}", "status": "closed", "type": "alert"}}, "type": "alert"},
			},
		},
		{
			name: "UpdateTicketType",
			args: args{method: "PUT", url: "/api/tickettypes/alert", data: map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-bell", "id": "alert", "name": "Alerts"}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"default_playbooks": []interface{}{}, "default_template": "default", "icon": "mdi-bell", "id": "alert", "name": "Alerts"},
			},
		},
		{
			name: "UpdateUser",
			args: args{method: "PUT", url: "/api/users/bob", data: map[string]interface{}{"roles": []interface{}{"analyst", "admin"}}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"apikey": false, "blocked": false, "id": "bob", "roles": []interface{}{"admin:backup:read", "admin:backup:restore", "admin:group:write", "admin:job:read", "admin:job:write", "admin:log:read", "admin:ticket:delete", "admin:user:write", "admin:userdata:read", "admin:userdata:write", "analyst:automation:read", "analyst:currentsettings:write", "analyst:currentuser:read", "analyst:currentuserdata:read", "analyst:file", "analyst:group:read", "analyst:playbook:read", "analyst:rule:read", "analyst:settings:read", "analyst:template:read", "analyst:ticket:read", "analyst:ticket:write", "analyst:tickettype:read", "analyst:user:read", "engineer:automation:write", "engineer:playbook:write", "engineer:rule:write", "engineer:template:write", "engineer:tickettype:write"}},
			},
		},
		{
			name: "UpdateUserData",
			args: args{method: "PUT", url: "/api/userdata/bob", data: map[string]interface{}{"blocked": false, "email": "bob@example.org", "name": "Bob Bad"}},
			want: want{
				status: 200,
				body:   map[string]interface{}{"email": "bob@example.org", "id": "bob", "name": "Bob Bad"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, _, _, _, db, _, server, cleanup, err := test.Server(t)
			if err != nil {
				t.Fatal(err)
			}
			defer cleanup()

			if err := test.SetupTestData(ctx, db); err != nil {
				t.Fatal(err)
			}

			setUser := func(context *gin.Context) {
				busdb.SetContext(context, test.Bob)
			}
			server.ApiGroup.Use(setUser)

			server.ConfigureRoutes()
			w := httptest.NewRecorder()

			// setup request
			var req *http.Request
			if tt.args.data != nil {
				b, err := json.Marshal(tt.args.data)
				if err != nil {
					t.Fatal(err)
				}

				req = httptest.NewRequest(tt.args.method, tt.args.url, bytes.NewBuffer(b))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.args.method, tt.args.url, nil)
			}

			// run request
			server.ServeHTTP(w, req)

			result := w.Result()

			// assert results
			if result.StatusCode != tt.want.status {
				msg, _ := io.ReadAll(result.Body)

				t.Fatalf("Status got = %v, want %v: %s", result.Status, tt.want.status, msg)
			}
			if tt.want.status != http.StatusNoContent {
				jsonEqual(t, result.Body, tt.want.body)
			}
		})
	}
}

func jsonEqual(t *testing.T, got io.Reader, want interface{}) {
	var gotObject, wantObject interface{}

	// load bytes
	wantBytes, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	gotBytes, err := io.ReadAll(got)
	if err != nil {
		t.Fatal(err)
	}

	fields := []string{"secret"}
	for _, field := range fields {
		gField := gjson.GetBytes(wantBytes, field)
		if gField.Exists() && gjson.GetBytes(gotBytes, field).Exists() {
			gotBytes, err = sjson.SetBytes(gotBytes, field, gField.Value())
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	// normalize bytes
	if err = json.Unmarshal(wantBytes, &wantObject); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(gotBytes, &gotObject); err != nil {
		t.Fatal(string(gotBytes), err)
	}

	// compare
	assert.Equal(t, wantObject, gotObject)
}
