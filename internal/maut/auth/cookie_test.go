package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
)

func TestJar_deleteUserSession(t *testing.T) {
	t.Parallel()
	type fields struct {
		store sessions.Store
	}
	type args struct {
		r *http.Request
		w http.ResponseWriter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"delete user session", fields{store: sessions.NewCookieStore([]byte("secret"))}, args{r: httptest.NewRequest("GET", "/", nil), w: httptest.NewRecorder()}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			j := &Jar{store: tt.fields.store}
			j.deleteUserSession(tt.args.r, tt.args.w)
		})
	}
}

func TestJar_setStateSession(t *testing.T) {
	t.Parallel()
	type fields struct {
		store sessions.Store
	}
	type args struct {
		r     *http.Request
		w     http.ResponseWriter
		state string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"set state session", fields{store: sessions.NewCookieStore([]byte("secret"))}, args{r: httptest.NewRequest("GET", "/", nil), w: httptest.NewRecorder(), state: "state"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			j := &Jar{store: tt.fields.store}
			j.setStateSession(tt.args.r, tt.args.w, tt.args.state)
		})
	}
}

func TestJar_setUserSession(t *testing.T) {
	t.Parallel()
	type fields struct {
		store sessions.Store
	}
	type args struct {
		r      *http.Request
		w      http.ResponseWriter
		userID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"set user session", fields{store: sessions.NewCookieStore([]byte("secret"))}, args{r: httptest.NewRequest("GET", "/", nil), w: httptest.NewRecorder(), userID: "user"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			j := &Jar{store: tt.fields.store}
			j.setUserSession(tt.args.r, tt.args.w, tt.args.userID)
		})
	}
}

func TestJar_stateSession(t *testing.T) {
	t.Parallel()
	type fields struct {
		store     sessions.Store
		hasCookie bool
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantState string
		wantIsNew bool
	}{
		{"get state session", fields{store: sessions.NewCookieStore([]byte("secret")), hasCookie: true}, args{r: httptest.NewRequest("GET", "/", nil)}, "foo", false},
		{"get new state session", fields{store: sessions.NewCookieStore([]byte("secret")), hasCookie: false}, args{r: httptest.NewRequest("GET", "/", nil)}, "", true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			j := &Jar{store: tt.fields.store}
			if tt.fields.hasCookie {
				if err := addStateCookie(j.store, "foo", tt.args.r); err != nil {
					t.Error(err)
				}
			}
			gotState, gotIsNew := j.stateSession(tt.args.r)
			if gotState != tt.wantState {
				t.Errorf("stateSession() gotState = %v, want %v", gotState, tt.wantState)
			}
			if gotIsNew != tt.wantIsNew {
				t.Errorf("stateSession() gotIsNew = %v, want %v", gotIsNew, tt.wantIsNew)
			}
		})
	}
}

func TestJar_userSession(t *testing.T) {
	t.Parallel()
	type fields struct {
		store sessions.Store
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantUserID string
		wantIsNew  bool
	}{
		{"get user session", fields{store: sessions.NewCookieStore([]byte("secret"))}, args{r: httptest.NewRequest("GET", "/", nil)}, "alice", false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			j := &Jar{store: tt.fields.store}
			if err := addUserCookie(j.store, "alice", tt.args.r); err != nil {
				t.Error(err)
			}
			gotUserID, gotIsNew := j.userSession(tt.args.r)
			if gotUserID != tt.wantUserID {
				t.Errorf("userSession() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
			if gotIsNew != tt.wantIsNew {
				t.Errorf("userSession() gotIsNew = %v, want %v", gotIsNew, tt.wantIsNew)
			}
		})
	}
}
