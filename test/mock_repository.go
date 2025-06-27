package mocks

import (
	"errors"
	"reflect"
	"sync"
)

// MockRepository implements the ports.Repository interface for testing purposes.
// It uses in-memory storage and allows configuring predefined responses and errors.
// The mock is generic and works with any entity type T and comparable ID type ID.
type MockRepository[T any, ID comparable] struct {
	// entities stores the mock data
	entities map[ID]*T
	// mu ensures thread-safety for concurrent access
	mu sync.RWMutex
	// calls tracks method invocations for verification
	calls map[string][]interface{}
	// responses stores predefined responses for each method
	responses map[string]response
}

// response defines the expected response for a mocked method
type response struct {
	value interface{}
	err   error
}

// NewMockRepository creates a new instance of MockRepository
func NewMockRepository[T any, ID comparable]() *MockRepository[T, ID] {
	return &MockRepository[T, ID]{
		entities:  make(map[ID]*T),
		calls:     make(map[string][]interface{}),
		responses: make(map[string]response),
	}
}

// recordCall logs a method call with its arguments
func (m *MockRepository[T, ID]) recordCall(method string, args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls[method] = append(m.calls[method], args)
}

// setResponse sets a predefined response for a method
func (m *MockRepository[T, ID]) SetResponse(method string, value interface{}, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses[method] = response{value: value, err: err}
}

// Calls returns the recorded calls for a given method
func (m *MockRepository[T, ID]) Calls(method string) []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.calls[method]
}

// Create persists a new entity in the mock repository
func (m *MockRepository[T, ID]) Create(entity *T) (*T, error) {
	m.recordCall("Create", entity)

	// Check for predefined response
	m.mu.RLock()
	resp, exists := m.responses["Create"]
	m.mu.RUnlock()
	if exists {
		if resp.err != nil {
			return nil, resp.err
		}
		if val, ok := resp.value.(*T); ok {
			return val, nil
		}
		return nil, errors.New("invalid response type for Create")
	}

	// Default behavior: store the entity
	m.mu.Lock()
	defer m.mu.Unlock()

	// Extract ID using reflection (assuming ID is a field named "ID")
	v := reflect.ValueOf(*entity)
	idField := v.FieldByName("ID")
	if !idField.IsValid() {
		return nil, errors.New("entity must have an ID field")
	}
	id, ok := idField.Interface().(ID)
	if !ok {
		return nil, errors.New("ID field type does not match repository ID type")
	}

	if _, exists := m.entities[id]; exists {
		return nil, errors.New("entity already exists")
	}
	m.entities[id] = entity
	return entity, nil
}

// GetByID retrieves an entity by its ID
func (m *MockRepository[T, ID]) GetByID(id ID) (*T, error) {
	m.recordCall("GetByID", id)

	// Check for predefined response
	m.mu.RLock()
	resp, exists := m.responses["GetByID"]
	m.mu.RUnlock()
	if exists {
		if resp.err != nil {
			return nil, resp.err
		}
		if val, ok := resp.value.(*T); ok {
			return val, nil
		}
		return nil, errors.New("invalid response type for GetByID")
	}

	// Default behavior: retrieve from entities
	m.mu.RLock()
	defer m.mu.RUnlock()
	entity, exists := m.entities[id]
	if !exists {
		return nil, errors.New("entity not found")
	}
	return entity, nil
}

// GetAll retrieves all entities
func (m *MockRepository[T, ID]) GetAll() ([]T, error) {
	m.recordCall("GetAll")

	// Check for predefined response
	m.mu.RLock()
	resp, exists := m.responses["GetAll"]
	m.mu.RUnlock()
	if exists {
		if resp.err != nil {
			return nil, resp.err
		}
		if val, ok := resp.value.([]T); ok {
			return val, nil
		}
		return nil, errors.New("invalid response type for GetAll")
	}

	// Default behavior: return all entities
	m.mu.RLock()
	defer m.mu.RUnlock()
	entities := make([]T, 0, len(m.entities))
	for _, entity := range m.entities {
		entities = append(entities, *entity)
	}
	return entities, nil
}

// Update modifies an existing entity
func (m *MockRepository[T, ID]) Update(entity *T) (*T, error) {
	m.recordCall("Update", entity)

	// Check for predefined response
	m.mu.RLock()
	resp, exists := m.responses["Update"]
	m.mu.RUnlock()
	if exists {
		if resp.err != nil {
			return nil, resp.err
		}
		if val, ok := resp.value.(*T); ok {
			return val, nil
		}
		return nil, errors.New("invalid response type for Update")
	}

	// Default behavior: update the entity
	m.mu.Lock()
	defer m.mu.Unlock()

	v := reflect.ValueOf(*entity)
	idField := v.FieldByName("ID")
	if !idField.IsValid() {
		return nil, errors.New("entity must have an ID field")
	}
	id, ok := idField.Interface().(ID)
	if !ok {
		return nil, errors.New("ID field type does not match repository ID type")
	}

	if _, exists := m.entities[id]; !exists {
		return nil, errors.New("entity not found")
	}
	m.entities[id] = entity
	return entity, nil
}

// Delete removes an entity by its ID
func (m *MockRepository[T, ID]) Delete(id ID) error {
	m.recordCall("Delete", id)

	// Check for predefined response
	m.mu.RLock()
	resp, exists := m.responses["Delete"]
	m.mu.RUnlock()
	if exists {
		return resp.err
	}

	// Default behavior: delete the entity
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.entities[id]; !exists {
		return errors.New("entity not found")
	}
	delete(m.entities, id)
	return nil
}

// FindByField retrieves the first entity matching the field-value pair
func (m *MockRepository[T, ID]) FindByField(field string, value any) (*T, error) {
	m.recordCall("FindByField", field, value)

	// Check for predefined response
	m.mu.RLock()
	resp, exists := m.responses["FindByField"]
	m.mu.RUnlock()
	if exists {
		if resp.err != nil {
			return nil, resp.err
		}
		if val, ok := resp.value.(*T); ok {
			return val, nil
		}
		return nil, errors.New("invalid response type for FindByField")
	}

	// Default behavior: search using reflection
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, entity := range m.entities {
		v := reflect.ValueOf(*entity)
		f := v.FieldByName(field)
		if !f.IsValid() {
			return nil, errors.New("invalid field name")
		}
		if f.Interface() == value {
			return entity, nil
		}
	}
	return nil, errors.New("entity not found")
}

// Query executes a custom query and returns the result as interface{}.
// This method provides a flexible way to execute custom queries that don't fit the standard CRUD operations.
func (m *MockRepository[T, ID]) Query(query string, args ...interface{}) (interface{}, error) {
	m.recordCall("Query", append([]interface{}{query}, args...)...)

	// Check for predefined response
	m.mu.RLock()
	resp, exists := m.responses["Query"]
	m.mu.RUnlock()

	if exists {
		if resp.err != nil {
			return nil, resp.err
		}
		return resp.value, nil
	}

	// Default behavior: return nil
	return nil, nil
}
