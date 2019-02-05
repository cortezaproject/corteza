package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"encoding/json"
	"io/ioutil"

	"github.com/davecgh/go-spew/spew"
)

type SpecFile []*SpecEntry

type SpecEntry struct {
	Title          string
	Description    string
	Protocol       string
	Authentication []string
	Entrypoint     string
	Path           string
	Struct         interface{}
	Parameters     map[string]interface{}
	APIs           []*SpecAPI
}

func (s *SpecEntry) toOutFile() OutFile {
	file := OutFile{
		APIs: []*OutFileAPI{},
	}
	s.applyToOutFile(&file)
	return file
}

func (s *SpecEntry) applyToOutFile(o *OutFile) {
	// reset title/interface/path to spec data
	o.Title = s.Title
	o.Description = s.Description
	o.Parameters = s.Parameters
	o.Interface = strings.ToUpper(s.Entrypoint[0:1]) + s.Entrypoint[1:]
	o.Path = s.Path
	if o.Path == "" {
		o.Path = "/" + s.Entrypoint
	}
	o.Struct = s.Struct
	o.Protocol = s.Protocol
	o.Authentication = s.Authentication

	namedAPIs := o.NamedAPIs()

	for _, val := range s.APIs {
		path := val.Path
		if path == "" {
			path = "/" + val.Name
		}
		// add new API calls
		call, ok := namedAPIs[val.Name]
		if !ok {
			o.APIs = append(o.APIs, &OutFileAPI{
				Name:       val.Name,
				Method:     val.Method,
				Title:      val.Title,
				Path:       path,
				Parameters: val.Parameters,
			})
		} else {
			// update title/method/path of existing APIs
			call.Name = val.Name
			call.Title = val.Title
			call.Method = val.Method
			call.Path = path
			if val.Parameters != nil {
				call.Parameters = val.Parameters
			}
		}
	}
}

type SpecAPI struct {
	Name       string
	Method     string
	Title      string
	Path       string
	Parameters map[string]interface{}
}

type OutFile struct {
	Title          string
	Description    string `json:",omitempty"`
	Interface      string
	Struct         interface{}
	Parameters     map[string]interface{}
	Protocol       string
	Authentication []string
	Path           string
	APIs           []*OutFileAPI
}

func (o *OutFile) NamedAPIs() map[string]*OutFileAPI {
	apis := map[string]*OutFileAPI{}
	for _, api := range o.APIs {
		apis[api.Name] = api
	}
	return apis
}

type OutFileAPI struct {
	Name        string
	Method      string
	Title       string
	Description string `json:",omitempty"`
	Path        string
	Parameters  map[string]interface{}
}

func main() {
	debug := false

	raw, err := ioutil.ReadFile("./spec.json")
	if err != nil {
		log.Fatal(err)
	}

	var spec SpecFile
	json.Unmarshal(raw, &spec)

	os.Mkdir("./spec", 0755)
	for _, val := range spec {
		filename := val.Entrypoint + ".json"
		var file OutFile
		contents, err := ioutil.ReadFile("./" + filename)
		if err != nil {
			file = val.toOutFile()
		} else {
			err = json.Unmarshal(contents, &file)
			if err != nil {
				log.Fatal("Error parsing ", filename, ": ", err)
			}
			val.applyToOutFile(&file)
		}
		raw, _ := json.MarshalIndent(file, "", "  ")
		ioutil.WriteFile("./spec/"+filename, raw, 0644)
		fmt.Println(filename)
	}

	if debug {
		spew.Dump(spec)
	}
}
