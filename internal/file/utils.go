package file

import (
	"bytes"
	"encoding/gob"
	"p86l/internal/debug"
)

func (a *AppFS) DecodeData(appDebug *debug.Debug, b []byte) (Data, *debug.Error) {
	var d Data
	decoder := gob.NewDecoder(bytes.NewReader(b))
	if err := decoder.Decode(&d); err != nil {
		return d, appDebug.New(err, debug.FSError, debug.ErrDataLoad)
	}
	return d, nil
}

func (a *AppFS) DecodeCache(appDebug *debug.Debug, b []byte) (Cache, *debug.Error) {
	var c Cache
	decoder := gob.NewDecoder(bytes.NewReader(b))
	if err := decoder.Decode(&c); err != nil {
		return c, appDebug.New(err, debug.FSError, debug.ErrCacheLoad)
	}
	return c, nil
}
