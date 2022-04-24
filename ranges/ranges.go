package ranges

import (
	"fmt"
	"github.com/juju/errors"
	"regexp"
	"strconv"
	"strings"
)

func ResolveIndexedRanges(inputs []string, values []string) ([]string, error) {
	result := make([]string, 0)

	rangeSplitter := regexp.MustCompile(`^([^-]+?)\s*-\s*([^-]+)$`)

	for _, input := range inputs {
		if !strings.Contains(input, "-") {
			result = append(result, input)
			continue
		}

		// range split
		parts := rangeSplitter.FindStringSubmatch(input)
		if len(parts) != 3 {
			result = append(result, input)
			continue
		}
		start := parts[1]
		end := parts[2]

		validRange := true
		startIndex, err := getIndex(values, start)
		if err != nil {
			validRange = false
		}
		endIndex, err := getIndex(values, end)
		if err != nil {
			validRange = false
		}

		// don't allow empty range as range
		if endIndex == startIndex {
			validRange = false
		}

		if !validRange {
			result = append(result, input)
			continue
		}

		// prepare for wrapping around
		valueCount := len(values)
		if endIndex < startIndex {
			endIndex += valueCount
		}
		for i := startIndex; i <= endIndex; i++ {
			result = append(result, values[i%valueCount])
		}
	}

	return result, nil
}

func ResolveIntRanges(inputs []string) ([]int, error) {
	result := make([]int, 0)

	rangeSplitter := regexp.MustCompile(`^(-?\d+?)\s*-\s*(-?\d+)$`)

	var (
		value int
		err   error
	)
	for _, input := range inputs {
		if !strings.Contains(input, "-") {
			value, err = strconv.Atoi(input)
			if err != nil {
				return result, err
			}
			result = append(result, value)
			continue
		}

		// range split
		parts := rangeSplitter.FindStringSubmatch(input)
		if len(parts) != 3 {
			value, err = strconv.Atoi(input)
			if err != nil {
				return result, err
			}
			result = append(result, value)
			continue
		}
		start := parts[1]
		end := parts[2]

		startInt, err := strconv.Atoi(start)
		if err != nil {
			return result, err
		}

		endInt, err := strconv.Atoi(end)
		if err != nil {
			return result, err
		}

		validRange := true

		// don't allow reverse and empty ranges
		if endInt <= startInt {
			validRange = false
		}

		if !validRange {
			value, err = strconv.Atoi(input)
			if err != nil {
				return result, err
			}
			result = append(result, value)
			continue
		}

		for i := startInt; i <= endInt; i++ {
			result = append(result, i)
		}
	}

	return result, nil
}

func getIndex(values []string, search string) (int, error) {
	for i, v := range values {
		if v != search {
			continue
		}
		return i, nil
	}

	return -1, errors.New(fmt.Sprintf("value %s not found in list [%s]", search, strings.Join(values, ", ")))
}
