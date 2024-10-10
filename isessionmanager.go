package apisession

import (
	"context"
	"time"
)

type ISessionManager interface {
	// Records an API call: load session, update the session with the API call, validate it and save it back if api call is valid.
	//
	// Validating is performed by calling ValidateAPICall
	//
	// Params:
	//   - sessionId string: session id
	//   - owner string: owner of the session
	//   - url: url of the api call
	//
	// Returns:
	//   - session: updated session
	//   - error: nil if success, an error instance if any
	RecordAPICall(ctx context.Context, sessionId string, owner string, url string) (*APISession, error)
	// Validates an API call to a session. Only validating, doesn't perform any update to database.
	//
	// Params:
	//   - request *APIRequest: api call
	//   - session *APISession
	//   - now: current time in milliseconds
	//
	// Returns:
	//   - error: nil if success, an error instance if any
	ValidateAPICall(request *APIRequest, session *APISession, now time.Time) error

	// Loads session of a user from database
	GetSession(ctx context.Context, owner string) (*APISession, error)
	// Deletes session of a user from database
	DeleteSession(ctx context.Context, owner string) error
	// Saves session of a user to database
	SetSession(ctx context.Context, owner string, session *APISession) error

	StartSession(ctx context.Context, owner string) (string, error)

	GetRequestInterval() int64
	GetMaxCallPerWindow() int64
	GetWindowSize() int64
}
