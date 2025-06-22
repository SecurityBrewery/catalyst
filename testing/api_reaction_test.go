package testing

import (
	"net/http"
	"testing"

	"github.com/SecurityBrewery/catalyst/app/data"
)

func TestReactionsCollection(t *testing.T) {
	t.Parallel()

	testSets := []catalystTest{
		{
			baseTest: baseTest{
				Name:   "ListReactions",
				Method: http.MethodGet,
				URL:    "/api/reactions",
			},
			userTests: []userTest{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusUnauthorized,
					ExpectedContent: []string{
						`"invalid bearer token"`,
					},
					ExpectedEvents: map[string]int{},
				},
				{
					Name:           "Analyst",
					AuthRecord:     data.AnalystEmail,
					ExpectedStatus: http.StatusUnauthorized,
					ExpectedContent: []string{
						`"missing required scopes"`,
					},
				},
				{
					Name:           "Admin",
					Admin:          data.AdminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"r-test-webhook"`,
					},
					ExpectedHeaders: map[string]string{
						"X-Total-Count": "3",
					},
					NotExpectedContent: []string{
						`"items":[]`,
					},
					ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
				},
			},
		},
		{
			baseTest: baseTest{
				Name:           "CreateReaction",
				Method:         http.MethodPost,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/reactions",
				Body: s(map[string]any{
					"name":        "test",
					"trigger":     "webhook",
					"triggerdata": map[string]any{"path": "test"},
					"action":      "python",
					"actiondata":  map[string]any{"script": "print('Hello, World!')"},
				}),
			},
			userTests: []userTest{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusUnauthorized,
					ExpectedContent: []string{
						`"invalid bearer token"`,
					},
				},
				{
					Name:           "Analyst",
					AuthRecord:     data.AnalystEmail,
					ExpectedStatus: http.StatusUnauthorized,
					ExpectedContent: []string{
						`"missing required scopes"`,
					},
				},
				{
					Name:           "Admin",
					Admin:          data.AdminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"name":"test"`,
					},
					NotExpectedContent: []string{
						`"items":[]`,
					},
					ExpectedEvents: map[string]int{
						// "OnModelAfterCreate":          1,
						// "OnModelBeforeCreate":         1,
						"OnRecordAfterCreateRequest":  1,
						"OnRecordBeforeCreateRequest": 1,
					},
				},
			},
		},
		{
			baseTest: baseTest{
				Name:           "GetReaction",
				Method:         http.MethodGet,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/reactions/r-test-webhook",
			},
			userTests: []userTest{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusUnauthorized,
					ExpectedContent: []string{
						`"invalid bearer token"`,
					},
				},
				{
					Name:           "Analyst",
					AuthRecord:     data.AnalystEmail,
					ExpectedStatus: http.StatusUnauthorized,
					ExpectedContent: []string{
						`"missing required scopes"`,
					},
				},
				{
					Name:           "Admin",
					Admin:          data.AdminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"r-test-webhook"`,
					},
					ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
				},
			},
		},
		{
			baseTest: baseTest{
				Name:           "UpdateReaction",
				Method:         http.MethodPatch,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/reactions/r-test-webhook",
				Body:           s(map[string]any{"name": "update"}),
			},
			userTests: []userTest{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusUnauthorized,
					ExpectedContent: []string{
						`"invalid bearer token"`,
					},
				},
				{
					Name:           "Analyst",
					AuthRecord:     data.AnalystEmail,
					ExpectedStatus: http.StatusUnauthorized,
					ExpectedContent: []string{
						`"missing required scopes"`,
					},
				},
				{
					Name:           "Admin",
					Admin:          data.AdminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"r-test-webhook"`,
						`"name":"update"`,
					},
					ExpectedEvents: map[string]int{
						// "OnModelAfterUpdate":          1,
						// "OnModelBeforeUpdate":         1,
						"OnRecordAfterUpdateRequest":  1,
						"OnRecordBeforeUpdateRequest": 1,
					},
				},
			},
		},
		{
			baseTest: baseTest{
				Name:   "DeleteReaction",
				Method: http.MethodDelete,
				URL:    "/api/reactions/r-test-webhook",
			},
			userTests: []userTest{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusUnauthorized,
					ExpectedContent: []string{
						`"invalid bearer token"`,
					},
				},
				{
					Name:           "Analyst",
					AuthRecord:     data.AnalystEmail,
					ExpectedStatus: http.StatusUnauthorized,
					ExpectedContent: []string{
						`"missing required scopes"`,
					},
				},
				{
					Name:           "Admin",
					Admin:          data.AdminEmail,
					ExpectedStatus: http.StatusNoContent,
					ExpectedEvents: map[string]int{
						// "OnModelAfterDelete":          1,
						// "OnModelBeforeDelete":         1,
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
