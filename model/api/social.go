package api

type RegisterQ struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Captcha  int    `json:"captcha"`
}

type RegisterA struct {
	Message string `json:"message"`
}

type GetCaptchaQ struct {
	Email string `json:"email"`
}

type GetCaptchaA struct {
	Message string `json:"message"`
}
