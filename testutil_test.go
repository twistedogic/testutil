package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func mockResponseHandler(msg []byte) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write(msg)
	})
}

func mockJSONHandler(i any) http.Handler {
	b, _ := json.Marshal(i)
	return mockResponseHandler(b)
}

func mockStatusCodeHandler(code int) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(code)
	})
}

type mockTesting struct {
	msg string
}

func (t *mockTesting) Fatalf(msg string, args ...any) { t.msg = fmt.Sprintf(msg, args...) }
func (t *mockTesting) Helper()                        {}
func (m *mockTesting) check(t *testing.T, wantFatal bool) {
	hasError := m.msg != ""
	switch {
	case wantFatal && !hasError:
		t.Fatal("test should failed but it passed.")
	case !wantFatal && hasError:
		t.Fatalf("test should pass but got: %s", m.msg)
	}
}

func Test_MatchResponseStatusCode(t *testing.T) {
	cases := map[string]struct {
		want, got  int
		shouldFail bool
	}{
		"test should pass with status code match": {
			want: 200,
			got:  200,
		},
		"test should fail with status code mismatch": {
			want:       200,
			got:        500,
			shouldFail: true,
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			mockT := &mockTesting{}
			handler := mockStatusCodeHandler(tc.got)
			MatchResponseStatusCode(mockT, handler, &http.Request{}, tc.want)
			mockT.check(t, tc.shouldFail)
		})
	}
}

func Test_MatchResponseBody(t *testing.T) {
	cases := map[string]struct {
		want, got  []byte
		shouldFail bool
	}{
		"test should pass with status code match": {
			want: []byte("good"),
			got:  []byte("good"),
		},
		"test should fail with status code mismatch": {
			want:       []byte("good"),
			got:        []byte("bad"),
			shouldFail: true,
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			mockT := &mockTesting{}
			handler := mockResponseHandler(tc.got)
			MatchResponseBody(mockT, handler, &http.Request{}, tc.want)
			mockT.check(t, tc.shouldFail)
		})
	}
}

func Test_MatchResponseJSON(t *testing.T) {
	cases := map[string]struct {
		want, got  any
		shouldFail bool
	}{
		"test should pass with status code match": {
			want: struct {
				Key, Value string
			}{
				Key: "key", Value: "value",
			},
			got: struct {
				Key, Value string
			}{
				Key: "key", Value: "value",
			},
		},
		"test should fail with status code mismatch": {
			want: struct {
				Key, Value string
			}{
				Key: "key", Value: "value",
			},
			got: struct {
				Key, Value string
			}{
				Value: "value",
			},
			shouldFail: true,
		},
	}
	for name := range cases {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			mockT := &mockTesting{}
			handler := mockJSONHandler(tc.got)
			MatchResponseJSON(mockT, handler, &http.Request{}, tc.want)
			mockT.check(t, tc.shouldFail)
		})
	}
}
