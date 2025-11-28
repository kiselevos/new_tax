package handlers

type PageData struct {
	PageTemplate    string      // имя шаблона контента
	PageTitle       string      // title страницы
	MetaDescription string      // meta description
	MetaKeywords    string      // meta keywords
	Payload         interface{} // любые данные для конкретной страницы
	FeedbackEmail   string      // общий параметр, чтобы не тянуть из ENV в layout
	CurrentYear     int         // общий параметр (например для footer)
}
