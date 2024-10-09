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
	Payload any `json:"p" msgpack:"p"`

	//
	Meta map[string]any `json:"m" msgpack:"m"`
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
		Meta:    nil,
	}
}

func NewAPISessionFull(owner string, payload any, meta map[string]any) *APISession {
	return &APISession{
		Id:      GenerateSessionValue(owner),
		Records: make(map[string]*APICallRecord),
		Window:  0,
		Payload: payload,
		Meta:    meta,
	}
}
func (ses *APISession) SetWindow(window int64) {
	ses.Window = window
	for _, record := range ses.Records {
		record.Count = 0
	}
}
func (ses *APISession) SetMeta(key string, value any) {
	if ses.Meta == nil {
		ses.Meta = make(map[string]any)
	}
	ses.Meta[key] = value
}

func (ses *APISession) GetMeta(key string) any {
	if ses.Meta == nil {
		return nil
	}
	return ses.Meta[key]
}

// Returns metadata as string, empty string if not found
func (ses *APISession) GetMetaString(key string) string {
	if ses.Meta == nil {
		return ""
	}
	value := ses.Meta[key]
	if value != nil {
		strValue, isString := value.(string)
		if isString {
			return strValue
		}
	}
	return ""
}

// Returns metadata as int64, 0 if not found
func (ses *APISession) GetMetaInt64(key string) int64 {
	if ses.Meta == nil {
		return 0
	}
	value := ses.Meta[key]
	if value != nil {
		intValue, isString := value.(int64)
		if isString {
			return intValue
		}
	}
	return 0
}

// Returns metadata as int, 0 if not found
func (ses *APISession) GetMetaInt(key string) int {
	if ses.Meta == nil {
		return 0
	}
	value := ses.Meta[key]
	if value != nil {
		intValue, isString := value.(int)
		if isString {
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
