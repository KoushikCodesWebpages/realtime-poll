package utils

func Flatten(prefix string, input map[string]interface{}, out map[string]interface{}) {
	for k, v := range input {

		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		if nested, ok := v.(map[string]interface{}); ok {
			Flatten(key, nested, out)
		} else {
			out[key] = v
		}
	}
}