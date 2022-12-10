package consts

const (
	RU = "ru"
	EN = "en"

	SubscribeCommandFail = "subFail"
	NotValidLink         = "notValid"
)

var (
	LocText = map[string]map[string]string{
		RU: {
			StartCommand:         "Привет\\! \nЯ умею читать RSS\\-ленты и отправлять тебе обновления из твоих любимых источников\\! Если тебе нужна помощь по командам, пиши /help :\\)",
			HelpCommand:          "/sub <ссылка на RSS\\-ресурс\\> \\- добавить ссылку на интересующий RSS\\-ресурс\n/unsub <ссылка на RSS\\-ресурс\\> \\- отписаться от обновлений, используйте команду\n/upd \\- вывести последние обновления новостей\n/act <ссылка на RSS\\-ресурс\\> \\- включить получение обновлений от ресурса\n/deact <ссылка на RSS\\-ресурс\\> \\- отключить получение обновлений от ресурса, но оставить его в списке подписок\n/list \\- посмотреть cвои текущие подписки\n/help \\- получить информацию о всех командах",
			ListCommand:          "*Список подисок*",
			UpdateCommand:        "Обновления:",
			ActivateCommand:      "Пожалуйста, выбери подписку, которую ты бы %v:",
			UnsubscribeCommand:   "Пожалуйста, выбери подпику, от которой хотел бы отписаться:",
			SubscribeCommand:     "*Успешно подписан*\n[%v](%v\\)",
			SubscribeCommandFail: "*Подисаться не получилось*\n%v",
			NotValidLink:         "Ссылка не подходит",
		},
		EN: {
			StartCommand:         "Hi\\! \nI can read RSS feeds and send you updates from your favorite sources\\! If you need help, then write /help :\\)",
			HelpCommand:          "/sub <RSS resource link\\> \\- add a link to the RSS resource of interest\n/unsub <RSS resource link\\> \\- unsubscribe from updates\n/upd \\- display the latest news updates\n/act <RSS resource link\\> \\- enable receiving updates from the resource\n/deact <link to RSS resource\\> \\- disable receiving updates from the resource, but leave it in the list of subscriptions\n/list \\- view your current subscriptions\n/help \\- view information about all commands",
			ListCommand:          "*Subscription list*",
			UpdateCommand:        "Updates:",
			ActivateCommand:      "Please choose a subscription you'd like to %v:",
			UnsubscribeCommand:   "Please choose a subscription you'd like to unsubscribe from:",
			SubscribeCommand:     "*Successfully subscribed*\n[%v](%v\\)",
			SubscribeCommandFail: "*Failed to subscribe*\n%v",
			NotValidLink:         "The link is not valid",
		},
	}
)
