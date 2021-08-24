# logkit
A Golang package providing high-performance logging

## Getting start
```shell script
go get github.com:echotrue/logkit
```

## Usage
Some simple example to get started，or refer to the `example` directory
```go
    t1 := logkit.NewConsoleTarget()
    t2 := logkit.NewFileTarget()
    t2.FilePath = "./temp"
    
    log := logkit.NewLogKitByOptions(logkit.WithBuffer(1024), logkit.WithTargets(t1,t2), logkit.WithCallStackDepth(0))
    
    err := log.Open()
    if err != nil {
        log2.Fatal(err)
    }
    defer log.Close()

    log.Debug("%d--%s", time.Now().Unix(), "当前时间戳")
    log.Info("%d--%s", time.Now().Unix(), "当前时间戳")
    log.Notice("%d--%s", time.Now().Unix(), "当前时间戳")
    log.Warning("%d--%s", time.Now().Unix(), "当前时间戳")
    log.Error("%d--%s", time.Now().Unix(), "当前时间戳")
    log.Critical("%d--%s", time.Now().Unix(), "当前时间戳")
    log.Alert("%d--%s", time.Now().Unix(), "当前时间戳")
    log.Emergency("%d--%s", time.Now().Unix(), "当前时间戳")
```

## Benchmark
```shell script
$ go test -v -bench . -benchmem
=== RUN   TestNewConsoleTarget
--- PASS: TestNewConsoleTarget (0.03s)
=== RUN   TestNewFileTarget
--- PASS: TestNewFileTarget (0.00s)
goos: windows
goarch: amd64
pkg: logkit
BenchmarkNewConsoleTarget
BenchmarkNewConsoleTarget-12            1000000000             0 B/op          0 allocs/op
BenchmarkNewFileTarget
BenchmarkNewFileTarget-12               1000000000               0.00103 ns/op         0 B/op          0 allocs/op
PASS
ok      logkit  0.303s

```

## MIT License
Copyright (c) 2021 axlrose

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.