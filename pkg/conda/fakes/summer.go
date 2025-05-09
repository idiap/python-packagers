// SPDX-FileCopyrightText: Copyright (c) 2013-Present CloudFoundry.org Foundation, Inc. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0

package fakes

import "sync"

type Summer struct {
	SumCall struct {
		mutex     sync.Mutex
		CallCount int
		Receives  struct {
			Arg []string
		}
		Returns struct {
			String string
			Error  error
		}
		Stub func(...string) (string, error)
	}
}

func (f *Summer) Sum(param1 ...string) (string, error) {
	f.SumCall.mutex.Lock()
	defer f.SumCall.mutex.Unlock()
	f.SumCall.CallCount++
	f.SumCall.Receives.Arg = param1
	if f.SumCall.Stub != nil {
		return f.SumCall.Stub(param1...)
	}
	return f.SumCall.Returns.String, f.SumCall.Returns.Error
}
