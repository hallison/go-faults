// Copyright 2019 Hallison Batista. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package faults provides a small errors handler.
package faults

import (
	"fmt"
)

// Type function that runs a block code and returns error.
type BlockFunction func() error

// Faults is a basic structure to handle errors.
type Faults struct {
	lastMessage string           `json:"-"`
	failures    map[string]error `json:"errors"`
	allowStack  bool             `json:"-"`
	locked      bool             `json:"-"`
}

func (f *Faults) resetFailures() *Faults {
	f.failures = make(map[string]error)

	return f
}

func (f *Faults) addFailWithMessage(fail error, message string) *Faults {
	if !f.allowStack {
		f.resetFailures()
	}

	f.lastMessage = message
	f.failures[message] = fail

	return f
}

func (f *Faults) lock() *Faults {
	f.locked = !f.allowStack

	return f
}

// Reset set default values to faults handler.
func (f *Faults) Reset() *Faults {
	f.lastMessage = ""
	f.allowStack = false
	f.locked = false

	return f.resetFailures()
}

// LastMessage returns the last message from failures.
func (f *Faults) LastMessage() string {
	return f.lastMessage
}

// EnableStack if handle a stack of errors, this function enable this feature.
func (f *Faults) EnableStack() *Faults {
	f.allowStack = true

	return f
}

// DisableStack if do not handle a stack of errors, this function disable this feature.
func (f *Faults) DisableStack() *Faults {
	f.allowStack = false

	return f
}

// Check add a message and check if the block code should be runned testing if
// has no failures or allow stack or unlocked.
func (f *Faults) Check(message string, block BlockFunction) *Faults {
	if !f.locked {
		if fail := block(); fail != nil {
			f.addFailWithMessage(fail, message).lock()
		}
	}

	return f
}

// Add fail and message.
func (f *Faults) Add(fail error, message string) *Faults {
	return f.addFailWithMessage(fail, message)
}

// AddIf add fail and massage by condition.
func (f *Faults) AddIf(condition bool, message string) *Faults {
	if f.locked = condition; condition {
		f.failures[message] = fmt.Errorf(message)
	}
	return f
}

// GetLast returns the last error.
func (f *Faults) GetLast() error {
	if f.lastMessage != "" {
		return f.failures[f.lastMessage]
	}
	return nil
}

// GetAll returns all failures.
func (f *Faults) GetAll() map[string]error {
	return f.failures
}

// IsEmpty check if failures is empty.
func (f *Faults) IsEmpty() bool {
	return len(f.failures) == 0
}

// IsNotEmpty check if failures is not empty.
func (f *Faults) IsNotEmpty() bool {
	return !f.IsEmpty()
}

// New creates a new faults handler.
func New() (f *Faults) {
	return new(Faults).Reset()
}
