// Copyright (c) 2014 RightScale, Inc. - see LICENSE

package misc

import (
	"time"

	"gopkg.in/inconshreveable/log15.v2"
)

// TimeNow() is a replacement for time.Now() that allows testing
var TimeNow = time.Now

// TimeSleep() is a replacement for time.Sleep() that allows testing
var TimeSleep = time.Sleep

// FakeTime simplifies tests where time.Now() and time.Sleep() are called. It allows the
// test to set time.Now() and not have it change so tests can rely on one value. For example,
// if code calculates time.Now() + 2 seconds then the test can do that too with the same value for
// now. FakeTime.Sleep() sleeps a minimal amount to let other goroutines run and then adds the
// sleep time to the value for "now". In other words, time only advances when sleep is called.
// Use all this in normal execution by importing this package and using TimeNow and TimeSleep
// instead of time.Now() and time.Sleep()
// For testing use:
// func init() {
//      ft := NewFakeTime()
//	timeNow = func() *time.Time { return ft.Now() }
//	timeSleep = func(d time.Duration) { ft.Sleep(d) }
//}
type FakeTime struct {
	now   time.Time
	slept time.Duration
}

// NewFakeTime returns a fresh FakeTime struct initialized to the current time
func NewFakeTime() *FakeTime {
	ft := FakeTime{now: time.Now()}
	log15.Debug("FakeTime starts", "now", ft.now.Format(time.RFC3339Nano))
	return &ft
}

// Now implements time.Now() for FakeTime
func (ft *FakeTime) Now() time.Time {
	return ft.now
}

// Sleep implements time.Sleep for FakeTime. It simply advances the notion of "now" by the time
// being slept and does a real sleep for 100us to allow other goroutines to run a quantum
func (ft *FakeTime) Sleep(d time.Duration) {
	ft.now = ft.now.Add(d)
	ft.slept += d
	time.Sleep(time.Millisecond / 10)
}

// Slept returns a duration representing the amount of time we have slept since NewFakeTime.
// This is useful in tests to find out how much fake time has elapsed, which allows to check
// that things like timeouts and back-off have actually occurred.
func (ft *FakeTime) Slept() time.Duration { return ft.slept }
