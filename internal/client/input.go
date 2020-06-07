package client

import (
	"bufio"
	"fmt"
	"os"
)

func InpupChannel() chan chan interface{} {
	var (
		InputChannel = make(chan chan interface{})
		r            = bufio.NewReader(os.Stdin)
	)

	go func() {
		for {
			var s string
			fmt.Fscanf(r, "%s ", &s)
			select {
			case newInput := <-InputChannel:
				newInput <- s
			default:
			}
		}
	}()

	return InputChannel
}
