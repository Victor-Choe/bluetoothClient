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

	// 예상 시나리오
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

	successCounting := 0
	failureCounting := 0

	// 토탈 카운팅만큼 반복
	for counting := 0; counting < len(expectedSequence); counting++ { // 수정: totalCounting 대신 예상 시퀀스 길이 사용
		buf := make([]byte, 128)
		// 데이터 수신
		n, err := port.Read(buf)
		if err != nil {
			fmt.Printf("시리얼 데이터 읽기 오류: %v\n", err)
			return
		}

		// 수신한 데이터를 문자열로 변환
		receivedData := string(buf[:n])

		// 데이터 비교
		if receivedData == expectedSequence[counting] {
			successCounting++
			fmt.Printf("Success: %v\n", successCounting)
			fmt.Printf("Expected: %v\n", expectedSequence[counting])
			fmt.Printf("Received: %v\n", receivedData)
		} else {
			failureCounting++
			fmt.Printf("Failure: %v\n", failureCounting)
			fmt.Printf("Expected: %v\n", expectedSequence[counting])
			fmt.Printf("Received: %v\n", receivedData)
		}
	}

	// 결과 출력
	fmt.Printf("성공 횟수: %v\n", successCounting)
	fmt.Printf("실패 횟수: %v\n", failureCounting)
}
