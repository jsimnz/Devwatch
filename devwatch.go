/*
The BSD 3-Clause

Copyright (c) 2013 John-Alan Simmons

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

  Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

  Redistributions in binary form must reproduce the above copyright notice, this
  list of conditions and the following disclaimer in the documentation and/or
  other materials provided with the distribution.

  Neither the name of the {organization} nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"fmt"
	"github.com/howeyc/fsnotify"
	"log"
	"net/http"
	"os"
)

// Basic websocket message/broadcast platfotm based
// on http://gary.beagledreams.com/page/go-websocket-chat.html
type hub struct {
	//Registered connections
	connections map[*connection]bool

	//Inbound messages from the connections
	broadcast chan string

	//Register requests from the connections
	register chan *connection

	//Unregister requests from connections
	unregister chan *connection
}

//Register hub
var h = hub{
	broadcast:   make(chan string),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

//Hub main loop
func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
		case m := <-h.broadcast:
			log.Printf("Broadcasting: %s\n", m)
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					delete(h.connections, c)
					close(c.send)
					go c.ws.Close()
				}
			}
		}
	}
}

//Websocket connection wrapper
type connection struct {
	//Websocket connection
	ws *websocket.Conn

	//Channel for outbound messages.
	send chan string
}

//Write broadcast msg to websocket
func (c *connection) writer() {
	for message := range c.send {
		err := websocket.Message.Send(c.ws, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

//Command args defn
var port string

func init() {
	flag.StringVar(&port, "port", "8080", "The port to listen on")
}

func main() {
	log.Println("Getting current working directory...")
	pwd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	flag.Parse()

	//Enpoint definitions
	http.Handle("/", http.FileServer(http.Dir(pwd)))
	http.Handle("/ws/refresh", websocket.Handler(wsRefreshHandler))

	go h.run()
	go watchNRefresh(pwd)

	log.Printf("Starting server at localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// Watch the current directory (and direct child dirs) for changes
// and notify all connected clients via websocket to refresh browser
func watchNRefresh(pwd string) {
	directory, err := os.Open(pwd)
	if err != nil {
		log.Panic(err)
	}

	dirs, err := directory.Readdir(0)
	if err != nil {
		log.Panic(err)
	}

	done := make(chan bool)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panic(err)
	}

	// Process events
	go func() {
		for {
			select {
			//Send refresh signal
			case ev := <-watcher.Event:
				log.Println("event:", ev)
				h.broadcast <- "REFRESH"
			//send kill signal
			case err := <-watcher.Error:
				log.Println("error:", err)
				h.broadcast <- "KILL"
				done <- true
			}
		}
	}()

	//Register current directory and direct child dirs
	err = watcher.Watch(pwd)
	if err != nil {
		log.Panic(err)
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			err = watcher.Watch(dir.Name())
			if err != nil {
				log.Panic(err)
			}
		}
	}

	<-done
	watcher.Close()
}

//Websocket connection handler
func wsRefreshHandler(ws *websocket.Conn) {
	//This buffered channel is overkill
	c := &connection{send: make(chan string, 256), ws: ws}
	h.register <- c
	defer func() { h.unregister <- c }()
	c.writer()
}
