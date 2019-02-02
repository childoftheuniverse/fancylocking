package fancylocking

import (
  "testing"
  "time"
)

func TestTryLock_Locked(t *testing.T) {
  var m = NewMutexWithDeadline()
  m.Lock()
  defer m.Unlock()

  if m.TryLock() != false {
    t.Error("TryLock succeeded on a currently locked mutex")
  }
}

func TestTryLock_UnLocked(t *testing.T) {
  var m = NewMutexWithDeadline()

  if m.TryLock() != true {
    t.Error("TryLock failed on a currently unlocked mutex")
  } else {
    m.Unlock()
  }
}

func TestLockWithDeadline_Locked(t *testing.T) {
  var startTime = time.Now()
  var endTime time.Time
  var deadline = startTime.Add(3 * time.Second)
  var m = NewMutexWithDeadline()
  m.Lock()
  defer m.Unlock()

  if m.LockWithDeadline(deadline) != false {
    t.Error("LockWithDeadline succeeded on a currently locked mutex")
  }

  endTime = time.Now()
  if endTime.Before(deadline) {
    t.Errorf("LockWithDeadline stopped attempting before deadline " +
      "(expected %s, stopped at %s)", deadline, endTime)
  }
}

func TestLockWithDeadline_Unlocked(t *testing.T) {
  var startTime = time.Now()
  var endTime time.Time
  var deadline = startTime.Add(3 * time.Second)
  var m = NewMutexWithDeadline()

  if m.LockWithDeadline(deadline) != true {
    t.Error("LockWithDeadline failed on a currently unlocked mutex")
  } else {
    m.Unlock()
  }

  endTime = time.Now()
  if !endTime.Before(deadline) {
    t.Errorf("LockWithDeadline waited until after the deadline " +
      "(deadline %s, stopped at %s)", deadline, endTime)
  }
}
