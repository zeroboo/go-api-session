package apisession

import (
	"context"
	"errors"
	"os"
	"testing"

	redis "github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var redisClient *redis.Client
var sessionOwners []string = []string{}
var sessionPrefix string = "sess"

func TestMain(m *testing.M) {
	log.Infof("TestMain: Init done!!!")
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
		FullTimestamp:    false,
	})

	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	code := m.Run()

	CleanUpTest()
	os.Exit(code)

}

func CleanUpTest() {
	for _, owner := range sessionOwners {
		key := GetRedisSessionKey(sessionPrefix, owner)
		cmd := redisClient.Del(context.TODO(), key)
		if cmd.Err() != nil {
			log.Infof("Failed to delete session %s: %v", owner, cmd.Err())
		} else {
			log.Infof("Deleted session of %s, key=%v", owner, key)
		}
	}
}

func TestDeleteSession(t *testing.T) {
	owner := "user_" + t.Name()
	manager := NewRedisSessionManager(redisClient, sessionPrefix, 1000, 10000, 10, 5, false)

	sessionId, errNewSession := manager.StartSession(context.TODO(), owner)
	if errNewSession != nil {
		t.Fatalf("Create new session failed: %v", errNewSession)
	}
	sessionOwners = append(sessionOwners, owner)

	session, errGetSession := manager.GetSession(context.TODO(), owner)
	if errGetSession != nil {
		t.Fatalf("Get session failed: %v", errGetSession)
	}
	assert.Equal(t, session.Id, sessionId, "Session ID must match")

	errDelete := manager.DeleteSession(context.TODO(), owner)
	if errDelete != nil {
		t.Errorf("Delete session failed: %v", errDelete)
	}

	// Verify that the session is deleted
	session, errGet := manager.GetSession(context.TODO(), owner)
	assert.True(t, errors.Is(errGet, redis.Nil), "Get deleted session should return nil")
	assert.Nil(t, session, "Session should be nil after deletion")

	log.Infof("Successfully deleted session for owner %s with ID %s", owner, sessionId)
}
