package internal

import (
	"encoding/json"
)

// based of https://groups.google.com/g/golang-nuts/c/nLCy75zMlS8
// merge merges the two JSON-marshalable values x1 and x2,
// preferring x1 over x2 except where x1 and x2 are
// JSON objects, in which case the keys from both objects
// are included and their values merged recursively.
//
// It returns an error if x1 or x2 cannot be JSON-marshaled.
func MergeStructsToMap(x1, x2 interface{}) (interface{}, error) {
	data1, err := json.Marshal(x1)
	if err != nil {
		return nil, err
	}
	data2, err := json.Marshal(x2)
	if err != nil {
		return nil, err
	}
	var j1 interface{}
	err = json.Unmarshal(data1, &j1)
	if err != nil {
		return nil, err
	}
	var j2 interface{}
	err = json.Unmarshal(data2, &j2)
	if err != nil {
		return nil, err
	}
	return MergeMaps(j1, j2), nil
}

func MergeMaps(x1, x2 interface{}) interface{} {
	switch x1 := x1.(type) {
	case map[string]interface{}:
		x2, ok := x2.(map[string]interface{})
		if !ok {
			return x1
		}
		for k, v2 := range x2 {
			if v1, ok := x1[k]; ok {
				x1[k] = MergeMaps(v1, v2)
			} else {
				x1[k] = v2
			}
		}
	case nil:
		// merge(nil, map[string]interface{...}) -> map[string]interface{...}
		x2, ok := x2.(map[string]interface{})
		if ok {
			return x2
		}
	default:
		// log.Printf("Type went undetected %+v\n", x1)
	}
	return x1
}

func MergeJSONsToJSON(args ...[]byte) ([]byte, error) {
	if len(args) == 1 {
		return args[0], nil
	}
	var mergedMap interface{}
	json.Unmarshal(args[0], &mergedMap)
	for _, m := range args[1:] {
		var mapOther interface{}

		json.Unmarshal(m, &mapOther)

		mergedMap = MergeMaps(mergedMap, mapOther)
	}

	return json.Marshal(mergedMap)
}
func JSONToMap(data []byte) map[string]interface{} {
	var map1 map[string]interface{}
	json.Unmarshal(data, &map1)
	return map1
}
