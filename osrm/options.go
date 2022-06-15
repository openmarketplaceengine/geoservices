package osrm

import (
	"fmt"
	"net/url"
)

func EncodeUrlParam(key string, val []int) string {
	var buf []byte
	buf = append(buf, url.QueryEscape(key)...)
	buf = append(buf, '=')
	for i, v := range val {
		if i > 0 {
			buf = append(buf, ';')
		}
		buf = append(buf, fmt.Sprintf("%v", v)...)
	}
	return string(buf)
}
