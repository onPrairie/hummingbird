module Hummingbird

go 1.17

require (
	github.com/axgle/mahonia v0.0.0-20180208002826-3358181d7394
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/robertkrimen/otto v0.0.0-20210614181706-373ff5438452
	github.com/sirupsen/logrus v1.8.1
)

require (
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.12.0 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	github.com/subosito/gotenv v1.3.0 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0 // indirect
)

require (
	github.com/jmoiron/sqlx v1.3.4
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/lestrrat-go/strftime v1.0.5 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
)

require utils v0.0.0

replace utils => ./utils
