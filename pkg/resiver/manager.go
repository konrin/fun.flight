package resiver

import (
	"errors"
	"sync"
)

type Manager struct {
	mu       sync.RWMutex
	resivers map[string]*ResiversRuntime
}

func (m *Manager) Add(resiver *Resiver) error {
	if _, ok := m.resivers[resiver.Name()]; !ok {
		return errors.New("Resiver is already exist")
	}

	m.mu.Lock()
	m.resivers[resiver.Name()] = NewResiversRuntime(resiver)
	m.mu.Unlock()

	return nil
}

func (m *Manager) StartResiver(resiver *Resiver) error {
	runtime := m.getRuntime(resiver)
	if runtime == nil {
		return errors.New("Resiver is not found")
	}

	runtime.Start()

	return nil
}

func (m *Manager) StartAllResivers() error {
	for _, runtime := range m.resivers {
		err := m.StartResiver(runtime.GetResiver())
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) GetResivers() ([]*Resiver, error) {
	return nil, nil
}

func (m *Manager) getRuntime(resiver *Resiver) *ResiversRuntime {
	runtime, ok := m.resivers[resiver.Name()]
	if !ok {
		return nil
	}

	return runtime
}

type ResiversRuntime struct {
	resiver     *Resiver
	dataHandler *DataHandler
}

func NewResiversRuntime(resiver *Resiver) *ResiversRuntime {
	dataHandler := NewDataHandler(resiver)

	return &ResiversRuntime{
		resiver:     resiver,
		dataHandler: dataHandler,
	}
}

func (rr *ResiversRuntime) Start() error {
	return rr.dataHandler.Start()
}

func (rr *ResiversRuntime) Stop() error {
	return rr.dataHandler.Stop()
}

func (rr *ResiversRuntime) GetResiver() *Resiver {
	return rr.resiver
}
