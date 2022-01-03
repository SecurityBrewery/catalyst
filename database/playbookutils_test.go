package database

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SecurityBrewery/catalyst/generated/model"
)

var playbook2 = &model.Playbook{
	Name: "Phishing",
	Tasks: map[string]*model.Task{
		"board": {Next: map[string]string{
			"escalate":    "boardInvolved == true",
			"aquire-mail": "boardInvolved == false",
		}},
		"escalate": {},
		"aquire-mail": {Next: map[string]string{
			"extract-iocs":         "schemaKey == 'yes'",
			"block-sender":         "schemaKey == 'yes'",
			"search-email-gateway": "schemaKey == 'no'",
		}},
		"extract-iocs":         {Next: map[string]string{"fetch-iocs": ""}},
		"fetch-iocs":           {Next: map[string]string{"block-iocs": ""}},
		"search-email-gateway": {Next: map[string]string{"block-iocs": ""}},
		"block-sender":         {Next: map[string]string{"block-iocs": ""}},
		"block-iocs":           {Next: map[string]string{"block-ioc": ""}},
		"block-ioc":            {},
	},
}

var playbook3 = &model.Playbook{
	Name: "Phishing",
	Tasks: map[string]*model.Task{
		"board": {Next: map[string]string{
			"escalate":    "boardInvolved == true",
			"aquire-mail": "boardInvolved == false",
		}, Data: map[string]interface{}{"boardInvolved": true}, Done: true},
		"escalate": {},
		"aquire-mail": {Next: map[string]string{
			"extract-iocs":         "schemaKey == 'yes'",
			"block-sender":         "schemaKey == 'yes'",
			"search-email-gateway": "schemaKey == 'no'",
		}},
		"extract-iocs":         {Next: map[string]string{"fetch-iocs": ""}},
		"fetch-iocs":           {Next: map[string]string{"block-iocs": ""}},
		"search-email-gateway": {Next: map[string]string{"block-iocs": ""}},
		"block-sender":         {Next: map[string]string{"block-iocs": ""}},
		"block-iocs":           {Next: map[string]string{"block-ioc": ""}},
		"block-ioc":            {},
	},
}

var playbook4 = &model.Playbook{
	Name: "Malware",
	Tasks: map[string]*model.Task{
		"file-or-hash": {Next: map[string]string{
			"enter-hash": "file == 'Hash'",
			"upload":     "file == 'File'",
		}},
		"enter-hash": {Next: map[string]string{
			"virustotal": "hash != ''",
		}},
		"upload": {Next: map[string]string{
			"hash": "malware",
		}},
		"hash":       {Next: map[string]string{"virustotal": ""}},
		"virustotal": {},
	},
}

func Test_canBeCompleted(t *testing.T) {
	type args struct {
		playbook *model.Playbook
		taskID   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"playbook2 board", args{playbook: playbook2, taskID: "board"}, true, false},
		{"playbook2 escalate", args{playbook: playbook2, taskID: "escalate"}, false, false},
		{"playbook2 aquire-mail", args{playbook: playbook2, taskID: "aquire-mail"}, false, false},
		{"playbook2 block-ioc", args{playbook: playbook2, taskID: "block-ioc"}, false, false},
		{"playbook3 board", args{playbook: playbook3, taskID: "board"}, false, false},
		{"playbook3 escalate", args{playbook: playbook3, taskID: "escalate"}, true, false},
		{"playbook3 aquire-mail", args{playbook: playbook3, taskID: "aquire-mail"}, false, false},
		{"playbook3 block-ioc", args{playbook: playbook3, taskID: "block-ioc"}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := activePlaybook(tt.args.playbook, tt.args.taskID)
			if (err != nil) != tt.wantErr {
				t.Errorf("activePlaybook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("activePlaybook() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_playbookOrder(t *testing.T) {
	type args struct {
		playbook *model.Playbook
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"playbook4", args{playbook: playbook4}, []string{"file-or-hash", "enter-hash", "upload", "hash", "virustotal"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toPlaybookResponse(tt.args.playbook)
			if (err != nil) != tt.wantErr {
				t.Errorf("activePlaybook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			names := make([]string, len(got.Tasks))
			for name, task := range got.Tasks {
				names[task.Order] = name
			}

			assert.Equal(t, tt.want, names)
		})
	}
}
