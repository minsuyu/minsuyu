package app

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

func Run() {
	for {
		//pid := int32(os.Getpid())
		pid := int32(1040)
		process_info, err := process.NewProcess(pid)
		if err != nil {
			fmt.Println("get process info fail")
			return
		}

		cpu_precent, _ := process_info.CPUPercent()
		mem_info, _ := process_info.MemoryInfo()

		fmt.Printf("cpu : %.2f%%", cpu_precent)
		fmt.Printf("\t mem : %d MB\n", mem_info.RSS/1024/1024)

		time.Sleep(1 * time.Second)
	}
}
