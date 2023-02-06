package hybrid_message_logger

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rs/xid"
)

var globalIds []xid.ID
var hml *HybridMessageLogger

func TestMain(m *testing.M) {
	// os.Exit skips defer calls
	// so we need to call another function
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	// pseudo-code, some implementation excluded:
	//
	// 1. create test.db if it does not exist
	// 2. run our DDL statements to create the required tables if they do not exist
	// 3. run our tests
	// 4. truncate the test db tables

	dir1, err := ioutil.TempDir("", "streamer_test_data")
	if err != nil {
		return -1, fmt.Errorf("could not create 1 of 2 temp directories to set up databse: %w", err)
	}
	dir2, err := ioutil.TempDir("", "streamer_test_data")
	if err != nil {
		return -1, fmt.Errorf("could not create 2 of 2 temp directories to set up databse: %w", err)
	}

	hml, err = NewHybridMessageLogger(dir1, dir2)

	if err != nil {
		return -1, fmt.Errorf("could not create hybrid message logger for test suite: %w", err)
	}

	defer func() {
		defer os.Remove(dir1)
		defer os.Remove(dir2)
		err := hml.Cleanup()
		if err != nil {
			fmt.Println(err)
		}
	}()

	for i := 0; i < 3; i++ {
		globalIds = append(globalIds, xid.New())
	}

	for _, id := range globalIds {
		err := hml.AddEvent(id, []byte("data"))
		if err != nil {
			return -1, fmt.Errorf("unable to add event to testsuite database: %w", err)
		}
	}

	return m.Run(), nil
}

func TestLoggerAddsEventWithValidStringParameters(t *testing.T) {
	t.Parallel()

	err := hml.AddEvent(xid.New(), []byte("valid string"))
	if err != nil {
		t.Error(err)
	}
}

func TestLoggerReturnsErrorWhenAddingEventWithInvalidValueParameter(t *testing.T) {
	t.Parallel()

	err := hml.AddEvent(xid.New(), []byte(""))
	if err == nil {
		t.Error("Should not allow empty string as value")
	}
}

func TestLoggerRemovesEventWithValidParameters(t *testing.T) {
	t.Parallel()

	id := globalIds[0]
	err := hml.RemoveEvent(id)
	if err != nil {
		t.Error("Unable to remove event with id" + id.String())
	}
}

func TestLoggerReturnsErrorWhenRemovingEventWithInvalidParameters(t *testing.T) {
	t.SkipNow()
	t.Parallel()
	t.Error(false)
}
