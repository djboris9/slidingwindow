package stockutil

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestCreate(t *testing.T) {
	w := Window{}
	if err := w.Create(10, 3); err != nil {
		t.Error("Create(10, 3) should not fail")
	}
	if err := w.Create(3, 10); err == nil {
		t.Error("Create(3, 10) should fail")
	}
}

func TestAdd(t *testing.T) {
	w := Window{}
	w.Create(4, 2)
	if s := w.Slice(); len(s) != 0 {
		t.Fail()
	}

	// Check first append
	w.Add(1)
	if s := w.Slice(); len(s) != 1 || s[0] != 1 {
		t.Errorf("%v\n", s)
	}

	// Check window moving
	w.Add(2)
	if s := w.Slice(); len(s) != 2 || s[0] != 1 || s[1] != 2 {
		t.Errorf("%v\n", s)
	}

	w.Add(3)
	if s := w.Slice(); len(s) != 2 || s[0] != 2 || s[1] != 3 {
		t.Errorf("%v\n", s)
	}

	w.Add(4)
	if s := w.Slice(); len(s) != 2 || s[0] != 3 || s[1] != 4 {
		t.Errorf("%v\n", s)
	}

	// Check going over capacity
	if getSliceHeader(&w.base).Cap != 4 {
		t.Fail()
	}

	w.Add(5)
	if s := w.Slice(); len(s) != 2 || s[0] != 4 || s[1] != 5 {
		t.Errorf("%v\n", w.base)
	}

	if getSliceHeader(&w.base).Cap != 4 {
		t.Errorf("%v\n", w.base)
		t.Errorf("Capacity was extended, shouldn't be")
	}

	// And the next after rollover
	w.Add(6)
	if s := w.Slice(); len(s) != 2 || s[0] != 5 || s[1] != 6 {
		t.Errorf("%v\n", s)
	}
}

func TestClear(t *testing.T) {
	w := Window{}
	w.Create(4, 2)
	w.Add(1)
	w.Add(2)
	w.Add(3)
	if s := w.Slice(); len(s) != 2 {
		t.Errorf("%v\n", s)
	}

	w.Clear()
	if s := w.Slice(); len(s) != 0 {
		t.Errorf("%v should be empty\n", s)
	}
}

func TestRemove(t *testing.T) {
	w := Window{}
	w.Create(4, 2)
	w.Remove()
	if s := w.Slice(); len(s) != 0 {
		t.Errorf("Empty-1 should be empty\n", s)
	}

	w.Add(1)
	w.Add(2)
	w.Add(3)
	if s := w.Slice(); len(s) != 2 {
		t.Errorf("%v\n", s)
	}

	w.Remove()
	if s := w.Slice(); len(s) != 1 || s[0] != 2 {
		t.Errorf("%v should be [2]\n", s)
	}
}

func TestLoad(t *testing.T) {
	w := Window{}
	w.Create(4, 2)
	w.Add(1)
	w.Add(2)
	w.Add(3)

	w.Load([]float64{})
	if s := w.Slice(); len(s) != 2 || s[1] != 3 {
		t.Errorf("%v should be [2, 3]\n", s)
	}

	w.Load([]float64{4, 5})
	if s := w.Slice(); len(s) != 2 || s[0] != 4 || s[1] != 5 {
		t.Errorf("%v should be [4, 5]\n", s)
	}

	w.Load([]float64{5, 6, 7})
	if s := w.Slice(); len(s) != 2 || s[0] != 6 || s[1] != 7 {
		t.Errorf("%v should be [6, 7]\n", s)
	}
}

func BenchmarkSlice(b *testing.B) {
	w := Window{}
	w.Create(2, 10)
	w.Load([]float64{42, 43, 44})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w.Slice()
	}
}

func benchAdd(i, j int, b *testing.B) {
	w := Window{}
	w.Create(i, j)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w.Add(42)
	}
}

func BenchmarkAdd100X10(b *testing.B) {
	benchAdd(100, 10, b)
}
func BenchmarkAdd100X50(b *testing.B) {
	benchAdd(100, 50, b)
}
func BenchmarkAdd100X80(b *testing.B) {
	benchAdd(100, 80, b)
}
func BenchmarkAdd1000X20(b *testing.B) {
	benchAdd(1000, 20, b)
}
func BenchmarkAdd100000X20(b *testing.B) {
	benchAdd(100000, 20, b)
}
func BenchmarkAdd100000X200(b *testing.B) {
	benchAdd(100000, 200, b)
}

func BenchmarkLoadNormal(b *testing.B) {
	ds := make([]float64, 10)
	w := Window{}
	w.Create(100, 10000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w.Load(ds)
	}
}

func BenchmarkLoadRollover(b *testing.B) {
	ds := make([]float64, 1000)
	w := Window{}
	w.Create(1500, 3000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w.Load(ds)
	}
}

func getSliceHeader(x *[]float64) *reflect.SliceHeader {
	return (*reflect.SliceHeader)(unsafe.Pointer(x))
}
