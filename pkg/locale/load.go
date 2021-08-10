package locale

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

const serverApplication = "corteza-server"

var defaultLanguage = language.English

func load(log *zap.Logger, pp ...string) (ll []*Language, err error) {
	var (
		pattern string
		f       *os.File
		configs []string
	)

	// @todo reads language files from all paths
	for _, p := range pp {
		pattern = filepath.Join(p, "*", "config.yaml")
		configs, err = filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("%s glob failed under: %v", pattern, err)
		}

		for _, c := range configs {
			tag := language.Make(filepath.Base(filepath.Dir(c)))

			lang := &Language{
				Tag:      tag,
				internal: make(internal),
				external: make(external),
			}

			err = func() error {
				f, err = os.Open(c)
				if err != nil {
					return fmt.Errorf("could not read %s: %v", p, err)
				}

				defer f.Close()

				if err = yaml.NewDecoder(f).Decode(&lang); err != nil {
					return fmt.Errorf("could not decode %s: %v", p, err)
				}

				return loadTranslations(lang, filepath.Dir(c))
			}()

			if err != nil {
				return nil, err
			}

			log.Info("language loaded", zap.String("tag", lang.Tag.String()), zap.String("config", c))

			ll = append(ll, lang)
		}

	}

	return
}

func loadTranslations(lang *Language, dir string) (err error) {
	lang.internal = make(internal)
	lang.external = make(external)

	var (
		// this will help us collect structured keys from YAML files
		// for external translations
		auxExternal = make(map[string]map[string]map[string]interface{})
	)

	err = filepath.Walk(dir, func(p string, finfo fs.FileInfo, err error) error {
		if err != nil || finfo.IsDir() {
			return err
		}

		var (
			f        *os.File
			relPath  = p[len(dir)+1:]
			firstSep = strings.Index(relPath, string(filepath.Separator))
		)

		if firstSep == -1 {
			return nil
		}

		// From this point on, we're assuming we are at least 1 level deep
		// and the name of the 1st directory is the name of the application
		// that the content is for
		var (
			namespace string
			keyPath   string

			ext        = path.Ext(relPath)
			appDir     = relPath[:strings.Index(relPath, string(filepath.Separator))]
			subpath    = relPath[len(appDir)+1:]
			nsDirIndex = strings.Index(subpath, string(filepath.Separator))
			nsDotIndex = strings.LastIndex(subpath, ".")

			isYAML = ext == ".yaml" || ext == ".yml"
			isJSON = ext == ".json"
		)

		if !isYAML && !isJSON {
			// unsupported format, skipping
			return nil
		}

		if appDir == serverApplication && isJSON {
			return fmt.Errorf("expecting YAML translation files for server translations")
		}

		// namespace is whatever we put on the 1st level - either filename or directory
		if nsDirIndex > 0 {
			namespace = subpath[:nsDirIndex]

			keyPath = subpath[nsDirIndex+1:]
			keyPath = keyPath[:len(keyPath)-len(ext)]
			keyPath = strings.ReplaceAll(keyPath, string(filepath.Separator), ".")
		} else if nsDotIndex > 0 {
			namespace = subpath[:nsDotIndex]
		}

		f, err = os.Open(p)
		if err != nil {
			return fmt.Errorf("could not open %s: %w", p, err)
		}

		defer f.Close()

		if appDir == serverApplication {
			if lang.internal[namespace] == nil {
				lang.internal[namespace] = make(map[string]string)
			}

			if err = procInternal(lang.internal[namespace], keyPath, f); err != nil {
				return fmt.Errorf("could not process %s: %v", relPath, err)
			}
		} else {
			if isJSON {
				if lang.external[appDir] != nil {
					return fmt.Errorf("only one JSON file per namespace supported")
				}

				var buf = &bytes.Buffer{}
				if _, err = io.Copy(buf, f); err != nil {
					return fmt.Errorf("could not copy buffer: %v", err)
				}

				lang.external[appDir] = buf
			} else {
				if auxExternal[appDir] == nil {
					auxExternal[appDir] = make(map[string]map[string]interface{})
				}
				if auxExternal[appDir][namespace] == nil {
					auxExternal[appDir][namespace] = make(map[string]interface{})
				}

				if err = procExternal(auxExternal[appDir][namespace], keyPath, f); err != nil {
					return fmt.Errorf("could not process %s: %v", relPath, err)
				}
			}
		}

		return err
	})

	if err != nil {
		return
	}

	for app, nss := range auxExternal {
		var buf = &bytes.Buffer{}
		if err = json.NewEncoder(buf).Encode(nss); err != nil {
			return
		}
		lang.external[app] = buf
	}

	return
}

// procInternal reads internal YAML translation files and converts it into
// simple key/value structure
func procInternal(ns map[string]string, prefix string, f io.Reader) (err error) {
	var (
		aux = make(map[string]interface{})

		// this allows us to do recursive calls to lambda
		flatten func(prefix string, sub map[string]interface{})
	)

	if err = yaml.NewDecoder(f).Decode(&aux); err != nil {
		return
	}

	flatten = func(prefix string, sub map[string]interface{}) {
		if len(prefix) > 0 {
			prefix += "."
		}

		for k, v := range sub {
			switch val := v.(type) {
			case string:
				ns[prefix+k] = val
			case map[string]interface{}:
				flatten(prefix+k, val)
			}
		}
	}

	flatten(prefix, aux)

	return
}

func procExternal(kv map[string]interface{}, prefix string, f io.Reader) error {
	// get to the level we want to be:
	if len(prefix) > 0 {
		for _, p := range strings.Split(prefix, ".") {
			if kv[p] == nil {
				kv[p] = make(map[string]interface{})
				kv = kv[p].(map[string]interface{})
			}
		}
	}

	return yaml.NewDecoder(f).Decode(&kv)
}
