package orders

func ValidLuhn(number string) bool {
	var sum int
	alt := false
	n := len(number)

	for i := n - 1; i >= 0; i-- {
		digit := int(number[i] - '0')
		if alt {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alt = !alt
	}
	return sum%10 == 0
}
