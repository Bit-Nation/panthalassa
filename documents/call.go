package documents

import (
	"encoding/base64"
	"errors"
	"time"
)

type DocumentCreateCall struct {
	s *Storage
}

func NewDocumentCreateCall(s *Storage) *DocumentCreateCall {
	return &DocumentCreateCall{
		s: s,
	}
}

func (c *DocumentCreateCall) CallID() string {
	return "DOCUMENT:CREATE"
}

func (c *DocumentCreateCall) Validate(map[string]interface{}) error {
	return nil
}

func (c *DocumentCreateCall) Handle(data map[string]interface{}) (map[string]interface{}, error) {

	// get title
	title, k := data["title"].(string)
	if !k {
		return map[string]interface{}{}, errors.New("title must be a string")
	}

	// mime type
	mimeType, k := data["mime_type"].(string)
	if !k {
		return map[string]interface{}{}, errors.New("mime type must be a string")
	}

	// content
	content, k := data["content"].(string)
	if !k {
		return map[string]interface{}{}, errors.New("content type must be a string")
	}
	contentBuff, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return map[string]interface{}{}, err
	}

	// description
	description, k := data["description"].(string)
	if !k {
		return map[string]interface{}{}, errors.New("description type must be a string")
	}

	return map[string]interface{}{}, c.s.Save(&Document{
		Title:       title,
		Content:     contentBuff,
		Description: description,
		MimeType:    mimeType,
		CreatedAt:   time.Now().Unix(),
	})

}

type DocumentAllCall struct {
	s *Storage
}

func NewDocumentAllCall(db *Storage) *DocumentAllCall {
	return &DocumentAllCall{
		s: db,
	}
}

func (c *DocumentAllCall) CallID() string {
	return "DOCUMENT:ALL"
}

func (c *DocumentAllCall) Validate(map[string]interface{}) error {
	return nil
}

func (c *DocumentAllCall) Handle(data map[string]interface{}) (map[string]interface{}, error) {

	docs, err := c.s.All()
	if err != nil {
		return map[string]interface{}{}, nil
	}

	jsonDocs := []map[string]interface{}{}
	for _, d := range docs {
		jsonDocs = append(jsonDocs, map[string]interface{}{
			"id":          d.ID,
			"content":     d.Content,
			"mime_type":   d.MimeType,
			"description": d.Description,
			"title":       d.Title,
		})
	}
	return map[string]interface{}{
		"docs": jsonDocs,
	}, nil
}

type DocumentUpdateCall struct {
	s *Storage
}

func NewDocumentUpdateCall(s *Storage) *DocumentUpdateCall {
	return &DocumentUpdateCall{
		s: s,
	}
}

func (c *DocumentUpdateCall) CallID() string {
	return "DOCUMENT:UPDATE"
}

func (c *DocumentUpdateCall) Validate(map[string]interface{}) error {
	return nil
}

func (c *DocumentUpdateCall) Handle(data map[string]interface{}) (map[string]interface{}, error) {

	docID, k := data["doc_id"].(float64)
	if !k {
		return map[string]interface{}{}, errors.New("expect doc_id to be an integer")
	}

	title, k := data["title"].(string)
	if !k {
		return map[string]interface{}{}, errors.New("expect title to be an string")
	}

	description, k := data["description"].(string)
	if !k {
		return map[string]interface{}{}, errors.New("expect description to be an string")
	}

	var doc Document
	if err := c.s.db.One("ID", int(docID), &doc); err != nil {
		return map[string]interface{}{}, err
	}

	doc.Title = title
	doc.Description = description

	return map[string]interface{}{}, c.s.db.Update(&doc)
}

type DocumentDeleteCall struct {
	s *Storage
}

func NewDocumentDeleteCall(s *Storage) *DocumentDeleteCall {
	return &DocumentDeleteCall{
		s: s,
	}
}

func (d *DocumentDeleteCall) CallID() string {
	return "DOCUMENT:DELETE"
}

func (d *DocumentDeleteCall) Validate(map[string]interface{}) error {
	return nil
}

func (d *DocumentDeleteCall) Handle(data map[string]interface{}) (map[string]interface{}, error) {

	docID, k := data["doc_id"].(float64)
	if !k {
		return map[string]interface{}{}, errors.New("expect doc_id to be an integer")
	}

	var doc Document
	if err := d.s.db.Find("ID", int(docID), &doc); err != nil {
		return map[string]interface{}{}, err
	}

	return map[string]interface{}{}, d.s.db.DeleteStruct(&doc)

}
