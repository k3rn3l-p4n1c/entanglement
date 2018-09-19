package entanglement

import (
	"sync"
	"os"
	"github.com/hashicorp/raft"
	"github.com/k3rn3l-p4n1c/entanglement/httpd"
	"encoding/json"
	"net/http"
	"fmt"
	"bytes"
	"github.com/pborman/uuid"
	"errors"
	"time"
	"github.com/sirupsen/logrus"
)


type Config struct {
	RaftDir  string
	RaftAddr string
	HttpAddr string
	JoinAddr string
	MaxRetry int
	NodeID   string
	Log      *logrus.Logger
}

func DefaultConfig() Config {
	return Config{
		RaftDir:  "./",
		RaftAddr: ":12700",
		HttpAddr: ":12701",
		JoinAddr: "",
		MaxRetry: 3,
		NodeID:   uuid.New(),
		Log:      logrus.New(),
	}
}

type System struct {
	m map[string]*Entanglement

	logger *logrus.Logger
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

func Bootstrap(config Config) *System {
	os.MkdirAll(config.RaftDir, 0700)

	s := &System{
		m:      make(map[string]*Entanglement),
		logger: logrus.New(),
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
		config.Log.Debugf("trying to connect to %s", config.JoinAddr)
		err := join(config.JoinAddr, config.RaftAddr, config.NodeID)
		retry := 0
		for err != nil {
			if retry >= config.MaxRetry {
				config.Log.Fatalf("failed to join node at %s: %s after %d retries", config.JoinAddr, err.Error(), config.MaxRetry)
			}
			time.Sleep(1 * time.Second)
			config.Log.Debugf("retrying to connect to %s", config.JoinAddr)
			err = join(config.JoinAddr, config.RaftAddr, config.NodeID)
			retry ++
		}
	}

	config.Log.Info("hraft started successfully")

	return s
}

func (s *System) New(key string) *Entanglement {
	s.mu.Lock()
	defer s.mu.Unlock()

	e := s.m[key]
	if e == nil {
		e = &Entanglement{
			key:    key,
			data:   "",
			system: s,
			mu:     &sync.Mutex{},
		}
		s.m[key] = e
	}
	return e
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
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("unable to join. status_code=%d", resp.StatusCode))
	}
	defer resp.Body.Close()

	return nil
}
