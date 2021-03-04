module gitlab.com/inetmock/inetmock

go 1.16

require (
	github.com/bwmarrin/snowflake v0.3.0
	github.com/docker/go-connections v0.4.0
	github.com/golang/mock v1.5.0
	github.com/golang/protobuf v1.4.3
	github.com/google/gopacket v1.1.19
	github.com/google/uuid v1.2.0
	github.com/imdario/mergo v0.3.11
	github.com/jinzhu/copier v0.2.5
	github.com/miekg/dns v1.1.40
	github.com/mitchellh/mapstructure v1.4.1
	github.com/olekukonko/tablewriter v0.0.5
	github.com/prometheus/client_golang v1.9.0
	github.com/soheilhy/cmux v0.1.4
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/testcontainers/testcontainers-go v0.9.0
	go.uber.org/multierr v1.6.0
	go.uber.org/zap v1.16.0
	google.golang.org/grpc v1.36.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/elazarl/goproxy.v1 v1.0.0-20180725130230-947c36da3153
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/docker/docker => github.com/docker/docker v17.12.0-ce-rc1.0.20200916142827-bd33bbf0497b+incompatible
