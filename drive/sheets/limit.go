package sheets

import (
	"golang.org/x/time/rate"
)

const (
	// Default rate (without paying) is 1 per second per user.
	ratePerSecond = 1

	// As many as 100 requests can be sent at a time.
	maxBurstRate = 100
)

var (
	// Rate limiter for slowing down access to acceptable rate.
	// Prevents annoying "USER-100s quota limit error".
	Limiter = rate.NewLimiter(ratePerSecond, maxBurstRate)
)
