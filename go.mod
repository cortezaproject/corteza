module github.com/cortezaproject/corteza-server

go 1.16

// This is useful when testing changes on corteza-locale
// and you do not want to push on every change in the locale repo
// replace github.com/cortezaproject/corteza-locale => ../locale

require (
	github.com/766b/chi-prometheus v0.0.0-20180509160047-46ac2b31aa30
	github.com/99designs/basicauth-go v0.0.0-20160802081356-2a93ba0f464d
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/Masterminds/squirrel v1.1.1-0.20191017225151-12f2162c8d8d
	github.com/PaesslerAG/gval v1.1.1-0.20201104175134-7847ed0c7671
	github.com/PaesslerAG/jsonpath v0.1.1 // indirect
	github.com/SentimensRG/ctx v0.0.0-20180729130232-0bfd988c655d
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/brianvoe/gofakeit/v6 v6.5.0
	github.com/cortezaproject/corteza-locale v0.0.0-20211116171437-a53f20dbdbf9
	github.com/crewjam/saml v0.4.5
	github.com/crusttech/go-oidc v0.0.0-20180918092017-982855dad3e1
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dgryski/dgoogauth v0.0.0-20190221195224-5a805980a5f3
	github.com/disintegration/imaging v1.6.0
	github.com/dop251/goja v0.0.0-20210726224656-a55e4cfac4cf
	github.com/edwvee/exiffix v0.0.0-20180602190213-b57537c92a6b
	github.com/evanw/esbuild v0.12.16
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gabriel-vasile/mimetype v1.1.2
	github.com/getsentry/sentry-go v0.1.1
	github.com/go-chi/chi v3.3.4+incompatible
	github.com/go-chi/cors v1.0.0
	github.com/go-chi/httprate v0.4.0
	github.com/go-chi/jwtauth v0.0.0-20190109153619-47840abb19b3
	github.com/go-oauth2/oauth2/v4 v4.3.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.2
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/gorilla/csrf v1.7.0
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/securecookie v1.1.1
	github.com/gorilla/sessions v1.2.1
	github.com/gorilla/websocket v1.4.2
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/jmoiron/sqlx v1.2.0
	github.com/joho/godotenv v1.3.0
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0
	github.com/lestrrat-go/strftime v1.0.3
	github.com/lib/pq v1.1.0
	github.com/markbates/goth v1.67.1
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/mattn/go-sqlite3 v1.14.0
	github.com/microcosm-cc/bluemonday v1.0.16
	github.com/minio/minio-go/v6 v6.0.39
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/ngrok/sqlmw v0.0.0-20200129213757-d5c93a81bec6
	github.com/pkg/errors v0.9.1
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/prometheus/client_golang v0.9.3
	github.com/rakyll/gotest v0.0.6 // indirect
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd // indirect
	github.com/sony/sonyflake v0.0.0-20181109022403-6d5bd6181009
	github.com/spf13/afero v1.2.2
	github.com/spf13/cast v1.3.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3 // indirect
	github.com/steinfletcher/apitest v1.3.8
	github.com/steinfletcher/apitest-jsonpath v1.3.0
	github.com/stretchr/testify v1.7.0
	github.com/tebeka/strftime v0.1.5 // indirect
	go.uber.org/atomic v1.7.0
	go.uber.org/zap v1.19.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/image v0.0.0-20190910094157-69e4b8554b2a // indirect
	golang.org/x/oauth2 v0.0.0-20210628180205-a41e5a781914
	golang.org/x/sys v0.0.0-20211117180635-dee7805ff2e1 // indirect
	golang.org/x/text v0.3.7
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/ini.v1 v1.51.0 // indirect
	gopkg.in/mail.v2 v2.3.1
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	moul.io/zapfilter v1.6.1
	rsc.io/qr v0.2.0
)
