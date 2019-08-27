package panics

import (
	"fmt"
	"github.com/clcert/osr/mailer"
)

// An info struct defines a panic information container.
type Info struct {
	Text        string
	Err         error
	Attachments []mailer.Attachable
}

func (info *Info) Error() string {
	if info.Err != nil {
		return fmt.Sprintf("%s: %s", info.Text, info.Err.Error())
	} else {
		return info.Text
	}
}
