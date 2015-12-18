# Sliding Window
This module implements an sliding window. But without move back.
It's like a ringbuffer with the possibility to read an arbitary cell. And remove items from the front
[godoc](https://godoc.org/github.com/djboris9/slidingwindow)

## Explain it!
```go
   w := sw.Window{}
   w.Create(19, 3) 
// Creates the following:
// ┌--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┐
// |00|00|00|00|00|00|00|00|00|00|00|00|00|00|00|00|00|00|00|00|
// |__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|
// |----------------------CAPACITY OF 19-----------------------|

   w.Add(1)
   w.Add(2)
// ┌--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┐
// |04|02|00|00|00|00|00|00|00|00|00|00|00|00|00|00|00|00|00|00|
// |__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|
// |-----|
//  Sliding Window (Size 3, but not enough items in it)

   w.Add(3)
   w.Add(4)
// ┌--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┬--┐
// |01|02|03|04|00|00|00|00|00|00|00|00|00|00|00|00|00|00|00|00|
// |__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|__|
//    |--------|
//     Sliding Window (Size 3)
```
If the capacity is exceeded, it will simply rollover, which can be an more time consuming.

## Example
```go
package main

import (
	"fmt"
	sw "github.com/djboris9/slidingwindow"
)

func main() {
	// Create a new sliding window with capacity 10 and window size 3
	w := sw.Window{}
	w.Create(10, 3)

	// Print newly created slice
	fmt.Printf("%v\n", w.Slice()) // Prints: []

	// Add two new items
	w.Add(7)
	w.Add(8)
	fmt.Printf("%v\n", w.Slice()) // Prints: [7 8]

	// Remove the last item
	w.Remove()
	fmt.Printf("%v\n", w.Slice()) // Prints: [7]

	// Bulk load two new items
	w.Load([]float64{4, 5})
	fmt.Printf("%v\n", w.Slice()) // Prints: [7 4 5]

	// Clear window
	w.Clear()
	fmt.Printf("%v\n", w.Slice()) // Prints: []
}
```

## Testing and benchmarking
Run go test on your own hardware.
```
$ go test -cover -bench=. 
PASS
BenchmarkSlice-4        	50000000	        30.3 ns/op
BenchmarkAdd100X10-4    	30000000	        53.8 ns/op
BenchmarkAdd100X50-4    	30000000	        57.5 ns/op
BenchmarkAdd100X80-4    	20000000	        63.5 ns/op
BenchmarkAdd1000X20-4   	30000000	        53.7 ns/op
BenchmarkAdd100000X20-4 	30000000	        54.9 ns/op
BenchmarkAdd100000X200-4	20000000	        55.8 ns/op
BenchmarkLoadNormal-4   	  500000	      2809 ns/op
BenchmarkLoadRollover-4 	 3000000	       512 ns/op
coverage: 100.0% of statements
ok  	github.com/djboris9/slidingwindow	14.402s
```

