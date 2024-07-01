package apisession

import (
	"fmt"
)

type APISession struct {
	Session string `json:"s" msgpack:"s"`
	//Map of url to API call track
	Records map[string]*APICallRecord `json:"r" msgpack:"r"`

	//Current time window
	Window int64 `json:"w" msgpack:"w"`

	//Payload are extra data of session
	Payload map[string]interface{} `json:"p" msgpack:"p"`
}

// Tracks how an api is being called
type APICallRecord struct {

	//calls in current window
	Count int64 `json:"c" msgpack:"c"`

	//Last call in milliseconds
	Last int64 `json:"l" msgpack:"l"`
}

func NewAPICallRecord() *APICallRecord {
	return &APICallRecord{
		Count: 0,
		Last:  0,
	}

}
func (session *APISession) RecordCall(url string) error {
	return fmt.Errorf("not implemented")
}

func NewAPISession(owner string) *APISession {
	return &APISession{
		Session: GenerateSessionValue(owner),
		Records: make(map[string]*APICallRecord),
		Window:  0,
		Payload: make(map[string]interface{}),
	}
}

func (ses *APISession) SetWindow(window int64) {
	ses.Window = window
	for _, record := range ses.Records {
		record.Count = 0
	}
}

func (ses *APISession) GetCallRecord(url string) *APICallRecord {
	var record *APICallRecord
	record, exist := ses.Records[url]
	if !exist {
		record = NewAPICallRecord()
		ses.Records[url] = record
	}
	return record
}

func (ses *APISession) ValidateSession(session string) bool {
	return ses.Session == session
}
