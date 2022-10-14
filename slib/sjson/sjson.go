package sjson

import "github.com/bitly/go-simplejson"

func NewJson(body []byte) (*simplejson.Json, error) {
	return simplejson.NewJson(body)
}
