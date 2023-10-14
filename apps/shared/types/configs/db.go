package configs

type DB struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	User       string `json:"user"`
	Password   string `json:"password"`
	DBName     string `json:"db_name"`
	Connection struct {
		RetryCount   int `json:"retry_count"`
		RetryTimeout int `json:"retry_timeout"`
	} `json:"connection"`
}
