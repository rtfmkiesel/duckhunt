package speed

// InitHistory() will create a slice of the specified length
func InitHistory(length int, value int64) []int64 {
	slice := make([]int64, length)

	for i := range slice {
		// Double value here so that the
		// program will not instantly trigger
		slice[i] = value * 2
	}

	return slice
}

// CalcAvrgHistory() will calculate the average time interval of the last N keystrokes
func CalcAvrgHistory(history []int64) int64 {
	// Calc sum first
	var sum int64
	for _, element := range history {
		sum += element
	}

	// Calculate amount
	length := int64(len(history))

	// Sum/Elements = Avrg
	return sum / length
}
