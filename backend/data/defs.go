package data

import "flag"

var Port int

var FAST_STATION_COUNT int = 2
var SLOW_STATION_COUNT int = 3
var MAX_WAITING_SLOT int = 6
var MAX_STATION_QUEUE int = 2
var CALL_SCHEDULE int = 0	// will not implement

func init() {
	port_arg := flag.Int("port", 8080, "server address")
	fast_count_arg := flag.Int("fast", 2, "fast station count")
	slow_count_arg := flag.Int("slow", 3, "slow station count")
	flag.Parse()
	FAST_STATION_COUNT = *fast_count_arg
	SLOW_STATION_COUNT = *slow_count_arg
	Port = *port_arg
}
