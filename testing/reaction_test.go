package testing

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWebhookReactions(t *testing.T) {
	t.Parallel()

	server := NewRecordingServer()

	go http.ListenAndServe("127.0.0.1:12345", server) //nolint:gosec,errcheck

	if err := WaitForStatus("http://127.0.0.1:12345/health", http.StatusOK, 5*time.Second); err != nil {
		t.Fatal(err)
	}

	testSets := []catalystTest{
		{
			baseTest: BaseTest{
				Name:           "TriggerWebhookReaction",
				Method:         http.MethodGet,
				RequestHeaders: map[string]string{"Authorization": "Bearer 1234567890"},
				URL:            "/reaction/test",
			},
			userTests: []UserTest{
				{
					Name:            "Unauthorized",
					ExpectedStatus:  http.StatusOK,
					ExpectedContent: []string{`Hello, World!`},
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:   "TriggerWebhookReaction2",
				Method: http.MethodGet,
				URL:    "/reaction/test2",
			},
			userTests: []UserTest{
				{
					Name:            "Unauthorized",
					ExpectedStatus:  http.StatusOK,
					ExpectedContent: []string{`"test":true`},
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

func TestHookReactions(t *testing.T) {
	t.Parallel()

	server := NewRecordingServer()

	go http.ListenAndServe("127.0.0.1:12346", server) //nolint:gosec,errcheck

	if err := WaitForStatus("http://127.0.0.1:12346/health", http.StatusOK, 5*time.Second); err != nil {
		t.Fatal(err)
	}

	testSets := []catalystTest{
		{
			baseTest: BaseTest{
				Name:           "TriggerHookReaction",
				Method:         http.MethodPost,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/collections/tickets/records",
				Body: s(map[string]any{
					"name": "test",
				}),
			},
			userTests: []UserTest{
				// {
				// 	Name:            "Unauthorized",
				// 	ExpectedStatus:  http.StatusOK,
				// 	ExpectedContent: []string{`Hello, World!`},
				// },
				{
					Name:           "Analyst",
					AuthRecord:     analystEmail,
					ExpectedStatus: http.StatusOK,
					ExpectedContent: []string{
						`"collectionName":"tickets"`,
						`"name":"test"`,
					},
					ExpectedEvents: map[string]int{
						"OnModelAfterCreate":          1,
						"OnModelBeforeCreate":         1,
						"OnRecordAfterCreateRequest":  1,
						"OnRecordBeforeCreateRequest": 1,
					},
				},
				// {
				// 	Name:            "Admin",
				// 	RequestHeaders:  map[string]string{"Authorization": adminToken},
				// 	ExpectedStatus:  http.StatusOK,
				// 	ExpectedContent: []string{`Hello, World!`},
				// },
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

					require.NotEmpty(t, server.Entries)
				})
			}
		})
	}
}
