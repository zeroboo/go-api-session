package apisession

import (
	"context"
	"fmt"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// go test -timeout 30s -run ^TestValidateSession_ValidCall_NoError$ github.com/zeroboo/go-api-session -v
func TestValidateSession_ValidCall_NoError(t *testing.T) {
	owner := "user_" + t.Name()
	manager := CreateNewRedisSessionManager(redisClient, sessionPrefix, 1000, 10000, 10, 5)

	sessionId, errNewSession := manager.StartSession(context.TODO(), owner)
	assert.Nil(t, errNewSession, "Create new session, no error")
	sessionOwners = append(sessionOwners, owner)
	fmt.Printf("SessionId: %v\n", sessionId)
	session, errGetSession := manager.GetSession(context.TODO(), owner)
	assert.Nil(t, errGetSession, "Get session, no error")
	now := time.Now()

	errValidate := manager.ValidateAPICall(&APIRequest{
		Owner:     owner,
		SessionId: sessionId,
		URL:       "url1",
	}, session, now)
	assert.Nil(t, errValidate, "Valid owner, no error")

}

// go test -timeout 30s -run ^TestValidateSession_InvalidSessionValue_Error$ github.com/zeroboo/go-api-session
func TestValidateSession_InvalidSessionValue_Error(t *testing.T) {
	owner := "user_" + t.Name()
	manager := CreateNewRedisSessionManager(redisClient, sessionPrefix, 1000, 10000, 10, 5)
	sessionId, _ := manager.StartSession(context.TODO(), owner)
	sessionOwners = append(sessionOwners, owner)

	session, _ := manager.GetSession(context.TODO(), owner)
	now := time.Now()

	errValidate := manager.ValidateAPICall(&APIRequest{
		Owner:     owner,
		SessionId: "invalid_session_id",
		URL:       "url1",
	}, session, now)
	assert.Equal(t, ErrInvalidSession, errValidate, "Invalid owner, error")
	log.Infof("Session: %v", sessionId)
}

// go test -timeout 30s -run ^TestValidateSession_TooFast_Error$ github.com/zeroboo/go-api-session
func TestValidateSession_TooFast_Error(t *testing.T) {
	owner := "user_" + t.Name()
	manager := CreateNewRedisSessionManager(redisClient, sessionPrefix, 1000, 10000, 3, 10)
	sessionId, _ := manager.StartSession(context.TODO(), owner)
	session, _ := manager.GetSession(context.TODO(), owner)
	sessionOwners = append(sessionOwners, owner)

	now := time.Now()
	var errValidate error

	errValidate = manager.ValidateAPICall(&APIRequest{
		Owner:     owner,
		SessionId: sessionId,
		URL:       "url1",
	}, session, now)
	assert.Nil(t, errValidate, "Valid call, no error")

	errValidate = manager.ValidateAPICall(&APIRequest{
		Owner:     owner,
		SessionId: sessionId,
		URL:       "url1",
	}, session, now)
	assert.Equal(t, ErrTooFast, errValidate, "Too fast call, has error")
}

// go test -timeout 30s -run ^TestValidateSession_TooFrequently_Error$ github.com/zeroboo/go-api-session
func TestValidateSession_TooFrequently_Error(t *testing.T) {
	owner := "user_" + t.Name()
	manager := CreateNewRedisSessionManager(redisClient, sessionPrefix, 1000, 10000, 2, 0)
	sessionId, _ := manager.StartSession(context.TODO(), owner)
	session, _ := manager.GetSession(context.TODO(), owner)
	sessionOwners = append(sessionOwners, owner)
	now := time.Now()
	var errValidate error

	errValidate = manager.ValidateAPICall(&APIRequest{
		Owner:     owner,
		SessionId: sessionId,
		URL:       "url1",
	}, session, now)
	assert.Nil(t, errValidate, "Valid call, no error")

	errValidate = manager.ValidateAPICall(&APIRequest{
		Owner:     owner,
		SessionId: sessionId,
		URL:       "url1",
	}, session, now)
	assert.Equal(t, nil, errValidate, "second call, no error")

	errValidate = manager.ValidateAPICall(&APIRequest{
		Owner:     owner,
		SessionId: sessionId,
		URL:       "url1",
	}, session, now)
	assert.Equal(t, ErrTooMany, errValidate, "third call is too frequently, has error")
}

// go test -timeout 30s -run ^TestValidateSession_NewWindow_Correct$ github.com/zeroboo/go-api-session
func TestValidateSession_NewWindow_Correct(t *testing.T) {
	owner := "user_" + t.Name()
	interval := int64(10)
	manager := CreateNewRedisSessionManager(redisClient, sessionPrefix, 1000, 10000, 2, interval)
	sessionId, _ := manager.StartSession(context.TODO(), owner)
	sessionOwners = append(sessionOwners, owner)
	session, _ := manager.GetSession(context.TODO(), owner)

	now := time.Now()
	var errValidate error

	errValidate = manager.ValidateAPICall(&APIRequest{
		Owner:     owner,
		SessionId: sessionId,
		URL:       "url1",
	}, session, now)
	assert.Nil(t, errValidate, "Valid call, no error")

	errValidate = manager.ValidateAPICall(&APIRequest{
		Owner:     owner,
		SessionId: sessionId,
		URL:       "url1",
	}, session, now.Add(time.Duration(interval+1)*time.Millisecond))
	assert.Equal(t, nil, errValidate, "second call, no error")

	errValidate = manager.ValidateAPICall(&APIRequest{
		Owner:     owner,
		SessionId: sessionId,
		URL:       "url1",
	}, session, now.Add(time.Duration(2*interval+1)*time.Millisecond))

	assert.Equal(t, ErrTooMany, errValidate, "third call is too frequently, has error")

	session.SetWindow(session.Window + 1)
	errValidate = manager.ValidateAPICall(&APIRequest{
		Owner:     owner,
		SessionId: sessionId,
		URL:       "url1",
	}, session, now.Add(time.Duration(2*interval+1)*time.Millisecond))
	assert.Equal(t, nil, errValidate, "new windows, request valid")
}

// go test -timeout 30s -run ^TestGetMapFromSession_Correct$ github.com/zeroboo/go-api-session -v
func TestGetMapFromSession_Correct(t *testing.T) {
	session := NewAPISessionWithPayload("user1", map[string]any{})
	mapValue, ok := GetPayloadMap[string, int64](session, "key")
	assert.False(t, ok, "Key not found")
	assert.Nil(t, mapValue, "Value is nil")
	t.Logf("Empty map value: %v", mapValue)

	session.SetPayload("key", map[string]int64{"key1": 1, "key2": 2})
	mapValue, ok = GetPayloadMap[string, int64](session, "key")
	assert.True(t, ok, "Key found")
	assert.NotNil(t, mapValue, "Value is not nil")
	t.Logf("Valid map value: %v", mapValue)
	mapValue["key3"] = 3

	updatedMapValue, ok := GetPayloadMap[string, int64](session, "key")
	assert.True(t, ok, "Key found")
	assert.NotNil(t, mapValue, "Value is not nil")
	t.Logf("Updated map value: %v", updatedMapValue)

}

// go test -timeout 30s -run ^TestGetSliceFromSession_Correct$ github.com/zeroboo/go-api-session -v
func TestGetSliceFromSession_Correct(t *testing.T) {
	session := NewAPISessionWithPayload("user1", map[string]any{
		"key": []string{"1"},
	})
	slice, ok := GetOrCreatePayloadSlice[string](session, "key")
	assert.True(t, ok, "Key found")
	assert.NotNil(t, slice, "Value is nil")
	t.Logf("Empty slice value: %v", slice)

	slice, ok = GetPayloadSlice[string](session, "invalidkey")
	assert.False(t, ok, "Key found")
	assert.Nil(t, slice, "Value is nil")
	t.Logf("Empty slice value: %v", slice)

}
