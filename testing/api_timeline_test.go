package testing

import (
	"net/http"
	"testing"

	"github.com/SecurityBrewery/catalyst/app/data"
)

func TestTimelineCollection(t *testing.T) {
	t.Parallel()

	testSets := []catalystTest{
		{
			baseTest: baseTest{
				Name:   "ListTimeline",
				Method: http.MethodGet,
				URL:    "/api/timeline?ticket=test-ticket",
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
					ExpectedStatus: http.StatusOK,
					ExpectedHeaders: map[string]string{
						"X-Total-Count": "1",
					},
					ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
				},
				{
					Name:           "Admin",
					Admin:          data.AdminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedHeaders: map[string]string{
						"X-Total-Count": "1",
					},
					ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
				},
			},
		},
		{
			baseTest: baseTest{
				Name:           "CreateTimeline",
				Method:         http.MethodPost,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/timeline",
				Body: s(map[string]any{
					"ticket":  "test-ticket",
					"message": "new",
					"time":    "2023-01-01T00:00:00Z",
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
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"ticket":"test-ticket"`,
					},
					ExpectedEvents: map[string]int{
						"OnRecordAfterCreateRequest":  1,
						"OnRecordBeforeCreateRequest": 1,
					},
				},
				{
					Name:           "Admin",
					Admin:          data.AdminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"ticket":"test-ticket"`,
					},
					ExpectedEvents: map[string]int{
						"OnRecordAfterCreateRequest":  1,
						"OnRecordBeforeCreateRequest": 1,
					},
				},
			},
		},
		{
			baseTest: baseTest{
				Name:   "GetTimeline",
				Method: http.MethodGet,
				URL:    "/api/timeline/h_test_timeline",
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
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"h_test_timeline"`,
					},
					ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
				},
				{
					Name:           "Admin",
					Admin:          data.AdminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"h_test_timeline"`,
					},
					ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
				},
			},
		},
		{
			baseTest: baseTest{
				Name:           "UpdateTimeline",
				Method:         http.MethodPatch,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/timeline/h_test_timeline",
				Body:           s(map[string]any{"message": "update"}),
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
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"h_test_timeline"`,
						`"message":"update"`,
					},
					ExpectedEvents: map[string]int{
						"OnRecordAfterUpdateRequest":  1,
						"OnRecordBeforeUpdateRequest": 1,
					},
				},
				{
					Name:           "Admin",
					Admin:          data.AdminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"h_test_timeline"`,
						`"message":"update"`,
					},
					ExpectedEvents: map[string]int{
						"OnRecordAfterUpdateRequest":  1,
						"OnRecordBeforeUpdateRequest": 1,
					},
				},
			},
		},
		{
			baseTest: baseTest{
				Name:   "DeleteTimeline",
				Method: http.MethodDelete,
				URL:    "/api/timeline/h_test_timeline",
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
					ExpectedStatus: http.StatusNoContent,
					ExpectedEvents: map[string]int{
						"OnRecordAfterDeleteRequest":  1,
						"OnRecordBeforeDeleteRequest": 1,
					},
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
