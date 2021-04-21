package seeder

import (
	"fmt"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
)

type (
	faker struct {
		methods map[string]func() string
	}

	valueOptions struct {
	}
)

var (
	fakerMethods = map[string]func() string{
		"Name":        gofakeit.Name,
		"FirstName":   gofakeit.FirstName,
		"LastName":    gofakeit.LastName,
		"Title":       gofakeit.JobTitle,
		"Phone":       gofakeit.Phone,
		"MobilePhone": gofakeit.Phone,
		"Email":       gofakeit.Email,
	}
)

func Faker() *faker {
	return &faker{fakerMethods}
}

// seed to ensure randomization on initial
func (f faker) seed() {
	gofakeit.Seed(0)
	return
}

// generateValue generate value based on name or given type
func (f faker) fakeValueByName(name string) (val string, ok bool) {
	// Ensure randomization on initial
	f.seed()

	// Generate & return value from mapped methods
	method, ok := f.methods[name]
	if ok {
		return method(), ok
	}
	return
}

// @todo: currently its dependent on predefined methods but custom kind fake gen can be improved!
// generateValue generate value based on name or given type
func (f faker) fakeValue(name, kind string, opt valueOptions) (val string, err error) {
	// Generate & return value from mapped methods
	val, ok := f.fakeValueByName(name)
	if ok {
		return
	}

	// Ensure randomization on initial
	f.seed()
	valueKind := toLowerCase(kind)

	// Since, we don't have faker method for it,
	// we will generate the value based on kind(type)
	switch true {
	case valueKind == "string":
		val = gofakeit.LoremIpsumWord()
		break
	case valueKind == "int":
		break
	case valueKind == "Geometry":
		// @todo need to update once we support all type of point for geo location
		val = fmt.Sprintf("{\"coordinates\":[%f,%f]}", gofakeit.Longitude(), gofakeit.Latitude())
		break
	}
	return
}

func (f faker) fakeUserHandle(s string) string {
	return gofakeit.Generate("??????") + "_seeded"
}

func toLowerCase(s string) string {
	return strings.ToLower(s)
}
