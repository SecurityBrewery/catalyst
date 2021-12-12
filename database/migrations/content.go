package migrations

import _ "embed"

//go:embed templates/default.json
var DefaultTemplateSchema string

//go:embed automations/hash.sha1.py
var SHA1HashAutomation string

//go:embed automations/vt.hash.py
var VTHashAutomation string

//go:embed automations/thehive.py
var TheHiveAutomation string

//go:embed automations/comment.py
var CommentAutomation string

//go:embed playbooks/malware.yml
var MalwarePlaybook string

//go:embed playbooks/phishing.yml
var PhishingPlaybook string

//go:embed playbooks/simple.yaml
var SimplePlaybook string
