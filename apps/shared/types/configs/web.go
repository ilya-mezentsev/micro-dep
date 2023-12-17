package configs

type Web struct {
	Port         int    `json:"port"`
	Domain       string `json:"domain"`
	SecureCookie bool   `json:"secure_cookie"`
	HttpOnly     bool   `json:"http_only"`
}
