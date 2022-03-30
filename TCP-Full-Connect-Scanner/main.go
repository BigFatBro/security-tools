package main

import (
	"log"
	"os"

	"github.com/BigFatBro/security-tools/TCP-Full-Connect-Scanner/argParser"
	"github.com/BigFatBro/security-tools/TCP-Full-Connect-Scanner/scanner"
)

func main() {
	// 1.generate scan task
	// 2.split task
	// 3.exec task group by group
	// 4.show scan result
	if len(os.Args) == 3 {
		ipArgs := os.Args[1]
		portArgs := os.Args[2]
		ipList, err := argParser.GetIpList(ipArgs)
		if err != nil {
			log.Fatal(err)
		}
		portList, err := argParser.GetPorts(portArgs)
		if err != nil {
			log.Fatal(err)
		}

		tasks, _ := scanner.GenerateTask(ipList, portList)
		scanner.AssigningTask(tasks)
		scanner.PrintResult()

	} else {
		log.Fatal("please input complete args!")
	}
}
