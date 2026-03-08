package storage

import (
	"fmt"
	"maps"
	"slices"
	"sync"
)

type StorageEntry struct {
	Storage *Storage
	Alias   string
}

type StorageEntryConfig struct {
	Name    string
	File    string
	Alias   string
	Schemas []string
}

type StorageAggregator struct {
	storages map[string]*StorageEntry
	mu       sync.RWMutex
}

func NewAggregator() *StorageAggregator {
	return &StorageAggregator{make(map[string]*StorageEntry), sync.RWMutex{}}
}

func (sa *StorageAggregator) NewStorage(cfg StorageEntryConfig) (*StorageEntry, error) {
	sa.mu.Lock()
	defer sa.mu.Unlock()

	st, err := NewDB(cfg.File)
	if err != nil {
		return nil, err
	}

	if err := st.InitSchemas(cfg.Schemas...); err != nil {
		return nil, err
	}

	stEntry := &StorageEntry{st, cfg.Alias}

	sa.storages[cfg.Name] = stEntry

	return stEntry, nil
}

func (sa *StorageAggregator) GetEntry(name string) (*StorageEntry, error) {
	se, ok := sa.storages[name]
	if !ok {
		return nil, fmt.Errorf("StorageEntry not found with name: %s", name)
	}

	return se, nil
}

func (sa *StorageAggregator) GetEntries() []*StorageEntry {
	return slices.Collect(maps.Values(sa.storages))
}

func (sa *StorageAggregator) CloseAll() error {
	for _, se := range sa.storages {
		if err := se.Storage.Close(); err != nil {
			return err
		}
	}

	return nil
}
