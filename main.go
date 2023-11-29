package main

import (
	"fmt"
	"github.com/tarm/serial"
)

func main() {
	// 시리얼 포트 설정
	config := &serial.Config{
		Name: "/dev/ttyTHS1",
		Baud: 9600,
	}

	// 시리얼 포트 열기
	port, err := serial.OpenPort(config)
	if err != nil {
		fmt.Printf("시리얼 포트 열기 오류: %v\n", err)
		return
	}
	defer port.Close()
	//예상 시나리오
	expectedSequence := []string{
		"On\r\n",
		"Off\r\n",
		"Forward\r\n",
		"Backward\r\n",
		"Left\r\n",
		"Right\r\n",
		"Center\r\n",
		"Stop\r\n",
	}
	//버퍼 채널
	bufferChannel := make(chan string, 10)

	successCounting := 0
	failreCounting := 0
	//차후 os.Args[1]로 받아올 예정
	totalCounting := 10
	//시도 횟수
	counting := 0

	//종료 채널
	//signalChannel := make(chan os.Signal, 1)
	//signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	//토탈 카운팅만큼 반복
	for i := 0; i < totalCounting; i++ {
		buf := make([]byte, 128)
		// 데이터 수신
		n, err := port.Read(buf)
		if err != nil {
			fmt.Printf("시리얼 데이터 읽기 오류: %v\n", err)
			return
		}

		// 수신한 데이터를 문자열로 변환
		data := string(buf[:n])
		counting++
		//total
		bufferChannel <- data
	}

	receivedData := <-bufferChannel

	//데이터 비교
	if receivedData == expectedSequence[counting%totalCounting] {
		successCounting++
		fmt.Printf("Success: %v\n", successCounting)
		fmt.Printf("Expected: %v\n", expectedSequence[counting%totalCounting])
		fmt.Printf("Received: %v\n", receivedData)
	} else {
		failreCounting++
		fmt.Printf("Failure: %v\n", failreCounting)
		fmt.Printf("Expected: %v\n", expectedSequence[counting%totalCounting])
		fmt.Printf("Received: %v\n", receivedData)
	}
}
