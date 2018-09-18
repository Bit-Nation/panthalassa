package documents

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	keyManager "github.com/Bit-Nation/panthalassa/keyManager"
	"github.com/ethereum/go-ethereum/common"
	cid "github.com/ipfs/go-cid"
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
			"hash":        d.CID,
			"signature":   d.Signature,
			"tx_hash":     d.TransactionHash,
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
	if err := d.s.db.One("ID", int(docID), &doc); err != nil {
		return map[string]interface{}{}, err
	}

	return map[string]interface{}{}, d.s.db.DeleteStruct(&doc)

}

type DocumentSubmitCall struct {
	s          *Storage
	km         *keyManager.KeyManager
	n          *NotaryMulti
	notaryAddr common.Address
}

func NewDocumentNotariseCall(s *Storage, km *keyManager.KeyManager, n *NotaryMulti, notaryAddr common.Address) *DocumentSubmitCall {
	return &DocumentSubmitCall{
		s:          s,
		km:         km,
		n:          n,
		notaryAddr: notaryAddr,
	}
}

func (d *DocumentSubmitCall) CallID() string {
	return "DOCUMENT:NOTARISE"
}

func (d *DocumentSubmitCall) Validate(map[string]interface{}) error {
	return nil
}

func (d *DocumentSubmitCall) Handle(data map[string]interface{}) (map[string]interface{}, error) {

	docID, k := data["doc_id"].(float64)
	if !k {
		return map[string]interface{}{}, errors.New("expect doc_id to be an integer")
	}

	var doc Document
	if err := d.s.db.One("ID", int(docID), &doc); err != nil {
		return map[string]interface{}{}, err
	}

	// prepare request data
	reqData := map[string][]byte{
		"file": doc.Content,
	}
	rawReqData, err := json.Marshal(reqData)
	if err != nil {
		return map[string]interface{}{}, err
	}

	// upload to ipfs
	req, err := http.NewRequest("POST", "https://ipfs.infura.io:5001/api/v0/add", bytes.NewBuffer(rawReqData))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return map[string]interface{}{}, err
	}

	// exec request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{}, err
	}

	// make sure we got a valid status code back
	if resp.Status != "200" {
		return map[string]interface{}{}, errors.New("invalid status: " + resp.Status)
	}

	// read response
	var rawResp []byte
	if _, err := resp.Body.Read(rawResp); err != nil {
		return map[string]interface{}{}, err
	}

	// unmarshal response
	respMap := map[string]string{}
	if err := json.Unmarshal(rawResp, &respMap); err != nil {
		return map[string]interface{}{}, err
	}

	// make sure hash exist in response
	strCid, exist := respMap["Hash"]
	if !exist {
		return map[string]interface{}{}, errors.New("hash doesn't exist in response")
	}

	// cast string hash to CID
	c, err := cid.Cast([]byte(strCid))
	if err != nil {
		return map[string]interface{}{}, err
	}

	// attach cid to document
	doc.CID = c.Bytes()

	// sign cid
	cidSignature, err := d.km.IdentitySign(c.Bytes())
	if err != nil {
		return map[string]interface{}{}, err
	}

	// attach signature to doc
	doc.Signature = cidSignature

	// submit tx to chain
	tx, err := d.n.NotarizeTwo(nil, d.notaryAddr, c.Bytes(), cidSignature)
	if err != nil {
		return map[string]interface{}{}, err
	}

	doc.TransactionHash = tx.Hash().Hex()

	return map[string]interface{}{}, d.s.Save(&doc)

}
