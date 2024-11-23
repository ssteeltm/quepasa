package library

type LogFieldCollection struct {
	MessageId string
	ChatId    string
	EventId   string
	WId       string
	Token     string
	Url       string
}

var LogFields = LogFieldCollection{
	MessageId: "msgid",
	ChatId:    "chatid",
	EventId:   "eventid",
	WId:       "wid",
	Token:     "token",
	Url:       "url",
}
