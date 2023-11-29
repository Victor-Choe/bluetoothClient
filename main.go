package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/tarm/serial"
)

func main() {
	// 시리얼 포트 설정
	config := &serial.Config{
		Name: "/dev/ttyTHS1",
		Baud: 115200,
	}

	// 시리얼 포트 열기
	port, err := serial.OpenPort(config)
	if err != nil {
		fmt.Printf("시리얼 포트 열기 오류: %v\n", err)
		return
	}
	defer port.Close()

	expectedSequence := []string{
		"Open",
		"Close",
		"Left",
		"Right",
	}

	currentIndex := 0

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	//버퍼 채널
	buf := make([]byte, 128)

	for {
		// 시리얼 포트로부터 데이터 읽기
		n, err := port.Read(buf)
		if err != nil {
			fmt.Printf("시리얼 데이터 읽기 오류: %v\n", err)
			return
		}

		// 수신한 데이터를 문자열로 변환
		data := string(buf[:n])

		// 시리얼 데이터를 받았을 때 예상 시퀀스와 비교
		if currentIndex < len(expectedSequence) && data == expectedSequence[currentIndex] {
			fmt.Printf("수신한 시리얼 데이터와 일치하는 예상 동작: %s\n", data)
			currentIndex++
		} else {
			fmt.Printf("예상 동작과 일치하지 않는 시리얼 데이터: %s\n", data)
		}

		// 모든 예상 동작이 완료되면 프로그램 종료
		if currentIndex == len(expectedSequence) {
			fmt.Println("예상 동작 시퀀스가 모두 완료되었습니다.")
			break
		}
	}

	// 시그널을 받으면 프로그램 종료
	<-signalChannel
	fmt.Println("프로그램 종료")
}
