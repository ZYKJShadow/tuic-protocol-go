package fragment

import "sync"

type FBuffer struct {
	fragments map[uint8][]byte
	total     uint8
	size      uint16
}

type FCache struct {
	cache map[uint16]*FBuffer
	sync.Mutex
}

func NewFCache() *FCache {
	return &FCache{
		cache: make(map[uint16]*FBuffer),
	}
}

func NewFragmentBuffer(total uint8, size uint16) *FBuffer {
	return &FBuffer{
		fragments: make(map[uint8][]byte),
		total:     total,
		size:      size,
	}
}

func (fc *FCache) DelFragment(assocID uint16) {
	fc.Lock()
	defer fc.Unlock()

	delete(fc.cache, assocID)
}

func (fc *FCache) AddFragment(assocID uint16, fragID uint8, total uint8, size uint16, data []byte) []byte {
	fc.Lock()
	defer fc.Unlock()

	fb, ok := fc.cache[assocID]
	if !ok {
		fb = NewFragmentBuffer(total, size)
		fc.cache[assocID] = fb
	}

	if fb.SetFragData(fragID, data).IsComplete() {
		assembled := fb.Assemble()
		delete(fc.cache, assocID)
		return assembled
	}

	return nil
}

func (fb *FBuffer) SetFragData(fragID uint8, data []byte) *FBuffer {
	fb.fragments[fragID] = data
	return fb
}

func (fb *FBuffer) IsComplete() bool {
	return uint8(len(fb.fragments)) == fb.total
}

func (fb *FBuffer) Assemble() []byte {
	assembled := make([]byte, fb.size)
	for i := uint8(0); i < fb.total; i++ {
		copy(assembled[i*uint8(len(fb.fragments[i])):(i+1)*uint8(len(fb.fragments[i]))], fb.fragments[i])
	}

	return assembled
}
