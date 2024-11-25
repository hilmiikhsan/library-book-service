package helpers

import "encoding/json"

func MarshalJSON(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}

func UnmarshalJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
