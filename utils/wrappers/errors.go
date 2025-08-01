// Copyright (C) 2020-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package wrappers

type Errs struct{ Err error }

func (errs *Errs) Errored() bool {
	return errs.Err != nil
}

func (errs *Errs) Add(errors ...error) {
	if errs.Err == nil {
		for _, err := range errors {
			if err != nil {
				errs.Err = err
				break
			}
		}
	}
}
