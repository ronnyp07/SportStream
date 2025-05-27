package ccvalidation

// // Validator defines the interface for validation
// type Validator[T any] interface {
// 	Validate(src T) Result
// }

// // validator implements the Validator interface
// type validator[T any] struct {
// 	breakOnFailure bool
// 	validators
// }

// // Validate executes all validation steps
// func (v *validator[T]) Validate(src T) Result {
// 	result := Result{}
// 	if len(v.validators) == 0 {
// 		result.AddErrorMessage("no validator step defined")
// 		return result
// 	}

// 	for _, step := range v.validators {
// 		err := step.validator(src)

// 		if err != nil {
// 			result.AddError(err)
// 			if step.breakOnFailure || v.breakOnFailure {
// 				return result
// 			}
// 		}
// 	}
// 	return result
// }

// // NewValidator creates a new validator
// func NewValidator[T any]() *validator[T] {
// 	return &validator[T]{
// 		breakOnFailure: false,
// 		validators:     make([]validationStep[T], 0),
// 	}
// }

// // WithBreakOnFailure configures the validator to break on first failure
// func (v *validator[T]) WithBreakOnFailure() *validator[T] {
// 	v.breakOnFailure = true
// 	return v
// }

// // AddValidator adds a new validation step
// func (v *validator[T]) AddValidator(validator Validator[T]) *validator[T] {
// 	step := validationStep[T]{
// 		validator:      validator,
// 		breakOnFailure: false,
// 	}
// 	v.validators = append(v.validators, step)
// 	return v
// }

// // conditionalValidator validates based on a condition
// type conditionalValidator[Env any, T any] struct {
// 	validator Validator[T]
// 	condition func(Env) bool
// }

// // NewConditionalValidator creates a new conditional validator
// func NewConditionalValidator[Env any, T any](validator Validator[T]) *conditionalValidator[Env, T] {
// 	return &conditionalValidator[Env, T]{
// 		validator: validator,
// 	}
// }

// // WithCondition sets the condition function
// func (c *conditionalValidator[Env, T]) WithCondition(condition func(Env) bool) *conditionalValidator[Env, T] {
// 	c.condition = condition
// 	return c
// }

// Validate implements the Validator interface for conditional validation
// func (c *conditionalValidator[Env, T]) Validate(env Env, src T) Result {
// 	if c.condition != nil && !c.condition(env) {
// 		return NewResult() // Skip validation if condition is not met
// 	}
// 	return c.validator.Validate(src)
// }
