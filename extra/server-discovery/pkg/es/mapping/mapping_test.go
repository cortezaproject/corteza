package mapping

import (
	"fmt"
	"github.com/crusttech/human/server/pkg/options"
	"os"
	"strings"
	"testing"
)

type (
	config struct {
		httpAddr string
		es       struct {
			addresses []string
		}
		humanAuth         string
		humanDiscoveryAPI string
		schemas             []*schema
	}

	schema struct {
		indexPrefix  string
		clientKey    string
		clientSecret string
	}
)

const (
	envKeyHttpAddr = "HTTP_ADDR"
)

// @todo use it properly
func TestMappings(t *testing.T) {
	//var (
	//	ctx = context.Background()
	//	log = logger.Default()
	//	req = require.New(t)
	//)
	//
	//cfg, err := GetConfig()
	//req.NoError(err)

	//client, err := Es(cfg.es.addresses)
	//req.NoError(err)
	//
	//api, err := indexer.ApiClient(cfg.humanDiscoveryAPI, cfg.humanAuth, cfg.schemas[0].clientKey, cfg.schemas[0].clientSecret)
	//req.NoError(err)
	//
	//err = Mappings(ctx, log, client, api, "private")
	//req.NoError(err)
}

// @todo reuse existing one
func GetConfig() (*config, error) {
	c := &config{}
	return c, func() error {
		baseUrl := options.EnvString("HUMAN_SERVER_BASE_URL", "http://server:80")

		c.httpAddr = options.EnvString(envKeyHttpAddr, "127.0.0.1:3201")

		c.humanAuth = options.EnvString("HUMAN_SERVER_AUTH", baseUrl+"/auth")
		if c.humanAuth == "" {
			return fmt.Errorf("Human Auth endpoint value empty, set it directly with HUMAN_SERVER_AUTH or indirectly with HUMAN_SERVER_BASE_URL")
		}

		c.humanDiscoveryAPI = options.EnvString("HUMAN_SERVER_API_DISCOVERY", baseUrl+"/api/discovery")
		if c.humanDiscoveryAPI == "" {
			return fmt.Errorf("Human Discovery API endpoint value empty, set it directly with HUMAN_SERVER_AUTH or indirectly with HUMAN_SERVER_API_DISCOVERY")
		}

		for _, ar := range []string{"public", "protected", "private"} {
			var (
				has  bool
				ucAr = strings.ToUpper(ar)
				s    = &schema{indexPrefix: ar}

				keyEnv = ucAr + "_INDEX_CLIENT_KEY"
				secEnv = ucAr + "_INDEX_CLIENT_SECRET"
			)

			if s.clientKey, has = os.LookupEnv(keyEnv); !has {
				continue
			} else if s.clientKey == "" {
				return fmt.Errorf("client key (%s) for '%s' is empty or missing", keyEnv, s.indexPrefix)
			}

			if s.clientSecret = os.Getenv(secEnv); s.clientSecret == "" {
				return fmt.Errorf("client secret (%s) for '%s' is empty or missing", secEnv, s.indexPrefix)
			}

			c.schemas = append(c.schemas, s)
		}

		if len(c.schemas) == 0 {
			return fmt.Errorf("set at least one client secret pair using <PREFIX>_INDEX_CLIENT_KEY and <PREFIX>_INDEX_CLIENT_SECRET where prefix is one of 'public', 'protected' or 'private'")
		}

		for _, a := range strings.Split(options.EnvString("ES_ADDRESS", "http://es:9200"), " ") {
			if a = strings.TrimSpace(a); a != "" {
				c.es.addresses = append(c.es.addresses, a)
			}
		}
		return nil
	}()
}
