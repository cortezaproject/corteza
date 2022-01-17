package codegen

import (
	"strings"

	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

gig:
	[...schema.#codegen] &
	[
		for k in ["decoders", "preprocessors", "postprocessors"] {
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

				tasks: [ for k in ["decoders", "preprocessors", "postprocessors"] {
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
	]+

	[
		for k in ["decoders", "postprocessors"] {
			template: "gocode/gig/integration/task_$kind_test.go.tpl"
			output:   "tests/gig/task_\(k).gen_test.go"
			payload: {
				package: "gig"
				imports: []

				taskKind:  strings.TrimSuffix(k, "s")
				taskExpKind: strings.ToTitle(taskKind)
				taskConst: "Task" + strings.ToTitle(taskKind)

				test: string | *"test_\(taskKind)_tasks"
				expTest: string | *strings.ToTitle(test)
				testWorker: string | *"\(test)_worker"

				tasks:     app.corteza.gig[k]
			}
		},
	]+

	// Moving preprocessors into sepparate block as it does have some
	// extra logic to it
	[
		{
			template: "gocode/gig/integration/task_preprocessors_test.go.tpl"
			output:   "tests/gig/task_preprocessors.gen_test.go"
			payload: {
				package: "gig"
				imports: []

				taskKind:  strings.TrimSuffix("preprocessors", "s")
				taskExpKind: strings.ToTitle(taskKind)
				taskConst: "Task" + strings.ToTitle(taskKind)

				test: string | *"test_\(taskKind)_tasks"
				expTest: string | *strings.ToTitle(test)
				testWorker: string | *"\(test)_worker"

				workers: app.corteza.gig.workers
			}
		},
	]+

	// Placeholder
	[]
