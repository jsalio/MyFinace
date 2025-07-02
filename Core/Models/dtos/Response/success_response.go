package dtos

// SuccessResponse represents a standardized success response structure for API responses.
// It is a generic type that can wrap any data type, providing a consistent response format
// across the application. This helps maintain consistency in API responses and makes it
// easier for clients to handle successful responses.
//
// The generic type parameter TObject specifies the type of data being returned.
//
// Example usage:
//
//	response := SuccessResponse[User]{
//		Message: "User retrieved successfully",
//		Data:    user,
//	}
//
// swagger:model SuccessResponse
// @name SuccessResponse
// @description Standardized success response structure
// @property {string} message - Human-readable message describing the result
// @property {object} data - The actual data being returned (type depends on TObject)
type SuccessResponse[TObject any] struct {
	// Message provides a human-readable description of the operation's result.
	// Example: "Operation completed successfully"
	// Required: true
	Message string `json:"message"`

	// Data contains the actual payload of the response.
	// The structure of this field depends on the generic type parameter TObject.
	// Example: {"id": 1, "name": "John Doe"}
	// Required: true
	Data TObject `json:"data"`
}
