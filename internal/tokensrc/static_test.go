package tokensrc

import (
	"context"
	"testing"
	"time"
)

func TestFromEnvOrSession(t *testing.T) {
	ts := FromEnvOrSession(func(string) string { return "dev:u1" }, context.Background(), nil, nil, time.Now)
	tok, err := ts.Token()
	if err != nil || tok.AccessToken != "dev:u1" {
		t.Fatalf("static token: %v %v", tok, err)
	}
	if _, ok := FromEnvOrSession(func(string) string { return "" }, context.Background(), nil, nil, nil).(*Source); !ok {
		t.Fatal("expected the session-backed source")
	}
}
