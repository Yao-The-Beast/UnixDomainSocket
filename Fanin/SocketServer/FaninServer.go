package main

import (
    "log"
    "net"
    "os"
    "time"
    "strconv"
    "io/ioutil"
    "encoding/binary"
)

const WARMUP = 1000
const MESSAGES = 10000

func echoServer(c net.Conn, index int) {
    var results []byte
    var latencySum int64
    var counter int64
    for {
        buf := make([]byte, 1024)
        nr, err := c.Read(buf)
        if err != nil {
            return
        }

        data := buf[0:nr]
        sentTime, _ := binary.Varint(data)
        currentTime := time.Now().UnixNano()
        latency := currentTime - sentTime
        

        if counter >= WARMUP {
            latencySum += latency
            latencyByte := strconv.FormatInt(latency,10)
            results = append(results, latencyByte...)
            results = append(results, "\n"...)
        }
        
        if counter == MESSAGES + WARMUP && index == 8 {
            println("Average Latency is:",latencySum / MESSAGES)
            ioutil.WriteFile("Latency", results, 0777)
            break
        }
        counter++
    }
}

func main() {
    l, err := net.Listen("tcp", ":8080")  
    if err != nil{
        os.Exit(1)
    }
    index := 0
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