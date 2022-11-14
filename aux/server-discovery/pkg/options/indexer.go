package options

import (
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"os"
	"strings"
)

type (
	IndexerOpt struct {
		Enabled bool
		//HttpAddr             string
		CortezaServerBaseUrl string
		CortezaServerAuthUrl string
		CortezaDiscoveryAPI  string

		Schemas []*schema
	}

	schema struct {
		IndexPrefix  string
		ClientKey    string
		ClientSecret string
	}
)

const (
	discoveryIndexer     = "DISCOVERY_INDEXER_"
	indexerEnvKeyEnabled = discoveryIndexer + "ENABLED"
	//indexerEnvKeyHttpAddr = discoveryIndexer + "HTTP_ADDR"
)

func Indexer() (o *IndexerOpt, err error) {
	o = &IndexerOpt{}

	return o, func() error {
		o.Enabled = options.EnvBool(indexerEnvKeyEnabled, true)

		//baseUrl := options.EnvString(envKeyBaseUrl, "http://server:80")

		//o.HttpAddr = options.EnvString(indexerEnvKeyHttpAddr, "0.0.0.0:80")

		//o.CortezaServerAuthUrl = options.EnvString(envKeyAuthUrl, baseUrl+"/auth")
		//if o.CortezaServerAuthUrl == "" {
		//	return fmt.Errorf("corteza Auth endpoint value empty, set it directly with %s or indirectly with %s", envKeyAuthUrl, envKeyBaseUrl)
		//}
		//
		//o.CortezaDiscoveryAPI = options.EnvString(envKeyDiscoveryUrl, baseUrl+"/api/discovery")
		//if o.CortezaDiscoveryAPI == "" {
		//	return fmt.Errorf("corteza Discovery API endpoint value empty, set it directly with %s or indirectly with %s", envKeyDiscoveryUrl, envKeyBaseUrl)
		//}

		//o.IndexInterval = options.EnvInt(indexerEnvKeyIndexInterval, 30)

		for _, ar := range []string{"public", "protected", "private"} {
			var (
				has  bool
				ucAr = strings.ToUpper(ar)
				s    = &schema{IndexPrefix: ar}

				keyEnv = discoveryIndexer + ucAr + "_INDEX_CLIENT_KEY"
				secEnv = discoveryIndexer + ucAr + "_INDEX_CLIENT_SECRET"
			)

			if s.ClientKey, has = os.LookupEnv(keyEnv); !has {
				continue
			} else if s.ClientKey == "" {
				return fmt.Errorf("client key (%s) for '%s' is empty or missing", keyEnv, s.IndexPrefix)
			}

			if s.ClientSecret = os.Getenv(secEnv); s.ClientSecret == "" {
				return fmt.Errorf("client secret (%s) for '%s' is empty or missing", secEnv, s.IndexPrefix)
			}

			o.Schemas = append(o.Schemas, s)
		}

		if len(o.Schemas) == 0 {
			return fmt.Errorf("set at least one client secret pair using %s<PREFIX>_INDEX_CLIENT_KEY and <PREFIX>_INDEX_CLIENT_SECRET where prefix is one of 'public', 'protected' or 'private'", discoveryIndexer)
		}

		//for _, a := range strings.Split(options.EnvString(indexerEnvKeyEsAddr, "http://es:9200"), " ") {
		//	if a = strings.TrimSpace(a); a != "" {
		//		o.Es.Addresses = append(o.Es.Addresses, a)
		//	}
		//}

		return nil
	}()
}
