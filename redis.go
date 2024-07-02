package apisession

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	redis "github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

type RedisSessionManager struct {
	sessionKeyPrefix string
	redisClient      *redis.Client
	//In milliseconds
	sessionTTL time.Duration

	//In milliseconds
	windowSize int64

	//Max request in a time window
	maxCallPerWindow int64

	//minimum milliseconds between 2 request, 0 mean no limit
	requestInterval int64
}

// Create redis session manager
// Params:
// - redisClient: redis client
//
// - sessionKeyPrefix: prefix for session key
//
// - sessionTTL: session time to live in milliseconds
//
// - windowSize: time window in milliseconds
//
// - maxCallPerWindow: max calls allowed per window
//
// - minRequestInterval: minimum milliseconds between 2 request, 0 mean no limit
func NewRedisSessionManager(redisClient *redis.Client,
	sessionKeyPrefix string,
	sessionTTL int64,
	windowSize int64,
	maxCallPerWindow int64,
	requestInterval int64) ISessionManager {
	return CreateNewRedisSessionManager(redisClient, sessionKeyPrefix, sessionTTL, windowSize, maxCallPerWindow, requestInterval)
}

func CreateNewRedisSessionManager(redisClient *redis.Client,
	sessionKeyPrefix string,
	sessionTTL int64,
	windowSize int64,
	maxCallPerWindow int64,
	requestInterval int64) *RedisSessionManager {
	sessManager := &RedisSessionManager{
		redisClient:      redisClient,
		sessionKeyPrefix: sessionKeyPrefix,
		sessionTTL:       time.Duration(sessionTTL) * time.Millisecond,
		windowSize:       windowSize,
		maxCallPerWindow: maxCallPerWindow,
		requestInterval:  requestInterval,
	}
	return sessManager
}

// returns sha256 hash of the value
func Hash(value string) string {
	h := sha256.New()
	h.Write([]byte(value))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func GenerateSessionValue(ownerId string) string {
	payload := fmt.Sprintf("%s-%d-%d", ownerId, time.Now().Unix(), rand.Int63())
	sessionId := Hash(payload)
	return sessionId
}

func (sm *RedisSessionManager) GetSessionKey(sessionId string) string {
	return GetRedisSessionKey(sm.sessionKeyPrefix, sessionId)
}

func GetRedisSessionKey(prefix string, sessionId string) string {
	return fmt.Sprintf("%v:%v", prefix, sessionId)
}

func (sm *RedisSessionManager) RecordAPICall(ctx context.Context, sessionValue string, owner string, url string) (*APISession, error) {
	session, errGet := sm.GetSession(ctx, owner)
	if errGet != nil {
		return nil, errGet
	}

	//Validate session
	now := time.Now().UnixMilli()
	errValidate := sm.ValidateAPICall(APIRequest{
		Owner:   owner,
		Session: sessionValue,
		URL:     url,
	}, session, now)
	if errValidate != nil {
		return nil, errValidate
	}

	errUpdate := sm.SetSession(ctx, owner, session)
	if errUpdate != nil {
		return nil, errUpdate
	}

	return session, nil
}

type APIRequest struct {
	Owner   string
	Session string
	URL     string
}

func (sm *RedisSessionManager) ValidateAPICall(request APIRequest, session *APISession, now int64) error {
	if session.Session != request.Session {
		return ErrInvalidSession
	}

	sm.UpdateSession(now, session)
	call := session.GetCallRecord(request.URL)

	if sm.requestInterval > 0 {
		if now-call.Last < sm.requestInterval {
			return ErrTooFast
		}
	}

	if call.Count+1 > sm.maxCallPerWindow {
		return ErrTooMany
	}
	call.Count++
	call.Last = now

	return nil

}

func (sm *RedisSessionManager) GetSession(ctx context.Context, owner string) (*APISession, error) {
	key := sm.GetSessionKey(owner)
	cmd := sm.redisClient.Get(ctx, key)
	payload, errRedis := cmd.Bytes()
	if errRedis != nil {
		return nil, errRedis
	}

	session := &APISession{}
	errUnmarshal := msgpack.Unmarshal(payload, session)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}
	return session, errUnmarshal
}

func (sm *RedisSessionManager) UpdateSession(currentMillis int64, session *APISession) {
	window := currentMillis / sm.windowSize
	if window != session.Window {
		session.SetWindow(window)
	}

}
func (sm *RedisSessionManager) SetSession(ctx context.Context, owner string, session *APISession) error {
	payload, errSerialize := msgpack.Marshal(session)
	if errSerialize != nil {
		return errSerialize
	}

	key := sm.GetSessionKey(owner)
	cmd := sm.redisClient.Set(ctx, key, payload, sm.sessionTTL)

	return cmd.Err()
}

// StartSession creates a new session for the owner and insert to db
//
// Returns:
//   - sessionId string: id of new session
//   - error: error if exists, nil is successful
func (sm *RedisSessionManager) StartSession(ctx context.Context, owner string) (string, error) {
	session := NewAPISession(owner)
	errSet := sm.SetSession(ctx, owner, session)
	if errSet != nil {
		return "", errSet
	}
	return session.Session, nil
}

func (sm *RedisSessionManager) DeleteSession(ctx context.Context, owner string) error {
	key := sm.GetSessionKey(owner)
	cmd := sm.redisClient.Del(ctx, key)
	return cmd.Err()
}
