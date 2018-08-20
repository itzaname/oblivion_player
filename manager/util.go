package manager

import (
	"math"
	"strconv"
	"time"
)

func isInBounds(x, y, w, h, mx, my int32) bool {
	if mx >= x && my >= y && mx <= (x+w) && my <= (y+h) {
		return true
	}
	return false
}

func easing(t, b, c, d float64) int32 {
	t /= d
	return int32(c*t*t*t + b)
}

func animate(duration int, starttime float64, start int32, end int32) int32 {
	now := float64(time.Now().UnixNano()) / float64(time.Second)
	starttime = now - starttime
	if starttime > float64(duration) {
		starttime = float64(duration)
	}
	return easing(float64(starttime), float64(start), float64(end), float64(duration))
}

func getHHMMSS(sec_num float64) string {
	var minutes = math.Floor(sec_num / 60)
	var seconds = sec_num - (minutes * 60)

	if seconds < 10 {
		return strconv.Itoa(int(minutes)) + ":0" + strconv.Itoa(int(seconds))
	}
	return strconv.Itoa(int(minutes)) + ":" + strconv.Itoa(int(seconds))
}
