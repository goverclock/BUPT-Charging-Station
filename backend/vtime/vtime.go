package vtime

import (
	"log"
	"sync"
	"time"
)

const Xrate = 20

type vtime struct {
	ts     int64
	freeze bool
	mu     sync.Mutex
}

var vt vtime

func init() {
	now := time.Now().Unix()
	hour := time.Now().Hour()
	minu := time.Now().Minute()
	seco := time.Now().Second()
	now -= int64((hour*60+minu)*60 + seco)
	now += 6 * 60 * 60 // begining hour = 6
	now -= 5 * Xrate   // real 5s for 1st operation
	vt.ts = now

	go ticker()
}

func Now() time.Time {
	vt.mu.Lock()
	defer vt.mu.Unlock()
	return time.Unix(vt.ts, 0)
}

func ShouldFreeze() bool {
	vt.mu.Lock()
	defer vt.mu.Unlock()
	return vt.freeze
}

func UnFreeze() {
	vt.mu.Lock()
	defer vt.mu.Unlock()
	log.Println("UnFreeze!")
	vt.freeze = false
}

func ticker() {
	for {
		time.Sleep(time.Second)
		vt.mu.Lock()

		if vt.freeze {
			vt.mu.Unlock()
			continue
		}
		vt.ts += Xrate
		minu := time.Unix(vt.ts, 0).Minute()
		secs := time.Unix(vt.ts, 0).Second()
		if secs == 0 && (minu == 30 || minu == 0) {
			vt.freeze = true
		}

		vt.mu.Unlock()
	}
}
