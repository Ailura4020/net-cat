# Netcat Server
A simple chat server written in Go.

## Features
* Handles multiple clients
* Broadcasts messages to all connected clients
* Logs messages and user actions (join, leave, rename)
* Supports renaming of users
* Limits the number of connected clients to 10
## Usage
To run the server, simply execute the following command:

```sh
# clone the project
git clone https://zone01normandie.org/git/clecart/net-cat.git
# build application
go build -o TCPChat
# launch server
./TCPChat [Port number]
# in a new terminal window, launch the client
nc localhost [Port number]
```
Replace `[Port number]` with the desired port number.

## Code Structure
The code is organized into several packages:

`server`: contains the main server logic, including the Server and Client structs, as well as functions for handling connections, broadcasting messages, and logging.
`main`: contains the main function, which creates a new Server instance and starts it.
## Functions
`AsciiArt()`
Returns a string containing an ASCII art logo.

`Broadcast(client, message, messagetype)`
Broadcasts a message to all connected clients. The `messagetype` parameter determines the type of message (join, leave, or regular message).

`HandleConnection(client)`
Handles incoming messages from a client and broadcasts them to all other clients.

`User(conn)`
Checks if a username is already taken and prompts the user to enter a new name if necessary.

`GestionErreur(err)`
Handles errors by printing an error message and returning.

`Run()`
Starts the server and listens for incoming connections.

## Structs
`Server`
Represents the server instance, containing a slice of Client structs and a mutex for synchronization.

`Client`
Represents a connected client, containing a `net.Conn` object, a username, and a slice of messages.

`Historic`
Represents a log entry, containing a timestamp, username, and message.

## Constants
`IP`
The IP address to listen on (localhost).

`PORT`
The port number to listen on (8081).

 ### MIT License

Copyright (c) [2024] [pukei-pukei]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.