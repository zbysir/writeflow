module github.com/zbysir/writeflow

go 1.20

require (
	github.com/docker/libkv v0.2.1
	github.com/dop251/goja v0.0.0-20230605162241-28ee0ee714f3
	github.com/gin-gonic/gin v1.9.0
	github.com/go-git/go-billy/v5 v5.4.1
	github.com/go-git/go-git/v5 v5.7.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/samber/lo v1.38.1
	github.com/sashabaranov/go-openai v1.11.2
	github.com/spf13/cast v1.5.0
	github.com/spf13/cobra v1.7.0
	github.com/spf13/viper v1.15.0
	github.com/stretchr/testify v1.8.2
	github.com/thoas/go-funk v0.9.3
	github.com/tmc/langchaingo v0.0.0-20230522045238-97426d911826
	github.com/traefik/yaegi v0.15.1
	github.com/zbysir/gojsx v0.4.8
	github.com/zbysir/writeflow-ui v0.0.0-20230619073658-fb18c2e01d4c
	go.uber.org/zap v1.21.0
)

replace (
	// remove replace if this issue (https://github.com/traefik/yaegi/issues/1571) is fixed
	github.com/traefik/yaegi v0.15.1 => "./outpkg/yaegi"
)

require (
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20230518184743-7afd39499903 // indirect
	github.com/acomagu/bufpipe v1.0.4 // indirect
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/bytedance/sonic v1.8.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/cloudflare/circl v1.3.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dlclark/regexp2 v1.7.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/evanw/esbuild v0.14.51 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.11.2 // indirect
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/goccy/go-json v0.10.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/pprof v0.0.0-20230207041349-798e818bf904 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/imdario/mergo v0.3.15 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jolestar/go-commons-pool/v2 v2.1.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/skeema/knownhosts v1.1.1 // indirect
	github.com/spf13/afero v1.9.3 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stoewer/go-strcase v1.2.0 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/tdewolff/parse/v2 v2.6.5 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.9 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	github.com/yuin/goldmark v1.5.3 // indirect
	github.com/yuin/goldmark-meta v1.1.0 // indirect
	go.abhg.dev/goldmark/mermaid v0.4.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/arch v0.0.0-20210923205945-b76863e36670 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/exp v0.0.0-20220303212507-bbda1eaf7a17 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
