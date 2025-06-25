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

	//Track online users
	trackOnlineUsers bool
	onlineUserKey    string
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
	requestInterval int64,
	trackOnlineUsers bool) *RedisSessionManager {
	sessManager := &RedisSessionManager{
		redisClient:      redisClient,
		sessionKeyPrefix: sessionKeyPrefix,
		sessionTTL:       time.Duration(sessionTTL) * time.Millisecond,
		windowSize:       windowSize,
		maxCallPerWindow: maxCallPerWindow,
		requestInterval:  requestInterval,
		trackOnlineUsers: trackOnlineUsers,
	}

	if trackOnlineUsers {
		sessManager.onlineUserKey = fmt.Sprintf("online:%s", sessionKeyPrefix)
	} else {
		sessManager.onlineUserKey = ""
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
	now := time.Now()
	errValidate := sm.ValidateAPICall(&APIRequest{
		Owner:     owner,
		SessionId: sessionValue,
		URL:       url,
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
	Owner     string
	SessionId string
	URL       string
}

func (sm *RedisSessionManager) ValidateAPICall(request *APIRequest, session *APISession, currentTime time.Time) error {
	if session.Id != request.SessionId {
		return ErrInvalidSession
	}
	now := currentTime.UnixMilli()
	err := sm.UpdateSession(now, session)
	if err != nil {
		return err
	}
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
	bytes, errRedis := cmd.Bytes()
	if errRedis != nil {
		return nil, errRedis
	}

	session := &APISession{}
	errUnmarshal := msgpack.Unmarshal(bytes, session)

	if errUnmarshal != nil {
		return nil, errUnmarshal
	}
	return session, errUnmarshal
}

func (sm *RedisSessionManager) UpdateSession(currentMillis int64, session *APISession) error {
	window := currentMillis / sm.windowSize
	if window != session.Window {
		session.SetWindow(window)
	}
	session.Updated = currentMillis

	return sm.SetSession(context.Background(), session.Owner, session)
}

func (sm *RedisSessionManager) SetSession(ctx context.Context, owner string, session *APISession) error {
	session.Updated = time.Now().UnixMilli()
	payload, errSerialize := msgpack.Marshal(session)
	if errSerialize != nil {
		return errSerialize
	}

	key := sm.GetSessionKey(owner)
	cmd := sm.redisClient.Set(ctx, key, payload, sm.sessionTTL)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	if sm.trackOnlineUsers {
		// Update online user tracking
		cmd := sm.redisClient.ZAdd(context.Background(), sm.onlineUserKey, redis.Z{
			Score:  float64(session.Updated),
			Member: session.Owner,
		})
		if cmd.Err() != nil {
			return cmd.Err()
		}
	}
	return nil
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
	return session.Id, nil
}

func (sm *RedisSessionManager) DeleteSession(ctx context.Context, owner string) error {
	key := sm.GetSessionKey(owner)
	cmd := sm.redisClient.Del(ctx, key)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	if sm.trackOnlineUsers {
		// Remove from online users tracking
		cmd := sm.redisClient.ZRem(ctx, sm.onlineUserKey, owner)
		if cmd.Err() != nil {
			return cmd.Err()
		}
	}
	return nil
}

func (sm *RedisSessionManager) GetRequestInterval() int64 {
	return sm.requestInterval
}

func (sm *RedisSessionManager) GetMaxCallPerWindow() int64 {
	return sm.maxCallPerWindow
}

func (sm *RedisSessionManager) GetWindowSize() int64 {
	return sm.windowSize
}

func (sm *RedisSessionManager) GetOnlineUsers(ctx context.Context) (map[string]int64, error) {
	if !sm.trackOnlineUsers {
		return nil, fmt.Errorf("online users tracking is disabled")
	}

	cmd := sm.redisClient.ZRangeWithScores(ctx, sm.onlineUserKey, 0, -1)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	onlineUsers := make(map[string]int64)
	for _, z := range cmd.Val() {
		onlineUsers[z.Member.(string)] = int64(z.Score)
	}
	return onlineUsers, nil
}
