package api

type CommonA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type RegisterQ struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
}

type GetCaptchaQ struct {
	Email string `json:"email"`
}

type LoginQ struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
	ID      uint64 `json:"id"`
}

type GetProfileA struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Name    string `json:"name"`
}
