package locale

import "golang.org/x/text/language"

type (
	ResourceTranslation struct {
		Resource string `json:"resource"`
		Lang     string `json:"lang"`
		Key      string `json:"key"`
		Msg      string `json:"message"`
	}

	ResourceTranslationIndex map[string]*ResourceTranslation

	ResourceTranslationSet []*ResourceTranslation
)

func ContentID(cID uint64, i int) uint64 {
	if cID == 0 && i > 0 {
		return uint64(i)
	}
	return cID
}

func (rr ResourceTranslationSet) SetLanguage(tag language.Tag) {
	str := tag.String()
	for _, r := range rr {
		r.Lang = str
	}
}

// Returns true if resource translation set  contains foreign (non-native) languages
func (rr ResourceTranslationSet) ContainsForeign(native language.Tag) bool {
	str := native.String()

	for _, r := range rr {
		if r.Lang != str {
			return true
		}
	}

	return false
}

func (rx ResourceTranslationIndex) FindByKey(k string) *ResourceTranslation {
	return rx[k]
}
