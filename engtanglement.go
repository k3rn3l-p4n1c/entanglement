package entanglement

import (
	"sync"
	"log"
	"os"
	"github.com/hashicorp/raft"
	"github.com/k3rn3l-p4n1c/entanglement/httpd"
	"encoding/json"
	"net/http"
	"fmt"
	"bytes"
	"github.com/pborman/uuid"
)

var (
	instance *System
	once     *sync.Once
)

type Config struct {
	RaftDir  string
	RaftAddr string
	HttpAddr string
	JoinAddr string
	NodeID   string
	Log      *log.Logger
}

func DefaultConfig() Config {
	return Config{
		RaftDir:  "./",
		RaftAddr: ":12700",
		HttpAddr: ":12701",
		JoinAddr: "",
		NodeID:   uuid.New(),
		Log:      log.New(os.Stdout, log.Prefix()+"[entanglement]", log.Flags()),
	}
}

type System struct {
	m map[string]*Entanglement

	logger *log.Logger
	mu     *sync.Mutex
	once   sync.Once

	RaftDir  string
	RaftBind string
	raft     *raft.Raft
}

type Entanglement struct {
	key    string
	data   string
	system *System
	mu     *sync.Mutex
}

func GetInstance(config Config) *System {
	once.Do(func() {
		os.MkdirAll(config.RaftDir, 0700)

		s := &System{
			logger: log.New(os.Stderr, "[store] ", log.LstdFlags),
			mu:     &sync.Mutex{},
		}

		s.RaftDir = config.RaftDir
		s.RaftBind = config.RaftAddr
		if err := s.Open(config.JoinAddr == "", config.NodeID); err != nil {
			config.Log.Fatalf("failed to open store: %s", err.Error())
		}

		h := httpd.New(config.HttpAddr, s)
		if err := h.Start(); err != nil {
			config.Log.Fatalf("failed to start HTTP service: %s", err.Error())
		}

		// If join was specified, make the join request.
		if config.JoinAddr != "" {
			if err := join(config.JoinAddr, config.RaftAddr, config.NodeID); err != nil {
				config.Log.Fatalf("failed to join node at %s: %s", config.JoinAddr, err.Error())
			}
		}

		config.Log.Println("hraft started successfully")

		instance = s
	})

	return instance
}

func (s *System) New(key string) *Entanglement {
	return &Entanglement{
		key:    key,
		data:   "",
		system: s,
		mu:     &sync.Mutex{},
	}
}

// Set sets the value for the given key.
func (e *Entanglement) Set(value string) error {
	return e.system.Set(e.key, value)
}

// Get returns the value for the given key.
func (e *Entanglement) Get() (string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.data, nil
}

func join(joinAddr, raftAddr, nodeID string) error {
	b, err := json.Marshal(map[string]string{"addr": raftAddr, "id": nodeID})
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/join", joinAddr), "application-type/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
