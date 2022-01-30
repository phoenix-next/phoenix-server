package api

type NormalResponse struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

type RegisterQ struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
}

type RegisterA struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type CaptchaValidQ struct {
	Email string `json:"email"`
}
