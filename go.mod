module github.com/cortezaproject/corteza-server

go 1.16

// This is useful when testing changes on corteza-locale
// and you do not want to push on every change in the locale repo
// replace github.com/cortezaproject/corteza-locale => ../locale

require (
	github.com/766b/chi-prometheus v0.0.0-20211217152057-87afa9aa2ca8
	github.com/99designs/basicauth-go v0.0.0-20160802081356-2a93ba0f464d
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/Masterminds/squirrel v1.5.2
	github.com/PaesslerAG/gval v1.1.2
	github.com/SentimensRG/ctx v0.0.0-20180729130232-0bfd988c655d
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/brianvoe/gofakeit/v6 v6.12.1
	github.com/cortezaproject/corteza-locale v0.0.0-20220208175500-968c388d516c
	github.com/crewjam/saml v0.4.6
	github.com/crusttech/go-oidc v0.0.0-20180918092017-982855dad3e1
	github.com/davecgh/go-spew v1.1.1
	github.com/dgryski/dgoogauth v0.0.0-20190221195224-5a805980a5f3
	github.com/disintegration/imaging v1.6.2
	github.com/dop251/goja v0.0.0-20220110113543-261677941f3c
	github.com/edwvee/exiffix v0.0.0-20210922235313-0f6cbda5e58f
	github.com/evanw/esbuild v0.14.11
	github.com/fsnotify/fsnotify v1.5.1
	github.com/gabriel-vasile/mimetype v1.4.0
	github.com/getsentry/sentry-go v0.12.0
	github.com/go-chi/chi v3.3.4+incompatible // indirect
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/cors v1.2.0
	github.com/go-chi/httprate v0.5.2
	github.com/go-chi/jwtauth v1.2.0
	github.com/go-oauth2/oauth2/v4 v4.4.2
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang-jwt/jwt/v4 v4.2.0
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/gorilla/csrf v1.7.1
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.1
	github.com/gorilla/websocket v1.4.2
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/jmoiron/sqlx v1.3.4
	github.com/joho/godotenv v1.4.0
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0
	github.com/lestrrat-go/jwx v1.2.15
	github.com/lestrrat-go/strftime v1.0.5
	github.com/lib/pq v1.10.4
	github.com/markbates/goth v1.68.0
	github.com/mattn/go-sqlite3 v1.14.10
	github.com/microcosm-cc/bluemonday v1.0.17
	github.com/minio/minio-go/v6 v6.0.57
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/ngrok/sqlmw v0.0.0-20211220175533-9d16fdc47b31
	github.com/pkg/errors v0.9.1
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/prometheus/client_golang v1.11.0
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd // indirect
	github.com/sony/sonyflake v1.0.0
	github.com/spf13/afero v1.8.0
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.3.0
	github.com/steinfletcher/apitest v1.5.11
	github.com/steinfletcher/apitest-jsonpath v1.7.1
	github.com/stretchr/testify v1.7.0
	go.uber.org/atomic v1.9.0
	go.uber.org/zap v1.20.0
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	golang.org/x/text v0.3.7
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/mail.v2 v2.3.1
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	moul.io/zapfilter v1.6.1
	rsc.io/qr v0.2.0
)
