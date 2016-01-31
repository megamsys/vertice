package alerts

import (
	"github.com/Bowery/slack"
)

type slacker struct {
	token string
	chnl  string
}

func NewSlack(tok string, cl string) Notifier {
	return &slacker{token: tok, chnl: cl}
}

func (s *slacker) satisfied() bool {
	return true
}

func (s *slacker) Notify(evt string) error {
	if !s.satisfied() {
		return nil
	}
	if err := slack.NewClient(s.token).SendMessage("#"+s.chnl, evt, "megamio"); err != nil {
		return err
	}
	return nil
}
