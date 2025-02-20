package parsing

import "strconv"

// IntMustParse Returns 0 if parsing fails or int if everything is ok
func IntMustParse(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func Int32PointerMustParse(s string) *int32 {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return nil
	}
	r := int32(i)
	return &r
}

// Float64MustParse Returns 0 if parsing fails or float64 if everything is ok
func Float64MustParse(s string) float64 {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return i
}
