package fakes

import (
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
)

const (
	Origin = "uaa"
	Schema = "urn:scim:schemas:core:1.0"
)

var Schemas = []string{Schema}

type ServerConfig struct {
	PublicKey string
}

type UAAServer struct {
	server    *httptest.Server
	users     *Users
	clients   *Clients
	tokenizer Tokenizer

	defaultScopes []string
	publicKey     string
}

func NewUAAServer(config ServerConfig) *UAAServer {
	router := mux.NewRouter()
	server := &UAAServer{
		server: httptest.NewUnstartedServer(router),
		defaultScopes: []string{
			"scim.read",
			"cloudcontroller.admin",
			"password.write",
			"scim.write",
			"openid",
			"cloud_controller.write",
			"cloud_controller.read",
			"doppler.firehose",
		},
		publicKey: config.PublicKey,
		tokenizer: NewTokenizer("this is the encryption key"),
		users:     NewUsers(),
		clients:   NewClients(),
	}

	router.HandleFunc("/Users", server.CreateUser).Methods("POST")
	router.HandleFunc("/Users", server.FindUsers).Methods("GET")
	router.HandleFunc("/Users/{guid}", server.GetUser).Methods("GET")
	router.HandleFunc("/Users/{guid}", server.DeleteUser).Methods("DELETE")
	router.HandleFunc("/Users/{guid}", server.UpdateUser).Methods("PUT")
	router.HandleFunc("/Users/{guid}/password", server.UpdateUserPassword).Methods("PUT")

	router.HandleFunc("/oauth/clients", server.CreateClient).Methods("POST")
	router.HandleFunc("/oauth/clients/{guid}", server.GetClient).Methods("GET")
	router.HandleFunc("/oauth/clients/{guid}", server.DeleteClient).Methods("DELETE")

	router.HandleFunc("/oauth/token", server.OAuthToken).Methods("POST")
	router.HandleFunc("/oauth/authorize", server.OAuthAuthorize).Methods("POST")

	router.HandleFunc("/token_key", server.GetTokenKey).Methods("GET")

	return server
}

func (s *UAAServer) Start() {
	s.server.Start()
}

func (s *UAAServer) Close() {
	s.server.Close()
}

func (s *UAAServer) Reset() {
	s.users.Clear()
	s.clients.Clear()
}

func (s *UAAServer) URL() string {
	return s.server.URL
}

func (s *UAAServer) SetDefaultScopes(scopes []string) {
	s.defaultScopes = scopes
}

func (s *UAAServer) ClientTokenFor(clientID string, scopes, audiences []string) string {
	return s.tokenizer.Encrypt(Token{
		ClientID:  clientID,
		Scopes:    scopes,
		Audiences: audiences,
	})
}

func (s *UAAServer) UserTokenFor(userID string, scopes, audiences []string) string {
	return s.tokenizer.Encrypt(Token{
		UserID:    userID,
		Scopes:    scopes,
		Audiences: audiences,
	})
}

func (s *UAAServer) ValidateToken(encryptedToken string, audiences, scopes []string) bool {
	token := s.tokenizer.Decrypt(encryptedToken)

	return s.tokenizer.Validate(token, Token{
		Audiences: audiences,
		Scopes:    scopes,
	})
}

func GenerateID() string {
	guid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return guid.String()
}
