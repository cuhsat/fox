package flags

import (
	"errors"
)

type BagMode string

const (
	BagName = "evidence"
)

const (
	BagModeNone   BagMode = "none"
	BagModeText   BagMode = "text"
	BagModeJson   BagMode = "json"
	BagModeJsonl  BagMode = "jsonl"
	BagModeXml    BagMode = "xml"
	BagModeSqlite BagMode = "sqlite"
)

const (
	BagUrlLogstash = "http://localhost:8080"
)

func (bm *BagMode) String() string {
	return string(*bm)
}

func (bm *BagMode) Type() string {
	return "BagMode"
}

func (bm *BagMode) Set(v string) error {
	switch v {
	case string(BagModeNone):
		fallthrough
	case string(BagModeText):
		fallthrough
	case string(BagModeJson):
		fallthrough
	case string(BagModeJsonl):
		fallthrough
	case string(BagModeXml):
		fallthrough
	case string(BagModeSqlite):
		*bm = BagMode(v)
		return nil

	default:
		return errors.New("mode not recognized")
	}
}
