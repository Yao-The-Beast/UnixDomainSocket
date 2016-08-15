package main

import (
    "io"
    "log"
    "net"
    "time"
    "os"
    "encoding/binary"
)

const WARMUP = 1000
const MESSAGES = 10000

func reader(r io.Reader) {
    buf := make([]byte, 1024)
    for {
        n, err := r.Read(buf[:])
        if err != nil {
            return
        }
        println("Client got:", string(buf[0:n]))
    }
}

func main() {
    //client address
    lAddr := os.Args[1]
    lAddr = "/tmp/client" + lAddr
    localAddr := net.UnixAddr{Name:lAddr, Net:"unix"}
    //server address
    sAddr := "/tmp/echoServer"
    c, err := net.DialUnix("unix", &localAddr, &net.UnixAddr{Name:sAddr, Net:"unix"})

    if err != nil {
        panic(err)
    }

    defer os.Remove(lAddr)

    time.Sleep(5 * time.Second)

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
    println("Client exists: ", os.Args[1])
    c.Close()
}