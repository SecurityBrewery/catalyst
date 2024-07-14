package testing

import (
	"net/http"
	"testing"
)

func TestReactions(t *testing.T) {
	baseApp, adminToken, analystToken, baseAppCleanup := BaseApp(t)
	defer baseAppCleanup()

	testSets := []authMatrixText{
		{
			baseTest: BaseTest{
				Name:           "TriggerReaction",
				Method:         http.MethodGet,
				URL:            "/reaction/test",
				TestAppFactory: AppFactory(baseApp),
			},
			authBasedExpectations: []AuthBasedExpectation{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`Hello, World!`,
					},
				},
				{
					Name:           "Analyst",
					RequestHeaders: map[string]string{"Authorization": analystToken},
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`Hello, World!`,
					},
				},
				{
					Name:           "Admin",
					RequestHeaders: map[string]string{"Authorization": adminToken},
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`Hello, World!`,
					},
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:           "TriggerReaction2",
				Method:         http.MethodGet,
				URL:            "/reaction/test2",
				TestAppFactory: AppFactory(baseApp),
			},
			authBasedExpectations: []AuthBasedExpectation{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusInternalServerError,
					ExpectedContent: []string{
						`"message":"Post \"http://localhost:8080/test\": dial tcp [::1]:8080: connect: connection refused."`,
					},
				},
				{
					Name:           "Analyst",
					RequestHeaders: map[string]string{"Authorization": analystToken},
					ExpectedStatus: http.StatusInternalServerError,
					ExpectedContent: []string{
						`"message":"Post \"http://localhost:8080/test\": dial tcp [::1]:8080: connect: connection refused."`,
					},
				},
				{
					Name:           "Admin",
					RequestHeaders: map[string]string{"Authorization": adminToken},
					ExpectedStatus: http.StatusInternalServerError,
					ExpectedContent: []string{
						`"message":"Post \"http://localhost:8080/test\": dial tcp [::1]:8080: connect: connection refused."`,
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
