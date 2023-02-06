package hybrid_message_logger

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cockroachdb/pebble"
	"github.com/rs/xid"
)

type HybridMessageLogger struct {
	transientStateDB *pebble.DB
	failedStateDB    *pebble.DB
}

func NewHybridMessageLogger(transientStateDBPath, failedStateDBPath string) (*HybridMessageLogger, error) {
	transientStateDB, err := pebble.Open(transientStateDBPath, &pebble.Options{})
	if err != nil {
		return nil, err
	}
	failedStateDB, err := pebble.Open(failedStateDBPath, &pebble.Options{})
	if err != nil {
		return nil, err
	}
	return &HybridMessageLogger{
		transientStateDB: transientStateDB,
		failedStateDB:    failedStateDB,
	}, nil
}

var ErrNotFound = errors.New("pebble: not found")
var ErrInvalidValueParameter = errors.New("invalid event value being persisted")

func (h *HybridMessageLogger) AddEvent(key xid.ID, val []byte) error {
	fmt.Println("Adding id '" + key.String() + "' with val\n" + string(val) + "\nto transient database.")

	if string(val) == "" {
		return ErrInvalidValueParameter
	}

	value, closer, err := h.transientStateDB.Get(key.Bytes())

	if string(value) != "" {
		closer.Close()
		return errors.New(key.String() + " already exists in the database.")
	}

	if err != nil && err.Error() == ErrNotFound.Error() {
		if err := h.transientStateDB.Set(key.Bytes(), val, pebble.Sync); err != nil {
			log.Println("Error setting data for\nkey: " + key.String() + " \nval: " + string(val))
			return err
		}
	}

	return nil
}

func (h *HybridMessageLogger) MoveToFailed(key xid.ID) {
	fmt.Println("Moving id '" + key.String() + "' to failed database.")
	if err := h.failedStateDB.Set(key.Bytes(), []byte("failed"), pebble.Sync); err != nil {
		log.Println(err)
		return
	}
}

func (h *HybridMessageLogger) RemoveEvent(key xid.ID) error {
	fmt.Println("Removing id '" + key.String() + "' from transient database.")

	if err := h.transientStateDB.Delete(key.Bytes(), pebble.Sync); err != nil {
		return err
	}
	return nil
}

func (h *HybridMessageLogger) Cleanup() error {
	errs := make([]string, 2)

	err := h.transientStateDB.Close()

	if err != nil {
		errs = append(errs, err.Error())
	}

	err = h.failedStateDB.Close()

	if err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}
