package shared

func Clamp(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func MapToRanges(percent float64, ranges ...[2]float64) []int {
	result := make([]int, len(ranges))
	for i, r := range ranges {
		diff := r[1] - r[0]
		val := r[0] + (percent/100.0)*diff
		result[i] = int(val)
	}
	return result
}
