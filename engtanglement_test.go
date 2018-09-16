package entanglement

import (
	"testing"
	"log"
	"os"
)

func TestGetAndSetLocal(t *testing.T) {
	conf := Config{
		RaftDir:  "./node1",
		RaftAddr: ":12001",
		HttpAddr: ":11001",
		JoinAddr: "",
		NodeID:   "node1",
		Log:      log.New(os.Stdout, log.Prefix()+"[entanglement1]", log.Flags()),
	}
	s := GetInstance(conf)

	entanglement := s.New("key")

	value := "sample data"
	entanglement.Set(value)

	receivedValue, _ := entanglement.Get()

	if value != receivedValue {
		t.Errorf("Setted value is not equal to getted value")
	}
}
