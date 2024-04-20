package shared

import (
	"encoding/json"
	"io"
	"os"
	"sync"
)

type JsonKvStore struct {
	path   string
	values map[string]interface{}
	mutex  sync.RWMutex
}

func NewJsonKvStore(path string) (*JsonKvStore, error) {
	js := &JsonKvStore{
		path:   path,
		values: make(map[string]interface{}),
	}

	if err := js.Load(); err != nil {
		return js, err
	}
	return js, nil
}

func (js *JsonKvStore) Load() error {
	js.mutex.Lock()
	defer js.mutex.Unlock()

	file, err := os.Open(js.path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &js.values); err != nil {
		return err
	}

	return nil
}

func (js *JsonKvStore) Get(key string) interface{} {
	js.mutex.RLock()
	defer js.mutex.RUnlock()
	return js.values[key]
}

func (js *JsonKvStore) Set(key string, value interface{}) {
	js.mutex.Lock()
	defer js.mutex.Unlock()
	js.values[key] = value
}

func (js *JsonKvStore) Save() error {
	js.mutex.Lock()
	defer js.mutex.Unlock()

	data, err := json.MarshalIndent(js.values, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(js.path, data, 0644)
}
