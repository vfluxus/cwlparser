package libs

import (
	"fmt"
)

func CastToString(data interface{}) (str string, err error) {
	switch dataCast := data.(type) {
	case string:
		str = dataCast
		return str, nil

	default:
		return "", fmt.Errorf("Can not cast to string. Data: %v. Type: %T", data, data)
	}
}
