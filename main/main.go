package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	apisession "github.com/zeroboo/go-api-session"
)

func main() {
	//Create redis client
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	//Create session manager
	sessionManager := apisession.NewRedisSessionManager(client,
		"session", //keys will have format of `session:<sessionId>``
		86400000,  //session last for 1 day
		60000,     //Time window is 1 minute
		10,        //Max 10 calls per minute
		1000)      //2 calls must be at least 1 second apart

	owner := "user1"

	//User starts a session: create new
	sessionId, errGet := sessionManager.StartSession(context.TODO(), owner)
	if errGet != nil {
		log.Printf("Failed to get session: %v", errGet)
	}
	//...

	//...
	//Update in api call
	session, errSession := sessionManager.RecordAPICall(context.TODO(), sessionId, owner, "url1")
	if errSession == nil {
		//Valid api call
		//Next processing...
		//...
	} else {
		//Invalid api call
		log.Fatalf("Failed to update session: %v", errGet)
	}
	log.Printf("Session: %#v", session)
}
