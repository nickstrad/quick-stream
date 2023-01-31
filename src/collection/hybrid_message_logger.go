package collection

import (
	"fmt"
	"log"

	"github.com/cockroachdb/pebble"
)

type HybridMessageLogger struct {
	transientStateDB *pebble.DB
	failedStateDB    *pebble.DB
}

func NewHybridMessageLogger() *HybridMessageLogger {
	options := &pebble.Options{}
	transientStateDB, err := pebble.Open("data/transient", options)
	if err != nil {
		log.Fatal(err)
	}
	failedStateDB, err := pebble.Open("data/failed", options)
	if err != nil {
		log.Fatal(err)
	}

	return &HybridMessageLogger{
		transientStateDB: transientStateDB,
		failedStateDB:    failedStateDB,
	}
}

func (logger *HybridMessageLogger) WriteTransientStateMessage(key, val string) {
	if err := logger.transientStateDB.Set([]byte(key), []byte(val), pebble.Sync); err != nil {
		log.Fatal(err)
	}
}
func (logger *HybridMessageLogger) WriteFailedStateMessage() {
	key := []byte("hello")
	if err := logger.failedStateDB.Set(key, []byte("world"), pebble.Sync); err != nil {
		log.Fatal(err)
	}
}

func (logger *HybridMessageLogger) DeleteTransientStateMessage() {
	err := logger.transientStateDB.Delete([]byte("world"), pebble.Sync)
	if err != nil {
		fmt.Println(err)
	}

}
func (logger *HybridMessageLogger) DeleteFailedStateMessage(messageKey string) {
	err := logger.failedStateDB.Delete([]byte("world"), pebble.Sync)
	if err != nil {
		fmt.Println(err)
	}
}

func (logger *HybridMessageLogger) CleanUp() {
	err := logger.transientStateDB.Close()
	err2 := logger.failedStateDB.Close()

	if err != nil {
		log.Fatal(err)
	}

	if err2 != nil {
		log.Fatal(err2)
	}
}
