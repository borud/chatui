# chatui

`chatui` is a simple chat-like user interface for typical command/response type CLI applications. 

You communicate with the UI via two channels: the `OutputCh` which is used to display output in the
main display area and `CommandCh` which is used to receive commands entered by the user in the input area of the UI.

## Sample usage

```go
package main

import (
    "log"
    "strings"
    "time"

    "github.com/borud/chatui"
)

func main() {
    outputCh := make(chan string, 10)
    commandCh := make(chan string)
   
    chatui := chatui.New(chatui.Config{
        OutputCh:     outputCh,
        CommandCh:    commandCh,
        DynamicColor: false,
        BlockCtrlC:   true,
        HistorySize:  10,
    })
   
    go func() {
        for {
            outputCh <- "hey there"
            time.Sleep(time.Second)
        }
    }()
   
    go func() {
        for command := range commandCh {
            if strings.ToLower(command) == "/quit" {
                chatui.Stop()
            }
            outputCh <- "command was: " + command
            chatui.SetStatus("last command was: " + command)
        }
    }()
   
    go func() {
        // this is done in a goroutine because it will block if the UI is not running.
        chatui.SetStatus("type /quit to exit")
    }()
   
    err := chatui.Run()
    if err != nil {
        log.Fatal(err)
    }
}   
   
```
