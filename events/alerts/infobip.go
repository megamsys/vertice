package alerts

type infobip struct {
	url            string
	username       string
	password       string
	api_key        string
	application_id string
	message_id     string
}

func NewInfobip(m map[string]string) Notifier {
	return &infobip{
		url:            "https://infobip.com/v2",
		username:       m[USERNAME],
		password:       m[PASSWORD],
		api_key:        m[API_KEY],
		application_id: m[APPLICATION_ID],
		message_id:     m[MESSAGE_ID],
	}
}

func (i *infobip) satisfied() bool {
	return false
}

func (i *infobip) Notify(eva EventAction, m map[string]string) error {
	return nil
}
