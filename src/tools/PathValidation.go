package tools

import "net"

// написать валидацию пути в отдельной функции
func PathValidation(path string) bool {
	if len(path) < 3 || len(path) > 63 {
		return false
	}
	if path[0] == '-' || path[len(path)-1] == '-' {
		return false
	}

	for i := 0; i < len(path)-1; i++ {
		if path[i] == path[i+1] && (path[i] == '.' || path[i] == '-') {
			return false
		}
	}

	if net.ParseIP(path) != nil {
		return false
	}

	return true
}
