package api

type RegisterQ struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
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

type LoginQ struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginA struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
