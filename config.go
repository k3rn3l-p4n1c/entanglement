package entangle

type Config struct {
	listenPort int
	interfaceHost string
	bufferSize int
	sendBufferSize int
}

func GetDefaultConfig() Config {
	return Config{
		listenPort: 6543,
		interfaceHost: "0.0.0.0",
		bufferSize: 1024,
		sendBufferSize: 10,
	}
}
