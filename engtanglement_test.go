package entangle

import "testing"

func TestGetAndSetLocal(t *testing.T) {
	entanglement := New()

	value := "sample data"
	entanglement.Set(value)

	receivedValue := entanglement.Get()

	if value != receivedValue {
		t.Errorf("Setted value is not equal to getted value")
	}
}
