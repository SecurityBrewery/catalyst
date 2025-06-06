package testing

import (
	"net/http"
	"testing"

	"github.com/SecurityBrewery/catalyst/app/database"
)

func TestTasksCollection(t *testing.T) {
	t.Parallel()

	testSets := []catalystTest{
		{
			baseTest: BaseTest{
				Name:   "ListTasks",
				Method: http.MethodGet,
				URL:    "/api/tasks?ticket=test-ticket",
			},
			userTests: []UserTest{
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
					AuthRecord:     database.AnalystEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedHeaders: map[string]string{
						"X-Total-Count": "1",
					},
					ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
				},
				{
					Name:           "Admin",
					Admin:          database.AdminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedHeaders: map[string]string{
						"X-Total-Count": "1",
					},
					ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:           "CreateTask",
				Method:         http.MethodPost,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/tasks",
				Body: s(map[string]any{
					"ticket": "test-ticket",
					"name":   "new",
					"open":   true,
					"owner":  "u_bob_analyst",
				}),
			},
			userTests: []UserTest{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusUnauthorized,
					ExpectedContent: []string{
						`"invalid bearer token"`,
					},
				},
				{
					Name:           "Analyst",
					AuthRecord:     database.AnalystEmail,
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
					Admin:          database.AdminEmail,
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
			baseTest: BaseTest{
				Name:   "GetTask",
				Method: http.MethodGet,
				URL:    "/api/tasks/k_test_task",
			},
			userTests: []UserTest{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusUnauthorized,
					ExpectedContent: []string{
						`"invalid bearer token"`,
					},
				},
				{
					Name:           "Analyst",
					AuthRecord:     database.AnalystEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"k_test_task"`,
					},
					ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
				},
				{
					Name:           "Admin",
					Admin:          database.AdminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"k_test_task"`,
					},
					ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:           "UpdateTask",
				Method:         http.MethodPatch,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/tasks/k_test_task",
				Body:           s(map[string]any{"name": "update"}),
			},
			userTests: []UserTest{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusUnauthorized,
					ExpectedContent: []string{
						`"invalid bearer token"`,
					},
				},
				{
					Name:           "Analyst",
					AuthRecord:     database.AnalystEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"k_test_task"`,
						`"name":"update"`,
					},
					ExpectedEvents: map[string]int{
						"OnRecordAfterUpdateRequest":  1,
						"OnRecordBeforeUpdateRequest": 1,
					},
				},
				{
					Name:           "Admin",
					Admin:          database.AdminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"k_test_task"`,
						`"name":"update"`,
					},
					ExpectedEvents: map[string]int{
						"OnRecordAfterUpdateRequest":  1,
						"OnRecordBeforeUpdateRequest": 1,
					},
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:   "DeleteTask",
				Method: http.MethodDelete,
				URL:    "/api/tasks/k_test_task",
			},
			userTests: []UserTest{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusUnauthorized,
					ExpectedContent: []string{
						`"invalid bearer token"`,
					},
				},
				{
					Name:           "Analyst",
					AuthRecord:     database.AnalystEmail,
					ExpectedStatus: http.StatusNoContent,
					ExpectedEvents: map[string]int{
						"OnRecordAfterDeleteRequest":  1,
						"OnRecordBeforeDeleteRequest": 1,
					},
				},
				{
					Name:           "Admin",
					Admin:          database.AdminEmail,
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
