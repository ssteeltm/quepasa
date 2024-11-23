package library

type LogFieldCollection struct {
	MessageId string
	ChatId    string
	EventId   string
}

var LogFields = LogFieldCollection{
	MessageId: "msgid",
	ChatId:    "chatid",
	EventId:   "eventid",
}
