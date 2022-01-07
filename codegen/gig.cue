package codegen

import (
	"strings"

	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

_taskKinds: ["decoders", "preprocessors", "postprocessors"]

gig:
	[...schema.#codegen] &
	[
		for k in _taskKinds {
			template: "gocode/gig/task_$kind.go.tpl"
			output:   "pkg/gig/task_\(k).gen.go"
			payload: {
				package: "gig"
				imports: []

				taskKind:  strings.TrimSuffix(k, "s")
				taskConst: "Task" + strings.ToTitle(taskKind)
				tasks:     app.corteza.gig[k]
			}
		},
	]+

	[
		{
			template: "gocode/gig/worker.go.tpl"
			output:   "pkg/gig/worker.gen.go"
			payload: {
				package: "gig"
				imports: []
				workers: app.corteza.gig.workers
			}
		},
	]+

	[
		{
			template: "gocode/gig/conv.go.tpl"
			output:   "system/rest/conv/gig.gen.go"
			payload: {
				package: "conv"
				imports: []

				tasks: [ for k in _taskKinds {
					kind:    strings.TrimSuffix(k, "s")
					expKind: strings.ToTitle(kind)

					unwrapFunc:    "Unwrap\(expKind)"
					unwrapSetFunc: "\(unwrapFunc)Set"

					wrapFunc:    "Wrap\(expKind)"
					wrapSetFunc: "\(wrapFunc)Set"

					set: app.corteza.gig[k]
				}]
				workers: app.corteza.gig.workers
			}
		},
	]
