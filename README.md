# pubsub

Small and synchronous pubsub lib written in Go.

## Usage example

```go
package main

import (
	"fmt"
	"github.com/xrash/pubsub"
	"time"
)

func main() {
	e := pubsub.NewEngine[int](1024)

	go e.Start()
	defer e.Stop()

	subscriber := e.Subscribe("any_topic", 1024)
	defer e.Unsubscribe("any_topic", subscriber)

	go func() {
		for message := range subscriber {
			fmt.Println(message)
		}
	}()

	e.Publish("any_topic", 111222333)
	e.Publish("any_topic", 369)

	time.Sleep(time.Second * 1)
}
```
