package tasks

import (
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/mailer"
	"github.com/sirupsen/logrus"
)

// Subject for Task Done mail notification
const TaskDoneMailSubject = "[OSR] Tarea Finalizada"

// Body for Task Done mail notification
const TaskDoneMailTemplate = `
Hola, 

Te informamos que se terminaron de ejecutar una o más procesos incluidos en una tarea.

Nombre Tarea:					{{ .Name }}
Descripción Tarea:				{{ .Description }}
Número Tarea:					{{ .TaskSession.ID }}
Procesos satisfactorios:		{{ .GetSucceeded }}
Procesos con error:				{{ .GetFailed }}
Fecha Inicio: 					{{ .TaskSession.StartDate.Format "02-01-2006 15:04:05 -0700"}}
Fecha Fin:						{{ .TaskSession.EndDate.Format "02-01-2006 15:04:05 -0700"}}
Estado Final:					{{ .TaskSession.GetStatus }}
Abortar en Error:				{{ .AbortOnError }}

{{if .HasErrors}}
A continuación, mostramos una lista con los errores encontrados por tarea:

{{range $name, $err := .Failed }}
Tarea: {{ $name }}
Error: {{ $err }}

{{end}}
{{end}}

Adjuntamos también los archivos de log asociados a esta sesión de importación.

Saludos!

Observatorio de Seguridad de la Red
CLCERT Universidad de Chile.
`

// Checks if it should notify about an ended savers
// And notifies about it if it's the case.
func notify(task *Task) {
	fields := logrus.Fields{
		"has_error":     task.HasErrors(),
		"taskID":        task.TaskSession.ID,
		"start_date":    task.TaskSession.StartDate,
		"end_date":      task.TaskSession.EndDate,
		"import_status": task.TaskSession.GetStatus(),
	}

	if len(task.GetSucceeded()) > 0 {
		fields["succeeded"] = task.GetSucceeded()
	}
	if len(task.GetFailed()) > 0 {
		fields["failed"] = task.GetFailed()
	}
	logs.Log.WithFields(fields).Info("Task finished, sending mail...")
	var level mailer.NotifyLevel
	if task.HasErrors() {
		level = mailer.ERROR
	} else {
		level = mailer.INFO
	}
	mail := mailer.NotifyMail{
		Title:       TaskDoneMailSubject,
		Template:    TaskDoneMailTemplate,
		Values:      task,
		Attachments: []mailer.Attachable{logs.Log, task},
		Level:       level,
	}
	err := mail.Send()
	if err != nil {
		logs.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Couldn't send mail.")
	}
}
