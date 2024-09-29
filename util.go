package nestcsv

func shallowCopyMap(m map[string]any) map[string]any {
	clone := make(map[string]any)
	for k, v := range m {
		clone[k] = v
	}
	return clone
}
