package testing

import (
	"net/http"
	"testing"
)

func TestReactionsCollection(t *testing.T) {
	t.Parallel()

	testSets := []catalystTest{
		{
			baseTest: BaseTest{
				Name:   "ListReactions",
				Method: http.MethodGet,
				URL:    "/api/collections/reactions/records",
			},
			userTests: []UserTest{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"totalItems":0`,
						`"items":[]`,
					},
					ExpectedEvents: []string{"OnRecordsListRequest"},
				},
				{
					Name:           "Analyst",
					AuthRecord:     analystEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"totalItems":3`,
						`"id":"r_reaction"`,
					},
					NotExpectedContent: []string{
						`"items":[]`,
					},
					ExpectedEvents: []string{"OnRecordsListRequest"},
				},
				{
					Name:           "Admin",
					Admin:          adminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"totalItems":3`,
						`"id":"r_reaction"`,
					},
					NotExpectedContent: []string{
						`"items":[]`,
					},
					ExpectedEvents: []string{"OnRecordsListRequest"},
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
			},
			userTests: []UserTest{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusBadRequest,
					ExpectedContent: []string{
						`"message":"Failed to create record."`,
					},
					ExpectedEvents: []string{"OnRecordCreateRequest"},
				},
				{
					Name:           "Analyst",
					AuthRecord:     analystEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"name":"test"`,
					},
					NotExpectedContent: []string{
						`"items":[]`,
					},
					ExpectedEvents: []string{
						"OnModelAfterCreateSuccess",
						"OnRecordCreateRequest",
					},
				},
				{
					Name:           "Admin",
					Admin:          adminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"name":"test"`,
					},
					NotExpectedContent: []string{
						`"items":[]`,
					},
					ExpectedEvents: []string{
						"OnModelAfterCreateSuccess",
						"OnRecordCreateRequest",
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
			},
			userTests: []UserTest{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusNotFound,
					ExpectedContent: []string{
						`"message":"The requested resource wasn't found."`,
					},
				},
				{
					Name:           "Analyst",
					AuthRecord:     analystEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"r_reaction"`,
					},
					ExpectedEvents: []string{"OnRecordViewRequest"},
				},
				{
					Name:           "Admin",
					Admin:          adminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"r_reaction"`,
					},
					ExpectedEvents: []string{"OnRecordViewRequest"},
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
			},
			userTests: []UserTest{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusNotFound,
					ExpectedContent: []string{
						`"message":"The requested resource wasn't found."`,
					},
				},
				{
					Name:           "Analyst",
					AuthRecord:     analystEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"r_reaction"`,
						`"name":"update"`,
					},
					ExpectedEvents: []string{
						"OnModelAfterUpdateSuccess",
						"OnRecordUpdateRequest",
					},
				},
				{
					Name:           "Admin",
					Admin:          adminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"r_reaction"`,
						`"name":"update"`,
					},
					ExpectedEvents: []string{
						"OnModelAfterUpdateSuccess",
						"OnRecordUpdateRequest",
					},
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:   "DeleteReaction",
				Method: http.MethodDelete,
				URL:    "/api/collections/reactions/records/r_reaction",
			},
			userTests: []UserTest{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusNotFound,
					ExpectedContent: []string{
						`"message":"The requested resource wasn't found."`,
					},
				},
				{
					Name:           "Analyst",
					AuthRecord:     analystEmail,
					ExpectedStatus: http.StatusNoContent,
					ExpectedEvents: []string{
						"OnModelAfterDeleteSuccess",
						"OnRecordDeleteRequest",
					},
				},
				{
					Name:           "Admin",
					Admin:          adminEmail,
					ExpectedStatus: http.StatusNoContent,
					ExpectedEvents: []string{
						"OnModelAfterDeleteSuccess",
						"OnRecordDeleteRequest",
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
