package alerts

type infobip struct {
	url            string
	username       string
	password       string
	api_key        string
	application_id string
	message_id     string
}

func NewInfobip(us string, pa string, ak string, aid string, mid string) Notifier {
	return &infobip{
		url:            "https://infobip.com/v2",
		username:       us,
		password:       pa,
		api_key:        ak,
		application_id: aid,
		message_id:     mid,
	}
}

func (i *infobip) satisfied() bool {
	return false
}

func (i *infobip) Notify(evt string) error {
	return nil
}
