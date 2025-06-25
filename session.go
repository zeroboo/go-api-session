package apisession

import (
	"fmt"
	"time"
)

type APISession struct {
	Id    string `json:"i" msgpack:"i"`
	Owner string `json:"o" msgpack:"o"`
	//Map of url to API call track
	Records map[string]*APICallRecord `json:"r" msgpack:"r"`

	//Current time window
	Window int64 `json:"w" msgpack:"w"`

	//Payload are extra data of session
	Payload map[string]any `json:"p" msgpack:"p"`

	Created int64 `json:"c" msgpack:"c"` //Created time in milliseconds
	Updated int64 `json:"u" msgpack:"u"` //Updated time in milliseconds
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
		Owner:   owner,
		Records: make(map[string]*APICallRecord),
		Window:  0,
		Payload: nil,
		Created: time.Now().UnixMilli(),
		Updated: time.Now().UnixMilli(),
	}
}

func NewAPISessionWithPayload(owner string, payload map[string]any) *APISession {
	return &APISession{
		Id:      GenerateSessionValue(owner),
		Owner:   owner,
		Records: make(map[string]*APICallRecord),
		Window:  0,
		Payload: payload,
		Created: time.Now().UnixMilli(),
		Updated: time.Now().UnixMilli(),
	}
}

func (ses *APISession) SetPayload(key string, value any) {
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

// GetPayloadMap returns a map from session payload, if not exist, return nil
func GetPayloadMap[K comparable, V any](sess *APISession, key string) map[K]V {
	value, exist := sess.Payload[key]
	if !exist {
		return nil
	}

	typedValue, ok := value.(map[K]V)
	if !ok {
		return nil
	}
	return typedValue
}

// GetOrCreatePayloadMap returns a map from session payload, if not exist, create a new map.
//   - If the value of key is not a map, replace it with a new map.
//   - If any value of the map is not the correct type, it will be ignored.
//
// Return the map and a boolean indicate if the map is created
func GetOrCreatePayloadMap[K comparable, V any](sess *APISession, key string) (map[K]V, bool) {
	value, exist := sess.Payload[key]
	if !exist {
		newMap := make(map[K]V)
		sess.Payload[key] = newMap
		return newMap, true
	}

	//Check if the value is correct map type
	typedValue, ok := value.(map[K]V)
	if ok {
		return typedValue, false
	}

	//Check if the value is map of any
	mapAny, ok := value.(map[K]any)
	if ok {
		newMap := make(map[K]V)
		for k, v := range mapAny {
			mapValue, valueOk := v.(V)
			if valueOk {
				newMap[k] = mapValue
			}
		}
		sess.Payload[key] = newMap
		return newMap, false
	}

	//Casting failed, replace with new map
	newMap := make(map[K]V)
	sess.Payload[key] = newMap
	return newMap, true

}
func GetPayloadSlice[V any](sess *APISession, key string) ([]V, bool) {
	value, exist := sess.Payload[key]
	if !exist {
		return nil, false
	}
	retValues := []V{}
	for _, v := range value.([]any) {
		typedValue, ok := v.(V)
		if !ok {
			return nil, false
		}
		retValues = append(retValues, typedValue)
	}

	return retValues, true
}

func GetOrCreatePayloadSlice[V any](sess *APISession, key string) ([]V, bool) {
	value, exist := sess.Payload[key]
	if !exist {
		newSlice := make([]V, 0)
		sess.Payload[key] = value
		return newSlice, false
	}

	typedValue, ok := value.([]V)
	return typedValue, ok
}

// SetPayloadMap init a map in session payload
func SetPayloadMap[K comparable, V any](sess *APISession, key string, value map[K]V) {
	if sess.Payload == nil {
		sess.Payload = make(map[string]any)
	}
	sess.Payload[key] = value
}
