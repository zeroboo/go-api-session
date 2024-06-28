package apisession

import (
	"fmt"
	"strings"
)

type APISession struct {
	Id string `json:"i" msgpack:"i"`

	Owner string `json:"o" msgpack:"o"`

	//Map of url to API call track
	Records map[string]*APICallRecord `json:"r" msgpack:"r"`

	//Current time window
	Window int64 `json:"w" msgpack:"w"`
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
		Owner:   owner,
		Records: make(map[string]*APICallRecord),
		Window:  0,
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

func (ses *APISession) ValidateOwner(owner string) bool {

	return strings.EqualFold(ses.Owner, owner)
}
