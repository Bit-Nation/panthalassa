package contacts

import (
	"errors"
)

type ContactCreateCall struct {
	s *Storage
}

func NewContactCreateCall(s *Storage) *ContactCreateCall {
	return &ContactCreateCall{
		s: s,
	}
}

func (c *ContactCreateCall) CallID() string {
	return "CONTACT:CREATE"
}

func (c *ContactCreateCall) Validate(map[string]interface{}) error {
	return nil
}

func (c *ContactCreateCall) Handle(data map[string]interface{}) (map[string]interface{}, error) {

	// get title
	identityKey, k := data["identity_key"].(string)
	if !k {
		return map[string]interface{}{}, errors.New("identity_key must be astring")
	}

	return map[string]interface{}{}, c.s.Save(&Contact{
		IdentityKey: identityKey,
	})

}

type ContactAllCall struct {
	s *Storage
}

func NewContactAllCall(db *Storage) *ContactAllCall {
	return &ContactAllCall{
		s: db,
	}
}

func (c *ContactAllCall) CallID() string {
	return "CONTACT:ALL"
}

func (c *ContactAllCall) Validate(map[string]interface{}) error {
	return nil
}

func (c *ContactAllCall) Handle(data map[string]interface{}) (map[string]interface{}, error) {

	contacts, err := c.s.All()
	if err != nil {
		return map[string]interface{}{}, nil
	}

	jsonContacts := []map[string]interface{}{}
	for _, c := range contacts {
		jsonContacts = append(jsonContacts, map[string]interface{}{
			"identity_key": c.IdentityKey,
		})
	}
	return map[string]interface{}{
		"contacts": jsonContacts,
	}, nil
}
