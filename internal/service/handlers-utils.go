package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ReadJsonBody(body io.ReadCloser, out interface{}) error {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(body)

	if err != nil {
		return fmt.Errorf("read request body error %w", err)
	}

	if err = json.Unmarshal(buf.Bytes(), &out); err != nil {
		return fmt.Errorf("requset body json error %w", err)
	}

	return nil
}

func WriteJsonBody(w http.ResponseWriter, data interface{}) error {
	resp, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("response body json error %w", err)
	}

	_, err = w.Write(resp)
	if err != nil {
		return fmt.Errorf("response body write error %w", err)
	}

	return nil
}
