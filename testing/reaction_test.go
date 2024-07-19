package testing

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWebhookReactions(t *testing.T) {
	baseApp, adminToken, analystToken, baseAppCleanup := BaseApp(t)
	defer baseAppCleanup()

	server := testWebhookServer()
	defer server.Close()

	go server.ListenAndServe() //nolint:errcheck

	testSets := []authMatrixText{
		{
			baseTest: BaseTest{
				Name:           "TriggerWebhookReaction",
				Method:         http.MethodGet,
				URL:            "/reaction/test",
				TestAppFactory: AppFactory(baseApp),
			},
			authBasedExpectations: []AuthBasedExpectation{
				{
					Name:            "Unauthorized",
					ExpectedStatus:  http.StatusOK,
					ExpectedContent: []string{`Hello, World!`},
				},
				{
					Name:            "Analyst",
					RequestHeaders:  map[string]string{"Authorization": analystToken},
					ExpectedStatus:  http.StatusOK,
					ExpectedContent: []string{`Hello, World!`},
				},
				{
					Name:            "Admin",
					RequestHeaders:  map[string]string{"Authorization": adminToken},
					ExpectedStatus:  http.StatusOK,
					ExpectedContent: []string{`Hello, World!`},
				},
			},
		},
		{
			baseTest: BaseTest{
				Name:           "TriggerWebhookReaction2",
				Method:         http.MethodGet,
				URL:            "/reaction/test2",
				TestAppFactory: AppFactory(baseApp),
			},
			authBasedExpectations: []AuthBasedExpectation{
				{
					Name:            "Unauthorized",
					ExpectedStatus:  http.StatusOK,
					ExpectedContent: []string{`"test":true`},
				},
				{
					Name:            "Analyst",
					RequestHeaders:  map[string]string{"Authorization": analystToken},
					ExpectedStatus:  http.StatusOK,
					ExpectedContent: []string{`"test":true`},
				},
				{
					Name:            "Admin",
					RequestHeaders:  map[string]string{"Authorization": adminToken},
					ExpectedStatus:  http.StatusOK,
					ExpectedContent: []string{`"test":true`},
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

func testWebhookServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"test":true}`)) //nolint:errcheck
	})

	return &http.Server{
		Addr:              "127.0.0.1:12345",
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}
}

type RecordingServer struct {
	Entries []string
}

func NewRecordingServer() *RecordingServer {
	return &RecordingServer{}
}

func (s *RecordingServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Entries = append(s.Entries, r.URL.Path)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"test":true}`)) //nolint:errcheck
}

func TestHookReactions(t *testing.T) {
	baseApp, _, analystToken, baseAppCleanup := BaseApp(t)
	defer baseAppCleanup()

	server := NewRecordingServer()

	go http.ListenAndServe("127.0.0.1:12346", server) //nolint:gosec,errcheck

	testSets := []authMatrixText{
		{
			baseTest: BaseTest{
				Name:           "TriggerHookReaction",
				Method:         http.MethodPost,
				RequestHeaders: map[string]string{"Content-Type": "application/json"},
				URL:            "/api/collections/tickets/records",
				Body: s(map[string]any{
					"name": "test",
				}),
				TestAppFactory: AppFactory(baseApp),
			},
			authBasedExpectations: []AuthBasedExpectation{
				// {
				// 	Name:            "Unauthorized",
				// 	ExpectedStatus:  http.StatusOK,
				// 	ExpectedContent: []string{`Hello, World!`},
				// },
				{
					Name:           "Analyst",
					RequestHeaders: map[string]string{"Authorization": analystToken},
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
			for _, authBasedExpectation := range testSet.authBasedExpectations {
				scenario := mergeScenario(testSet.baseTest, authBasedExpectation)
				scenario.Test(t)
			}

			require.NotEmpty(t, server.Entries)
		})
	}
}
