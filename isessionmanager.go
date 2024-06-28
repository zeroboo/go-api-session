package apisession

import "context"

type ISessionManager interface {
	RecordAPICall(ctx context.Context, sessionId string, owner string, url string) (*APISession, error)
	GetSession(ctx context.Context, sessionId string) (*APISession, error)
	DeleteSession(ctx context.Context, sessionId string) error
	SetSession(ctx context.Context, sessionId string, session *APISession) error
	CreateNewSession(ctx context.Context, owner string) (string, error)
}
