package testing

import (
	"net/http"
	"testing"

	"github.com/SecurityBrewery/catalyst/app/database"
)

func TestFilesCollection(t *testing.T) {
	t.Parallel()

	testSets := []catalystTest{
		{
			baseTest: BaseTest{
				Name:   "ListFiles",
				Method: http.MethodGet,
				URL:    "/api/files?ticket=test-ticket",
			},
			userTests: []UserTest{
				{
					Name:            "Unauthorized",
					ExpectedStatus:  http.StatusUnauthorized,
					ExpectedContent: []string{`"invalid bearer token"`},
				},
				{
					Name:            "Analyst",
					AuthRecord:      database.AnalystEmail,
					ExpectedStatus:  http.StatusOK,
					ExpectedHeaders: map[string]string{"X-Total-Count": "1"},
					ExpectedEvents:  map[string]int{},
				},
				{
					Name:            "Admin",
					Admin:           database.AdminEmail,
					ExpectedStatus:  http.StatusOK,
					ExpectedHeaders: map[string]string{"X-Total-Count": "1"},
					ExpectedEvents:  map[string]int{},
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:           "CreateFile",
				Method:         http.MethodPost,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/files",
				Body: s(map[string]any{
					"ticket": "test-ticket",
					"name":   "new.txt",
					"size":   3,
					"blob":   "data:text/plain;base64,bmV3",
				}),
			},
			userTests: []UserTest{
				{
					Name:            "Unauthorized",
					ExpectedStatus:  http.StatusUnauthorized,
					ExpectedContent: []string{`"invalid bearer token"`},
				},
				{
					Name:            "Analyst",
					AuthRecord:      database.AnalystEmail,
					ExpectedStatus:  http.StatusUnauthorized,
					ExpectedContent: []string{`"missing required scopes"`},
				},
				{
					Name:            "Admin",
					Admin:           database.AdminEmail,
					ExpectedStatus:  http.StatusOK,
					ExpectedContent: []string{`"ticket":"test-ticket"`},
					ExpectedEvents: map[string]int{
						"OnRecordAfterCreateRequest":  1,
						"OnRecordBeforeCreateRequest": 1,
					},
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:   "GetFile",
				Method: http.MethodGet,
				URL:    "/api/files/b_test_file",
			},
			userTests: []UserTest{
				{
					Name:            "Unauthorized",
					ExpectedStatus:  http.StatusUnauthorized,
					ExpectedContent: []string{`"invalid bearer token"`},
				},
				{
					Name:            "Analyst",
					AuthRecord:      database.AnalystEmail,
					ExpectedStatus:  http.StatusOK,
					ExpectedContent: []string{`"id":"b_test_file"`},
					ExpectedEvents:  map[string]int{"OnRecordViewRequest": 1},
				},
				{
					Name:            "Admin",
					Admin:           database.AdminEmail,
					ExpectedStatus:  http.StatusOK,
					ExpectedContent: []string{`"id":"b_test_file"`},
					ExpectedEvents:  map[string]int{"OnRecordViewRequest": 1},
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:           "UpdateFile",
				Method:         http.MethodPatch,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/files/b_test_file",
				Body:           s(map[string]any{"name": "update.txt"}),
			},
			userTests: []UserTest{
				{
					Name:            "Unauthorized",
					ExpectedStatus:  http.StatusUnauthorized,
					ExpectedContent: []string{`"invalid bearer token"`},
				},
				{
					Name:            "Analyst",
					AuthRecord:      database.AnalystEmail,
					ExpectedStatus:  http.StatusUnauthorized,
					ExpectedContent: []string{`"missing required scopes"`},
				},
				{
					Name:            "Admin",
					Admin:           database.AdminEmail,
					ExpectedStatus:  http.StatusOK,
					ExpectedContent: []string{`"name":"update.txt"`},
					ExpectedEvents: map[string]int{
						"OnRecordAfterUpdateRequest":  1,
						"OnRecordBeforeUpdateRequest": 1,
					},
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:   "DeleteFile",
				Method: http.MethodDelete,
				URL:    "/api/files/b_test_file",
			},
			userTests: []UserTest{
				{
					Name:            "Unauthorized",
					ExpectedStatus:  http.StatusUnauthorized,
					ExpectedContent: []string{`"invalid bearer token"`},
				},
				{
					Name:            "Analyst",
					AuthRecord:      database.AnalystEmail,
					ExpectedStatus:  http.StatusUnauthorized,
					ExpectedContent: []string{`"missing required scopes"`},
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
