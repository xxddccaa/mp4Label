package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/xd/mp4label/pkg/server"
)

//go:embed web
var webFS embed.FS

// 版本信息（编译时注入，默认为 v0.1）
var version = "v0.1"

func main() {
	// 如果没有参数，显示使用说明
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	// 解析子命令
	switch os.Args[1] {
	case "web":
		runWebServer()
	case "version", "--version", "-v":
		printVersion()
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Printf("未知命令: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

// 打印版本信息
func printVersion() {
	fmt.Printf("mp4Label %s\n", version)
	fmt.Println("视频标注工具 - 用于标注视频剪辑教程")
}

// 打印使用说明
func printUsage() {
	fmt.Println("mp4Label - 视频标注工具")
	fmt.Printf("版本: %s\n\n", version)
	fmt.Println("使用方式:")
	fmt.Println("  mp4label web [选项]    启动 Web 服务器")
	fmt.Println("  mp4label version       显示版本信息")
	fmt.Println("  mp4label help          显示此帮助信息")
	fmt.Println()
	fmt.Println("Web 服务器选项:")
	fmt.Println("  -port string           服务器端口 (默认: 8080)")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  mp4label web           # 在默认端口 8080 启动")
	fmt.Println("  mp4label web -port 3000  # 在端口 3000 启动")
	fmt.Println("  mp4label version       # 显示版本")
}

// 运行 Web 服务器
func runWebServer() {
	// 创建 web 子命令的 flag set
	webCmd := flag.NewFlagSet("web", flag.ExitOnError)
	port := webCmd.String("port", "8080", "服务器端口")
	
	// 解析 web 子命令的参数
	webCmd.Parse(os.Args[2:])

	// 创建服务器
	srv, err := server.NewServer(webFS)
	if err != nil {
		log.Fatalf("创建服务器失败: %v", err)
	}

	// 显示启动信息
	fmt.Printf("mp4Label %s\n", version)
	fmt.Printf("服务器启动在 http://localhost:%s\n", *port)
	fmt.Println("按 Ctrl+C 停止服务器")
	fmt.Println()

	// 启动服务器
	if err := srv.Start(*port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
