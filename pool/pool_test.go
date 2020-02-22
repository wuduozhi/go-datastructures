package pool

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestPool_Process(t *testing.T) {
	pool := NewFunc(4, func(in interface{}) interface{} {
		intVal := in.(int)
		return intVal * 2
	})
	for i := 0; i < 10; i++ {
		ret := pool.Process(i)
		if exp, act := i*2, ret.(int); exp != act {
			t.Errorf("Wrong result: %v != %v", act, exp)
		}
		time.Sleep(1)
	}
	pool.Close()

	time.Sleep(100000)
}

func TestPool_ProcessTimed(t *testing.T) {
	pool := NewFunc(4, func(in interface{}) interface{} {
		intVal := in.(int)
		return intVal * 2
	})
	defer pool.Close()

	for i := 0; i < 10; i++ {
		ret, err := pool.ProcessTimed(i, time.Millisecond)
		if err != nil {
			t.Fatalf("Failed to process %v: %v", i, err)
		}
		if exp, act := i*2, ret.(int); exp != act {
			t.Errorf("Wrong result: %v != %v", act, exp)
		}
	}
}

func TestTimeout(t *testing.T) {
	pool := NewFunc(1, func(in interface{}) interface{} {
		intVal := in.(int)
		<-time.After(time.Millisecond)
		return intVal * 2
	})
	defer pool.Close()

	_, act := pool.ProcessTimed(10, time.Duration(1))
	if exp := ErrJobTimedOut; exp != act {
		t.Errorf("Wrong error returned: %v != %v", act, exp)
	}
}

func TestParallelJobs(t *testing.T) {
	nWorkers := 10

	pool := NewFunc(nWorkers, func(in interface{}) interface{} {
		intVal := in.(int)
		return intVal * 2
	})
	defer pool.Close()

	for i := 0; i < nWorkers*8; i++ {
		go func() {
			ret := pool.Process(10)
			if exp, act := 20, ret.(int); exp != act {
				t.Errorf("Wrong result: %v != %v", act, exp)
			}
		}()
	}

	time.Sleep(1000000)
}

func TestCallbackJob(t *testing.T) {
	pool := NewCallback(10)
	defer pool.Close()

	var counter int32
	for i := 0; i < 10; i++ {
		ret := pool.Process(func() {
			atomic.AddInt32(&counter, 1)
		})
		if ret != nil {
			t.Errorf("Non-nil callback response: %v", ret)
		}
	}

	ret := pool.Process("foo")
	if exp, act := ErrJobNotFunc, ret; exp != act {
		t.Errorf("Wrong result from non-func: %v != %v", act, exp)
	}

	if exp, act := int32(10), counter; exp != act {
		t.Errorf("Wrong result: %v != %v", act, exp)
	}
}
