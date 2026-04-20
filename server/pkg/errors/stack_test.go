package errors

import (
	"reflect"
	"testing"
)

func Test_convertNodeStack(t *testing.T) {
	converted := convertNodeStack([]string{
		"Object.exec (/dev/human/server-corredor/usr/testing/server-scripts/foo.js:12:11)",
		"Object.<anonymous> (/dev/human/server-corredor/node_modules/@planetcrust/human-js/src/corredor/exec.ts:26:35)",
		"Generator.next (<anonymous>)",
		"/dev/human/server-corredor/node_modules/@planetcrust/human-js/dist/index.js:277:71",
		"new Promise (<anonymous>)",
		"__awaiter (/dev/human/server-corredor/node_modules/@planetcrust/human-js/dist/index.js:273:12)",
		"Object.Exec (/dev/human/server-corredor/node_modules/@planetcrust/human-js/dist/index.js:483:12)",
		"ServerScripts.exec (/dev/human/server-corredor/src/services/server-scripts.ts:86:17)",
		"Object.Exec (/dev/human/server-corredor/src/grpc-handlers/server-scripts.ts:142:11)",
		"/dev/human/server-corredor/node_modules/grpc/src/server.js:593:13",
	})

	expecting := []*frame{
		{"Object.exec", "/dev/human/server-corredor/usr/testing/server-scripts/foo.js", 12},
		{"Object.<anonymous>", "/dev/human/server-corredor/node_modules/@planetcrust/human-js/src/corredor/exec.ts", 26},
		{"Generator.next", "<anonymous>", 0},
		{"", "/dev/human/server-corredor/node_modules/@planetcrust/human-js/dist/index.js", 277},
		{"new Promise", "<anonymous>", 0},
		{"__awaiter", "/dev/human/server-corredor/node_modules/@planetcrust/human-js/dist/index.js", 273},
		{"Object.Exec", "/dev/human/server-corredor/node_modules/@planetcrust/human-js/dist/index.js", 483},
		{"ServerScripts.exec", "/dev/human/server-corredor/src/services/server-scripts.ts", 86},
		{"Object.Exec", "/dev/human/server-corredor/src/grpc-handlers/server-scripts.ts", 142},
		{"", "/dev/human/server-corredor/node_modules/grpc/src/server.js", 593},
	}

	if len(expecting) != len(converted) {
		t.Errorf("converted stack length (%d) not match expectations (%d)", len(converted), len(expecting))
	}

	for i := range expecting {
		if !reflect.DeepEqual(converted[i], expecting[i]) {
			t.Errorf("converted stack frame %d does not match\nexpecting: %v\ngot:       %v", i, expecting[i], converted[i])
		} else {
			t.Logf("converted stack frame %d ok", i)
		}
	}

}
