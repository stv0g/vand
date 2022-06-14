module github.com/stv0g/vand

go 1.18

require (
	gioui.org v0.0.0-20220517073035-79f037f983a7
	github.com/adrianmo/go-nmea v1.7.0
	github.com/dgraph-io/badger/v3 v3.2103.2
	github.com/eclipse/paho.mqtt.golang v1.3.5
	github.com/gin-contrib/static v0.0.1
	github.com/gin-gonic/gin v1.7.7
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/websocket v1.5.0
	github.com/imdario/mergo v0.3.13
	github.com/mitchellh/mapstructure v1.5.0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/rzetterberg/elmobd v0.0.0-20200309135549-334e700512dd
	github.com/simonvetter/modbus v1.5.0
	github.com/spf13/cobra v1.4.0
	github.com/spf13/viper v1.11.0
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07
	github.com/tdewolff/canvas v0.0.0-00010101000000-000000000000
	google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v2 v2.4.0
	periph.io/x/conn/v3 v3.6.10
)

require (
	gioui.org/cpu v0.0.0-20220412190645-f1e9e8c3b1f7 // indirect
	gioui.org/shader v1.0.6 // indirect
	github.com/ByteArena/poly2tri-go v0.0.0-20170716161910-d102ad91854f // indirect
	github.com/adrg/strutil v0.3.0 // indirect
	github.com/adrg/sysfont v0.1.2 // indirect
	github.com/adrg/xdg v0.4.0 // indirect
	github.com/benoitkugler/textlayout v0.1.3-0.20220520134940-234370ccc6fe // indirect
	github.com/benoitkugler/textprocessing v0.0.2 // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/dgraph-io/ristretto v0.1.0 // indirect
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gioui/uax v0.2.1-0.20220325163150-e3d987515a12 // indirect
	github.com/go-latex/latex v0.0.0-20210823091927-c0d11ff05a81 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.11.0 // indirect
	github.com/go-text/typesetting v0.0.0-20220411150340-35994bc27a7b // indirect
	github.com/goburrow/serial v0.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/google/flatbuffers v1.12.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.12.3 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/tdewolff/minify/v2 v2.11.2 // indirect
	github.com/tdewolff/parse/v2 v2.5.29 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	golang.org/x/crypto v0.0.0-20220518034528-6f7dac969898 // indirect
	golang.org/x/exp/shiny v0.0.0-20220518171630-0b5c67f07fdf // indirect
	golang.org/x/image v0.0.0-20220413100746-70e8d0d3baa9 // indirect
	golang.org/x/net v0.0.0-20220520000938-2e3eb7b945c2 // indirect
	golang.org/x/sys v0.0.0-20220519141025-dcacdad47464 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/yaml.v3 v3.0.0 // indirect
)

replace github.com/dgraph-io/ristretto => github.com/dgraph-io/ristretto v0.1.1-0.20220403145359-8e850b710d6d

replace github.com/tdewolff/canvas => github.com/stv0g/canvas v0.0.0-20220520160803-c53607df80eb

replace github.com/rzetterberg/elmobd => github.com/stv0g/elmobd v0.0.0-20220614123457-e89412c28604
