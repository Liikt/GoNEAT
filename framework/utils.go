package framework

import (
	"bytes"
	"encoding/gob"
)

func DeepCopy(dst, src interface{}) error {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(&buffer).Decode(dst)
}
