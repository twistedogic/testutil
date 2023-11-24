package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
)

type TestingT interface {
	Helper()
	Fatalf(string, ...any)
}

func result(t TestingT, handler http.Handler, req *http.Request) *http.Response {
	t.Helper()
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	return res.Result()
}

func body(t TestingT, handler http.Handler, req *http.Request) []byte {
	res := result(t, handler, req)
	b, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("fail to read response: %v", err)
	}
	return b
}

func MatchResponseStatusCode(t TestingT, handler http.Handler, req *http.Request, want int) {
	res := result(t, handler, req)
	if got := res.StatusCode; got != want {
		t.Fatalf("want: %v, got: %v", want, got)
	}
}

func MatchResponseBody(t TestingT, handler http.Handler, req *http.Request, want []byte) {
	if got := body(t, handler, req); !bytes.Equal(got, want) {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}

func MatchResponseJSON(t TestingT, handler http.Handler, req *http.Request, want any) {
	wantByte, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("fail to marshal %v: %v", want, err)
	}
	MatchResponseBody(t, handler, req, wantByte)
}

func ReadFile(t TestingT, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("fail to read %v: %v", path, err)
	}
	return b
}

func TempWriter(t TestingT, dir, pattern string) io.WriteCloser {
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		t.Fatalf("fail to create temp file %v/%v-*: %v", dir, pattern, err)
	}
	return f
}
