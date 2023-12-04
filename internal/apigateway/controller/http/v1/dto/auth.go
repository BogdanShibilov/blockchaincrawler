package dto

type UserCreds struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JwtToken struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RenewTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type SendConfirmCodeRequest struct {
	Email string `json:"email"`
}

type ConfirmUserRequest struct {
	Code string `json:"code"`
}
