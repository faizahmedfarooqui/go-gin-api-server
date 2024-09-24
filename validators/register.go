//go:generate go run ../cmd/generate_validators/main.go

package validators

var ValidatorRegistry = map[string]interface{}{}

// RegisterValidator registers a validator by name dynamically.
func RegisterValidator(name string, v interface{}) {
	ValidatorRegistry[name] = v
}
