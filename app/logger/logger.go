package logger

import (
	"fmt"
	"staploy-cli/app/consts"
)

var NoColor bool

func InitLogger(noColor bool) {
	NoColor = noColor
}

const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorDim    = "\033[2m"
)

func Process(format string, a ...interface{}) {
	if NoColor {
		fmt.Printf(consts.ProcessPrefix+" "+format+"\n", a...)
		return
	}

	prefix := colorCyan + consts.ProcessPrefix + colorReset + " "
	fmt.Printf(prefix+format+"\n", a...)
}

func Info(format string, a ...interface{}) {
	if NoColor {
		fmt.Printf(consts.InfoPrefix+" "+format+"\n", a...)
		return
	}

	prefix := colorGreen + consts.InfoPrefix + colorReset + " "
	fmt.Printf(prefix+format+"\n", a...)
}

func Tip(format string, a ...interface{}) {
	if NoColor {
		fmt.Printf(consts.SkipPrefix+" "+format+"\n", a...)
		return
	}

	prefix := colorDim + consts.SkipPrefix + colorReset + " "
	fmt.Printf(prefix+colorDim+format+colorReset+"\n", a...)
}

func Warn(format string, a ...interface{}) {
	if NoColor {
		fmt.Printf(consts.WarningPrefix+" "+format+"\n", a...)
		return
	}

	prefix := colorYellow + consts.WarningPrefix + colorReset + " "
	fmt.Printf(prefix+colorYellow+format+colorReset+"\n", a...)
}

func Error(format string, a ...interface{}) {
	if NoColor {
		fmt.Printf(consts.ErrorPrefix+" "+format+"\n", a...)
		return
	}

	prefix := colorRed + colorBold + consts.ErrorPrefix + colorReset + " "
	fmt.Printf(prefix+colorRed+format+colorReset+"\n", a...)
}
