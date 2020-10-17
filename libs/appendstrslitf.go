package libs

import "fmt"

func AppendStringSliceWithInterface(strSl *[]string, data interface{}) (err error) {
	switch castedData := data.(type) {
	case string:
		*strSl = append(*strSl, castedData)
		return nil

	case []string:
		*strSl = append(*strSl, castedData...)
		return nil

	case []interface{}:
		for castedDataIndex := range castedData {
			switch casted := castedData[castedDataIndex].(type) {
			case string:
				*strSl = append(*strSl, casted)

			default:
				return fmt.Errorf("Can not casted []interfaces. Data: %v, Type: %T", casted, casted)
			}
		}

	default:
		return fmt.Errorf("can not append to string slice. Data: %v, Type: %T", data, data)
	}

	return nil
}
