module github.com/cortezaproject/corteza/server

go 1.19

// This is useful when testing changes on corteza-locale
// and you do not want to push on every change in the locale repo
// replace github.com/cortezaproject/corteza-locale => ../locale

require (
	github.com/766b/chi-prometheus v0.0.0-20211217152057-87afa9aa2ca8
	github.com/99designs/basicauth-go v0.0.0-20160802081356-2a93ba0f464d
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/PaesslerAG/gval v1.2.1
	github.com/SentimensRG/ctx v0.0.0-20180729130232-0bfd988c655d
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/bep/godartsass/v2 v2.0.0
	github.com/brianvoe/gofakeit/v6 v6.16.0
	github.com/cespare/xxhash/v2 v2.2.0
	github.com/crewjam/saml v0.4.14
	github.com/crusttech/go-oidc v0.0.0-20180918092017-982855dad3e1
	github.com/davecgh/go-spew v1.1.1
	github.com/denisenkom/go-mssqldb v0.12.3
	github.com/dgryski/dgoogauth v0.0.0-20190221195224-5a805980a5f3
	github.com/disintegration/imaging v1.6.2
	github.com/dop251/goja v0.0.0-20230304130813-e2f543bf4b4c
	github.com/doug-martin/goqu/v9 v9.18.0
	github.com/edwvee/exiffix v0.0.0-20210922235313-0f6cbda5e58f
	github.com/evanw/esbuild v0.14.38
	github.com/fogleman/gg v1.3.0
	github.com/fsnotify/fsnotify v1.5.3
	github.com/gabriel-vasile/mimetype v1.4.0
	github.com/getsentry/sentry-go v0.13.0
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/cors v1.2.1
	github.com/go-chi/httprate v0.5.3
	github.com/go-chi/jwtauth v1.2.0
	github.com/go-oauth2/oauth2/v4 v4.4.3
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang-jwt/jwt/v4 v4.4.3
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.3
	github.com/google/uuid v1.6.0
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/gorilla/csrf v1.7.1
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.1
	github.com/gorilla/websocket v1.5.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/joho/godotenv v1.4.0
	github.com/lestrrat-go/jwx v1.2.28
	github.com/lestrrat-go/strftime v1.0.6
	github.com/lib/pq v1.10.5
	github.com/markbates/goth v1.71.1
	github.com/mat/besticon v3.12.0+incompatible
	github.com/mattn/go-sqlite3 v1.14.12
	github.com/microcosm-cc/bluemonday v1.0.18
	github.com/minio/minio-go/v6 v6.0.57
	github.com/modern-go/reflect2 v1.0.2
	github.com/ngrok/sqlmw v0.0.0-20211220175533-9d16fdc47b31
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.1
	github.com/sony/sonyflake v1.0.0
	github.com/spf13/afero v1.8.2
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.4.0
	github.com/steinfletcher/apitest v1.5.11
	github.com/steinfletcher/apitest-jsonpath v1.7.1
	github.com/stretchr/testify v1.8.4
	github.com/tidwall/btree v1.3.1
	github.com/valyala/fastjson v1.6.3
	go.uber.org/atomic v1.9.0
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.19.0
	golang.org/x/image v0.15.0
	golang.org/x/net v0.21.0
	golang.org/x/oauth2 v0.16.0
	golang.org/x/text v0.14.0
	google.golang.org/grpc v1.61.1
	google.golang.org/protobuf v1.32.0
	gopkg.in/mail.v2 v2.3.1
	gopkg.in/yaml.v3 v3.0.1
	moul.io/zapfilter v1.7.0
	rsc.io/qr v0.2.0
)

require (
	cloud.google.com/go/compute v1.24.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/PaesslerAG/jsonpath v0.1.1 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/beevik/etree v1.3.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cli/safeexec v1.0.1 // indirect
	github.com/crewjam/httperr v0.2.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/dlclark/regexp2 v1.8.1 // indirect
	github.com/go-chi/chi v3.3.4+incompatible // indirect
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/google/pprof v0.0.0-20230228050547-1710fef4ab10 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jonboulle/clockwork v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid v1.2.3 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/mattermost/xml-roundtrip-validator v0.1.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/minio/md5-simd v1.1.0 // indirect
	github.com/minio/sha256-simd v0.1.1 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/russellhaering/goxmldsig v1.4.0 // indirect
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/valyala/fasthttp v1.35.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/term v0.17.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240213162025-012b6fc9bca9 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/ini.v1 v1.66.2 // indirect
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
)

replace github.com/doug-martin/goqu/v9 => github.com/cortezaproject/goqu/v9 v9.18.4

replace github.com/PaesslerAG/gval => github.com/cortezaproject/gval v1.2.4
