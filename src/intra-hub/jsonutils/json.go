package jsonutils
import "encoding/json"

func MarshalUnsafe(val interface{}) string {
    js, err := json.Marshal(val)
    if err != nil {
        return "{}"
    }
    return string(js)
}