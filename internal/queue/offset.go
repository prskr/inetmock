package queue

import "math"

type Offset struct {
	CurrentOffset int
}

func (o *Offset) Next() (offset int, overflow bool) {
	if o.CurrentOffset == math.MaxInt64 {
		o.CurrentOffset = 0
		overflow = true
	}
	offset = o.CurrentOffset
	o.CurrentOffset++
	return
}

func (o *Offset) Inc(val int) (newVal int, overflow bool) {
	var maxOffset = math.MaxInt64 - o.CurrentOffset
	if val >= maxOffset {
		val -= maxOffset
		o.CurrentOffset = 0
		overflow = true
	}
	o.CurrentOffset += val
	newVal = o.CurrentOffset
	return
}
