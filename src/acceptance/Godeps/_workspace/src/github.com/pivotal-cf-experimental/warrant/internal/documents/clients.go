package documents

type ClientResponse struct {
	ClientID             string   `json:"client_id"`
	Scope                []string `json:"scope"`
	ResourceIDs          []string `json:"resource_ids"`
	Authorities          []string `json:"authorities"`
	AuthorizedGrantTypes []string `json:"authorized_grant_types"`
	AccessTokenValidity  int      `json:"access_token_validity"`
}

type CreateClientRequest struct {
	ClientID             string   `json:"client_id"`
	ClientSecret         string   `json:"client_secret"`
	Scope                []string `json:"scope"`
	ResourceIDs          []string `json:"resource_ids"`
	Authorities          []string `json:"authorities"`
	AuthorizedGrantTypes []string `json:"authorized_grant_types"`
	AccessTokenValidity  int      `json:"access_token_validity"`
}
