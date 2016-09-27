package main

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"golang.org/x/net/websocket"
)

func TestBridgeServer(t *testing.T) {
	s := httptest.NewServer(websocket.Handler(BridgeServer))
	defer s.Close()
	t.Logf("URL %s", s.URL)
	u, err := url.Parse(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	u.Scheme = "ws"
	conn, err := websocket.Dial(u.String(), "", s.URL)
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("hello, world\n")
	if _, err := conn.Write(msg); err != nil {
		t.Errorf("Write: %v", err)
	}
}
