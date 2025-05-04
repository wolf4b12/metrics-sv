package storage

import (

    "sync"


)

// KVStorage — базовое хранилище ключ-значение
type KVStorage struct {
    data map[string]any
    mu   sync.RWMutex
}

// NewKVStorage создаёт новое KV-хранилище
func NewKVStorage() *KVStorage {
    return &KVStorage{
        data: make(map[string]any),
    }
}

// Set устанавливает значение по ключу
func (s *KVStorage) Set(key string, value any) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.data[key] = value
}

// Get получает значение по ключу
func (s *KVStorage) Get(key string) (any, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    value, exists := s.data[key]
    return value, exists
}

// Delete удаляет значение по ключу
func (s *KVStorage) Delete(key string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    delete(s.data, key)
}

// All возвращает все данные
func (s *KVStorage) All() map[string]any {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.data
}

