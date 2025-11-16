// Package main 定义了一个可独立执行的程序
// 在Go语言中，包含main函数的包必须命名为package main
package main

// import语句用于导入标准库或第三方包
// fmt包提供了格式化I/O功能，类似于C语言中的printf和scanf
import "fmt"

// main函数是程序的入口点
// 当程序运行时，Go会自动调用main函数
// main函数没有参数，也没有返回值
func main() {
	// fmt.Println函数用于向标准输出打印一行文本
	// 字符串"Hello, World!"会被打印到控制台
	// Println会在输出后自动添加换行符
	fmt.Println("Hello, World!")

	// 以下是Go语言基本语法的示例注释：

	// 变量声明示例：
	// 使用var关键字声明变量，可以指定类型，也可以让Go自动推断类型
	// var name string = "Go语言"
	// age := 10  // 简短声明，自动推断类型为int

	// 常量声明示例：
	// const PI = 3.14159

	// 函数声明示例：
	// func add(a int, b int) int {
	//     return a + b
	// }

	// 控制结构示例：
	// if-else条件语句：
	// if age >= 18 {
	//     fmt.Println("成年人")
	// } else {
	//     fmt.Println("未成年人")
	// }

	// for循环：
	// for i := 0; i < 5; i++ {
	//     fmt.Println(i)
	// }

	/*
		=== Go程序运行方式说明 ===

		Go语言提供两种主要的方式来运行程序：

		1. 直接运行源码（开发调试阶段推荐）：
		   go run main.go
		   // 作用：直接编译并执行Go源文件，不生成可执行文件
		   // 优点：快速测试，无需手动编译，适合开发阶段
		   // 缺点：每次运行都重新编译，执行速度稍慢

		2. 先编译再运行（生产环境推荐）：
		   go build main.go        // 编译生成可执行文件（Windows下生成main.exe，Linux/Mac下生成main）
		   ./main.exe              // 运行可执行文件（Windows）或 ./main（Linux/Mac）
		                           // 或者运行：go build -o hello.exe main.go 指定输出文件名
		                           // 然后运行：hello.exe
		   // 作用：将Go源代码编译成机器码的可执行文件
		   // 优点：执行速度快，可独立运行，无需Go环境
		   // 缺点：需要额外的编译步骤

		3. 其他常用的Go命令：
		   go mod tidy           // 下载并整理项目依赖
		   go install            // 编译并安装到GOPATH/bin目录
		   go test               // 运行测试文件
		   go fmt                // 格式化代码
		   go vet                // 静态分析检查潜在问题

		4. 跨平台编译：
		   GOOS=linux GOARCH=amd64 go build main.go     // 编译为Linux 64位可执行文件
		   GOOS=windows GOARCH=amd64 go build main.go   // 编译为Windows 64位可执行文件
		   GOOS=darwin GOARCH=amd64 go build main.go    // 编译为macOS 64位可执行文件
	*/
}