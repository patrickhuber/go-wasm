package collections

func Select[TSource, TResult any](source []TSource, transform func(source TSource) (TResult, error)) ([]TResult, error) {
	results := []TResult{}
	for _, src := range source {
		result, err := transform(src)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
