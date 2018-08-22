package dyncall

import (
	"github.com/kataras/iris/core/errors"
	"github.com/stretchr/testify/require"
	"testing"
)

type testCallModule struct {
	callID   string
	validate func(map[string]interface{}) error
	handle   func(map[string]interface{}) (map[string]interface{}, error)
}

func (m *testCallModule) CallID() string {
	return m.callID
}

func (m *testCallModule) Validate(payload map[string]interface{}) error {
	return m.validate(payload)
}

func (m *testCallModule) Handle(payload map[string]interface{}) (map[string]interface{}, error) {
	return m.handle(payload)
}

func TestRegistry_Register(t *testing.T) {

	callModule := testCallModule{
		callID: "MY:MODULE",
	}

	reg := New()

	// should register without an error
	require.Nil(t, reg.Register(&callModule))

	// should fail since the same module can't be registered twice
	require.EqualError(t, reg.Register(&callModule), "a call modules with id MY:MODULE has already been registered")

}

// make sure call gets validated
func TestRegistry_CallValidate(t *testing.T) {

	callModule := testCallModule{
		callID: "MY:MODULE",
		validate: func(payload map[string]interface{}) error {
			return errors.New("invalid payload")
		},
	}

	reg := New()
	require.Nil(t, reg.Register(&callModule))

	// call should be validated
	_, err := reg.Call("MY:MODULE", map[string]interface{}{})
	require.EqualError(t, err, "invalid payload")

}

func TestRegistry_Call(t *testing.T) {

	callModule := testCallModule{
		callID: "MY:MODULE",
		validate: func(payload map[string]interface{}) error {
			return nil
		},
		handle: func(payload map[string]interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{
				"key": "value",
			}, nil
		},
	}

	reg := New()
	require.Nil(t, reg.Register(&callModule))

	// call should be validated
	result, err := reg.Call("MY:MODULE", map[string]interface{}{})
	require.Nil(t, err)
	require.Equal(t, map[string]interface{}{"key": "value"}, result)

}
