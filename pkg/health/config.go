package health

type Server struct {
	IP    string
	Port  uint16
	Proto string
}

type ClientsConfig struct {
	HTTP  Server
	HTTPS Server
	DNS   Server
	DoT   Server
}

type ValidationRule struct {
	Name string
	Rule string
}

type Config struct {
	Client ClientsConfig
	Rules  []ValidationRule
}
