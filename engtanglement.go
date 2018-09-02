package entangle

import "sync"

type Entanglement struct {
	data string

	mutex *sync.Mutex
}

func New() *Entanglement {
	return &Entanglement{
		mutex: &sync.Mutex{},
	}
}

func (e *Entanglement) Set(value string)  {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.data = value
}

func (e *Entanglement) Get() string {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.data
}
