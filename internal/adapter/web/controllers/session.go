package controllers

import "strconv"

func sessionUserID(value any) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case uint:
		return int(v), true
	case uint64:
		return int(v), true
	case float64:
		return int(v), true
	case string:
		id, err := strconv.Atoi(v)
		if err != nil {
			return 0, false
		}
		return id, true
	default:
		return 0, false
	}
}
