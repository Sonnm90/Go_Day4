package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var wg sync.WaitGroup

//func responseSize(url string) {
//
//	defer wg.Done()
//	fmt.Println("Step1: ", url)
//	response, err := http.Get(url)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println("Step2: ", url)
//	defer response.Body.Close()
//
//	fmt.Println("Step3: ", url)
//	body, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Step4: ", len(body))
//}

func Producer(factor int, out chan<- int) {
	for i := 0; ; i++ {
		out <- i * factor
	}
}

func Consumer(in <-chan int) {
	for v := range in {
		fmt.Println(v)
	}
}

//func worker(queue chan int, workNumber int, done chan bool) {
//	for j := range queue {
//		fmt.Println("worker", workNumber, "finished job", j)
//		done <- true
//	}
//}

func worker(queue chan int, workNumber int, done, ks chan bool) {
	for true {
		// dùng select để chờ cùng lúc trên cả 2 channel
		select {
		// xử lý job trong channel queue
		case k := <-queue:
			fmt.Println("doing work!", k, "workNumber", workNumber)
			done <- true

		// nếu nhận được kill signal thì return
		case <-ks:
			fmt.Println("worker halted, number", workNumber)
			return
		}
	}
}

func CalculateValue(c chan int) {
	value := rand.Intn(10)
	fmt.Println("Calculated Random Value: {}", value)
	//time.Sleep(1000 * time.Millisecond)
	c <- value
	//fmt.Println(value)
	fmt.Println("Only Executes after another goroutine performs a receive on the channel")
}
func responseSize(url string, nums chan int) {
	// Gọi hàm Done của WaitGroup để thông báo rằng goroutine hoàn thành.
	defer wg.Done()

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	// gửi giá trị cho unbuffered channel
	nums <- len(body)
}

//	func CalculateValue(values chan int) {
//		value := rand.Intn(10)
//		fmt.Println("Calculated Random Value: {}", value)
//		values <- value
//	}
func main() {
	//Goroutine
	//go fmt.Println("Xin chào goroutine")
	//fmt.Println("Xin chào main goroutine ")

	//time.Sleep(time.Second)
	//wg.Add(20)
	//fmt.Println("Start Goroutines")
	//
	//go responseSize("https://vngeeks.com")
	//go responseSize("https://coderwall.com")
	//go responseSize("https://stackoverflow.com")
	//
	////time.Sleep(10 * time.Second)
	//
	//wg.Wait()
	//fmt.Println("Terminating Program")

	//channel
	//nums := make(chan int) // Khai báo một unbuffered channel
	//wg.Add(1)
	//go responseSize("https://vngeeks.com", nums)
	//fmt.Println(<-nums) // Đọc giá trị từ unbuffered channel
	//close(nums) // Đóng channel

	//fmt.Println("Go Channel Tutorial")
	//
	//values := make(chan int)
	//defer close(values)
	//
	//go CalculateValue(values)
	//
	//value := <-values
	//fmt.Println(value)

	fmt.Println("Go Channel Tutorial")

	valueChannel := make(chan int, 3)
	defer close(valueChannel)

	//go CalculateValue(valueChannel)
	//
	//for i := 0; i < 3; i++ {
	//	go CalculateValue(valueChannel)
	//}

	//values := <-valueChannel
	//values1 := <-valueChannel
	//values2 := <-valueChannel
	//for i := 0; i < 3; i++ {
	//	fmt.Println(<-valueChannel)
	//}
	//fmt.Println(values1)
	//	fmt.Println(values2)
	//time.Sleep(time.Second)

	// hàng đợi
	//ch := make(chan int, 64)
	//
	//// tạo một chuỗi số với bội số 3
	//go Producer(3, ch)
	//
	//// tạo một chuỗi số với bội số 5
	//go Producer(5, ch)
	//
	//// tạo consumer
	//go Consumer(ch)
	//
	//// thoát ra sau khi chạy trong một khoảng thời gian nhất định
	//time.Sleep(5 * time.Second)

	killsignal := make(chan bool)
	// queue of jobs
	workers := []string{"q1", "q2", "q3", "q4"}
	totalWork := make(chan int, 20)
	for i := 0; i < 20; i++ {
		totalWork <- i
	}

	taskList := []chan int{}
	for i := 0; i < len(workers); i++ {
		taskList = append(taskList, make(chan int, len(totalWork)/len(workers)))
	}

	// done channel lấy ra kết quả của jobs
	done := make(chan bool)

	for i := 0; i < len(taskList); i++ {
		for j := 0; j < cap(taskList[i]); j++ {
			taskList[i] <- <-totalWork
		}
	}

	for i := 0; i < len(workers); i++ {
		for j := 0; j < cap(taskList[i]); j++ {
			go worker(taskList[i], i, done, killsignal)
			time.Sleep(time.Second)
		}
	}

	// số lượng worker trong pool
	//numberOfWorkers := 4
	//for i := 0; i < numberOfWorkers; i++ {
	//	go worker(q, i, done, killsignal)
	//}
	//go worker(q1, 1, done, killsignal)
	//go worker(q2, 2, done, killsignal)
	//go worker(q3, 3, done, killsignal)
	//go worker(q4, 4, done, killsignal)

	// đưa job vào queue
	//numberOfJobs := 16
	//
	//
	//
	//for j := 0; j < 4; j++ {
	//	go func(j int) {
	//		q1 <- j
	//	}(j)
	//}
	//for j := 4; j < 8; j++ {
	//	go func(j int) {
	//		q2 <- j
	//	}(j)
	//}
	//for j := 8; j < 12; j++ {
	//	go func(j int) {
	//		q3 <- j
	//	}(j)
	//}
	//for j := 12; j < 16; j++ {
	//	go func(j int) {
	//		q4 <- j
	//	}(j)
	//}

	// chờ nhận đủ kết quả
	for c := 0; c < len(totalWork); c++ {
		<-done
	}
	time.Sleep(time.Second * 10)
	//wg.Wait()
	close(killsignal)

}
