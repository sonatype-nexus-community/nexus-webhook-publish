/* MIT License

Copyright (C) RJ Zaworski <rj@rjzaworski.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE. */

package webhook

import (
	"encoding/hex"
	"crypto/hmac"
	"crypto/sha1"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const secret = "duckduckgoose"

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	
	computed.Write(body)

	return []byte(computed.Sum(nil))
}

func verifySignature(secret []byte, signature string, body []byte) bool {
	actual := make([]byte, 20)
	hex.Decode(actual, []byte(signature))

	return hmac.Equal(actual, signBody(secret, body))
}

type HookContext struct {
	Signature string
	Payload   []byte
}

func ParseHook(secret []byte, req *http.Request) (*HookContext, error) {
	hc := HookContext{}

	if hc.Signature = req.Header.Get("x-nexus-webhook-signature"); len(hc.Signature) == 0 {
		return nil, errors.New("No signature!")
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return nil, err
	}

	if !verifySignature(secret, hc.Signature, body) {
		return nil, errors.New("Invalid signature")
	}

	hc.Payload = body

	return &hc, nil
}

func Handler(w http.ResponseWriter, r *http.Request) (body []byte, err error) {
	hc, err := ParseHook([]byte(secret), r)

	w.Header().Set("Content-type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Failed processing hook! ('%s')", err)
		io.WriteString(w, "{}")
		return nil, err
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "{}")
	return hc.Payload, nil
}
