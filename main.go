package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/sparrc/go-ping"
)

type node struct {
	Name  string `json: "name"`
	IP    string `json: "ip"`
	Nodes []node `json: "nodes"`
	ID    int    `json: "id"`
	X     int    `json: "x"`
	Y     int    `json: "y"`
}

func main() {
	nodesFile, err := os.Open("nodes.json")
	if err != nil {
		fmt.Println("Не могу прочитать файл: ", err)
	} else {
		defer nodesFile.Close()
	}
	byteValue, _ := ioutil.ReadAll(nodesFile)

	var mainNode node

	json.Unmarshal(byteValue, &mainNode)

	//fmt.Println(mainNode.Nodes[1])

	//for _, n := range mainNode.Nodes {
	//	fmt.Println(n.Name)
	//}
	var wg sync.WaitGroup

	c := make(chan *ping.Statistics, 10)

	for _, n := range mainNode.Nodes {
		wg.Add(1)
		go pingNode(n, c, &wg)
	}

	go func(c chan *ping.Statistics) {
		for s := range c {
			fmt.Println(s)
		}
	}(c)

	wg.Wait()

}

func pingNode(n node, c chan *ping.Statistics, wg *sync.WaitGroup) {
	pinger, err := ping.NewPinger(n.IP)
	if err != nil {
		fmt.Println(err)
		return
	}
	pinger.Count = 1
	pinger.Timeout = time.Second * 2
	pinger.SetPrivileged(true)
	pinger.Run()
	stats := pinger.Statistics()
	c <- stats
	wg.Done()

}
