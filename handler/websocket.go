package handler

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var WebsocketHandler = func(res http.ResponseWriter, req *http.Request){

    var  websocketKey bytes.Buffer
    h := sha1.New()
    log.Println(req.Header.Get("Sec-Websocket-Key"))
	io.WriteString(h, fmt.Sprintf("%v258EAFA5-E914-47DA-95CA-C5AB0DC85B11", req.Header.Get("Sec-Websocket-Key")))
	encoder := base64.NewEncoder(base64.StdEncoding, &websocketKey)
	encoder.Write(h.Sum(nil))
	encoder.Close()
    log.Println(websocketKey.String())
    hj, ok := res.(http.Hijacker)
    if !ok{
        log.Printf("error hijacker")
        return
    }
    conn, bufrw, err := hj.Hijack()
    if err != nil{
        log.Println("error hijack")
        return
    }
    defer conn.Close()
    lines := []string{
        "HTTP/1.1 101 Web Socket Protocol Handshake",
        "Upgrade: WebSocket",
        "Connection: Upgrade",
        "Sec-WebSocket-Accept: " + websocketKey.String(),
        "",
        "", // required for extra CRLF 
    }
    bufrw.Write([]byte(strings.Join(lines, "\r\n")))
    bufrw.WriteString("Hello World")
    //frame := []byte{0x81, 0x05, 0x48, 0x65, 0x6c, 0x6c, 0x6f "Hello"}
    //log.Println(string(frame))
    bufrw.Write([]byte("%x1%x0%x0%x0%x1%x0%x00-7DHello"))
    bufrw.Flush()
    
}
