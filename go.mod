module my.com/secrets

go 1.22

toolchain go1.22.0

replace github.com/gorilla/sessions => github.com/ory/sessions v1.2.2-0.20220110165800-b09c17334dc2

require (
	github.com/Conight/go-googletrans v0.0.0-20200929083318-176776d061cb
	github.com/Masterminds/sprig/v3 v3.2.3
	github.com/Masterminds/squirrel v1.5.2
	github.com/arbovm/levenshtein v0.0.0-20160628152529-48b4e1c0c4d0
	github.com/avast/retry-go/v3 v3.1.1
	github.com/bradleyjkemp/cupaloy/v2 v2.8.0
	github.com/bwmarrin/discordgo v0.28.1
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/cockroachdb/cockroach-go/v2 v2.3.5
	github.com/coreos/go-oidc/v3 v3.10.0
	github.com/davecgh/go-spew v1.1.1
	github.com/dgraph-io/ristretto v0.1.1
	github.com/docker/go-connections v0.4.0
	github.com/fatih/color v1.13.0
	github.com/ghodss/yaml v1.0.0
	github.com/gin-gonic/gin v1.7.7
	github.com/go-crypt/crypt v0.2.21
	github.com/go-faker/faker/v4 v4.4.1
	github.com/go-openapi/strfmt v0.23.0
	github.com/go-playground/validator/v10 v10.9.0
	github.com/go-webauthn/webauthn v0.10.2
	github.com/gobuffalo/httptest v1.5.2
	github.com/gobuffalo/pop/v6 v6.1.2-0.20230318123913-c85387acc9a0
	github.com/gofrs/uuid v4.3.1+incompatible
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/golang-migrate/migrate/v4 v4.15.1
	github.com/golang/gddo v0.0.0-20210115222349-20d68f94ee1f
	github.com/golang/mock v1.6.0
	github.com/google/go-github/v27 v27.0.6
	github.com/google/go-github/v38 v38.1.0
	github.com/google/go-jsonnet v0.20.0
	github.com/google/uuid v1.6.0
	github.com/google/wire v0.5.0
	github.com/gorilla/sessions v1.2.1
	github.com/gtank/cryptopasta v0.0.0-20170601214702-1f550f6f2f69
	github.com/hashicorp/go-retryablehttp v0.7.5
	github.com/hashicorp/golang-lru v0.5.4
	github.com/ilyakaznacheev/cleanenv v1.2.6
	github.com/imdario/mergo v0.3.13
	github.com/inhies/go-bytesize v0.0.0-20220417184213-4913239db9cf
	github.com/jackc/pgx/v4 v4.18.2
	github.com/jarcoal/httpmock v1.3.1
	github.com/jmoiron/sqlx v1.3.5
	github.com/julienschmidt/httprouter v1.3.0
	github.com/knadh/koanf/parsers/json v0.1.0
	github.com/laher/mergefs v0.1.2-0.20230223191438-d16611b2f4e7
	github.com/lestrrat-go/jwx/v2 v2.0.21
	github.com/lib/pq v1.10.9
	github.com/luna-duclos/instrumentedsql v1.1.3
	github.com/mattes/migrate v3.0.1+incompatible
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe
	github.com/ory/analytics-go/v5 v5.0.1
	github.com/ory/client-go v1.9.0
	github.com/ory/dockertest/v3 v3.10.0
	github.com/ory/graceful v0.1.4-0.20230301144740-e222150c51d0
	github.com/ory/herodot v0.10.3-0.20230626083119-d7e5192f0d88
	github.com/ory/hydra-client-go/v2 v2.2.0-rc.3.0.20240202131107-1c7b57df3bb0
	github.com/ory/jsonschema/v3 v3.0.8
	github.com/ory/mail/v3 v3.0.0
	github.com/ory/nosurf v1.2.7
	github.com/ory/x v0.0.628
	github.com/peterhellberg/link v1.2.0
	github.com/phayes/freeport v0.0.0-20220201140144-74d24b5ae9f5
	github.com/pkg/errors v0.9.1
	github.com/pquerna/otp v1.4.0
	github.com/prometheus/client_golang v1.13.0
	github.com/rakutentech/jwk-go v1.1.3
	github.com/rs/cors v1.11.0
	github.com/rs/zerolog v1.26.1
	github.com/samber/lo v1.39.0
	github.com/sirupsen/logrus v1.9.0
	github.com/slack-go/slack v0.12.5
	github.com/spf13/cobra v1.7.0
	github.com/spf13/pflag v1.0.5
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.9.0
	github.com/swaggo/swag v1.7.6
	github.com/testcontainers/testcontainers-go v0.20.0
	github.com/testcontainers/testcontainers-go/modules/postgres v0.20.0
	github.com/tidwall/gjson v1.17.1
	github.com/tidwall/sjson v1.2.5
	github.com/urfave/negroni v1.0.0
	github.com/zmb3/spotify/v2 v2.4.2
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.51.0
	go.opentelemetry.io/otel v1.26.0
	go.opentelemetry.io/otel/sdk v1.21.0
	go.opentelemetry.io/otel/trace v1.26.0
	golang.org/x/crypto v0.22.0
	golang.org/x/exp v0.0.0-20240416160154-fe59bbe5cc7f
	golang.org/x/net v0.24.0
	golang.org/x/oauth2 v0.18.0
	golang.org/x/sync v0.7.0
	golang.org/x/text v0.14.0
	google.golang.org/grpc v1.59.0
)

require (
	code.dny.dev/ssrf v0.2.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/BurntSushi/toml v0.4.1 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.2.0 // indirect
	github.com/Microsoft/go-winio v0.6.0 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/avast/retry-go/v4 v4.3.0 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/containerd/containerd v1.6.19 // indirect
	github.com/containerd/continuity v0.3.0 // indirect
	github.com/cpuguy83/dockercfg v0.3.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/docker/cli v20.10.21+incompatible // indirect
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/docker/docker v23.0.5+incompatible // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/felixge/fgprof v0.9.3 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/fxamacker/cbor/v2 v2.6.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-crypt/x v0.2.14 // indirect
	github.com/go-jose/go-jose/v3 v3.0.3 // indirect
	github.com/go-jose/go-jose/v4 v4.0.1 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/errors v0.22.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/spec v0.20.9 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/go-webauthn/x v0.1.9 // indirect
	github.com/gobuffalo/envy v1.10.2 // indirect
	github.com/gobuffalo/fizz v1.14.4 // indirect
	github.com/gobuffalo/flect v1.0.0 // indirect
	github.com/gobuffalo/github_flavored_markdown v1.1.3 // indirect
	github.com/gobuffalo/helpers v0.6.7 // indirect
	github.com/gobuffalo/nulls v0.4.2 // indirect
	github.com/gobuffalo/plush/v4 v4.1.18 // indirect
	github.com/gobuffalo/tags/v3 v3.1.4 // indirect
	github.com/gobuffalo/validate/v3 v3.3.3 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/goccy/go-yaml v1.9.6 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.1.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/google/go-tpm v0.9.0 // indirect
	github.com/google/pprof v0.0.0-20221010195024-131d412537ea // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.18.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/huandu/xstrings v1.3.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/jackc/puddle/v2 v2.1.2 // indirect
	github.com/jandelgado/gcov2lcov v1.0.5 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/klauspost/compress v1.16.6 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/knadh/koanf/parsers/toml v0.1.0 // indirect
	github.com/knadh/koanf/parsers/yaml v0.1.0 // indirect
	github.com/knadh/koanf/providers/posflag v0.1.0 // indirect
	github.com/knadh/koanf/v2 v2.0.1 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.5 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/jwx v1.2.29 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/microcosm-cc/bluemonday v1.0.21 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/patternmatcher v0.5.0 // indirect
	github.com/moby/sys/sequential v0.5.0 // indirect
	github.com/moby/term v0.0.0-20221128092401-c43b287e0e0f // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/nyaruka/phonenumbers v1.1.6 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc2 // indirect
	github.com/opencontainers/runc v1.1.12 // indirect
	github.com/openzipkin/zipkin-go v0.4.2 // indirect
	github.com/ory/go-acc v0.2.9-0.20230103102148-6b1c9a70dbbe // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/pkg/profile v1.7.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/seatgeek/logrus-gelf-formatter v0.0.0-20210414080842-5b05eb8ff761 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/segmentio/backo-go v1.0.1 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/sourcegraph/annotate v0.0.0-20160123013949-f4cad6c6324d // indirect
	github.com/sourcegraph/syntaxhighlight v0.0.0-20170531221838-bd320f5d308e // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.16.0 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/ugorji/go/codec v1.2.6 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/xtgo/uuid v0.0.0-20140804021211-a0b114877d4c // indirect
	go.mongodb.org/mongo-driver v1.14.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.47.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.21.0 // indirect
	go.opentelemetry.io/contrib/propagators/jaeger v1.21.1 // indirect
	go.opentelemetry.io/contrib/samplers/jaegerremote v0.15.1 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/zipkin v1.21.0 // indirect
	go.opentelemetry.io/otel/metric v1.26.0 // indirect
	go.opentelemetry.io/proto/otlp v1.0.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/tools v0.20.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto v0.0.0-20231106174013-bbf56f31fb17 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20231106174013-bbf56f31fb17 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231106174013-bbf56f31fb17 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
