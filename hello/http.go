package main

import (
	"encoding/json"
	"fmt"
	"github.com/o1egl/paseto/v2"
	"net"
	"net/http"
	"time"
)

type Server struct {
	*http.ServeMux
}

func NewHTTPServer() *Server {
	server := &Server{
		ServeMux: http.NewServeMux(),
	}
	fs := http.FileServer(http.Dir("./docs/swagger"))
	server.Handle("/token", http.HandlerFunc(server.GetTokenHandler))
	server.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
	return server
}

func (s *Server) GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	symmetricKey := []byte("YELLOW SUBMARINE, BLACK WIZARDRY") // Must be 32 bytes
	now := time.Now()
	exp := now.Add(24 * time.Hour)
	nbt := now

	jsonToken := paseto.JSONToken{
		Audience:   "test",
		Issuer:     "test_service",
		Jti:        "123",
		Subject:    "test_subject",
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  nbt,
	}
	jsonToken.Set("data", "this is a signed message")
	footer := "some footer"

	// Encrypt data
	token, err := paseto.Encrypt(symmetricKey, jsonToken, footer)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"error": err,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"access_token": token,
	})
	return
}

func (s *Server) Serve(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("cannot create http network listener:%w", err)
	}
	if err = http.Serve(listener, s); err != nil {
		return fmt.Errorf("cannot start HTTP server: %w", err)
	}
	return nil
}
