package lang

type ServiceLanguage string

const (
	ServiceLanguageUnknown ServiceLanguage = "unknown"
	ServiceLanguageGo      ServiceLanguage = "go"
	ServiceLanguageJs      ServiceLanguage = "js"
)

var LanguagesList = []ServiceLanguage{
	ServiceLanguageGo,
	ServiceLanguageJs,
}
