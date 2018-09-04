package entangle

import (
	"net"
	"fmt"
	"log"
)

type CommunicationConfig struct {
	bufferSize int
}

type Communication struct {
	config      CommunicationConfig
	connection  net.PacketConn
	sendChannel *chan string
}

func NewCommunication(config Config) (*Communication, error) {
	lp, err := net.ListenPacket("udp", fmt.Sprintf("%s:%d", config.interfaceHost, config.listenPort))
	if err != nil {
		return nil, err
	}
	return &Communication{
		connection: lp,
		config: CommunicationConfig{
			bufferSize: config.bufferSize,
		},
	}, nil
}

func (c *Communication) startToListen() {
	for {
		buf := make([]byte, c.config.bufferSize)

		n, addr, err := c.connection.ReadFrom(buf)
		if err != nil {
			log.Print(err)
		}
		log.Printf("get %d byte from %s", n, addr)
	}

}

func (c *Communication) startSend() {
	for value := range *c.sendChannel {
		f
	}
}
func (c *Communication) BroadcastNewValue(value string) {
	*c.sendChannel <- value
}
