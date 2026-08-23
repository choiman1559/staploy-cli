package build

import (
	"fmt"
	"math"
	"regexp"
	"staploy-cli/app/proto"
	"strconv"
	"strings"
)

const (
	UNIVERSE_FILTER_NOT  = "!"
	NUMERIC_FILTER_RANGE = ".."
	STRING_FILTER_REGEX  = "regex:"
	STRING_FILTER_PREFIX = "prefix:"
	STRING_FILTER_SUFFIX = "suffix:"
	STRING_FILTER_EXACT  = "exact:"
)

var memoryPattern = regexp.MustCompile(`(?i)^([0-9]+(?:\.[0-9]+)?)([kmgt]?i?b?)?$`)

func (a *StaFileTask) evalFilter(workerInfo *proto.WorkerInfo, where *Where) (bool, error) {
	if workerInfo == nil || where == nil {
		return false, fmt.Errorf("worker info or where is nil")
	}

	if where.Name != "" {
		ok, err := evalStringFilter(workerInfo.GetWorkerName(), where.Name)
		if err != nil {
			return false, fmt.Errorf("name filter: %w", err)
		}
		if !ok {
			return false, nil
		}
	}

	if where.WorkingDir != "" {
		ok, err := evalStringFilter(workerInfo.GetBinLocation(), where.WorkingDir)
		if err != nil {
			return false, fmt.Errorf("workdir filter: %w", err)
		}
		if !ok {
			return false, nil
		}
	}

	if len(where.Arch) > 0 {
		arch := workerInfo.GetCpuArch().String()

		matched := false
		for _, filter := range where.Arch {
			ok, err := evalStringFilter(arch, filter)
			if err != nil {
				return false, fmt.Errorf("arch filter %q: %w", filter, err)
			}

			if ok {
				matched = true
				break
			}
		}

		if !matched {
			return false, nil
		}
	}

	if where.Cpu != "" {
		ok, err := evalNumericFilter(workerInfo.GetCpuCoreCount(), where.Cpu)
		if err != nil {
			return false, fmt.Errorf("cpu filter: %w", err)
		}
		if !ok {
			return false, nil
		}
	}

	if where.Memory != "" {
		ok, err := evalMemoryFilter(workerInfo.GetMemoryInBytes(), where.Memory)
		if err != nil {
			return false, fmt.Errorf("memory filter: %w", err)
		}
		if !ok {
			return false, nil
		}
	}

	if where.ShellEnabled != nil {
		if workerInfo.GetWorkerFlags().GetUSE_REMOTE_SHELL() != *where.ShellEnabled {
			return false, nil
		}
	}

	if where.SkipIntegrityCheck != nil {
		integrityCheck := !workerInfo.GetWorkerFlags().GetSKIP_HASH_VERIFICATION()
		if integrityCheck != *where.SkipIntegrityCheck {
			return false, nil
		}
	}

	return true, nil
}

func evalStringFilter(value string, filter string) (bool, error) {
	negative := false

	if strings.HasPrefix(filter, UNIVERSE_FILTER_NOT) {
		negative = true
		filter = strings.TrimPrefix(filter, UNIVERSE_FILTER_NOT)

		if filter == "" {
			return false, fmt.Errorf("empty filter after %q", UNIVERSE_FILTER_NOT)
		}
	}

	matched, err := evalStringFilterPositive(value, filter)
	if err != nil {
		return false, err
	}

	if negative {
		return !matched, nil
	}

	return matched, nil
}

func evalStringFilterPositive(value string, filter string) (bool, error) {
	switch {
	case strings.HasPrefix(filter, STRING_FILTER_REGEX):
		expr := strings.TrimPrefix(filter, STRING_FILTER_REGEX)

		re, err := regexp.Compile(expr)
		if err != nil {
			return false, fmt.Errorf("invalid regex %q: %w", expr, err)
		}

		return re.MatchString(value), nil

	case strings.HasPrefix(filter, STRING_FILTER_PREFIX):
		prefix := strings.TrimPrefix(filter, STRING_FILTER_PREFIX)
		return strings.HasPrefix(value, prefix), nil

	case strings.HasPrefix(filter, STRING_FILTER_SUFFIX):
		suffix := strings.TrimPrefix(filter, STRING_FILTER_SUFFIX)
		return strings.HasSuffix(value, suffix), nil

	case strings.HasPrefix(filter, STRING_FILTER_EXACT):
		exact := strings.TrimPrefix(filter, STRING_FILTER_EXACT)
		return value == exact, nil

	default:
		return value == filter, nil
	}
}

func evalNumericFilter(value int64, filter string) (bool, error) {
	negative := false

	if strings.HasPrefix(filter, UNIVERSE_FILTER_NOT) {
		negative = true
		filter = strings.TrimPrefix(filter, UNIVERSE_FILTER_NOT)

		if filter == "" {
			return false, fmt.Errorf("empty numeric filter")
		}
	}

	rangeMin, rangeMax, err := parseNumericRange(filter)
	if err != nil {
		return false, err
	}

	matched := true

	if rangeMin != nil && value < *rangeMin {
		matched = false
	}

	if rangeMax != nil && value > *rangeMax {
		matched = false
	}

	if negative {
		return !matched, nil
	}

	return matched, nil
}

func parseNumericRange(filter string) (*int64, *int64, error) {
	if strings.Contains(filter, NUMERIC_FILTER_RANGE) {
		parts := strings.Split(filter, NUMERIC_FILTER_RANGE)

		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("invalid numeric range: %q", filter)
		}

		var rangeMin *int64
		var rangeMax *int64

		if parts[0] != "" {
			value, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid lower bound %q", parts[0])
			}
			rangeMin = &value
		}

		if parts[1] != "" {
			value, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid upper bound %q", parts[1])
			}
			rangeMax = &value
		}

		if rangeMin == nil && rangeMax == nil {
			return nil, nil, fmt.Errorf("empty numeric range: %q", filter)
		}

		if rangeMin != nil && rangeMax != nil && *rangeMin > *rangeMax {
			return nil, nil, fmt.Errorf(
				"lower bound %d is greater than upper bound %d",
				*rangeMin,
				*rangeMax,
			)
		}

		return rangeMin, rangeMax, nil
	}

	value, err := strconv.ParseInt(filter, 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid numeric value %q", filter)
	}

	return &value, &value, nil
}

func evalMemoryFilter(value int64, filter string) (bool, error) {
	negative := false

	if strings.HasPrefix(filter, UNIVERSE_FILTER_NOT) {
		negative = true
		filter = strings.TrimPrefix(filter, UNIVERSE_FILTER_NOT)

		if filter == "" {
			return false, fmt.Errorf("empty memory filter after %q", UNIVERSE_FILTER_NOT)
		}
	}

	rangeMin, rangeMax, err := parseMemoryRange(filter)
	if err != nil {
		return false, err
	}

	matched := true

	if rangeMin != nil && value < *rangeMin {
		matched = false
	}

	if rangeMax != nil && value > *rangeMax {
		matched = false
	}

	if negative {
		return !matched, nil
	}

	return matched, nil
}

func parseMemoryValue(value string) (int64, error) {
	value = strings.TrimSpace(value)

	matches := memoryPattern.FindStringSubmatch(value)
	if matches == nil {
		return 0, fmt.Errorf("invalid memory value %q", value)
	}

	number, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, err
	}

	unit := strings.ToUpper(matches[2])

	var multiplier float64 = 1

	switch unit {
	case "", "B":
		multiplier = 1

	case "K", "KB":
		multiplier = 1000

	case "M", "MB":
		multiplier = 1000 * 1000

	case "G", "GB":
		multiplier = 1000 * 1000 * 1000

	case "T", "TB":
		multiplier = 1000 * 1000 * 1000 * 1000

	case "KI", "KIB":
		multiplier = 1024

	case "MI", "MIB":
		multiplier = 1024 * 1024

	case "GI", "GIB":
		multiplier = 1024 * 1024 * 1024

	case "TI", "TIB":
		multiplier = 1024 * 1024 * 1024 * 1024

	default:
		return 0, fmt.Errorf("unsupported memory unit %q", unit)
	}

	result := number * multiplier

	if result > math.MaxInt64 {
		return 0, fmt.Errorf("memory value overflow: %q", value)
	}

	return int64(result), nil
}

func parseMemoryRange(filter string) (*int64, *int64, error) {
	if !strings.Contains(filter, NUMERIC_FILTER_RANGE) {
		value, err := parseMemoryValue(filter)
		if err != nil {
			return nil, nil, err
		}
		return &value, &value, nil
	}

	parts := strings.Split(filter, NUMERIC_FILTER_RANGE)
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid memory range %q", filter)
	}

	var rangeMin *int64
	var rangeMax *int64

	if strings.TrimSpace(parts[0]) != "" {
		value, err := parseMemoryValue(parts[0])
		if err != nil {
			return nil, nil, fmt.Errorf("invalid lower bound: %w", err)
		}
		rangeMin = &value
	}

	if strings.TrimSpace(parts[1]) != "" {
		value, err := parseMemoryValue(parts[1])
		if err != nil {
			return nil, nil, fmt.Errorf("invalid upper bound: %w", err)
		}
		rangeMax = &value
	}

	if rangeMin == nil && rangeMax == nil {
		return nil, nil, fmt.Errorf("empty memory range %q", filter)
	}

	if rangeMin != nil && rangeMax != nil && *rangeMin > *rangeMax {
		return nil, nil, fmt.Errorf("lower bound exceeds upper bound")
	}

	return rangeMin, rangeMax, nil
}
