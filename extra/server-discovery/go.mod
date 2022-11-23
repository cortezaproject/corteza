module github.com/cortezaproject/corteza/extra/server-discovery

go 1.18

require (
	github.com/cortezaproject/corteza/server v0.0.0-00010101000000-000000000000
	github.com/davecgh/go-spew v1.1.1
	github.com/elastic/go-elasticsearch/v7 v7.12.0
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/cors v1.2.1
	github.com/go-chi/jwtauth v1.2.0
	github.com/go-oauth2/oauth2/v4 v4.4.3
	github.com/jmoiron/sqlx v1.3.5
	github.com/lestrrat-go/jwx v1.2.23
	github.com/microcosm-cc/bluemonday v1.0.18
	github.com/spf13/cast v1.4.1
	github.com/stretchr/testify v1.7.1
	go.uber.org/atomic v1.9.0
	go.uber.org/zap v1.21.0
)

require (
	github.com/PaesslerAG/gval v1.2.1 // indirect
	github.com/SentimensRG/ctx v0.0.0-20180729130232-0bfd988c655d // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/cortezaproject/corteza-locale v0.0.0-20221108130701-3981db126651 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.0-20210816181553-5444fa50b93d // indirect
	github.com/goccy/go-json v0.9.6 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.0 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/iter v1.0.1 // indirect
	github.com/lestrrat-go/option v1.0.0 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/spf13/cobra v1.4.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/net v0.2.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	moul.io/zapfilter v1.7.0 // indirect
)

replace github.com/cortezaproject/corteza/server => ../../server
