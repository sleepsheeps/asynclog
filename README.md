# asynclog

[![Go Reference](https://pkg.go.dev/badge/github.com/yourusername/asynclog.svg)](https://pkg.go.dev/github.com/yourusername/asynclog)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/asynclog)](https://goreportcard.com/report/github.com/yourusername/asynclog)

asynclog 是一个基于 Go 1.21+ `slog` 包开发的异步日志库。它提供了异步写入功能，可以显著提高日志记录性能，同时保持了 `slog` 的简单易用特性。

## 特性

- 基于 Go 标准库 `slog` 开发
- 异步写入，提高性能
- 支持多种输出方式
  - 标准输出
  - 文件输出
  - 自定义输出
- 优雅关闭，确保日志完整写入
- 兼容 `slog.Handler` 接口

## 快速开始

````go
package main
import (
"github.com/yourusername/asynclog"
"log/slog"
)
func main() {
// 创建异步日志处理器
handler := asynclog.New(asynclog.NewStdWriter())
defer handler.Close()
// 设置为默认日志记录器
slog.SetDefault(slog.New(handler))
// 记录日志
slog.Info("Hello, asynclog!")
slog.Error("Something went wrong", "error", err)
```
````

```go
package main
import (
"github.com/yourusername/asynclog"
"log/slog"
)
func main() {
// 创建文件写入器
writer, err := asynclog.NewFileWriter("app.log")
if err != nil {
  panic(err)
}
// 创建异步日志处理器
handler := asynclog.New(writer)
defer handler.Close()
// 设置为默认日志记录器
slog.SetDefault(slog.New(handler))
// 记录日志
slog.Info("This will be written to app.log")
```
