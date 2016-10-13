package director

import (
	"github.com/boltdb/bolt"
	"encoding/json"
	"log"
	"neubot/common"
	"time"
)

type Measurement struct {
	Status     string
	StdoutPath string
	StderrPath string
	Timestamp  time.Time
	TestName   string
	TestId     string
	Workdir    string
	CmdLine    []byte
}

func MeasurementsAppend(measurement *Measurement) error {

	database, err := bolt.Open(common.DefaultMeasurementsDb(), 0600, nil)
	if err != nil {
		return err
	}
	defer database.Close()

	log.Printf("database %s openned\n", common.DefaultMeasurementsDb())

	database.Update(func (tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("Measurements"))
		if err != nil {
			log.Printf("failed to create bucket")
			return err
		}
		value, err := json.Marshal(measurement)
		if err != nil {
			log.Printf("failed to marshal measurement value")
			return err
		}
		key, err := json.Marshal(measurement.Timestamp.Unix())
		if err != nil {
			log.Printf("failed to marshal measurement key")
			return err
		}
		log.Printf("about to put measurement: %s => %s", key, value)
		err = bucket.Put(key, value)
		if err != nil {
			log.Printf("failed to put measurement in database")
			return err
		}
		return nil
	})

	log.Printf("measurement appended to database\n")
	return nil
}
