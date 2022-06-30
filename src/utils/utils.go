package utils

func IsElementInArray[V comparable](arr []V, element V) bool {
	for _, i := range arr {
		if i == element {
			return true
		}
	}
	return false
}

func FilterArray[V any](arr []V, check func(V) bool) []V {
	var uniqueArr []V
	for _, elem := range arr {
		if check(elem) {
			uniqueArr = append(uniqueArr, elem)
		}
	}
	return uniqueArr
}

func FilterUnique[V any](arr []V, compare func(V, V) bool) []V {
	var uniqueArr []V

	checkElemFitsInArr := func(a []V, elem V) bool {
		for _, i := range a {
			if compare(i, elem) {
				return true
			}
		}
		return false
	}

	for _, elem := range arr {
		if !checkElemFitsInArr(uniqueArr, elem) {
			uniqueArr = append(uniqueArr, elem)
		}
	}
	return uniqueArr
}

func Map[V any, R any](arr []V, apply func(V) R) []R {
	var result []R
	for _, elem := range arr {
		result = append(result, apply(elem))
	}
	return result
}
