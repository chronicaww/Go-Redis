//   Copyright 2009-2012 Joubin Houshyar
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package redis

import (
	"fmt"
	"log"
)

// ----------------------------------------------------------------------------
// redis.Error
// ----------------------------------------------------------------------------

// Go-Redis API level error type
//
type Error interface {
	error
	// if true Error is a RedisError
	IsRedisError() bool
}

// ----------------------------------------------------------------------
// Go-Redis System Errors or Bugs
// ----------------------------------------------------------------------

// A system level error, ranging from connectivity issues to
// detected bugs.  Basically anything other than an Redis Server ERR.
//
// System errors can (and typically do) have an underlying Go std. lib
// or 3rd party lib error cause, such as net.Error, etc.
type SystemError interface {
	Cause() error
}

// supports SystemError interface
type systemError struct {
	msg   string
	cause error
}

func newSystemError(msg string) Error {
	return newSystemErrorWithCause(msg, nil)
}

func newSystemErrorWithCause(msg string, cause error) Error {
	e := &systemError{
		msg:   msg,
		cause: cause,
	}
	return e
}

// See: redis.Error#IsRedisError()
func (e *systemError) IsRedisError() bool { return false }

func (e *systemError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" [cause: %s]", e.cause.Error())
	}
	return fmt.Sprintf("SYSTEM_ERROR - %s%s", e.msg, cause)
}

func (e *systemError) Cause() error {
	return e.cause
}

// ----------------------------------------------------------------------
// Redis Server Errors
// ----------------------------------------------------------------------

// ERR Errors returned by the Redis server e.g. for bad AUTH password.
type RedisError interface {
	Message() string
}

type redisError2 struct {
	msg string
}

func newRedisError(msg string) Error {
	e := &redisError2{
		msg: msg,
	}
	return e
}

// See: redis.Error#IsRedisError()
func (e *redisError2) IsRedisError() bool { return true }

func (e *redisError2) Error() string {
	return fmt.Sprintf("REDIS_ERROR - %s", e.msg)
}

// ----------------------------------------------------------------------
// temp legacy junk
// ----------------------------------------------------------------------

// REVU - debug level should be handled in log
// TODO - redo log
// TODO - all errors should log on new
// utility function emits log if _debug flag /debug() is true
// Error is returned.
// usage:
//      foo, e := FooBar()
//      if e != nil {
//          return withError(e)
//      }
func withError(e Error) Error {
	if debug() {
		log.Println(e)
	}
	return e
}
