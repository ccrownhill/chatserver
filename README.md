# Go chatting server and client using UDP

(Using [tview](https://github.com/rivo/tview) for the TUI)

This is a demo project for learning about Go and socket programming.

## Usage

Start server in one terminal window (or in one of your terminal multiplexer's windows):

```
go run server.go
```

In other windows start as many clients as you like giving them a name:

```
go run client.go someName
```

(it will not disallow same names)

Inside a client type in a message and then press ENTER to send it to all other clients.

Press ESC to go scroll through the chat with Vim keys (`j` for down, `k` for up).
