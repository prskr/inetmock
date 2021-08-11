module gitlab.com/inetmock/inetmock

go 1.16

require (
	github.com/alecthomas/participle/v2 v2.0.0-alpha6
	github.com/bwmarrin/snowflake v0.3.0
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/go-connections v0.4.0
	github.com/elazarl/goproxy v0.0.0-20210110162100-a92cc753f88e
	github.com/golang/mock v1.6.0
	github.com/google/gopacket v1.1.19
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v0.0.0-00010101000000-000000000000 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/imdario/mergo v0.3.12
	github.com/jinzhu/copier v0.3.2
	github.com/maxatome/go-testdeep v1.9.2
	github.com/miekg/dns v1.1.43
	github.com/mitchellh/mapstructure v1.4.1
	github.com/olekukonko/tablewriter v0.0.5
	github.com/prometheus/client_golang v1.11.0
	github.com/soheilhy/cmux v0.1.5
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	github.com/testcontainers/testcontainers-go v0.11.1
	go.uber.org/multierr v1.7.0
	go.uber.org/zap v1.18.1
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.5.2
	github.com/docker/docker => github.com/docker/docker v20.10.7+incompatible
	github.com/google/gopacket => github.com/baez90/gopacket v1.1.20-0.20210811071216-a2eed1ae149e
	github.com/gorilla/mux => github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
	golang.org/x/text => golang.org/x/text v0.3.6
)
