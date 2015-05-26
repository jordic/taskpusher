package taskpusher

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

// A backend is a way to persist tasks on disk
type Backend interface {
	Save(t Tasker)
	Load(status int) []Tasker
}

var (
	BoltBucket = []byte("tasks")
)

// BoltBack is a Backend for Manager that will store tasks on a
// boltdb db.
type BoltBack struct {
	*bolt.DB
}

func (db *BoltBack) Open(path string, mode os.FileMode) error {

	d, err := bolt.Open(path, mode, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	db.DB = d
	return nil
}

// Save a task to disk
func (db *BoltBack) Save(t Tasker) {

	err := db.DB.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucketIfNotExists(BoltBucket)
		enc, err := json.Marshal(t)

		err = b.Put([]byte(t.UID()), enc)
		return err
	})
	if err != nil {
		log.Println("Error saving task")
	}
}

// Load tasks from disk. status == -1 should load all tasks
func (db *BoltBack) Load(status int) []Tasker {
	var dest []Tasker

	err := db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(BoltBucket)
		if b == nil {
			return nil
		}
		err := b.ForEach(func(k, v []byte) error {
			
			var temp map[string]interface{}
			err := json.Unmarshal(v, &temp)
			if err != nil {
				return err
			}
			
			state := int(temp["Status"].(float64))
			if state == -1 || state == status {
				obj := FactoryTask(temp["Type"].(string))
				err = json.Unmarshal(v, obj)
				obj.SetStatus(state)
				if err != nil {
					return err
				}
				dest = append(dest, obj)
			}
			return nil
		})
		return err
	})
	if err != nil {
		log.Printf("Unable to decode tasks %s", err)
	}

	return dest

}