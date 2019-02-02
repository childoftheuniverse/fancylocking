package fancylocking

import (
  "context"
  "time"
)

/*
MutexWithDeadline implements Mutex-like operations on top of a channel.
Unlike the base Mutex implementation in sync, MutexWithDeadline supports
the operations TryLock() (attempt to acquire a lock immediately and bail
if that fails) and LockWithDeadline() (attempt to acquire a lock until
certain time and fail if the deadline expires).

The mutexes are not recursive.

There is a potential overhead to this as it is using a channel in the
background, so only use this type of Mutex if you require the
LockWithDeadline() function. If you only need TryLock() functionality,
it is potentially better to implement that as an atomic operation.
*/
type MutexWithDeadline chan struct{}

/*
NewMutexWithDeadline initializes the MutexWithDeadline. This is
unfortunately necessary as the MutexWithDeadline is backed by a
channel.
*/
func NewMutexWithDeadline() MutexWithDeadline {
  return make(chan struct{}, 1)
}

/*
Lock unconditionally attempts to lock the mutex. This function will
potentially wait forever.
*/
func (m MutexWithDeadline) Lock() {
  m <- struct{}{}
}

/*
TryLock attempts to lock the mutex. If it is currently locked, TryLock
will return false immediately.
*/
func (m MutexWithDeadline) TryLock() bool {
  select {
  case m <- struct{}{}:
    return true
  default:
    return false
  }
}

/*
LockWithDeadline attempts to lock the mutex. If the lock can be acquired
before the given deadline, it returns true, otherwise returns false as
the given deadline expired.
*/
func (m MutexWithDeadline) LockWithDeadline(when time.Time) bool {
  select {
  case m <- struct{}{}:
    return true
  case <-time.After(time.Until(when)):
    return false
  }
}

/*
LockWithContext locks a mutex if the lock can be acquired while the
context passed is still active.
*/
func (m MutexWithDeadline) LockWithContext(ctx context.Context) bool {
  select {
  case m <- struct{}{}:
    return true
  case <-ctx.Done():
    return false
  }
}

/*
Unlock releases the lock currently held on the mutex. There is no check
that the locker and the unlocker are in any way related. Running Unlock
on a mutex that is not currently locked will block.
*/
func (m MutexWithDeadline) Unlock() {
  <-m
}
