package main

import (
	"fmt"
)

func main() {
	fmt.Printf("Time: %s\n", GetCurrentTime())
	fmt.Printf("OS: %s\n", GetOSInfo())
	fmt.Printf("CPU: %s\n", GetCPUModel())
	fmt.Printf("Memory: %s\n", GetMemorySize())

	gpuInfos := GetGPUInfo()
	if len(gpuInfos) > 0 {
		fmt.Println("GPU:")
		for _, info := range gpuInfos {
			fmt.Printf("  - %s\n", info)
		}
	} else {
		fmt.Println("GPU not found")
	}

	// 检测 .NET Runtimes
	PrintDotNetRuntimes()

	// 检测 VC++ Runtimes (仅 Windows)
	PrintVCRuntimes()
}
