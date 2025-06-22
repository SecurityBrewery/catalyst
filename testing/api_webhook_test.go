package testing

import (
	"net/http"
	"testing"

	"github.com/SecurityBrewery/catalyst/app/data"
)

func TestWebhooksCollection(t *testing.T) {
	t.Parallel()

	testSets := []catalystTest{
		{
			baseTest: baseTest{
				Name:   "ListWebhooks",
				Method: http.MethodGet,
				URL:    "/api/webhooks",
			},
			userTests: []userTest{
				{
					Name:            "Unauthorized",
					ExpectedStatus:  http.StatusUnauthorized,
					ExpectedContent: []string{`"invalid bearer token"`},
					ExpectedEvents:  map[string]int{},
				},
				{
					Name:            "Analyst",
					AuthRecord:      data.AnalystEmail,
					ExpectedStatus:  http.StatusUnauthorized,
					ExpectedContent: []string{`"missing required scopes"`},
				},
				{
					Name:            "Admin",
					Admin:           data.AdminEmail,
					ExpectedStatus:  http.StatusOK,
					ExpectedHeaders: map[string]string{"X-Total-Count": "1"},
					ExpectedContent: []string{`"id":"w_test_webhook"`},
					ExpectedEvents:  map[string]int{},
				},
			},
		},
		{
			baseTest: baseTest{
				Name:           "CreateWebhook",
				Method:         http.MethodPost,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/webhooks",
				Body: s(map[string]any{
					"name":        "new",
					"collection":  "tickets",
					"destination": "https://example.com/new",
				}),
			},
			userTests: []userTest{
				{
					Name:            "Unauthorized",
					ExpectedStatus:  http.StatusUnauthorized,
					ExpectedContent: []string{`"invalid bearer token"`},
				},
				{
					Name:            "Analyst",
					AuthRecord:      data.AnalystEmail,
					ExpectedStatus:  http.StatusUnauthorized,
					ExpectedContent: []string{`"missing required scopes"`},
				},
				{
					Name:            "Admin",
					Admin:           data.AdminEmail,
					ExpectedStatus:  http.StatusOK,
					ExpectedContent: []string{`"name":"new"`},
					ExpectedEvents: map[string]int{
						"OnRecordAfterCreateRequest":  1,
						"OnRecordBeforeCreateRequest": 1,
					},
				},
			},
		},
		{
			baseTest: baseTest{
				Name:   "GetWebhook",
				Method: http.MethodGet,
				URL:    "/api/webhooks/w_test_webhook",
			},
			userTests: []userTest{
				{
					Name:            "Unauthorized",
					ExpectedStatus:  http.StatusUnauthorized,
					ExpectedContent: []string{`"invalid bearer token"`},
				},
				{
					Name:            "Analyst",
					AuthRecord:      data.AnalystEmail,
					ExpectedStatus:  http.StatusUnauthorized,
					ExpectedContent: []string{`"missing required scopes"`},
				},
				{
					Name:            "Admin",
					Admin:           data.AdminEmail,
					ExpectedStatus:  http.StatusOK,
					ExpectedContent: []string{`"id":"w_test_webhook"`},
					ExpectedEvents:  map[string]int{"OnRecordViewRequest": 1},
				},
			},
		},
		{
			baseTest: baseTest{
				Name:           "UpdateWebhook",
				Method:         http.MethodPatch,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/webhooks/w_test_webhook",
				Body:           s(map[string]any{"name": "update"}),
			},
			userTests: []userTest{
				{
					Name:            "Unauthorized",
					ExpectedStatus:  http.StatusUnauthorized,
					ExpectedContent: []string{`"invalid bearer token"`},
				},
				{
					Name:            "Analyst",
					AuthRecord:      data.AnalystEmail,
					ExpectedStatus:  http.StatusUnauthorized,
					ExpectedContent: []string{`"missing required scopes"`},
				},
				{
					Name:            "Admin",
					Admin:           data.AdminEmail,
					ExpectedStatus:  http.StatusOK,
					ExpectedContent: []string{`"id":"w_test_webhook"`, `"name":"update"`},
					ExpectedEvents: map[string]int{
						"OnRecordAfterUpdateRequest":  1,
						"OnRecordBeforeUpdateRequest": 1,
					},
				},
			},
		},
		{
			baseTest: baseTest{
				Name:   "DeleteWebhook",
				Method: http.MethodDelete,
				URL:    "/api/webhooks/w_test_webhook",
			},
			userTests: []userTest{
				{
					Name:            "Unauthorized",
					ExpectedStatus:  http.StatusUnauthorized,
					ExpectedContent: []string{`"invalid bearer token"`},
				},
				{
					Name:            "Analyst",
					AuthRecord:      data.AnalystEmail,
					ExpectedStatus:  http.StatusUnauthorized,
					ExpectedContent: []string{`"missing required scopes"`},
				},
				{
					Name:           "Admin",
					Admin:          data.AdminEmail,
					ExpectedStatus: http.StatusNoContent,
					ExpectedEvents: map[string]int{
						"OnRecordAfterDeleteRequest":  1,
						"OnRecordBeforeDeleteRequest": 1,
					},
				},
			},
		},
	}

	for _, testSet := range testSets {
		t.Run(testSet.baseTest.Name, func(t *testing.T) {
			t.Parallel()

			for _, userTest := range testSet.userTests {
				t.Run(userTest.Name, func(t *testing.T) {
					t.Parallel()

					runMatrixTest(t, testSet.baseTest, userTest)
				})
			}
		})
	}
}
