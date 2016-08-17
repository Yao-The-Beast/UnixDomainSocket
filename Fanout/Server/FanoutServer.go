package main

import (
    "log"
    "net"
    "os"
    "time"
    "encoding/binary"
)

const WARMUP = 1000
const MESSAGES = 10000

func echoServer(c net.Conn, index int) {
    //listen to the hello message
    buf := make([]byte, 1024)
    _, err := c.Read(buf)
    if err != nil {
        return
    }

    time.Sleep(5*time.Second)

    //send message to the client
    counter := 0
    for {
        buffer := make([]byte, 1024)
        binary.PutVarint(buffer, time.Now().UnixNano())
        _, err := c.Write(buffer)
        if err != nil {
            log.Fatal("write error:", err)
            break
        }
        time.Sleep(4 * time.Millisecond)

        if counter == WARMUP + MESSAGES {
            break
        }
        counter++
    }
    println("Done")
}

func main() {
    localAddress := "/tmp/echoServer"
    l, err := net.ListenUnix("unix",  &net.UnixAddr{Name:localAddress, Net:"unix"})
    if err != nil {
        log.Fatal("listen error:", err)
    }
    defer os.Remove(localAddress)
    index := 0
    println("Fanout Server")
    for {
        fd, err := l.Accept()
        println("Grab one client")
        if err != nil {
            log.Fatal("accept error:", err)
        }

        go echoServer(fd, index)
        index++
    }
}