package scores

func GetIndexOfCoincidence(str string) float64 {
	var freq [26]int
	for _, c := range str {
		freq[c-'A']++
	}
	var length float64 = float64(len(str))
	var ioc float64
	for _, value := range freq {
		ioc += float64((value * (value - 1)))
	}
	ioc = ioc / (length * (length - 1))
	return ioc
}
