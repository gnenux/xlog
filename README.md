## introduction
The log has made some simple encapsulation of the standard library log.

## Usage

to file:
```
package main

import (
	"github.com/gnenux/xlog"
)

func main() {
	logger := xlog.NewLoggerFromFileName("xlog.log", "", xlog.DefaultLogFlag)
	logger.Warn("this is warn")
	logger.Info("this is info")
	logger.Error("this is error")
	logger.Fatal("this is fatal")
}

```
output:

xlog.log

```
2018/10/16 17:51:26.863411 xlog.go:66: [warn] this is warn
2018/10/16 17:51:26.863487 xlog.go:88: [info] this is info
2018/10/16 17:51:26.863494 xlog.go:110: [error] this is error
2018/10/16 17:51:26.863499 xlog.go:132: [fatal] this is fatal

```

to stdout:
```
package main

import (
	"os"

	"github.com/gnenux/xlog"
)

func main() {
	logger := xlog.NewLogger(os.Stdout, "", xlog.DefaultLogFlag)
	logger.Warn("this is warn")
	logger.Info("this is info")
	logger.Error("this is error")
	logger.Fatal("this is fatal")
}

```
output:
```
2018/10/16 17:53:45.024912 xlog.go:66: [warn] this is warn
2018/10/16 17:53:45.024975 xlog.go:88: [info] this is info
2018/10/16 17:53:45.024981 xlog.go:110: [error] this is error
2018/10/16 17:53:45.024989 xlog.go:132: [fatal] this is fatal
exit status 1
```
## Doc

xlog:https://godoc.org/xlog
