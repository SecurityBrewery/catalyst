package test

import (
	"context"
	"time"

	"github.com/SecurityBrewery/catalyst/database"
	"github.com/SecurityBrewery/catalyst/database/migrations"
	"github.com/SecurityBrewery/catalyst/generated/models"
	"github.com/SecurityBrewery/catalyst/pointer"
)

var bobSetting = &models.UserData{Email: pointer.String("bob@example.org"), Name: pointer.String("Bob Bad")}
var bobForm = &models.UserForm{ID: "bob", Blocked: false, Roles: []string{"admin"}}
var Bob = &models.UserResponse{ID: "bob", Blocked: false, Roles: []string{"admin"}}

func SetupTestData(ctx context.Context, db *database.Database) error {
	if err := db.UserDataCreate(ctx, "bob", bobSetting); err != nil {
		return err
	}

	if _, err := db.UserCreate(ctx, bobForm); err != nil {
		return err
	}
	if _, err := db.UserCreate(ctx, &models.UserForm{ID: "script", Roles: []string{"engineer"}, Apikey: true}); err != nil {
		return err
	}

	if _, err := db.TicketBatchCreate(ctx, []*models.TicketForm{
		{
			ID:         pointer.Int64(8125),
			Created:    parse("2021-10-02T18:04:59.078186+02:00"),
			Modified:   parse("2021-10-02T18:04:59.078186+02:00"),
			Name:       "phishing from selenafadel@von.com detected",
			Owner:      pointer.String("demo"),
			References: []*models.Reference{{Href: "https://www.seniorleading-edge.name/users/efficient", Name: "recovery"}, {Href: "http://www.dynamicseamless.com/clicks-and-mortar", Name: "force"}, {Href: "http://www.leadscalable.biz/envisioneer", Name: "fund"}},
			Schema:     pointer.String("{}"),
			Status:     "closed",
			Type:       "alert",
		}, {
			ID:         pointer.Int64(8126),
			Created:    parse("2021-10-02T18:04:59.078186+02:00"),
			Modified:   parse("2021-10-02T18:04:59.078186+02:00"),
			Name:       "Surfaceintroduce virus detected",
			Owner:      pointer.String("demo"),
			References: []*models.Reference{{Href: "http://www.centralworld-class.io/synthesize", Name: "university"}, {Href: "https://www.futurevirtual.org/supply-chains/markets/sticky/iterate", Name: "goal"}, {Href: "http://www.chiefsyndicate.io/action-items", Name: "unemployment"}},
			Schema:     pointer.String("{}"),
			Status:     "closed",
			Type:       "alert",
		}, {
			ID:       pointer.Int64(8123),
			Created:  parse("2021-10-02T18:04:59.078206+02:00"),
			Modified: parse("2021-10-02T18:04:59.078206+02:00"),
			Artifacts: []*models.Artifact{
				{Name: "94d5cab6f5fe3422a447ab15436e7a672bc0c09a", Status: pointer.String("unknown")},
				{Name: "http://www.customerviral.io/scalable/vertical/killer", Status: pointer.String("clean")},
				{Name: "leadreintermediate.io", Status: pointer.String("malicious")},
			},
			Name:       "live zebra",
			Owner:      pointer.String("demo"),
			References: []*models.Reference{{Href: "https://www.leadmaximize.net/e-services/back-end", Name: "performance"}, {Href: "http://www.corporateinteractive.name/rich", Name: "autumn"}, {Href: "https://www.corporateintuitive.org/intuitive/platforms/integrate", Name: "suggest"}},
			Schema:     pointer.String("{\n  \"definitions\": {},\n  \"$schema\": \"http://json-schema.org/draft-07/schema#\",\n  \"$id\": \"https://example.com/object1618746510.json\",\n  \"title\": \"Event\",\n  \"type\": \"object\",\n  \"required\": [\n    \"severity\",\n    \"description\",\n    \"tlp\"\n  ],\n  \"properties\": {\n    \"severity\": {\n      \"$id\": \"#root/severity\",\n      \"title\": \"Severity\",\n      \"type\": \"string\",\n      \"default\": \"Medium\",\n      \"nx-enum\": [\n        \"Low\",\n        \"Medium\",\n        \"High\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"Low\",\n          \"title\": \"Low\",\n          \"icon\": \"mdi-chevron-up\"\n        },\n        {\n          \"const\": \"Medium\",\n          \"title\": \"Medium\",\n          \"icon\": \"mdi-chevron-double-up\"\n        },\n        {\n          \"const\": \"High\",\n          \"title\": \"High\",\n          \"icon\": \"mdi-chevron-triple-up\"\n        }\n      ]\n    },\n    \"tlp\": {\n      \"$id\": \"#root/tlp\",\n      \"title\": \"TLP\",\n      \"type\": \"string\",\n      \"nx-enum\": [\n        \"White\",\n        \"Green\",\n        \"Amber\",\n        \"Red\"\n      ],\n      \"x-cols\": 6,\n      \"x-class\": \"pr-2\",\n      \"x-display\": \"icon\",\n      \"x-itemIcon\": \"icon\",\n      \"oneOf\": [\n        {\n          \"const\": \"White\",\n          \"title\": \"White\",\n          \"icon\": \"mdi-alpha-w\"\n        },\n        {\n          \"const\": \"Green\",\n          \"title\": \"Green\",\n          \"icon\": \"mdi-alpha-g\"\n        },\n        {\n          \"const\": \"Amber\",\n          \"title\": \"Amber\",\n          \"icon\": \"mdi-alpha-a\"\n        },\n        {\n          \"const\": \"Red\",\n          \"title\": \"Red\",\n          \"icon\": \"mdi-alpha-r\"\n        }\n      ]\n    },\n    \"description\": {\n      \"$id\": \"#root/description\",\n      \"title\": \"Description\",\n      \"type\": \"string\",\n      \"x-display\": \"textarea\",\n      \"x-class\": \"pr-2\"\n    }\n  }\n}\n"),
			Status:     "closed",
			Type:       "incident",
			Playbooks: []*models.PlaybookTemplateForm{
				{Yaml: migrations.PhishingPlaybook},
			},
		},
	}); err != nil {
		return err
	}

	if err := db.RelatedCreate(ctx, 8125, 8126); err != nil {
		return err
	}
	if _, err := db.PlaybookCreate(ctx, &models.PlaybookTemplateForm{Yaml: "name: Simple\ntasks:\n  input:\n    name: Enter something to hash\n    type: input\n    schema:\n      title: Something\n      type: object\n      properties:\n        something:\n          type: string\n          title: Something\n          default: \"\"\n    next:\n      hash: \"something != ''\"\n\n  hash:\n    name: Hash the something\n    type: automation\n    automation: hash.sha1\n    payload:\n      default: \"playbook.tasks['input'].data['something']\"\n    next:\n      comment: \"hash != ''\"\n\n  comment:\n    name: Comment the hash\n    type: automation\n    automation: comment\n    payload:\n      default: \"playbook.tasks['hash'].data['hash']\"\n    next:\n      done: \"done\"\n\n  done:\n    name: You can close this case now\n    type: task\n"}); err != nil {
		return err
	}

	if _, err := db.LogCreate(ctx, "tickets/294511", "Fail run account resist lend solve incident centre priority temperature. Cause change distribution examine location technique shape partner milk customer. Rail tea plate soil report cook railway interpretation breath action. Exercise dream accept park conclusion addition shoot assistance may answer. Gold writer link stop combine hear power name commitment operation. Determine lifespan support grow degree henry exclude detail set religion. Direct library policy convention chain retain discover ride walk student. Gather proposal select march aspect play noise avoid encourage employ. Assessment preserve transport combine wish influence income guess run stand. Charge limit crime ignore statement foundation study issue stop claim."); err != nil {
		return err
	}

	return nil
}

func parse(s string) *time.Time {
	modified, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return &modified
}
