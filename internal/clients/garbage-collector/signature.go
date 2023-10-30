package garbagecollector

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
)

func (c *client) AddSignature(r *http.Request, JSONBody []byte) {
	hash := hmac.New(sha256.New, []byte(c.config.APIKey))
	hash.Write(JSONBody)

	r.Header.Add("HashSHA256", base64.URLEncoding.EncodeToString(hash.Sum(nil)))
}
