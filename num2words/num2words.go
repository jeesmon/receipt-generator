package num2words

import (
	"strings"
)

var lowNames = []string{"zero", "One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen"}

var tensNames = []string{"Twenty", "Thirty", "Forty", "Fifty", "Sixty", "Seventy", "Eighty", "Ninety"}

var bigNames = []string{"Thousand", "Million", "Billion"}

func convert999(num int) string {
	s1 := lowNames[num/100] + " Hundred"
	s2 := convert99(num % 100)
	if num <= 99 {
		return s2
	}
	if num%100 == 0 {
		return s1
	} else {
		return s1 + " " + s2
	}
}

func convert99(num int) string {
	if num < 20 {
		return lowNames[num]
	}
	s := tensNames[num/10-2]
	if num%10 == 0 {
		return s
	}
	return s + "-" + lowNames[num%10]
}

func ConvertNum2Words(num int) string {
	if num < 0 {
		return strings.TrimSpace("negative " + ConvertNum2Words(-num))
	}

	if num <= 999 {
		return strings.TrimSpace(convert999(num))
	}

	s := ""
	t := 0
	for num > 0 {
		if num != 0 {
			s2 := convert999(num % 1000)
			if t > 0 {
				s2 = s2 + " " + bigNames[t-1]
			}
			if s == "" {
				s = s2
			} else {
				s = s2 + " " + s
			}
		}
		num /= 1000
		t++
	}
	return strings.TrimSpace(s)
}
