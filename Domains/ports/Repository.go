package ports

// Repository is a generic interface that defines the standard CRUD operations for domain entities.
// It uses Go generics to work with any entity type (T) and any comparable ID type (ID).
//
// Type parameters:
//   - T:  The domain entity type this repository manages
//   - ID: The type of the unique identifier for the entity (must be comparable)
//
// Implementations of this interface should handle data persistence and retrieval
// while abstracting away the underlying storage details.
type Repository[T any, ID comparable] interface {
	// Create persists a new entity in the repository.
	// This method is responsible for creating a new record in the underlying storage.
	//
	// Parameters:
	//   - entity: A pointer to the entity to be created
	//
	// Returns:
	//   - *T:    A pointer to the created entity with any system-generated fields populated
	//   - error: Error if the operation fails (e.g., validation error, storage error)
	//
	// Note: The entity parameter will be modified with any generated fields (like ID, timestamps, etc.)
	Create(entity *T) (*T, error)

	// GetByID retrieves an entity by its unique identifier.
	// This is the primary method for fetching a single entity when you have its ID.
	//
	// Parameters:
	//   - id: The unique identifier of the entity to retrieve
	//
	// Returns:
	//   - *T:    A pointer to the found entity, or nil if not found
	//   - error: Error if the operation fails (e.g., invalid ID format, database error)
	GetByID(id ID) (*T, error)

	// GetAll retrieves all entities of type T from the repository.
	// Use with caution on large datasets as it loads all records into memory.
	//
	// Returns:
	//   - []T:  A slice containing all entities, or an empty slice if none exist
	//   - error: Error if the operation fails (e.g., database connection error)
	//
	// Note: For large datasets, consider implementing pagination or streaming
	GetAll() ([]T, error)

	// Update modifies an existing entity in the repository.
	// The entity must have a valid ID that exists in the repository.
	//
	// Parameters:
	//   - entity: A pointer to the entity with updated fields
	//
	// Returns:
	//   - *T:    A pointer to the updated entity
	//   - error: Error if the operation fails (e.g., entity not found, validation error)
	//
	// Note: Implementations should perform optimistic locking if concurrent updates are a concern
	Update(entity *T) (*T, error)

	// Delete removes an entity from the repository by its ID.
	// This operation is permanent and cannot be undone.
	//
	// Parameters:
	//   - id: The unique identifier of the entity to delete
	//
	// Returns:
	//   - error: Error if the operation fails (e.g., entity not found, permission denied)
	//
	// Note: Some implementations might implement soft-delete instead of physical deletion
	Delete(id ID) error

	// FindByField retrieves the first entity that matches the given field-value pair.
	// This method is useful for looking up entities by non-primary key fields.
	//
	// Parameters:
	//   - field: The name of the field to search by (must be a valid field of type T)
	//   - value: The value to match against the specified field
	//
	// Returns:
	//   - *T:    A pointer to the found entity, or nil if no match is found
	//   - error: Error if the operation fails (e.g., invalid field name, database error)
	//
	// Note: The field name is case-sensitive and must match the struct field name exactly.
	// Only the first matching entity is returned. If multiple matches exist, consider
	// implementing a separate method to handle multiple results.
	FindByField(field string, value any) (*T, error)
}
