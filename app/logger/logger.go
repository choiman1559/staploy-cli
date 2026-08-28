package logger

import (
	"fmt"
	"staploy-cli/app/consts"
)

var NoColor bool
var IsTree bool
var LastTree bool

func InitLogger(noColor bool) {
	NoColor = noColor
	IsTree = false
}

const (
	colorReset   = "\033[0m"
	colorBold    = "\033[1m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorDim     = "\033[2m"
)

func EnableTree() {
	IsTree = true
}

func DisableTree(lastTree bool) {
	IsTree = false
	LastTree = lastTree
}

func printTree() string {
	if IsTree {
		return " ├─ "
	} else if LastTree {
		LastTree = false
		return " └─ "
	}
	return ""
}

func printf(format string, a ...interface{}) {
	fmt.Println(printTree() + fmt.Sprintf(format, a...))
}

func Task(format string, a ...interface{}) {
	if NoColor {
		printf(consts.TaskPrefix+" "+format, a...)
		return
	}

	prefix := colorBlue + consts.TaskPrefix + colorReset + " "
	printf(prefix+format, a...)
}

func Process(format string, a ...interface{}) {
	if NoColor {
		printf(consts.ProcessPrefix+" "+format, a...)
		return
	}

	prefix := colorCyan + consts.ProcessPrefix + colorReset + " "
	printf(prefix+format, a...)
}

func Info(format string, a ...interface{}) {
	if NoColor {
		printf(consts.InfoPrefix+" "+format, a...)
		return
	}

	prefix := colorGreen + consts.InfoPrefix + colorReset + " "
	printf(prefix+format, a...)
}

func Tip(format string, a ...interface{}) {
	if NoColor {
		printf(consts.SkipPrefix+" "+format, a...)
		return
	}

	prefix := colorDim + consts.SkipPrefix + colorReset + " "
	printf(prefix+colorDim+format+colorReset, a...)
}

func Hook(format string, a ...interface{}) {
	if NoColor {
		printf(consts.HookPrefix+" "+format, a...)
		return
	}

	prefix := colorMagenta + consts.HookPrefix + colorReset + " "
	printf(prefix+format, a...)
}

func Warn(format string, a ...interface{}) {
	if NoColor {
		printf(consts.WarningPrefix+" "+format, a...)
		return
	}

	prefix := colorYellow + consts.WarningPrefix + colorReset + " "
	printf(prefix+colorYellow+format+colorReset, a...)
}

func Error(format string, a ...interface{}) {
	if NoColor {
		printf(consts.ErrorPrefix+" "+format, a...)
		return
	}

	prefix := colorRed + colorBold + consts.ErrorPrefix + colorReset + " "
	printf(prefix+colorRed+format+colorReset, a...)
}
