package testing

import (
	"net/http"
	"testing"
)

func TestReactionsCollection(t *testing.T) {
	baseApp, adminToken, analystToken, baseAppCleanup := BaseApp(t)
	defer baseAppCleanup()

	testSets := []authMatrixText{
		{
			baseTest: BaseTest{
				Name:           "ListReactions",
				Method:         http.MethodGet,
				URL:            "/api/collections/reactions/records",
				TestAppFactory: AppFactory(baseApp),
			},
			authBasedExpectations: []AuthBasedExpectation{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"totalItems":0`,
						`"items":[]`,
					},
					ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
				},
				{
					Name:           "Analyst",
					RequestHeaders: map[string]string{"Authorization": analystToken},
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"totalItems":3`,
						`"id":"r_reaction"`,
					},
					NotExpectedContent: []string{
						`"items":[]`,
					},
					ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
				},
				{
					Name:           "Admin",
					RequestHeaders: map[string]string{"Authorization": adminToken},
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"totalItems":3`,
						`"id":"r_reaction"`,
					},
					NotExpectedContent: []string{
						`"items":[]`,
					},
					ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:           "CreateReaction",
				Method:         http.MethodPost,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/collections/reactions/records",
				Body: s(map[string]any{
					"name":        "test",
					"trigger":     "webhook",
					"triggerdata": map[string]any{"path": "test"},
					"action":      "python",
					"actiondata":  map[string]any{"script": "print('Hello, World!')"},
				}),
				TestAppFactory: AppFactory(baseApp),
			},
			authBasedExpectations: []AuthBasedExpectation{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusBadRequest,
					ExpectedContent: []string{
						`"message":"Failed to create record."`,
					},
				},
				{
					Name:           "Analyst",
					RequestHeaders: map[string]string{"Authorization": analystToken},
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"name":"test"`,
					},
					NotExpectedContent: []string{
						`"items":[]`,
					},
					ExpectedEvents: map[string]int{
						"OnModelAfterCreate":          1,
						"OnModelBeforeCreate":         1,
						"OnRecordAfterCreateRequest":  1,
						"OnRecordBeforeCreateRequest": 1,
					},
				},
				{
					Name:           "Admin",
					RequestHeaders: map[string]string{"Authorization": adminToken},
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"name":"test"`,
					},
					NotExpectedContent: []string{
						`"items":[]`,
					},
					ExpectedEvents: map[string]int{
						"OnModelAfterCreate":          1,
						"OnModelBeforeCreate":         1,
						"OnRecordAfterCreateRequest":  1,
						"OnRecordBeforeCreateRequest": 1,
					},
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:           "GetReaction",
				Method:         http.MethodGet,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/collections/reactions/records/r_reaction",
				TestAppFactory: AppFactory(baseApp),
			},
			authBasedExpectations: []AuthBasedExpectation{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusNotFound,
					ExpectedContent: []string{
						`"message":"The requested resource wasn't found."`,
					},
				},
				{
					Name:           "Analyst",
					RequestHeaders: map[string]string{"Authorization": analystToken},
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"r_reaction"`,
					},
					ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
				},
				{
					Name:           "Admin",
					RequestHeaders: map[string]string{"Authorization": adminToken},
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"r_reaction"`,
					},
					ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:           "UpdateReaction",
				Method:         http.MethodPatch,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/collections/reactions/records/r_reaction",
				Body:           s(map[string]any{"name": "update"}),
				TestAppFactory: AppFactory(baseApp),
			},
			authBasedExpectations: []AuthBasedExpectation{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusNotFound,
					ExpectedContent: []string{
						`"message":"The requested resource wasn't found."`,
					},
				},
				{
					Name:           "Analyst",
					RequestHeaders: map[string]string{"Authorization": analystToken},
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"r_reaction"`,
						`"name":"update"`,
					},
					ExpectedEvents: map[string]int{
						"OnModelAfterUpdate":          1,
						"OnModelBeforeUpdate":         1,
						"OnRecordAfterUpdateRequest":  1,
						"OnRecordBeforeUpdateRequest": 1,
					},
				},
				{
					Name:           "Admin",
					RequestHeaders: map[string]string{"Authorization": adminToken},
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"r_reaction"`,
						`"name":"update"`,
					},
					ExpectedEvents: map[string]int{
						"OnModelAfterUpdate":          1,
						"OnModelBeforeUpdate":         1,
						"OnRecordAfterUpdateRequest":  1,
						"OnRecordBeforeUpdateRequest": 1,
					},
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:           "DeleteReaction",
				Method:         http.MethodDelete,
				URL:            "/api/collections/reactions/records/r_reaction",
				TestAppFactory: AppFactory(baseApp),
			},
			authBasedExpectations: []AuthBasedExpectation{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusNotFound,
					ExpectedContent: []string{
						`"message":"The requested resource wasn't found."`,
					},
				},
				{
					Name:           "Analyst",
					RequestHeaders: map[string]string{"Authorization": analystToken},
					ExpectedStatus: http.StatusNoContent,
					ExpectedEvents: map[string]int{
						"OnModelAfterDelete":          1,
						"OnModelBeforeDelete":         1,
						"OnRecordAfterDeleteRequest":  1,
						"OnRecordBeforeDeleteRequest": 1,
					},
				},
				{
					Name:           "Admin",
					RequestHeaders: map[string]string{"Authorization": adminToken},
					ExpectedStatus: http.StatusNoContent,
					ExpectedEvents: map[string]int{
						"OnModelAfterDelete":          1,
						"OnModelBeforeDelete":         1,
						"OnRecordAfterDeleteRequest":  1,
						"OnRecordBeforeDeleteRequest": 1,
					},
				},
			},
		},
	}
	for _, testSet := range testSets {
		t.Run(testSet.baseTest.Name, func(t *testing.T) {
			for _, authBasedExpectation := range testSet.authBasedExpectations {
				scenario := mergeScenario(testSet.baseTest, authBasedExpectation)
				scenario.Test(t)
			}
		})
	}
}
