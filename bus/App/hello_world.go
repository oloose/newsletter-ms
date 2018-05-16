package main	

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("hello world")

    for {
        time.Sleep(time.Millisecond*5000)
    }
}
