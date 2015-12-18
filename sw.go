package stockutil

import (
	"errors"
	"sync"
)

// Window implements an sliding window
type Window struct {
	base       []float64
	start      int
	Len        int // Len <= Windowsize
	Windowsize int
	mx         *sync.RWMutex
}

// Create initializes the window with a capacity and an window size
//  X O O O O O X X X
//  |---CAPACITY----|
//    |-SIZE--|
//    |-start
func (w *Window) Create(capacity, windowsize int) error {
	w.mx = new(sync.RWMutex)
	w.base = make([]float64, 0, capacity)
	w.Windowsize = windowsize

	if windowsize > capacity {
		return errors.New("capacity can't be smaller than windowsize")
	}

	return nil
}

// Add appends an item to the window
func (w *Window) Add(i float64) {
	w.mx.Lock()
	w.add(i)
	w.mx.Unlock()
}

// add is like Add, but without locking
func (w *Window) add(i float64) {
	// Move all values to front, if would reach end of base
	if w.start+w.Len+1 > cap(w.base) {
		for j := 0; j < w.Len-1; j++ {
			w.base[j] = w.base[w.start+j+1]
		}
		w.start = 0
		w.Len--
	}

	// Check capacity and append if needed
	if len(w.base) < w.start+w.Len+1 {
		w.base = append(w.base, i)
	} else {
		w.base[w.start+w.Len] = i
	}

	// If window is "full" => Move one
	if w.Len == w.Windowsize {
		w.start++
	} else {
		w.Len++
	}
}

// Clear removes all items from the sliding window. Very efficient
func (w *Window) Clear() {
	w.mx.Lock()
	w.start = 0
	w.Len = 0
	w.mx.Unlock()
}

// Remove will remove the last item from the window. If the window is empty, nothing happens
func (w *Window) Remove() {
	w.mx.Lock()
	if w.Len > 0 {
		w.Len--
	}

	if w.Len == 0 {
		w.start = 0
	}
	w.mx.Unlock()
}

// Load loads all data from an slice to the window. The slice can be larger then the window
func (w *Window) Load(x []float64) {
	// Nothing to load
	if len(x) == 0 {
		return
	}

	w.mx.Lock()
	// Move to front if possible (reduces unneccessairy rollovers
	if len(x) >= w.Windowsize {
		w.start = 0
		w.Len = 0
	}

	loadlen := len(x)
	if loadlen >= w.Windowsize {
		loadlen = w.Windowsize
	}

	for _, i := range x[len(x)-loadlen : len(x)] {
		w.add(i)
	}
	w.mx.Unlock()
}

// Slice returns an slice to the window
func (w *Window) Slice() []float64 {
	w.mx.RLock()
	// 4 Times faster than "defer Unlock"
	ret := w.base[w.start : w.start+w.Len]
	w.mx.RUnlock()
	return ret
}
