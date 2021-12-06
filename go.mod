module gitlab.com/inetmock/inetmock

go 1.17

require (
	github.com/alecthomas/participle/v2 v2.0.0-alpha7
	github.com/bwmarrin/snowflake v0.3.0
	github.com/docker/go-connections v0.4.0
	github.com/elazarl/goproxy v0.0.0-20211114080932-d06c3be7c11b
	github.com/golang/mock v1.6.0
	github.com/google/gopacket v1.1.19
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/imdario/mergo v0.3.12
	github.com/jinzhu/copier v0.3.4
	github.com/maxatome/go-testdeep v1.10.1
	github.com/miekg/dns v1.1.43
	github.com/mitchellh/mapstructure v1.4.3
	github.com/olekukonko/tablewriter v0.0.5
	github.com/prometheus/client_golang v1.11.0
	github.com/soheilhy/cmux v0.1.5
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.9.0
	github.com/testcontainers/testcontainers-go v0.12.0
	go.uber.org/multierr v1.7.0
	go.uber.org/zap v1.19.1
	golang.org/x/net v0.0.0-20211205041911-012df41ee64c
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/Microsoft/hcsshim v0.8.23 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/containerd/cgroups v1.0.1 // indirect
	github.com/containerd/containerd v1.5.5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v20.10.11+incompatible // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gorilla/mux v0.0.0-00010101000000-000000000000 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/moby/sys/mount v0.2.0 // indirect
	github.com/moby/sys/mountinfo v0.5.0 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v1.0.2 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.30.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	golang.org/x/sys v0.0.0-20211109184856-51b60fd695b3 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210909211513-a8c4777a87af // indirect
	gopkg.in/ini.v1 v1.63.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.5.8
	github.com/docker/docker => github.com/docker/docker v20.10.11+incompatible
	github.com/google/gopacket => github.com/baez90/gopacket v1.1.20-0.20211124091904-a10d380ef3cd
	github.com/gorilla/mux => github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket => github.com/gorilla/websocket v1.4.2
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20211117183948-ae814b36b871
	golang.org/x/text => golang.org/x/text v0.3.7
	gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.4.0
)
