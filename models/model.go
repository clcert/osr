// models contains all the models defined and usable
// natively in the main Postgresql database.
package models

import (
	"fmt"

	"github.com/clcert/osr/databases"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/query"
	"github.com/fatih/structs"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ModelConfig defines a configuration list in
type ModelConfig struct {
	BeforeCreate query.Configs // Statements to execute before model creation
	AfterCreate  query.Configs // Statements to execute after model creation
}

// Model defines abstractly a new model in the application.
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

// ExecBefore executes the queries defined in BeforeCreate config file
func (conf *ModelConfig) ExecBefore(db *pg.DB) error {
	if conf == nil {
		return nil
	}
	for _, queryConfig := range conf.BeforeCreate {
		queries, err := queryConfig.Open()
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

// ExecAfter executes the queries defined in AfterCreate config file
func (conf *ModelConfig) ExecAfter(db *pg.DB) error {
	if conf == nil {
		return nil
	}
	for _, queryConfig := range conf.AfterCreate {
		queries, err := queryConfig.Open()
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

// GetConfig unmarshals the config defined in viper for a model as a ModelConfig struct.
func (m *Model) GetConfig() (*ModelConfig, error) {
	var modelConf *ModelConfig
	structName := structs.Name(m.StructType)
	err := viper.UnmarshalKey(fmt.Sprintf("models.%s", structName), &modelConf)
	return modelConf, err
}

// BeforeCreateTable executes the statements in BeforeCreateStmts
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
		}
	}
	// From config
	if conf, err := m.GetConfig(); err == nil {
		logs.Log.WithFields(logrus.Fields{
			"model": m.Name,
		}).Info("Executing before statements from config...")
		if err := conf.ExecBefore(db); err != nil {
			return err
		}
	} else {
		logs.Log.WithFields(logrus.Fields{
			"model": m.Name,
		}).Error("Error getting after statements: %s", err)
	}
	return nil
}

// AfterCreateTable executes the statements in AfterCreateStmts
func (m *Model) AfterCreateTable(db *pg.DB) error {
	// Hardcoded Statements
	if m.AfterCreateStmts != nil {
		for _, stmt := range m.AfterCreateStmts {
			logs.Log.WithFields(logrus.Fields{
				"model": m.Name,
				"stmt":  stmt,
			}).Info("Executing hardcoded after statements...")
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
		}).Info("Executing after statements from config...")
		if err := conf.ExecAfter(db); err != nil {
			return err
		}
	} else {
		logs.Log.WithFields(logrus.Fields{
			"model": m.Name,
		}).Error("Error getting after statements: %s", err)
	}
	return nil
}

// Append appends a model to a model list.
func (s *ModelsList) Append(m Model) {
	s.Models = append(s.Models, m)
}

// CreateTables executes the statements necessary to create all the tables of the Models Group.
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
		err = m.BeforeCreateTable(db)
		if err != nil {
			logs.Log.WithFields(logrus.Fields{
				"model": m.Name,
				"error": err,
			}).Errorf("Error executing the statements before model creation: %s", err)
			return err
		}
		logs.Log.WithFields(logrus.Fields{
			"model": m.Name,
		}).Info("Creating table for model...")
		err := db.Model(m.StructType).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: false,
		})
		if err != nil {
			logs.Log.WithFields(logrus.Fields{
				"model": m.Name,
			}).Errorf("Error creating table for model: %s", err)
			return err
		}
		err = m.AfterCreateTable(db)
		if err != nil {
			logs.Log.WithFields(logrus.Fields{
				"model": m.Name,
			}).Errorf("Error executing the statements after model creation: %s", err)
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

// execStmts executes a list of statements.
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
