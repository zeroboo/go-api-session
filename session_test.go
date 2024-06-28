package apisession

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// go test -timeout 30s -run ^TestValidateSession_ValidCall_NoError$ github.com/zeroboo/go-api-session
func TestValidateSession_ValidCall_NoError(t *testing.T) {
	owner := "user1"
	manager := CreateNewRedisSessionManager(redisClient, sessionPrefix, 1000, 10000, 10, 5)

	sessionId, _ := manager.CreateNewSession(context.TODO(), owner)
	sessionIds = append(sessionIds, sessionId)

	session, _ := manager.GetSession(context.TODO(), sessionId)
	now := time.Now().Unix()

	errValidate := manager.ValidateCall(owner, "url1", session, now)
	assert.Nil(t, errValidate, "Valid owner, no error")

}

// go test -timeout 30s -run ^TestValidateSession_InvalidOwner_Error$ github.com/zeroboo/go-api-session
func TestValidateSession_InvalidOwner_Error(t *testing.T) {
	owner := "user1"
	manager := CreateNewRedisSessionManager(redisClient, sessionPrefix, 1000, 10000, 10, 5)
	sessionId, _ := manager.CreateNewSession(context.TODO(), owner)
	sessionIds = append(sessionIds, sessionId)

	session, _ := manager.GetSession(context.TODO(), sessionId)
	now := time.Now().Unix()

	errValidate := manager.ValidateCall("nobody", "url1", session, now)
	assert.Equal(t, ErrInvalidSession, errValidate, "Invalid owner, error")

}

// go test -timeout 30s -run ^TestValidateSession_TooFast_Error$ github.com/zeroboo/go-api-session
func TestValidateSession_TooFast_Error(t *testing.T) {
	owner := "user1"
	manager := CreateNewRedisSessionManager(redisClient, sessionPrefix, 1000, 10000, 3, 10)
	sessionId, _ := manager.CreateNewSession(context.TODO(), owner)
	session, _ := manager.GetSession(context.TODO(), sessionId)
	sessionIds = append(sessionIds, sessionId)

	now := time.Now().Unix()
	var errValidate error

	errValidate = manager.ValidateCall(owner, "url1", session, now)
	assert.Nil(t, errValidate, "Valid call, no error")

	errValidate = manager.ValidateCall(owner, "url1", session, now)
	assert.Equal(t, ErrTooFast, errValidate, "Too fast call, has error")
}

// go test -timeout 30s -run ^TestValidateSession_TooFrequently_Error$ github.com/zeroboo/go-api-session
func TestValidateSession_TooFrequently_Error(t *testing.T) {
	owner := "user1"
	manager := CreateNewRedisSessionManager(redisClient, sessionPrefix, 1000, 10000, 2, 0)
	sessionId, _ := manager.CreateNewSession(context.TODO(), owner)
	session, _ := manager.GetSession(context.TODO(), sessionId)
	sessionIds = append(sessionIds, sessionId)
	now := time.Now().Unix()
	var errValidate error

	errValidate = manager.ValidateCall(owner, "url1", session, now)
	assert.Nil(t, errValidate, "Valid call, no error")

	errValidate = manager.ValidateCall(owner, "url1", session, now)
	assert.Equal(t, nil, errValidate, "second call, no error")

	errValidate = manager.ValidateCall(owner, "url1", session, now)
	assert.Equal(t, ErrTooMany, errValidate, "third call is too frequently, has error")
}

// go test -timeout 30s -run ^TestValidateSession_NewWindow_Correct$ github.com/zeroboo/go-api-session
func TestValidateSession_NewWindow_Correct(t *testing.T) {
	owner := "user1"
	interval := int64(10)
	manager := CreateNewRedisSessionManager(redisClient, sessionPrefix, 1000, 10000, 2, interval)
	sessionId, _ := manager.CreateNewSession(context.TODO(), owner)
	sessionIds = append(sessionIds, sessionId)
	session, _ := manager.GetSession(context.TODO(), sessionId)

	now := time.Now().Unix()
	var errValidate error

	errValidate = manager.ValidateCall(owner, "url1", session, now)
	assert.Nil(t, errValidate, "Valid call, no error")

	errValidate = manager.ValidateCall(owner, "url1", session, now+interval+1)
	assert.Equal(t, nil, errValidate, "second call, no error")

	errValidate = manager.ValidateCall(owner, "url1", session, now+2*interval+1)
	assert.Equal(t, ErrTooMany, errValidate, "third call is too frequently, has error")

	session.SetWindow(session.Window + 1)
	errValidate = manager.ValidateCall(owner, "url1", session, now+2*interval+1)
	assert.Equal(t, nil, errValidate, "new windows, request valid")
}
