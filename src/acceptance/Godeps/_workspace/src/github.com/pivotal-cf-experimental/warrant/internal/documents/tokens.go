package documents

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	JTI         string `json:"jti"`
}

type TokenKeyResponse struct {
	Alg   string `json:"alg"`
	Value string `json:"value"`
	Kty   string `json:"kty"`
	Use   string `json:"use"`
	N     string `json:"n"`
	E     string `json:"e"`
}
