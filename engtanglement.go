package entangle

import "sync"

type Entanglement struct {
	data          string
	config        Config
	mutex         *sync.RWMutex
	sendChan      *chan string
	communication *Communication
}

func New(config Config) (*Entanglement, error) {
	sendChannel := make(chan string, config.sendBufferSize)
	communication, err := NewCommunication(config)
	if err != nil {
		return nil, err
	}
	return &Entanglement{
		mutex:         &sync.RWMutex{},
		config:        config,
		sendChan:      &sendChannel,
		communication: communication,
	}, nil
}

func (e *Entanglement) Set(value string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.data = value
	e.communication.BroadcastNewValue(value)
}

func (e *Entanglement) Get() string {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	return e.data
}
