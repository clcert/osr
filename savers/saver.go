// Savers are the media where the scanned information is stored.
//
// Currently there are two different savers: SFTP and postgres
//
// SFTP saves the information into files in a remote server, in CSV files
//
// Postgres saves the information in the OSR configured Postgresql Database.

package savers

import (
	"fmt"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/mailer"
	"github.com/clcert/osr/panics"
	"github.com/clcert/osr/utils"
	"github.com/fatih/structs"
	"reflect"
)

// A saver is a device where I can save results of processes.
type Saver interface {
	mailer.Attachable // Attachments to mailer
	// GetName returns a given name of the saver, based on task and process association. Two instances with the same source conf
	// should have different Names.
	GetName() string
	// Start starts the saver. Returns error if the server could not be started.
	Start() error
	// Finish sends a message to the saver to finish all its routines. When it returns, you can be sure that
	// all the files were saved
	Finish() error
	// SendMessage allows to send a message to the saver. This is useful if the saver needs some signals when processing data.
	SendMessage(msg interface{}) error
	// GetErrors returns a list of errors the saver has produced.
	GetErrors() []error
	// Save saves an object and returns an error if the server returns an error.
	Save(objs ...interface{}) error
}

// Config defines a saver in a process. It must define only one from [SFTP, HTTP, Script, ...]
// If you want to extend the savers, you must add a new type of config for the new saver.
type Config struct {
	Type     string          // type of the config (sftp, postgres)
	SFTP     *SFTPConfig     // Config if type is sftp
	Postgres *PostgresConfig // Config if type is postgres
}

// A savable object is an object sent to being saved. It allows to add metainformation to the savable object, via a hashmap.
type Savable struct {
	Object interface{}       // Object related to savable
	Meta   map[string]string // Metainformation for object
}

// New creates a new Savable from a config file. If you want to extend the savers, you must add a new case, pointing
// to its constructor method.
func (config Config) New(name string, params utils.Params) (Saver, error) {
	switch {
	case config.SFTP != nil:
		return config.SFTP.New(name, params)
	case config.Postgres != nil:
		return config.Postgres.New(name, params)
	default:
		return nil, fmt.Errorf("invalid saver")
	}
}

// Returns true if savable has a struct
func (savable *Savable) IsStruct() bool {
	r := reflect.ValueOf(savable.Object)
	for {
		switch r.Kind() {
		case reflect.Struct:
			return true
		case reflect.Ptr:
			r = r.Elem()
			continue
		default:
			return false
		}
	}
}

// Returns true if savable has a map of strings to strings
func (savable *Savable) IsMap() bool {
	_, ok := savable.Object.(map[string]string)
	return ok
}

// StructName returns the name of the struct saved by the savable object.
func (savable *Savable) StructName() string {
	switch {
	case savable.IsStruct():
		s := structs.New(savable.Object)
		return s.Name()
	case savable.IsMap():
		if name, ok := savable.Meta["outID"]; ok {
			return name
		} else {
			panic(&panics.Info{
				Text:        "Cannot save map object without outID metaproperty using SFTP saver",
				Err:         fmt.Errorf("Cannot save map object without outID metaproperty using SFTP saver"),
				Attachments: []mailer.Attachable{logs.Log},
			})
		}
	default:
		panic(&panics.Info{
			Text:        "Cannot save object different than a struct or map on a SFTP saver",
			Err:         fmt.Errorf("trying to save on a sftp saver an object different than a struct or map"),
			Attachments: []mailer.Attachable{logs.Log},
		})
	}
}

// GetOutID returns the "out ID" of a savable object. If it's defined as a meta option
// its OutID is that value. If is not defined, the "out ID" is the struct name.
func (savable *Savable) GetOutID() string {
	if val, ok := savable.Meta["outID"]; ok {
		return val
	} else {
		return savable.StructName()
	}
}

// Fields transforms the structure contained by Savable in a map with values.
func (savable *Savable) Fields() map[string]interface{} {
	switch {
	case savable.IsStruct():
		s := structs.New(savable.Object)
		return s.Map()
	case savable.IsMap():
		savableMap := savable.Object.(map[string]string)
		fields := make(map[string]interface{}, len(savableMap))
		for f, v := range savableMap {
			fields[f] = v
		}
		return fields
	default:
		panic(&panics.Info{
			Text:        "Cannot save object different than a struct or map on a SFTP saver",
			Err:         fmt.Errorf("trying to save on a sftp saver an object different than a struct or map"),
			Attachments: []mailer.Attachable{logs.Log},
		})
	}
}

// FieldNames returns the names of the fields of the structure contained by Savable
func (savable *Savable) FieldNames() []string {
	switch {
	case savable.IsStruct():
		s := structs.New(savable.Object)
		return s.Names()
	case savable.IsMap():
		savableMap := savable.Object.(map[string]string)
		fieldNames := make([]string, len(savableMap))
		i := 0
		for f, _ := range savableMap {
			fieldNames[i] = f
			i++
		}
		return fieldNames
	default:
		panic(&panics.Info{
			Text:        "Cannot save object different than a struct or map on a SFTP saver",
			Err:         fmt.Errorf("trying to save on a sftp saver an object different than a struct or map"),
			Attachments: []mailer.Attachable{logs.Log},
		})
	}
}
