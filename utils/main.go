package utils

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/eymen-iron/insider-task/esrc"
)

func EncodeJSON(data esrc.Item) []byte {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(data)
	if err != nil {
		fmt.Println("JSON doesn't resolve:", err)
	}
	return buf.Bytes()
}
