package locale

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

const serverApplication = "corteza-server"

var (
	defaultLanguage = language.English
)

func loadConfigs(fsys fs.FS) (ll []*Language, err error) {
	var (
		sub     fs.FS
		lang    *Language
		configs = make([]string, 0)
	)

	configs, err = fs.Glob(fsys, "*/config.yaml")
	if err != nil {
		return nil, err
	}

	for _, config := range configs {
		if sub, err = fs.Sub(fsys, filepath.Dir(config)); err != nil {
			return nil, err
		}

		if lang, err = loadConfig(sub, filepath.Base(config)); err != nil {
			return nil, err
		}

		lang.Tag = language.Make(filepath.Base(filepath.Dir(config)))
		lang.src = filepath.Dir(config)

		ll = append(ll, lang)
	}

	return ll, nil
}

func loadConfig(fsys fs.FS, p string) (lang *Language, err error) {
	var f fs.File
	lang = &Language{
		Tag: language.Make(filepath.Base(filepath.Dir(p))),
		fs:  fsys,
	}

	f, err = fsys.Open(p)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %v", p, err)
	}

	defer f.Close()
	if err = yaml.NewDecoder(f).Decode(&lang); err != nil {
		return nil, fmt.Errorf("could not decode %s: %v", p, err)
	}

	return
}

func loadTranslations(lang *Language) (err error) {
	lang.internal = make(internal)
	lang.external = make(external)

	var (
		// this will help us collect structured keys from YAML files
		// for external translations
		auxExternal = make(map[string]map[string]map[string]interface{})
	)

	err = fs.WalkDir(lang.fs, ".", func(p string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}

		var (
			f        fs.File
			firstSep = strings.Index(p, string(filepath.Separator))
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

			ext        = path.Ext(p)
			appDir     = p[:firstSep]
			subpath    = p[len(appDir)+1:]
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

		if f, err = lang.fs.Open(p); err != nil {
			return fmt.Errorf("could not open %s: %w", p, err)
		}

		defer f.Close()

		if appDir == serverApplication {
			if lang.internal[namespace] == nil {
				lang.internal[namespace] = make(map[string]string)
			}

			if err = procInternal(lang.internal[namespace], keyPath, f); err != nil {
				return fmt.Errorf("could not process %s: %v", p, err)
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

				lang.external[appDir] = bytes.NewReader(buf.Bytes())
			} else {
				if auxExternal[appDir] == nil {
					auxExternal[appDir] = make(map[string]map[string]interface{})
				}
				if auxExternal[appDir][namespace] == nil {
					auxExternal[appDir][namespace] = make(map[string]interface{})
				}

				if err = procExternal(auxExternal[appDir][namespace], keyPath, f); err != nil {
					return fmt.Errorf("could not process %s: %v", p, err)
				}
			}
		}

		return err
	})

	if err != nil {
		return
	}

	// encode all external (for webapps, clientS) translations
	// into buffers that will be copied directly into responses
	for app, nss := range auxExternal {
		var (
			buf = &bytes.Buffer{}
		)

		if err = json.NewEncoder(buf).Encode(nss); err != nil {
			return
		}
		lang.external[app] = bytes.NewReader(buf.Bytes())
	}

	return
}

func (svc *service) loadResourceTranslations(ctx context.Context, lang *Language, tg language.Tag) (err error) {
	lang.resources, err = svc.s.TransformResource(ctx, tg)
	if err != nil {
		return err
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
