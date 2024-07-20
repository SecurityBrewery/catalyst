package testing

import (
	"net/http"
	"testing"
)

func Test_Routes(t *testing.T) {
	t.Parallel()

	testSets := []catalystTest{
		{
			baseTest: BaseTest{
				Name:   "Root",
				Method: http.MethodGet,
				URL:    "/",
			},
			userTests: []UserTest{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusFound,
				},
				{
					Name:           "Analyst",
					AuthRecord:     analystEmail,
					ExpectedStatus: http.StatusFound,
				},
				{
					Name:           "Admin",
					Admin:          adminEmail,
					ExpectedStatus: http.StatusFound,
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:   "Config",
				Method: http.MethodGet,
				URL:    "/api/config",
			},
			userTests: []UserTest{
				{
					Name:           "Unauthorized",
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"flags":null`,
					},
				},
				{
					Name:           "Analyst",
					AuthRecord:     analystEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"flags":null`,
					},
				},
				{
					Name:           "Admin",
					Admin:          adminEmail,
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
			t.Parallel()

			testSet.run(t)
		})
	}
}
