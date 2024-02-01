package canvas

import "strconv"

func StrStrOrNil(strPtr *string) string {
	if strPtr == nil {
		return "nil"
	}

	return *strPtr
}

func IntStrOrNil(intPtr *int) string {
	if intPtr == nil {
		return "nil"
	}

	return strconv.Itoa(*intPtr)
}

func BoolStrOrNil(boolPtr *bool) string {
	if boolPtr == nil {
		return "nil"
	}

	if *boolPtr {
		return "true"
	} else {
		return "false"
	}
}
