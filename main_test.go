package apisession

import (
	"context"
	"os"
	"testing"

	redis "github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
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
