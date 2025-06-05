package testing

import (
	"net/http"
	"testing"

	"github.com/SecurityBrewery/catalyst/app/database"
)

func TestTypesCollection(t *testing.T) {
	t.Parallel()

	testSets := []catalystTest{
		{
			baseTest: BaseTest{
				Name:   "ListTypes",
				Method: http.MethodGet,
				URL:    "/api/types",
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
						"X-Total-Count": "4",
					},
					ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
				},
				{
					Name:           "Admin",
					Admin:          database.AdminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedHeaders: map[string]string{
						"X-Total-Count": "4",
					},
					ExpectedEvents: map[string]int{"OnRecordsListRequest": 1},
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:           "CreateType",
				Method:         http.MethodPost,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/types",
				Body: s(map[string]any{
					"singular": "Example",
					"plural":   "Examples",
					"icon":     "Bug",
					"schema":   map[string]any{},
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
					ExpectedStatus: http.StatusUnauthorized,
					ExpectedContent: []string{
						`"missing required scopes"`,
					},
				},
				{
					Name:           "Admin",
					Admin:          database.AdminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"singular":"Example"`,
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
				Name:   "GetType",
				Method: http.MethodGet,
				URL:    "/api/types/test-type",
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
						`"id":"test-type"`,
					},
					ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
				},
				{
					Name:           "Admin",
					Admin:          database.AdminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"test-type"`,
					},
					ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:           "UpdateType",
				Method:         http.MethodPatch,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/types/test-type",
				Body:           s(map[string]any{"singular": "Update"}),
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
					ExpectedStatus: http.StatusUnauthorized,
					ExpectedContent: []string{
						`"missing required scopes"`,
					},
				},
				{
					Name:           "Admin",
					Admin:          database.AdminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"test-type"`,
						`"singular":"Update"`,
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
				Name:   "DeleteType",
				Method: http.MethodDelete,
				URL:    "/api/types/test-type",
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
					ExpectedStatus: http.StatusUnauthorized,
					ExpectedContent: []string{
						`"missing required scopes"`,
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
