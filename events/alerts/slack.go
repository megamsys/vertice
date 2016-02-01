package alerts

import (
	"github.com/Bowery/slack"
)

type slacker struct {
	token string
	chnl  string
}

func NewSlack(m map[string]string) Notifier {
	return &slacker{token: m[TOKEN], chnl: m[CHANNEL]}
}

func (s *slacker) satisfied() bool {
	return true
}

func (s *slacker) Notify(eva EventAction, m map[string]string) error {
	if !s.satisfied() {
		return nil
	}
	if err := slack.NewClient(s.token).SendMessage("#"+s.chnl, m["message"], "megamio"); err != nil {
		return err
	}
	return nil
}
