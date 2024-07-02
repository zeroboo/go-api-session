# go-api-session
Handle session for API, supports rate limitting


## Features
- Validating session 
- Rate limiting 
  - request too frequently (HTTP code 429 Too Many Requests) using Fixed Window algorithm
  - request too fast (HTTP code 425 Too Early)
## Usage
### Install
```shell
go get github.com/zeroboo/go-api-session
```
### Sample 
#### Setup
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
#### Create new session
```golang
	//User starts a session: create new
	sessionId, errGet := sessionManager.StartSession(context.TODO(), owner)
	if errGet != nil {
		log.Printf("Failed to get session: %v", errGet)
	}
	//...
```
#### Use session to verify API calls
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