package main

import (
	"fmt"
	"sync"
)

const max = 10

func main() {
	oddChan := make(chan struct{}, 1)
	evenChan := make(chan struct{}, 1)
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 1; i <= max; i += 2 {
			<-oddChan
			fmt.Println(i)
			evenChan <- struct{}{}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 2; i <= max; i += 2 {
			<-evenChan
			fmt.Println(i)
			oddChan <- struct{}{}
		}
	}()

	// Start the sequence
	oddChan <- struct{}{}

	wg.Wait()
}

/*

If you have use unbuffered channel and if either goroutine completes its loop before the other
(which happens in the context of the loops reaching their bounds),
the channels can become unbalanced, leading to one goroutine blocking forever waiting
for a signal that never comes, causing a deadlock.

By adding a buffer to the channels, we allow one signal to be "stored" in the channel,
preventing the goroutines from blocking when sending a signal:

Explanation
Buffered Channels:

oddChan := make(chan struct{}, 1)
evenChan := make(chan struct{}, 1)
These lines create buffered channels with a capacity of 1, meaning each channel can hold one message without blocking.
Initial Signal:

oddChan <- struct{}{}
This starts the sequence by sending the first signal to the oddChan.
Graceful Termination:

Each goroutine sends a signal to the other before completing its loop, ensuring the next number can be printed.
*/
