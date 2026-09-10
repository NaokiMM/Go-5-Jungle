package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func isolateAWS(t *testing.T) {
	t.Helper()
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if strings.HasPrefix(key, "AWS_") {
			t.Setenv(key, "")
		}
	}
	dir := t.TempDir()
	for _, name := range []string{"config", "credentials"} {
		if err := os.WriteFile(filepath.Join(dir, name), nil, 0600); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("AWS_CONFIG_FILE", filepath.Join(dir, "config"))
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", filepath.Join(dir, "credentials"))
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
}

func TestIdentity(t *testing.T) {
	isolateAWS(t)
	t.Setenv("AWS_ACCESS_KEY_ID", "test-key")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test-secret")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		if r.Method != "POST" || r.Form.Get("Action") != "GetCallerIdentity" {
			t.Errorf("unexpected request: %s %v", r.Method, r.Form)
		}
		if !strings.Contains(r.Header.Get("Authorization"), "Credential=test-key/") {
			t.Error("request was not signed with test credentials")
		}
		w.Header().Set("Content-Type", "text/xml")
		io.WriteString(w, `<GetCallerIdentityResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/"><GetCallerIdentityResult><Account>123456789012</Account><Arn>arn:aws:sts::123456789012:assumed-role/Developer/test</Arn><UserId>AROAEXAMPLE:test</UserId></GetCallerIdentityResult><ResponseMetadata><RequestId>test</RequestId></ResponseMetadata></GetCallerIdentityResponse>`)
	}))
	defer server.Close()
	t.Setenv("AWS_ENDPOINT_URL_STS", server.URL)
	var out bytes.Buffer
	if err := run(context.Background(), []string{"-region", "ap-northeast-1"}, &out, io.Discard); err != nil {
		t.Fatal(err)
	}
	want := "Account: 123456789012\nARN:     arn:aws:sts::123456789012:assumed-role/Developer/test\nUser ID: AROAEXAMPLE:test\n"
	if out.String() != want {
		t.Errorf("output = %q, want %q", out.String(), want)
	}
}

func TestErrors(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"missing region", nil, "リージョンが未設定"},
		{"missing credentials", []string{"-region", "ap-northeast-1"}, "認証情報の設定"},
		{"invalid timeout", []string{"-timeout", "0s"}, "0より大きい"},
		{"unexpected argument", []string{"extra"}, "位置引数"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			isolateAWS(t)
			var out bytes.Buffer
			err := run(context.Background(), tc.args, &out, io.Discard)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want %q", err, tc.want)
			}
			if out.Len() != 0 {
				t.Error("failure printed success output")
			}
		})
	}
}

func TestHelp(t *testing.T) {
	var out bytes.Buffer
	if err := run(context.Background(), []string{"-help"}, io.Discard, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "-profile") {
		t.Fatal("help missing profile option")
	}
}
