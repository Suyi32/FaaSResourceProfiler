package profiler

import (
	"time"
	"io/ioutil"
	"log"
	"strings"
	"strconv"
	"os"
	"bufio"
	"fmt"
)

type ResourceProfiler struct {
    lastCPU float64
    Memo float64
	lastTimeStamp float64
}

func (resourceProfiler *ResourceProfiler) ReadCPU() float64 {
	data, err := ioutil.ReadFile("/sys/fs/cgroup/cpuacct/cpuacct.usage")
	if err != nil{
			log.Fatal(err)
	}
	cpu_string := strings.TrimSuffix(string(data), "\n")
	cpu_usage, err := strconv.ParseFloat(cpu_string, 64)
	//fmt.Println(cpu_usage, err, reflect.TypeOf(cpu_usage))
	
	currentCPU := float64(cpu_usage)
	currentTimeStamp := float64(time.Now().UnixNano())

	var res float64 = (currentCPU - resourceProfiler.lastCPU) / (currentTimeStamp - resourceProfiler.lastTimeStamp)

	resourceProfiler.lastTimeStamp = currentTimeStamp
	resourceProfiler.lastCPU = currentCPU

	return res
}

func (resourceProfiler *ResourceProfiler) ReadMemo() float64 {
	var memo_usage_in_bytes float64
	var cache float64

	file, err := os.Open("/sys/fs/cgroup/memory/memory.usage_in_bytes")
	if err !=  nil{
			log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
			_, err := fmt.Sscan(scanner.Text(), &memo_usage_in_bytes)
			//strconv.ParseFloat(f, 64)
			if err != nil{
					log.Fatal(err)
			}
	}
	if err := scanner.Err(); err != nil{
			log.Fatal(err)
	}
	// fmt.Printf("memo=%f, type: %T\n", memo_usage_in_bytes, memo_usage_in_bytes)

	file, err = os.Open("/sys/fs/cgroup/memory/memory.stat")
	if err != nil{
			log.Fatal(err)
	}
	defer file.Close()
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), " ")
		// fmt.Println(split)
		if split[0] == "cache"{
			_, err = fmt.Sscan(split[1], &cache)
		}
	}
	// fmt.Printf("cache=%f, type: %T\n", cache, cache)
	var curMemo float64 = (memo_usage_in_bytes - cache) / 1048576.0
	resourceProfiler.Memo = curMemo
	return curMemo
}

func NewResourceProfiler() *ResourceProfiler {
	data, err := ioutil.ReadFile("/sys/fs/cgroup/cpuacct/cpuacct.usage")
	if err != nil{
			log.Fatal(err)
	}
	cpu_string := strings.TrimSuffix(string(data), "\n")
	cpu_usage, err := strconv.ParseFloat(cpu_string, 64)
	
	return &ResourceProfiler{
		lastCPU: float64(cpu_usage),
		Memo: float64(0.0),
		lastTimeStamp: float64(time.Now().UnixNano()),
	}
}