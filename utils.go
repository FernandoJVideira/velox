package velox

import (
	"fmt"
	"regexp"
	"runtime"
	"time"
)

// LoadTime calculates function execution time. To use, add
// defer c.LoadTime(time.Now()) to the function body
func (v *Velox) LoadTime(start time.Time) {
	elapsed := time.Since(start)
	pc, _, _, _ := runtime.Caller(1)
	funcObj := runtime.FuncForPC(pc)
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	funcName := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	v.InfoLog.Println(fmt.Sprintf("Time to load %s: %s", funcName, elapsed))
}
