package middlewares

import (
	"api-server/validators"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// getFieldMessage retrieves the custom error message from the 'message' struct tag.
func getFieldMessage(obj interface{}, fieldName string, tag string) string {
	// Get the reflect.Type of the struct
	rt := reflect.TypeOf(obj)

	// Loop through the fields to find the one that matches fieldName
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.Name == fieldName {
			// Return the value of the custom 'message' tag
			return field.Tag.Get(tag)
		}
	}
	return ""
}

// ErrorHandler is a middleware to handle validation and binding errors.
func ErrorHandler(c *gin.Context) {
	c.Next() // Process the request

	// Check if there were any errors in the context
	if len(c.Errors) > 0 {
		for _, err := range c.Errors {
			if err.Type == gin.ErrorTypeBind {
				handleBindingError(c, err)
				return
			}
		}
	}
}

func handleBindingError(c *gin.Context, err *gin.Error) {
	// Check if the error is a validation error
	if validationErrs, ok := err.Err.(validator.ValidationErrors); ok {
		handleValidationErrors(c, err, validationErrs)
		return
	}

	// For generic binding errors
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  http.StatusBadRequest,
		"message": "Invalid request body",
		"error":   err.Error(),
	})
}

func handleValidationErrors(c *gin.Context, err *gin.Error, validationErrs validator.ValidationErrors) {
	// Create a map to hold validation error details
	errors := make(map[string]string)

	// Identify the validator struct dynamically
	obj := identifyValidatorStruct(err)

	for _, fieldErr := range validationErrs {
		field := fieldErr.StructField() // Field name (e.g., "Username")

		// Try to get a custom message from the 'message' tag
		customMessage := getFieldMessage(obj, field, "message")

		if customMessage != "" {
			errors[fieldErr.Field()] = customMessage
		} else {
			// Fallback to the default validation error message
			errors[fieldErr.Field()] = fieldErr.Error()
		}
	}

	// Send a detailed validation error response
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  http.StatusBadRequest,
		"message": "Validation failed",
		"errors":  errors,
	})
}

// identifyValidatorStruct dynamically resolves the validator struct using ValidatorRegistry.
func identifyValidatorStruct(err *gin.Error) interface{} {
	for validatorName, obj := range validators.ValidatorRegistry {
		if containsValidatorName(err.Error(), validatorName) {
			return obj
		}
	}
	return nil
}

// containsValidatorName checks if the error string contains the validator name.
func containsValidatorName(errorStr string, validatorName string) bool {
	return strings.Contains(errorStr, validatorName)
}
