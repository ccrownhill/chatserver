package main

import (
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"net"
	"os"
	"strings"
	"time"
)

const SERVER_PORT = "12000"
const BUFSIZE = 1024

func udpReadTimeout(conn *net.UDPConn, b []byte) (int, error) {
	conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
	n, _, _ := conn.ReadFromUDP(b)
	if n == 0 {
		return 0, errors.New("timeout")
	} else {
		return n, nil
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: client <client name>")
		return
	}

	inCtrlKey := make(chan tcell.Key, 1)

	inputHandler := func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter, tcell.KeyEscape:
			go func() {
				inCtrlKey <- event.Key()
			}()
			return nil
		default:
			return event
		}
	}

	// chat nickname
	name := os.Args[1]

	app := tview.NewApplication()
	flex := tview.NewFlex()

	s, err := net.ResolveUDPAddr("udp4", "localhost:"+SERVER_PORT)
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	c.Write([]byte("client:" + name))

	textView := tview.NewTextView().
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	textView.SetBorder(true)

	inputField := tview.NewInputField().
		SetLabel(">> ").
		SetAcceptanceFunc(tview.InputFieldMaxLength(BUFSIZE))

	flex.SetDirection(tview.FlexRow).
		AddItem(textView, 0, 100, false).
		AddItem(inputField, 0, 1, true)

	app.SetInputCapture(inputHandler)

	go func() {
		for {
			buffer := make([]byte, BUFSIZE)
			n, err := udpReadTimeout(c, buffer)
			if err == nil { // got a message
				if !(string(buffer)[:len(name)] == name) {
					fmt.Fprintf(textView, "%s\n", string(buffer[:n]))
				} else {
					fmt.Fprintf(textView, ">>> %s\n", string(buffer[len(name)+2:n]))
				}
				textView.ScrollToEnd()
			}

			select {
			case key := <-inCtrlKey:
				switch key {
				case tcell.KeyEnter:
					if inputField.HasFocus() {
						input := inputField.GetText()
						if len(input) == 0 {
							break
						}
						data := []byte(name + ": " + input + "\n")
						if strings.TrimSpace(string(input)) == "STOP" {
							app.Stop()
							_, err = c.Write([]byte(name + ":STOP"))
							if err != nil {
								fmt.Println(err)
							}
							return
						}
						_, err = c.Write(data)
						if err != nil {
							app.Stop()
							fmt.Println(err)
							return
						}
						inputField.SetText("")
					} else {
						app.SetFocus(inputField)
					}
				case tcell.KeyESC:
					app.SetFocus(textView)
				}
				app.Draw()
			default:
			}

		}
	}()
	err = app.SetRoot(flex, true).Run()
	if err != nil {
		panic(err)
	}
	defer app.Stop()
}
