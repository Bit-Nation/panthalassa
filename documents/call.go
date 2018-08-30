package documents

import (
	"errors"
	"time"
	"encoding/base64"
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
	
	return map[string]interface{}{}, c.s.Save(&Document{
		Title: title,
		Content: contentBuff,
		MimeType: mimeType,
		CreatedAt: time.Now().Unix(),
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
	return "DOCUMENT:All"
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
			"id": d.ID,
			"content": d.Content,
			"mime_type": d.MimeType,
			"title": d.Title,
		})
	}
	return map[string]interface{}{
		"docs": jsonDocs,
	}, nil
}