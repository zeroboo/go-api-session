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

	// StartSession creates a new session for the owner and insert to db
	// Returns:
	//   - sessionId string: id of new session
	//   - error: error if exists, nil is successful
	StartSession(ctx context.Context, owner string) (string, error)

	// StartSession creates a new session for the owner and insert to db
	// Returns:
	//   - session *APISession: new session with payload
	//   - error: error if exists, nil is successful
	StartSessionWithPayload(ctx context.Context, owner string, payload map[string]any) (*APISession, error)

	//GetRequestInterval returns the interval in milliseconds for API requests
	GetRequestInterval() int64

	// GetMaxCallPerWindow returns the maximum number of API calls allowed per time window
	GetMaxCallPerWindow() int64

	// GetWindowSize returns the size of the time window in milliseconds
	GetWindowSize() int64

	// GetOnlineUsers returns a map of online users with their last activity timestamp
	GetOnlineUsers(ctx context.Context) (map[string]int64, error)
}
