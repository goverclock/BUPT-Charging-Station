package vtime

import (
	"sync"
	"time"
)

const Xrate = 20

type vtime struct {
	ts int64
	mu sync.Mutex
}

var vt vtime

func init() {
	now := time.Now().Unix()
	hour := time.Now().Hour()
	minu := time.Now().Minute()
	seco := time.Now().Second()
	now -= int64((hour*60+minu)*60 + seco)
	now += 6 * 60 * 60 // begining hour = 6
	now -= 30 * Xrate	// real 30s for 1st operation
	vt.ts = now

	go ticker()
}

func Now() time.Time {
	vt.mu.Lock()
	defer vt.mu.Unlock()
	return time.Unix(vt.ts, 0)
}

func ticker() {
	for {
		time.Sleep(time.Second)
		vt.mu.Lock()

		vt.ts += Xrate

		vt.mu.Unlock()
	}
}
