// This package contains all the models defined for using
// natively in the main Postgresql database.
package models

import (
	"fmt"
	"github.com/clcert/osr/databases"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/query"
	"github.com/fatih/structs"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Defines the default model list.
var DefaultModels = ModelsList{
	Name:   "default",
	Models: make([]Model, 0),
}

// ModelConfig defines a configuration list in
type ModelConfig struct {
	BeforeCreate query.FileMap   // Statements to execute before model creation
	AfterCreate  query.FileMap   // Statements to execute after model creation
}

// This struct defines abstractly a new model in the application.
type Model struct {
	Name                 string      // The human readable name for the model
	Description          string      // A description about the model
	StructType           interface{} // A pointer to the struct that defines the model
	BeforeCreateStmts    []string    // Statements to execute before the creation of the model
	AfterCreateStmts     []string    // Statements to execute after the creation of the model
	BeforeCreateFunction func(db *pg.DB) error
	AfterCreateFunction  func(db *pg.DB) error
}

// A ModelsList is a group of  Currently we use only
// one list, but that could change in the future if we want
// to put apart the logs models from the data
type ModelsList struct {
	Name              string   // A human readable name for the model group
	Models            []Model  // A list of models associated to the group
	BeforeCreateStmts []string // A list of SQL statements to execute before the model creation.
	AfterCreateStmts  []string // A list of SQL statements to execute after the model creation.
}

// Executes the queries defined in BeforeCreate config file
func (conf *ModelConfig) ExecBefore(db *pg.DB) error {
	if conf == nil {
		return nil
	}
	for queryFile, whitelist := range conf.BeforeCreate {
		queries, err := query.OpenFile(queryFile, whitelist...)
		if err != nil {
			return err
		}
		for _, q := range queries {
			_, err := q.Execute(db)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Executes the queries defined in AfterCreate config file
func (conf *ModelConfig) ExecAfter(db *pg.DB) error {
	if conf == nil {
		return nil
	}
	for queryFile, whitelist := range conf.AfterCreate {
		queries, err := query.OpenFile(queryFile, whitelist...)
		if err != nil {
			return err
		}
		for _, q := range queries {
			_, err := q.Execute(db)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Unmarshals the config defined in viper for a model as a ModelConfig struct.
func (m *Model) GetConfig() (*ModelConfig, error) {
	var modelConf *ModelConfig
	structName := structs.Name(m.StructType)
	err := viper.UnmarshalKey(fmt.Sprintf("models.%s", structName), &modelConf)
	return modelConf, err
}

// This function executes the statements in BeforeCreateStmts
func (m *Model) BeforeCreateTable(db *pg.DB) error {

	if m.BeforeCreateStmts != nil {
		logs.Log.WithFields(logrus.Fields{
			"model": m.Name,
		}).Info("Executing hardcoded before statements...")
		for _, stmt := range m.BeforeCreateStmts {
			_, err := db.Exec(stmt)
			if err != nil {
				return err
			}
		}
	}
	if m.BeforeCreateFunction != nil {
		logs.Log.WithFields(logrus.Fields{
			"model": m.Name,
		}).Info("Executing before function...")
		err := m.BeforeCreateFunction(db)
		if err != nil {
			return err
		}	}
	// From config
	if conf, err := m.GetConfig(); err == nil {
		logs.Log.WithFields(logrus.Fields{
			"model": m.Name,
			"conf": conf,
		}).Info("Executing before statements from config...")
		if err := conf.ExecBefore(db); err != nil {
			return err
		}
	}
	return nil
}

// This function executes the statements in AfterCreateStmts
func (m *Model) AfterCreateTable(db *pg.DB) error {
	// Hardcoded Statements
	if m.AfterCreateStmts != nil {
		logs.Log.WithFields(logrus.Fields{
			"model": m.Name,
		}).Info("Executing hardcoded after statements...")
		for _, stmt := range m.AfterCreateStmts {
			_, err := db.Model(m.StructType).Exec(stmt)
			if err != nil {
				return err
			}
		}
	}
	if m.AfterCreateFunction != nil {
		logs.Log.WithFields(logrus.Fields{
			"model": m.Name,
		}).Info("Executing after function...")
		err := m.AfterCreateFunction(db)
		if err != nil {
			return err
		}
	}
	// From config
	if conf, err := m.GetConfig(); err == nil {
		logs.Log.WithFields(logrus.Fields{
			"model": m.Name,
			"conf": conf,
		}).Info("Executing after statements from config...")
		if err := conf.ExecAfter(db); err != nil {
			return err
		}
	}
	return nil
}

// Appends a model to a model list.
func (s *ModelsList) Append(m Model) {
	s.Models = append(s.Models, m)
}

// Executes the statements necessary to create all the tables of the Models Group.
// First, it executes the statements present in "BeforeCreateStmts".
// Then, it creates the tables of the models group. Before and after each table created, it executes
// the "BeforeCreateTable" and "AfterCreateTable" function of the models, respectively.
// Finally, it executes the statements present in "AfterCreateStmts".
func (s *ModelsList) CreateTables() error {
	db, err := databases.GetPostgresWriter()
	if err != nil {
		return err
	}
	if err := s.execStmts(db, s.BeforeCreateStmts); err != nil {
		return err
	}
	for _, m := range s.Models {
		logs.Log.WithFields(logrus.Fields{
			"model": m.Name,
		}).Info("Creating table for model...")
		err = m.BeforeCreateTable(db)
		if err != nil {
			logs.Log.WithFields(logrus.Fields{
				"model": m.Name,
				"error": err,
			}).Error("Error executing the statements before model creation: %s", err)
			return err
		}
		err := db.CreateTable(m.StructType, &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: false,
		})
		if err != nil {
			logs.Log.WithFields(logrus.Fields{
				"model": m.Name,
			}).Error("Error creating table for model: %s", err)
			return err
		}
		err = m.AfterCreateTable(db)
		if err != nil {
			logs.Log.WithFields(logrus.Fields{
				"model": m.Name,
			}).Error("Error executing the statements after model creation: %s", err)
			return err
		}
	}
	if err := s.execStmts(db, s.AfterCreateStmts); err != nil {
		return err
	}
	// Give permission to reader
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}

// Executes a list of statements.
func (s *ModelsList) execStmts(db *pg.DB, stmtList []string) error {
	if stmtList != nil {
		for _, stmt := range stmtList {
			_, err := db.Exec(stmt)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
