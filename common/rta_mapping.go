package common

import (
	"github.com/RoaringBitmap/roaring"
	"sync"
)

const (
	MaxUserSize = 1 << 25
)

var _rtaMapping *RtaMap
var _rtaMapOnce sync.Once

func RtaMapInst() *RtaMap {
	_rtaMapOnce.Do(func() {
		_rtaMapping = &RtaMap{
			roaring.NewBitmap(),
		}
	})
	return _rtaMapping
}

type RtaMap struct {
	*roaring.Bitmap
}

func (m *RtaMap) InitByOneRta(rtaID int64, numbers []int) {
}
