package scalar

func stringSliceToInterfaceSlice(elems []string) []interface{} {
	interfaceElems := make([]interface{}, len(elems))
	for index, elem := range elems {
		interfaceElems[index] = elem
	}

	return interfaceElems
}
