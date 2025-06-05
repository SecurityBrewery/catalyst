package testing

import (
	"context"
	"net/http"
	"testing"

	"github.com/SecurityBrewery/catalyst/app"
	"github.com/SecurityBrewery/catalyst/app/database"
	"github.com/SecurityBrewery/catalyst/app/database/sqlc"
)

func TestCommentsCollection(t *testing.T) {
	t.Parallel()

	setup := func(t *testing.T, app *app.App) {
		_, err := app.Queries.CreateComment(context.Background(), sqlc.CreateCommentParams{
			ID:      "c_test_comment",
			Author:  "u_bob_analyst",
			Message: "comment",
			Ticket:  "test-ticket",
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	testSets := []catalystTest{
		{
			baseTest: BaseTest{
				Name:   "ListComments",
				Method: http.MethodGet,
				URL:    "/api/comments?ticket=test-ticket",
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
				Name:           "CreateComment",
				Method:         http.MethodPost,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/comments",
				Body: s(map[string]any{
					"author":  "u_bob_analyst",
					"message": "new",
					"ticket":  "test-ticket",
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
				Name:   "GetComment",
				Method: http.MethodGet,
				URL:    "/api/comments/c_test_comment",
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
						`"id":"c_test_comment"`,
					},
					ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
				},
				{
					Name:           "Admin",
					Admin:          database.AdminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"id":"c_test_comment"`,
					},
					ExpectedEvents: map[string]int{"OnRecordViewRequest": 1},
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:           "UpdateComment",
				Method:         http.MethodPatch,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/comments/c_test_comment",
				Body:           s(map[string]any{"message": "update"}),
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
						`"id":"c_test_comment"`,
						`"message":"update"`,
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
						`"id":"c_test_comment"`,
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
			baseTest: BaseTest{
				Name:   "DeleteComment",
				Method: http.MethodDelete,
				URL:    "/api/comments/c_test_comment",
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

					runMatrixTestWithSetup(t, testSet.baseTest, userTest, setup)
				})
			}
		})
	}
}
