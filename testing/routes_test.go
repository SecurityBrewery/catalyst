package testing

import (
	"net/http"
	"testing"
)

func Test_Routes(t *testing.T) {
	baseApp, adminToken, analystToken, baseAppCleanup := BaseApp(t)
	defer baseAppCleanup()

	testSets := []authMatrixText{
		{
			baseTest: BaseTest{
				Name:           "Root",
				Method:         http.MethodGet,
				URL:            "/",
				TestAppFactory: AppFactory(baseApp),
			},
			authBasedExpectations: []AuthBasedExpectation{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusFound,
				},
				{
					Name:           "Analyst",
					RequestHeaders: map[string]string{"Authorization": analystToken},
					ExpectedStatus: http.StatusFound,
				},
				{
					Name:           "Admin",
					RequestHeaders: map[string]string{"Authorization": adminToken},
					ExpectedStatus: http.StatusFound,
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:           "Config",
				Method:         http.MethodGet,
				URL:            "/api/config",
				TestAppFactory: AppFactory(baseApp),
			},
			authBasedExpectations: []AuthBasedExpectation{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"flags":null`,
					},
				},
				{
					Name:           "Analyst",
					RequestHeaders: map[string]string{"Authorization": analystToken},
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"flags":null`,
					},
				},
				{
					Name:           "Admin",
					RequestHeaders: map[string]string{"Authorization": adminToken},
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"flags":null`,
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
