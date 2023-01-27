package localization

import (
	"encoding/json"
	"fmt"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/model"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
)

var (
	localizer          *i18n.Localizer
	bundle             *i18n.Bundle
	SupportedLanguages = []model.Content{
		{Text: "English", Data: "en"},
		{Text: "Русский", Data: "ru"},
	}
)

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	for _, lang := range SupportedLanguages {
		path := fmt.Sprintf("configs/localization/dictionaries/%v.json", lang.Data)
		_, err := bundle.LoadMessageFile(path)
		if err != nil {
			log.Printf("no %v in %v", lang.Data, path)
		}
	}

	localizer = i18n.NewLocalizer(bundle, language.English.String())
}

// ChangeLanguage changes the app language to a new one and returns true if it is supported, otherwise returns false
func ChangeLanguage(newLang string) bool {
	for _, lang := range SupportedLanguages {
		if newLang == lang.Data {
			localizer = i18n.NewLocalizer(bundle, newLang)
			return true
		}
	}
	return false
}

// Message handles config by id and returns localized message
func Message(id string) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: id,
	})
}
