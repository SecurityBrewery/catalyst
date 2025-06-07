package testing

import (
	"net/http"
	"testing"

	"github.com/SecurityBrewery/catalyst/app/database"
)

func TestFeaturesConfig(t *testing.T) {
	t.Parallel()

	testSets := []catalystTest{
		{
			baseTest: BaseTest{
				Name:   "Config",
				Method: http.MethodGet,
				URL:    "/api/config",
			},
			userTests: []UserTest{
				{
					Name:           "NoAuth",
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"flags":`,
					},
				},
				{
					Name:           "Analyst",
					AuthRecord:     database.AnalystEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"flags":`,
					},
				},
				{
					Name:           "Admin",
					Admin:          database.AdminEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"flags":`,
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
