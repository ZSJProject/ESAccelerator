package main

import (
	Me "ESAccelerator/Core"
	"log"
	"os"
)

func main() {
	WaitForAppClosing := make(chan bool)
	Server := Me.OpenHTTPServer(":9200")
	Flag := <-WaitForAppClosing

	if Exception := Server.Shutdown(nil); Exception != nil {
		log.Fatalf("HTTP 서버를 종료하려던 중 예외가 발생했습니다: %s", Exception)
	}

	if Flag {
		os.Exit(0)

	} else {
		os.Exit(1)
	}
}
