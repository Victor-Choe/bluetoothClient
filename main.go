package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tarm/serial"
)

func unescapeString(s string) string {
	s = strings.Replace(s, "\\r", "\r", -1)
	s = strings.Replace(s, "\\n", "\n", -1)
	s = strings.Replace(s, "\\t", "\t", -1)
	return s
}

func main() {
	//매개변수로 받기
	if len(os.Args) != 2 {
		fmt.Println("사용법: 프로그램명 실행횟수")
		return
	}

	// 텍스트 파일 읽기
	data, err := os.ReadFile("scenario.txt") //
	if err != nil {
		fmt.Printf("파일 읽기 오류: %v\n", err)
		return
	}
	//파싱
	lines := strings.Split(string(data), "\n")
	expectedSequence := make([]string, len(lines))
	for i, line := range lines {
		expectedSequence[i] = unescapeString(line)
	}

	// 실행 횟수 받아오기
	totalCountingStr := os.Args[1]
	totalCounting, err := strconv.Atoi(totalCountingStr)
	if err != nil {
		fmt.Println("실행횟수를 정수로 변환할 수 없습니다.")
		return
	}
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

	successCounting := 0
	failureCounting := 0

	// 실행 횟수 만큼 반복
	for i := 1; i < totalCounting+1; i++ {
		// 현재 예상 시퀀스 선택
		currentSequence := expectedSequence[(i-1)%len(expectedSequence)]

		buf := make([]byte, 128)
		// 데이터 수신
		n, err := port.Read(buf)
		if err != nil {
			fmt.Printf("시리얼 데이터 읽기 오류: %v\n", err)
			return
		}

		// 수신한 데이터를 문자열로 변환(현재 \r,\n 도 포함되어있기에 -2로 조정)
		receivedData := string(buf[:n-2])

		// 데이터 비교
		if receivedData == currentSequence {
			successCounting++
			fmt.Printf("trying %v\n", i)
			fmt.Printf("Success: %v\n", successCounting)
			fmt.Printf("Expected: %v\n", currentSequence)
			fmt.Printf("Received: %v\n", receivedData)
			fmt.Println("")
		} else {
			failureCounting++
			fmt.Printf("trying %v\n", i)
			fmt.Printf("Failure: %v\n", failureCounting)
			fmt.Printf("Expected: %v\n", currentSequence)
			fmt.Printf("Received: %v\n", receivedData)
			fmt.Println("")
		}
	}

	// 결과 출력
	fmt.Printf("성공 횟수: %v\n", successCounting)
	fmt.Printf("실패 횟수: %v\n", failureCounting)
}
