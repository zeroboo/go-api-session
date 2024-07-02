package apisession

import "context"

type ISessionManager interface {
	RecordAPICall(ctx context.Context, sessionId string, owner string, url string) (*APISession, error)
	GetSession(ctx context.Context, owner string) (*APISession, error)
	DeleteSession(ctx context.Context, owner string) error
	SetSession(ctx context.Context, owner string, session *APISession) error
	StartSession(ctx context.Context, owner string) (string, error)
}
