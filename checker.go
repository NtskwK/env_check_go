package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// DotNetRuntime 表示一个 .NET Runtime 信息
type DotNetRuntime struct {
	Name    string // Runtime 名称 (如 Microsoft.NETCore.App, Microsoft.AspNetCore.App, Microsoft.WindowsDesktop.App)
	Version string // Runtime 版本
	Path    string // Runtime 安装路径
}

// GetDotNetRuntimes 获取当前环境中安装的所有 .NET Runtime（不包括 SDK）
// 支持 Windows、Linux 和 macOS
func GetDotNetRuntimes() ([]DotNetRuntime, error) {
	var runtimes []DotNetRuntime

	// 检查 dotnet 命令是否可用
	cmd := exec.Command("dotnet", "--list-runtimes")
	output, err := cmd.Output()

	if err != nil {
		// 如果命令执行失败，可能是 dotnet 未安装
		return nil, fmt.Errorf("无法执行 dotnet 命令，可能未安装 .NET\n%v", err)
	}

	// 解析输出
	// 输出格式示例:
	// Microsoft.AspNetCore.App 6.0.0 [C:\Program Files\dotnet\shared\Microsoft.AspNetCore.App]
	// Microsoft.NETCore.App 6.0.0 [C:\Program Files\dotnet\shared\Microsoft.NETCore.App]
	// Microsoft.WindowsDesktop.App 6.0.0 [C:\Program Files\dotnet\shared\Microsoft.WindowsDesktop.App]

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 解析每一行
		runtime, err := parseDotNetRuntimeLine(line)
		if err != nil {
			// 跳过无法解析的行
			continue
		}

		runtimes = append(runtimes, runtime)
	}

	return runtimes, nil
}

// parseDotNetRuntimeLine 解析 dotnet --list-runtimes 的单行输出
func parseDotNetRuntimeLine(line string) (DotNetRuntime, error) {
	// 查找路径部分（在方括号中）
	pathStart := strings.Index(line, "[")
	pathEnd := strings.Index(line, "]")

	if pathStart == -1 || pathEnd == -1 || pathEnd <= pathStart {
		return DotNetRuntime{}, fmt.Errorf("无效的 runtime 行格式: %s", line)
	}

	// 提取路径
	path := strings.TrimSpace(line[pathStart+1 : pathEnd])

	// 提取名称和版本（在路径之前的部分）
	nameVersionPart := strings.TrimSpace(line[:pathStart])
	parts := strings.Fields(nameVersionPart)

	if len(parts) < 2 {
		return DotNetRuntime{}, fmt.Errorf("无法解析 runtime 名称和版本: %s", line)
	}

	name := parts[0]
	version := parts[1]

	return DotNetRuntime{
		Name:    name,
		Version: version,
		Path:    path,
	}, nil
}

// GetDotNetRuntimesGrouped 获取按名称分组的 .NET Runtimes
// 返回一个 map，key 是 runtime 名称，value 是该 runtime 的所有版本
func GetDotNetRuntimesGrouped() (map[string][]DotNetRuntime, error) {
	runtimes, err := GetDotNetRuntimes()
	if err != nil {
		return nil, err
	}

	grouped := make(map[string][]DotNetRuntime)
	for _, rt := range runtimes {
		grouped[rt.Name] = append(grouped[rt.Name], rt)
	}

	return grouped, nil
}

// PrintDotNetRuntimes 打印所有 .NET Runtimes 信息
func PrintDotNetRuntimes() {
	fmt.Printf("\n=== .NET Runtimes ===\n")

	runtimes, err := GetDotNetRuntimes()
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	if len(runtimes) == 0 {
		fmt.Println("未检测到任何 .NET Runtime")
		return
	}

	// 按名称分组显示
	grouped := make(map[string][]DotNetRuntime)
	for _, rt := range runtimes {
		grouped[rt.Name] = append(grouped[rt.Name], rt)
	}

	for name, versions := range grouped {
		fmt.Printf("\n%s:\n", name)
		for _, rt := range versions {
			fmt.Printf("  - 版本: %s\n", rt.Version)
			fmt.Printf("    路径: %s\n", rt.Path)
		}
	}

	fmt.Printf("\n总计: %d 个 runtime\n", len(runtimes))
}

// VCRuntime 表示一个 Visual C++ Runtime 信息
type VCRuntime struct {
	Version      string // 版本号 (如 2015-2022, 2013, 2012, 2010, 2008, 2005)
	Architecture string // 架构 (x86, x64, ARM64)
	DLLPath      string // DLL 文件路径
}

// GetVCRuntimes 获取当前 Windows 系统中安装的 Visual C++ Runtimes
// 仅支持 Windows 系统
func GetVCRuntimes() ([]VCRuntime, error) {
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("此功能仅支持 Windows 系统")
	}

	var runtimes []VCRuntime

	// 定义要检查的 VC++ Runtime DLL 及其对应的版本
	// 格式: DLL名称 -> 版本标识
	vcDLLs := map[string]string{
		// VC++ 2015-2022 (这些版本共享相同的运行时)
		"vcruntime140.dll":   "2015-2022",
		"msvcp140.dll":       "2015-2022",
		"vcruntime140_1.dll": "2015-2022", // x64 版本额外的 DLL

		// VC++ 2013
		"msvcr120.dll": "2013",
		"msvcp120.dll": "2013",

		// VC++ 2012
		"msvcr110.dll": "2012",
		"msvcp110.dll": "2012",

		// VC++ 2010
		"msvcr100.dll": "2010",
		"msvcp100.dll": "2010",

		// VC++ 2008
		"msvcr90.dll": "2008",
		"msvcp90.dll": "2008",

		// VC++ 2005
		"msvcr80.dll": "2005",
		"msvcp80.dll": "2005",
	}

	// 要检查的系统目录
	systemDirs := []struct {
		path string
		arch string
	}{
		{getSystemDirectory("System32"), "x64"}, // 64位系统上的64位DLL
		{getSystemDirectory("SysWOW64"), "x86"}, // 64位系统上的32位DLL
	}

	// 用于去重的 map (版本+架构 -> VCRuntime)
	foundRuntimes := make(map[string]VCRuntime)

	for _, dir := range systemDirs {
		if dir.path == "" {
			continue
		}

		for dllName, version := range vcDLLs {
			dllPath := fmt.Sprintf("%s\\%s", dir.path, dllName)

			// 检查文件是否存在
			if fileExists(dllPath) {
				key := fmt.Sprintf("%s-%s", version, dir.arch)

				// 如果这个版本+架构组合还没有记录，或者找到了更具代表性的 DLL
				if _, exists := foundRuntimes[key]; !exists || isMainDLL(dllName) {
					foundRuntimes[key] = VCRuntime{
						Version:      version,
						Architecture: dir.arch,
						DLLPath:      dllPath,
					}
				}
			}
		}
	}

	// 将 map 转换为 slice
	for _, rt := range foundRuntimes {
		runtimes = append(runtimes, rt)
	}

	return runtimes, nil
}

// getSystemDirectory 获取 Windows 系统目录
func getSystemDirectory(subDir string) string {
	// 获取 Windows 目录 (通常是 C:\Windows)
	winDir := os.Getenv("WINDIR")
	if winDir == "" {
		winDir = "C:\\Windows"
	}

	return fmt.Sprintf("%s\\%s", winDir, subDir)
}

// fileExists 检查文件是否存在
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// isMainDLL 判断是否是主要的 DLL（用于优先显示）
func isMainDLL(dllName string) bool {
	mainDLLs := []string{
		"vcruntime140.dll",
		"msvcr120.dll",
		"msvcr110.dll",
		"msvcr100.dll",
		"msvcr90.dll",
		"msvcr80.dll",
	}

	for _, main := range mainDLLs {
		if dllName == main {
			return true
		}
	}
	return false
}

// PrintVCRuntimes 打印所有检测到的 VC++ Runtimes
func PrintVCRuntimes() {
	fmt.Printf("\n=== Visual C++ Runtimes (Windows Only) ===\n")

	if runtime.GOOS != "windows" {
		fmt.Println("此功能仅支持 Windows 系统")
		return
	}

	runtimes, err := GetVCRuntimes()
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	if len(runtimes) == 0 {
		fmt.Println("未检测到任何 Visual C++ Runtime")
		return
	}

	// 按版本分组
	grouped := make(map[string][]VCRuntime)
	for _, rt := range runtimes {
		grouped[rt.Version] = append(grouped[rt.Version], rt)
	}

	// 定义版本顺序（从新到旧）
	versionOrder := []string{"2015-2022", "2013", "2012", "2010", "2008", "2005"}

	for _, version := range versionOrder {
		if archs, exists := grouped[version]; exists {
			fmt.Printf("\nVisual C++ %s Redistributable:\n", version)
			for _, rt := range archs {
				fmt.Printf("  - 架构: %s\n", rt.Architecture)
				fmt.Printf("    DLL: %s\n", rt.DLLPath)
			}
		}
	}

	fmt.Printf("\n总计: %d 个 VC++ Runtime\n", len(runtimes))
}
