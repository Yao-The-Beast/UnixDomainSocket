package main

import (
    "io"
    "net"
    "time"
    "os"
    "encoding/binary"
    "strconv"
    "io/ioutil"
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

    //sent a hello to the server
    buffer := make([]byte, 1024)
    binary.PutVarint(buffer, time.Now().UnixNano())
    c.Write(buffer)

    //listen to the server
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
        
        if counter == MESSAGES + WARMUP {
            if os.Args[1] == "8" {
                println("Average Latency is:",latencySum / MESSAGES)
                ioutil.WriteFile("Latency", results, 0777)
            }
            break
        }
        counter++
    }

    println("Client exits: ", os.Args[1])
    c.Close()
}