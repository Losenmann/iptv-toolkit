package setup

import (
	"errors"
)

func StringToInt(s string) (int, error) {
    if len(s) == 0 {
        return 0, errors.New("empty string")
    }
    
    result := 0
    sign := 1
    start := 0
    
    // Обработка знака
    if s[0] == '-' {
        sign = -1
        start = 1
    } else if s[0] == '+' {
        start = 1
    }
    
    // Конвертация цифр
    for i := start; i < len(s); i++ {
        if s[i] < '0' || s[i] > '9' {
            return 0, errors.New("invalid character: " + string(s[i]))
        }
        result = result*10 + int(s[i]-'0')
    }
    
    return result * sign, nil
}