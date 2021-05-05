package health

type Server struct {
	IP   string
	Port uint16
}

type HTTPClientConfig struct {
	HTTP  Server
	HTTPS Server
}

type ValidationRule struct {
	Name string
	Rule string
}

type Config struct {
	Client HTTPClientConfig
	Rules  []ValidationRule
}
