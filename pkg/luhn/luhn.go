package luhn

func ValidateLuhn(number int) bool {
	if number < 0 {
		return false
	}

	var sum int
	var isEvenOrder bool
	for number > 0 {
		cur := number % 10

		if isEvenOrder {
			cur *= 2
			cur = cur/10 + cur%10
		}

		sum += cur
		number /= 10
		isEvenOrder = !isEvenOrder
	}
	return sum%10 == 0
}
