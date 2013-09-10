![logo](https://raw.github.com/jsimnz/Devwatch/master/Documentation/devwatch.png)

### A simple command line static web server utilizing websockets for quick development of web apps
=========
## Overview
Devwatch is a (very) simple static webserver. It uses the current working directory from which it was called to server its documents. However despite these common features devwatch has one special feature. It uses [WebSockets](http://www.w3.org/TR/2009/WD-websockets-20091222/)! 

Devwatch watches the directory it serves from (and direct child folders) to alert all connected websocket clients that the directory it serves from has recieved file changes, and to refresh their browser. It uses [fsnotify](https://github.com/howeyc/fsnotify) to watch the filesystem and alert on events.

#### Features
  - Easy to use
  - Content-type detection
  - Watches filesystem for changes
  - Alert browsers with websockets

#### Cross Browser
   - Chrome
   - Firefox
   - IE 9+
  - Cross Device (Mobile/Desktop)

#### Cross-Platform
  - Linux
  - Windows
  - OSX


## Installation
You need [go](http://golang.org) installed and `GOBIN` in your `PATH`. Once that is done, run the
command:
```shell
$ go get github.com/tsenart/vegeta
$ go install github.com/tsenart/vegeta
```

## Usage
Devwatch is made to be as simple to use as possible. Just `cd` into the directory you want to serve from and run devwatch.
```shell
$ cd /to/directory
$ devwatch
```

By defualt it runs on port `:8080` this can be changed with the `--port` flag. If you want to run on port `1234` run:
```shell
$ devwatch --port 1234
```

## Websocket
As mentioned above devwatch uses websockets to alert connected clients to refresh their browser when a notification has been detected. Devwatch has built-in support for websockets, however the connected client needs to connect to the websocket endpoint located at `ws://[host]:[port]/ws/refresh`

Devwatch comes with the file `js/devwatch.js` which should be used on the client to connect to the server. It is a simple script that connects and refreshs the browser when instructed to.

#### Usage
```
<script type='text/javascript' src='/js/path/devwatch.js'></script> //import script in head
```

Check the `example` folder for further information

##TODO
 * Use `--www` as an alternative for specifying serving directory
 * Watch all children of the serving directory. Currently is limited to children 1 level down.


## Licence
```
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
```
  

    