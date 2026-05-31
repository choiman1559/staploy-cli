package logger

import (
	"fmt"
	"math"
	"strings"
)

func ShortHash(hash string) string {
	if len(hash) > 10 {
		return hash[:10] + "..."
	}
	return hash
}

func SafeBlankStr(str string) string {
	if str == "" {
		return "(none)"
	}
	return str
}

func VersionNamePrefix(str string) string {
	if str == "" {
		return "(unknown)"
	} else if strings.HasPrefix(strings.ToLower(str), "v") {
		return str
	}

	return "v" + str
}

func FormatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}

	const base = 1024
	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}

	i := math.Floor(math.Log(float64(bytes)) / math.Log(base))
	if int(i) >= len(units) {
		i = float64(len(units) - 1)
	}

	value := float64(bytes) / math.Pow(base, i)
	if i == 0 {
		return fmt.Sprintf("%.0f %s", value, units[int(i)])
	}
	return fmt.Sprintf("%.2f %s", value, units[int(i)])
}

func FormatWithCommas(num int64) string {
	str := fmt.Sprintf("%d", num)
	length := len(str)
	if length <= 3 {
		return str
	}

	var result []byte
	remainder := length % 3
	if remainder == 0 {
		remainder = 3
	}

	result = append(result, str[:remainder]...)

	for i := remainder; i < length; i += 3 {
		result = append(result, ',')
		result = append(result, str[i:i+3]...)
	}

	return string(result)
}
