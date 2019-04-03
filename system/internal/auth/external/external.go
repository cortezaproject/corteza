package external

import (
	"log"

	"github.com/crusttech/crust/internal/settings"
)

func Init(settingsService settings.Service) {
	if eas, err := ExternalAuthSettings(settingsService); err != nil {
		log.Printf("failed load external authentication settings: %v", err)
	} else {
		setupGoth(eas)
	}
}
