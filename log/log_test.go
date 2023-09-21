package log_test

import (
	"errors"
	"testing"
	"time"

	"bharvest.io/axelmon/log"
)

func TestAllLog(t *testing.T) {
	// Thread safe test
	timeout := time.After(5*time.Second)
	done := make(chan bool)

	go func() {
		go log.Info("Info log test")
		go log.Debug("Debug log test")
		go log.Error(errors.New("Log test"))
		log.Error(errors.New("Log test"))
		log.Debug(struct{
			name string
			age int
		}{
			"Choi",
			26,
		})

		// Wait log printing
		time.Sleep(3*time.Second)

		done <- true
	}()

	select {
	case <- timeout:
		t.Fatal("Timeout")
	case <- done:
	}
}
