package example

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
)

/*
type throttleTripper struct {
	internal    http.RoundTripper
	ratelimiter *rate.Limiter
}

func (t *throttleTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	err := t.ratelimiter.Wait(request.Context())
	if err != nil {
		return nil, err
	}

	resp, err := t.internal.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type loggingTripper struct {
	internal http.RoundTripper
}

func (l *loggingTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	log.Println("request", request.Host, request.URL, request.Cookies())
	resp, err := l.internal.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	u, _ := resp.Location()
	if resp.Body == nil {
		log.Println("response", resp.StatusCode, u, resp.Cookies())
	} else {
		b, _ := io.ReadAll(resp.Body)
		log.Println("response", resp.StatusCode, u, resp.Cookies(), string(b))
		resp.Body = io.NopCloser(bytes.NewReader(b))
	}

	return resp, nil
}
*/

type unsecureJar struct {
	internal http.CookieJar
}

func (j *unsecureJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	unsecureCookies := make([]*http.Cookie, len(cookies))
	for i, cookie := range cookies {
		cookie.Secure = false
		unsecureCookies[i] = cookie
	}

	j.internal.SetCookies(u, cookies)
}

func (j *unsecureJar) Cookies(u *url.URL) []*http.Cookie {
	return j.internal.Cookies(u)
}

func testClient(t *testing.T) (*http.Client, error) {
	t.Helper()

	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return nil, err
	}

	// rl := rate.NewLimiter(rate.Every(time.Second), 5)
	return &http.Client{
		Jar: &unsecureJar{jar},
		// Transport: &throttleTripper{http.DefaultTransport, rl},
	}, nil
}
