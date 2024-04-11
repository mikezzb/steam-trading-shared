package shared

import (
	"encoding/json"
	"io"
	"os"
)

type PersisitedStore struct {
	path   string
	values map[string]interface{}
}

func NewPersisitedStore(path string) (*PersisitedStore, error) {
	ps := &PersisitedStore{
		path:   path,
		values: make(map[string]interface{}),
	}

	if err := ps.Load(); err != nil {
		return nil, err
	}
	return ps, nil
}

func (ps *PersisitedStore) Load() error {
	file, err := os.Open(ps.path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &ps.values); err != nil {
		return err
	}

	return nil
}

func (ps *PersisitedStore) Get(key string) interface{} {
	return ps.values[key]
}

func (ps *PersisitedStore) Set(key string, value interface{}) {
	ps.values[key] = value
}

func (ps *PersisitedStore) Save() error {
	data, err := json.MarshalIndent(ps.values, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ps.path, data, 0644)
}
