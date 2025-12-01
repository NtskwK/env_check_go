package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetOSInfo() string {
	hInfo, err := host.Info()
	if err != nil {
		return fmt.Sprintf("获取操作系统信息失败: %v", err)
	}
	return fmt.Sprintf("%s %s (Arch: %s)", hInfo.Platform, hInfo.PlatformVersion, runtime.GOARCH)
}

func GetCPUModel() string {
	cInfos, err := cpu.Info()
	if err != nil {
		return fmt.Sprintf("获取CPU信息失败: %v", err)
	}
	if len(cInfos) > 0 {
		return cInfos[0].ModelName
	}
	return "未知CPU型号"
}

func GetMemorySize() string {
	vMem, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Sprintf("获取内存信息失败: %v", err)
	}
	return fmt.Sprintf("%.2f GB", float64(vMem.Total)/1024/1024/1024)
}

func GetGPUInfo() []string {
	var gpuInfos []string

	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("wmic", "path", "win32_VideoController", "get", "name")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && line != "Name" {
					gpuInfos = append(gpuInfos, line)
				}
			}
		}
	case "darwin":
		cmd := exec.Command("system_profiler", "SPDisplaysDataType")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.Contains(line, "Chipset Model:") {
					parts := strings.Split(line, ":")
					if len(parts) > 1 {
						gpuInfos = append(gpuInfos, strings.TrimSpace(parts[1]))
					}
				}
			}
		}
	case "linux":
		cmd := exec.Command("lspci")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "VGA") || strings.Contains(line, "3D") || strings.Contains(line, "Display") {
					// lspci output format: "00:02.0 VGA compatible controller: Intel Corporation ..."
					// We want the part after the address and type
					parts := strings.SplitN(line, ": ", 2)
					if len(parts) > 1 {
						gpuInfos = append(gpuInfos, strings.TrimSpace(parts[1]))
					} else {
						gpuInfos = append(gpuInfos, strings.TrimSpace(line))
					}
				}
			}
		}
	}

	return gpuInfos
}
