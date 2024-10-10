# go-api-session
Handle session for API, supports rate limitting


## 1. Features
- Validating session 
- Rate limiting 
  - request too frequently (HTTP code 429 Too Many Requests) using Fixed Window algorithm
  - request too fast (HTTP code 425 Too Early)
- Session payload: session can store extra data
## 2. Usage
### Install
```shell
go get github.com/zeroboo/go-api-session
```

### Setup
```golang
	//Create session manager
	sessionManager := apisession.NewRedisSessionManager(client,
		"session", //keys will have format of `session:<sessionId>``
		86400000,  //session last for 1 day
		60000,     //Time window is 1 minute
		10,        //Max 10 calls per minute
		1000)      //2 calls must be at least 1 second apart

	owner := "user1"
```
### Create new session
```golang
	//User starts a session: create new
	sessionId, errGet := sessionManager.StartSession(context.TODO(), owner)
	if errGet != nil {
		log.Printf("Failed to get session: %v", errGet)
	}
	//...
```
### Use session to verify API calls
```golang
	//...
	//Update in api call
	session, errSession := sessionManager.RecordAPICall(context.TODO(), sessionValue, owner, "url1")
	if errSession == nil {
		//Valid api call
		//Next processing...
		//...
	} else {
		//Invalid api call
		log.Printf("Failed to update session: %v", errGet)
	}
	
```

### Session payload
```golang
//Init session with extra data
session := NewAPISessionWithPayload("user1", map[string]any{
		"nickname": "Jaian",
		"oldTokens":      []string{"token1", "token2"},
	})
//... or set payload directly
session.SetPayload("nickname","Jaian")


//Retrieve old tokens from session by generic helper
knownTokens, ok := GetPayloadSlice[string](session, "oldTokens")//value: []string{"token1", "token2"}

//Retrieve nickname by type assertion
value := session.GetPayload("nickname")
nickname, ok := value.(string)


```

