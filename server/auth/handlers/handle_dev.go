package handlers

// Handlers dev-view for templates & scenarios
//
// See auth/README.adoc for details on how this works

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cortezaproject/corteza/server/auth/request"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"gopkg.in/yaml.v3"
)

type (
	devScenariosDoc []*devTemplate
	devTemplate     struct {
		Template string
		Scenes   devScenes
	}

	devScenes []*devScene
	devScene  struct {
		Name     string
		Template string
		Data     map[string]interface{}
	}
)

func (h *AuthHandlers) devView(req *request.AuthReq) (err error) {
	req.Template = "template-dev.html.tpl"
	req.Data["templates"], err = getScenes(h.Opt.AssetsPath)
	return
}

func (h *AuthHandlers) devSceneView(w http.ResponseWriter, r *http.Request) {
	s, err := findScenario(
		h.Opt.AssetsPath,
		r.URL.Query().Get("template"),
		r.URL.Query().Get("scene"),
	)

	if err == nil && s != nil {
		err = h.Templates.ExecuteTemplate(w, s.Template+".html.tpl", s.Data)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func findScenario(assetsPath, template, scene string) (*devScene, error) {
	templates, err := getScenes(assetsPath)
	if err != nil {
		return nil, err
	}

	for _, t := range templates {
		if t.Template != template {
			continue
		}

		for _, s := range t.Scenes {
			if s.Name != scene {
				continue
			}

			return s, nil
		}
	}

	return nil, nil
}

func getScenes(assetsPath string) (devScenariosDoc, error) {
	if assetsPath == "" {
		return nil, fmt.Errorf("can not load scenarios without AUTH_ASSETS_PATH and exported assets")
	}

	f, err := os.Open(assetsPath + "/templates/scenarios.yaml")
	if err != nil {
		return nil, err
	}
	aux := devScenariosDoc{}
	return aux, yaml.NewDecoder(f).Decode(&aux)
}

func (doc *devScenariosDoc) UnmarshalYAML(n *yaml.Node) error {
	return y7s.EachMap(n, func(k *yaml.Node, v *yaml.Node) (err error) {
		dt := &devTemplate{Template: k.Value, Scenes: devScenes{}}

		err = y7s.EachMap(v, func(k *yaml.Node, v *yaml.Node) (err error) {
			s := &devScene{
				Name:     k.Value,
				Template: dt.Template,
				Data:     make(map[string]interface{}),
			}

			if err = v.Decode(s.Data); err != nil {
				return
			}

			dt.Scenes = append(dt.Scenes, s)
			return nil
		})

		if err != nil {
			return
		}

		*doc = append(*doc, dt)
		return nil
	})
}
