package clock

import (
	"errors"
	"time"

	"github.com/beevik/ntp"
	"golang.org/x/xerrors"
)

// ErrNTPQueryFailed is returned if an NTP query failed.
var ErrNTPQueryFailed = errors.New("NTP query failed")

// difference between network time and node's local time.
var offset time.Duration

// FetchTimeOffset establishes the difference in local vs network time.
// This difference is stored in offset so that it can be used to adjust the local clock.
func FetchTimeOffset(host string) error {
	resp, err := ntp.Query(host)
	if err != nil {
		return xerrors.Errorf("NTP query error (%v): %w", err, ErrNTPQueryFailed)
	}
	offset = resp.ClockOffset

	return nil
}

// SyncedTime gets the synchronized time (according to the network) of a node.
func SyncedTime() time.Time {
	return time.Now().Add(offset)
}
