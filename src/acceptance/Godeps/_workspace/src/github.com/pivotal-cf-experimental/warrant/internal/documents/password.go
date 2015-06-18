package documents

type SetPasswordRequest struct {
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	Password    string `json:"password"`
	OldPassword string `json:"oldPassword"`
}
