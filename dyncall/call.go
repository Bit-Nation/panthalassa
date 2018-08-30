package dyncall

type CallModule interface {
	// the call id this module relates to
	// e.g. "Messages:Fetch"
	CallID() string
	// validate the given parameters - might return an error when received invalid parameters
	Validate(map[string]interface{}) error
	// will handle the call and return the result
	Handle(map[string]interface{}) (map[string]interface{}, error)
}
