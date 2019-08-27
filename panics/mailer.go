package panics

import (
	"fmt"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/mailer"
	"os"
	"strings"
)

const panicMailSubject = "[OSR] Aplicación terminó de forma inesperada"

const panicMailTemplate = `
Hola, 

Te informamos que la aplicación de OSR acaba de caerse con panics.

Comando: {{.Command}}
Motivo: {{.Motive}}

Adjuntamos también los archivos de log asociados a la ejecución del comando.

Saludos!

Observatorio de Seguridad de la Red
CLCERT Universidad de Chile.
`

// Defines the args used in a panics mail.
type MailArgs struct {
	Motive  string
	Command string
}

// Notifies if notifylevel is not set to NONE
func NotifyPanic() {
	r := recover()
	if r != nil {
		info, ok := r.(*Info)
		if !ok {
			info = &Info{
				Err: fmt.Errorf("%s", r),
				Text: "Unhandled panic",
			}
		}
		logs.Log.Errorf("panic! %s", info)
		motive := fmt.Sprintf("%s: %s", info.Text, info.Error())
		panicArgs := &MailArgs{
			Motive:  motive,
			Command: strings.Join(os.Args, " "),
		}
		mail := mailer.NotifyMail{
			Title:       panicMailSubject,
			Template:    panicMailTemplate,
			Values:      panicArgs,
			Attachments: info.Attachments,
			Level:       mailer.PANIC,
		}
		err := mail.Send()
		if err != nil {
			fmt.Printf("%s\n", err)
		}

	}
}
