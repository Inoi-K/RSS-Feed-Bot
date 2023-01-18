package localization

import (
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
)

var (
	localizer          *i18n.Localizer
	bundle             *i18n.Bundle
	supportedLanguages = map[string]struct{}{
		"en": {},
		"ru": {},
	}
)

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	for lang := range supportedLanguages {
		path := fmt.Sprintf("configs/localization/dictionaries/%v.json", lang)
		_, err := bundle.LoadMessageFile(path)
		if err != nil {
			log.Printf("no %v in %v", lang, path)
		}
	}

	localizer = i18n.NewLocalizer(bundle, language.English.String())
}

func ChangeLanguage(lang string) bool {
	_, ok := supportedLanguages[lang]
	if ok {
		localizer = i18n.NewLocalizer(bundle, lang)
	}
	return ok
}

func Message(id string) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: id,
	})
}
