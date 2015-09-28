// Copyright (c) 2014 RightScale, Inc. - see LICENSE

package misc

import (
	"time"

	"gopkg.in/inconshreveable/log15.v2"
)

var TimeNow = time.Now
var TimeSleep = time.Sleep

//===== timeNow / timeSleep stubs

// fakeTime simplifies tests where time.Now() is called as well as time.Sleep(). It allows the
// test to set time.Now() and not have it change so tests can rely on one value. For example,
// if code calculates time.Now() + 2 seconds then the test can do that too with the same value for
// now. fakeTime.Sleep() sleeps a minimal amount ot let other goroutines run and then adds the
// sleep time to the value for "now". In other words, time only advances when sleep is called.
// Use all this in normal execution by importing this package and using timeNow and timeSleep
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

func NewFakeTime() *FakeTime {
	ft := FakeTime{now: time.Now()}
	log15.Debug("FakeTime starts", "now", ft.now.Format(time.RFC3339Nano))
	return &ft
}

func (ft *FakeTime) Now() time.Time {
	return ft.now
}

func (ft *FakeTime) Sleep(d time.Duration) {
	ft.now = ft.now.Add(d)
	ft.slept += d
	time.Sleep(time.Millisecond / 10)
}

func (ft *FakeTime) Slept() time.Duration { return ft.slept }
