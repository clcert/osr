package remote

import (
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/mailer"
	"github.com/sirupsen/logrus"
)

// Subject for Disk Full mail notification
const DiskFullMailSubject = "[OSR] ¡Un disco está casi lleno!"

// Form for Disk Full mail notification
const DiskFullMailTemplate = `
Hola, 

Te informamos que un servidor tiene problemas de espacio en disco:

Servidor:           {{.Name}}
Address:            {{.Address}}
Fecha Revisión:     {{.Date}}


{{.Disks.String}}

Adjuntamos también los archivos de log asociados a la última revisión.

Saludos!

Observatorio de Seguridad de la Red
CLCERT Universidad de Chile.
`

// Notifies if a server has exceeded the danger threshold in a disk
// and the flag MailNotify is set.
func notify(server *Server, info *ServerInfo) {
	maxCap := info.Disks.getMaxCapacity()
	logFields := logrus.Fields{
		"server":       info.Name,
		"ip":           info.Address,
		"max_capacity": maxCap,
		"cap_warning":  CapacityWarning,
	}
	var level mailer.NotifyLevel
	if maxCap >= CapacityWarning {
		level = mailer.WARN
		logs.Log.WithFields(logFields).Info("capacity exceeded in a disk!")
		mail := mailer.NotifyMail{
			Title:       DiskFullMailSubject,
			Template:    DiskFullMailTemplate,
			Values:      info,
			Attachments: []mailer.Attachable{server.Output, logs.Log},
			Level:       level,
		}
		err := mail.Send()
		if err != nil {
			logs.Log.WithFields(logFields).Info("couldn't send mail: %s", err)
		}
	}
}
