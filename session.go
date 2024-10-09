package apisession

import (
	"fmt"
)

type APISession struct {
	Id string `json:"i" msgpack:"i"`
	//Map of url to API call track
	Records map[string]*APICallRecord `json:"r" msgpack:"r"`

	//Current time window
	Window int64 `json:"w" msgpack:"w"`

	//Payload are extra data of session
	Payload map[string]any `json:"p" msgpack:"p"`
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
		Id:      GenerateSessionValue(owner),
		Records: make(map[string]*APICallRecord),
		Window:  0,
		Payload: nil,
	}
}

func NewAPISessionWithPayload(owner string, payload map[string]any) *APISession {
	return &APISession{
		Id:      GenerateSessionValue(owner),
		Records: make(map[string]*APICallRecord),
		Window:  0,
		Payload: payload,
	}
}

func (ses *APISession) AddPayload(key string, value any) {
	if ses.Payload == nil {
		ses.Payload = make(map[string]any)
	}
	ses.Payload[key] = value
}

func (ses *APISession) GetPayload(key string) any {
	if ses.Payload == nil {
		return nil
	}
	return ses.Payload[key]
}

func (ses *APISession) SetWindow(window int64) {
	ses.Window = window
	for _, record := range ses.Records {
		record.Count = 0
	}
}

// Returns metadata as string, empty string if not found
func (ses *APISession) GetPayloadString(key string) string {
	if ses.Payload == nil {
		return ""
	}
	value := ses.Payload[key]
	if value != nil {
		strValue, isString := value.(string)
		if isString {
			return strValue
		}
	}
	return ""
}

// Returns metadata as int64, 0 if not found
func (ses *APISession) GetPayloadInt64(key string) int64 {
	if ses.Payload == nil {
		return 0
	}
	value := ses.Payload[key]
	if value != nil {
		intValue, isInt64 := value.(int64)
		if isInt64 {
			return intValue
		}
	}
	return 0
}

// Returns metadata as int, 0 if not found
func (ses *APISession) GetPayloadInt(key string) int {
	if ses.Payload == nil {
		return 0
	}
	value := ses.Payload[key]
	if value != nil {
		intValue, isInt := value.(int)
		if isInt {
			return intValue
		}
	}
	return 0
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
	return ses.Id == session
}
