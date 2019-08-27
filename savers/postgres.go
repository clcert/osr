package savers

import (
	"fmt"
	"github.com/clcert/osr/databases"
	"github.com/clcert/osr/logs"
	"github.com/clcert/osr/utils"
	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

// PostgresConfig defines a basic configuration for a Postgres saver.
type PostgresConfig struct {
	Buffer        int                      // Number of objects by object type to hold before commiting to database.
	InsertTimeout time.Duration            // Number of seconds to wait to autocommit if there are no new objects received.
	Mutex         bool                     // Use a mutex to insert items.
	InsertConfig  map[string]*InsertConfig // Config specific for each object or outID
}

// InsertConfig defines a configuration related to a specific type of object, to being saved
// in the database. For now, it only specifies "onconflict" rules.
type InsertConfig struct {
	OnConflict string   // The statement for "on conflict". It generally defines the conflict fields and if the SQL sentence should do nothing or update.
	Set        []string // A list of fields to update in case of conflict. If onconflict is "do nothing", this is not necessary.
}

// PostgresSaver defines a saver structure that connects with the default Postgres database.
type PostgresSaver struct {
	*PostgresConfig                       // Configuration related to the saver
	name        string                    // Name for the saver instance
	chanMutex   sync.Mutex                // Mutex to edit maps concurrently
	insertMutex sync.Mutex                // Mutex for inserting if Config has mutex = true
	db          *pg.DB                    // Pointer to Postgres DB connection
	wg          sync.WaitGroup            // Wait Group to wait the finish of all channels
	channels    map[string]chan<- Savable // A map of channels identified by the outID of the objects
	inserted    int                       // Number of inserted elements (FYI)
	errors      []error                   // List of errors
	log         *logs.OSRLog              // Saver log
}

// New creates a new PostgresSaver based on a PostgresConfig struct.
func (config *PostgresConfig) New(name string, params utils.Params) (*PostgresSaver, error) {
	if config.InsertConfig == nil {
		config.InsertConfig = make(map[string]*InsertConfig, 0)
	}
	err := config.Format(params)
	if err != nil {
		return nil, err
	}
	if config.InsertTimeout == 0 {
		config.InsertTimeout = 10 // default timeout
	}
	if config.Buffer == 0 {
		config.Buffer = 1024 // Default buffer
	}
	log, err := logs.NewLog(name)
	if err != nil {
		return nil, err
	}
	return &PostgresSaver{
		name:           name,
		PostgresConfig: config,
		channels:       make(map[string]chan<- Savable),
		errors:         make([]error, 0),
		inserted:       0,
		log:            log,
	}, nil
}

func (config *PostgresConfig) Format(params utils.Params) error {
	if params == nil {
		return nil
	}
	newInsertConfig := make(map[string]*InsertConfig)
	for k, v := range config.InsertConfig {
		newInsertConfig[params.FormatString(k)] = v.Format(params)
	}
	config.InsertConfig = newInsertConfig
	return nil
}

func (saver *PostgresSaver) Start() error {
	db, err := databases.GetPostgresWriter()
	if err != nil {
		return err
	}
	saver.db = db
	return nil
}

func (saver *PostgresSaver) SendMessage(msg interface{}) error {
	msgMap, ok := msg.(map[string]string)
	if !ok {
		return fmt.Errorf("message not understood")
	}
	// check close messages
	closeID, ok := msgMap["close"]
	if !ok {
		return fmt.Errorf("message not understood")
	}
	saver.chanMutex.Lock()
	defer saver.chanMutex.Unlock()
	objChan, ok := saver.channels[closeID]
	if !ok {
		return fmt.Errorf("channel not found: %s", closeID)
	}
	close(objChan)
	delete(saver.channels, closeID)
	return nil
}

func (saver *PostgresSaver) Save(objs ...interface{}) error {
	for _, obj := range objs {
		var savable Savable
		switch obj.(type) {
		case Savable:
			savable = obj.(Savable)
		default:
			savable = Savable{Object: obj}
		}
		saver.getChannel(savable.GetOutID()) <- savable
	}
	return nil
}

func (saver *PostgresSaver) Finish() error {
	saver.chanMutex.Lock()
	for _, output := range saver.channels {
		close(output)
	}
	saver.chanMutex.Unlock()
	saver.wg.Wait()
	return saver.db.Close()
}

func (saver *PostgresSaver) GetErrors() []error {
	return saver.errors
}

// TODO: add queries as log
func (saver *PostgresSaver) GetAttachments() []string {
	return []string{saver.log.Path}
}

func (saver *PostgresSaver) GetName() string {
	return saver.name
}

func (saver *PostgresSaver) getInsertConfig(object Savable) *InsertConfig {
	conf, ok := saver.InsertConfig[object.GetOutID()]
	if !ok {
		conf, ok = saver.InsertConfig[object.StructName()]
		if !ok {
			conf = &InsertConfig{}
			saver.InsertConfig[object.GetOutID()] = conf
		}
	}
	return conf
}

func (saver *PostgresSaver) startChannel(name string, channel <-chan Savable) {
	saver.wg.Add(1)
	var config *InsertConfig
	objectList := make([]interface{}, 0)
	for newObject := range channel {
		if config == nil {
			config = saver.getInsertConfig(newObject)
		}
		objectList = append(objectList, newObject.Object)
		if len(objectList) >= saver.Buffer {
			saver.insertToDatabase(objectList, config)
			objectList = make([]interface{}, 0, saver.Buffer)
		}
	}
	saver.insertToDatabase(objectList, config)
	saver.log.WithFields(logrus.Fields{
		"outID": name,
	}).Info("Done, deleting saver branch...")
	saver.chanMutex.Lock()
	delete(saver.channels, name)
	saver.chanMutex.Unlock()
	saver.wg.Done()
}

func (saver *PostgresSaver) insertToDatabase(objList []interface{}, config *InsertConfig) {
	if len(objList) > 0 {
		query := saver.db.Model(objList...)
		if len(config.OnConflict) > 0 {
			query = query.OnConflict(config.OnConflict)
			for _, setOption := range config.Set {
				query = query.Set(setOption)
			}
		}
		saver.log.WithFields(logrus.Fields{
			"number":   len(objList),
			"inserted": saver.inserted,
			"mutex":    saver.Mutex,
		}).Info("Inserting entries into database...")
		if saver.Mutex {
			saver.insertMutex.Lock()
			defer saver.insertMutex.Unlock()
		}
		result, err := query.Insert()
		if err != nil {
			saver.log.Errorf("Couldn't insert entries to database: %v", err)
			saver.errors = append(saver.errors, err)
			return
		}
		saver.inserted += result.RowsAffected()
	}
}

func (saver *PostgresSaver) getChannel(name string) chan<- Savable {
	saver.chanMutex.Lock()
	defer saver.chanMutex.Unlock()
	channel, ok := saver.channels[name]
	if ok {
		return channel
	}
	newChannel := make(chan Savable)
	saver.channels[name] = newChannel
	go saver.startChannel(name, newChannel)
	return newChannel
}

func (insertconfig *InsertConfig) Format(params utils.Params) *InsertConfig {
	newFileConfig := &InsertConfig{
		OnConflict: params.FormatString(insertconfig.OnConflict),
		Set:        make([]string, len(insertconfig.Set)),
	}
	for i, v := range insertconfig.Set {
		newFileConfig.Set[i] = params.FormatString(v)
	}
	return newFileConfig
}
