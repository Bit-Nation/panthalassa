package api

// This package provider a kind of api to the client.
// In our case the client is our pangea app

// NOTICE! The "device" package is deprecated

// The API struct provides a set of utils to make requests to the api
// If you plan to create a new api call you need to add you messages
// to the request and response protobufs
// you can then extend the api struct
// if you e.g. would like to implement the DHT CRUD you can create a new file
// called `dht.go` in the api folder and start to implement the requests