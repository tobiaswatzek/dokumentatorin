package util

func Map[TSource any, TDestination any](source []TSource, mapper func(item TSource) TDestination) []TDestination {
	dest := make([]TDestination, len(source))
	for i := range source {
		dest[i] = mapper(source[i])
	}
	return dest
}
