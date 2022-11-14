package types

import (
	"reflect"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/locale"
	"golang.org/x/text/language"
)

func TestNew(t *testing.T) {
	lang := Lang{Tag: language.English}

	cases := []struct {
		name     string
		existing ResourceTranslationSet
		assert   locale.ResourceTranslationSet
		expect   ResourceTranslationSet
	}{{
		name:     "all empty",
		existing: nil,
		assert:   nil,
		expect:   nil,
	}, {
		name:     "no existing",
		existing: nil,
		assert:   locale.ResourceTranslationSet{{Lang: "en", Key: "k1", Msg: "k1 trans"}},
		expect:   ResourceTranslationSet{{Lang: lang, K: "k1", Message: "k1 trans"}},
	}, {
		name:     "determine diff",
		existing: ResourceTranslationSet{{Lang: lang, K: "k1", Message: "k1 trans"}, {Lang: lang, K: "k2", Message: "k2 trans"}},
		assert:   locale.ResourceTranslationSet{{Lang: "en", Key: "k1", Msg: "k1 trans"}, {Lang: "en", Key: "k2", Msg: "k2 trans"}, {Lang: "en", Key: "k3", Msg: "k3 trans"}},
		expect:   ResourceTranslationSet{{Lang: lang, K: "k3", Message: "k3 trans"}},
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := c.existing.New(c.assert)
			if !reflect.DeepEqual(c.expect, out) {
				t.Errorf("gotP = %v, want %v", out, c.expect)
			}
		})
	}
}

func TestOld(t *testing.T) {
	lang := Lang{Tag: language.English}

	cases := []struct {
		name     string
		existing ResourceTranslationSet
		assert   locale.ResourceTranslationSet
		expect   [][2]*ResourceTranslation
	}{{
		name:     "all empty",
		existing: nil,
		assert:   nil,
		expect:   nil,
	}, {
		name:     "no existing",
		existing: nil,
		assert:   locale.ResourceTranslationSet{{Lang: "en", Key: "k1", Msg: "k1 trans"}},
		expect:   nil,
	}, {
		name:     "determine diff",
		existing: ResourceTranslationSet{{Lang: lang, K: "k1", Message: "k1 trans"}, {Lang: lang, K: "k2", Message: "k2 trans"}},
		assert:   locale.ResourceTranslationSet{{Lang: "en", Key: "k1", Msg: "k1 trans (edit)"}, {Lang: "en", Key: "k2", Msg: "k2 trans (edit)"}, {Lang: "en", Key: "k3", Msg: "k3 trans (edit)"}},
		expect:   [][2]*ResourceTranslation{{{Lang: lang, K: "k1", Message: "k1 trans"}, {Lang: lang, K: "k1", Message: "k1 trans (edit)"}}, {{Lang: lang, K: "k2", Message: "k2 trans"}, {Lang: lang, K: "k2", Message: "k2 trans (edit)"}}},
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := c.existing.Old(c.assert)
			if !reflect.DeepEqual(c.expect, out) {
				t.Errorf("gotP = %v, want %v", out, c.expect)
			}
		})
	}
}
