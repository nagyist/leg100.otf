/*
Package http provides an HTTP interface allowing HTTP clients to interact with OTF.
*/
package http

import "github.com/gorilla/schema"

var (
	// Query schema Encoder, caches structs, and safe for sharing
	encoder = schema.NewEncoder()
)
