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
2018/10/25 09:59:23.387468 main.go:9: [warn] this is warn
2018/10/25 09:59:23.387540 main.go:10: [info] this is info
2018/10/25 09:59:23.387546 main.go:11: [error] this is error
2018/10/25 09:59:23.387550 main.go:12: [fatal] this is fatal

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
2018/10/25 10:00:15.779688 main.go:11: [warn] this is warn
2018/10/25 10:00:15.779750 main.go:12: [info] this is info
2018/10/25 10:00:15.779757 main.go:13: [error] this is error
2018/10/25 10:00:15.779761 main.go:14: [fatal] this is fatal
exit status 1
```
## Doc

xlog:https://godoc.org/github.com/gnenux/xlog
