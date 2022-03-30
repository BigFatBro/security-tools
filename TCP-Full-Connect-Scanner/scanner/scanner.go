package scanner

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/BigFatBro/security-tools/TCP-Full-Connect-Scanner/vars"
)

func Connect(ip string, port int) (string, int, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ip, port), 3*time.Second)
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	return ip, port, err
}

func GenerateTask(ipList []net.IP, ports []int) ([]map[string]int, int) {
	tasks := make([]map[string]int, 0)
	for _, ip := range ipList {
		for _, port := range ports {
			ipPort := map[string]int{ip.String(): port}
			tasks = append(tasks, ipPort)
		}
	}

	return tasks, len(tasks)

}

func AssigningTask(tasks []map[string]int) {
	scanBatch := len(tasks) / vars.ThreadNum
	for i := 0; i < scanBatch; i++ {
		curTask := tasks[vars.ThreadNum*i : vars.ThreadNum*(i+1)]
		RunTask(curTask)
	}
	if len(tasks)%vars.ThreadNum > 0 {
		lastTask := tasks[vars.ThreadNum*scanBatch:]
		RunTask(lastTask)
	}
}

func RunTask(tasks []map[string]int) {
	wg := &sync.WaitGroup{}
	taskChan := make(chan map[string]int, vars.ThreadNum)
	for i := 0; i < vars.ThreadNum; i++ {
		go Scan(taskChan, wg)
	}

	for _, task := range tasks {
		wg.Add(1)
		taskChan <- task
	}

	close(taskChan)
	wg.Wait()

}

func Scan(taskChan chan map[string]int, wg *sync.WaitGroup) {
	for task := range taskChan {
		for ip, port := range task {
			err := SaveResult(Connect(ip, port))
			_ = err
			wg.Done()
		}
	}
}

func SaveResult(ip string, port int, err error) error {
	if err != nil {
		return err
	}
	v, ok := vars.Result.Load(ip)
	if ok {
		ports, ok1 := v.([]int)
		if ok1 {
			ports = append(ports, port)
			vars.Result.Store(ip, ports)
		}

	} else {
		ports := make([]int, 0)
		ports = append(ports, port)
		vars.Result.Store(ip, ports)

	}

	return nil

}

func PrintResult() {
	vars.Result.Range(func(key, value interface{}) bool {
		fmt.Printf("ip:%v\n", key)
		fmt.Printf("ports: %v\n", value)
		fmt.Println(strings.Repeat("-", 100))
		return true
	})
}
