package mocks

import "github.com/leg100/ots"

type Subscription struct {
	c chan ots.Event
}

func NewSubscription(size int) *Subscription {
	return &Subscription{c: make(chan ots.Event, size)}
}

func (s *Subscription) SendEvent(ev ots.Event) { s.c <- ev }

func (s *Subscription) C() <-chan ots.Event { return s.c }

func (s *Subscription) Close() error { return nil }
