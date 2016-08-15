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
    address := "127.0.0.1:8080"
    c, err := net.Dial("tcp",address)
    if err != nil{
        println("CONNECTION ERROR ", err.Error())
        os.Exit(1)
    }

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