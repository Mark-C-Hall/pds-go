package model

import (
	"time"

	"github.com/bluesky-social/indigo/atproto/syntax"
)

type Account struct {
	DID       syntax.DID
	Handle    syntax.Handle
	Email     string
	CreatedAt time.Time
}
