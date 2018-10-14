package entanglement

import (
	"testing"
	"github.com/sirupsen/logrus"
	"time"
)

func TestGetAndSetLocal(t *testing.T) {
	conf := Config{
		RaftDir:  "./node1",
		RaftAddr: ":12001",
		HttpAddr: ":11001",
		JoinAddr: "",
		NodeID:   "node1",
		Log:      logrus.New(),
	}
	s := Bootstrap(conf)

	entanglement := s.New("key")

	value := "sample data"
	err := entanglement.Set(value)
	if err != nil {
		t.Fatal(err.Error())
	}
	time.Sleep(1 * time.Second)
	receivedValue, _ := entanglement.Get()

	if value != receivedValue {
		t.Errorf("Setted value is not equal to getted value: \"%s\" != \"%s\"", value, receivedValue)
	}
}
