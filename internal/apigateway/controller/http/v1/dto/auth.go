package dto

type UserCreds struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5"`
}

type JwtToken struct {
	AccessToken  string `json:"accessToken" binding:"required"`
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type RenewTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type SendConfirmCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ConfirmUserRequest struct {
	Code string `json:"code" binding:"required,min=4,max=4"`
}
