package middleware

import "net/http"

//Func ...
type Func func(http.Handler) http.Handler

//Cover ...
func Cover(f http.HandlerFunc, m ...Func) http.Handler {
	// if there are no more middlewares, we just return the
	// handlerfunc, as we are done recursing.
	if len(m) == 0 {
		return f
	}
	// otherwise pop the middleware from the list,
	// and call build chain recursively as it's parameter
	return m[0](Cover(f, m[1:cap(m)]...))
}
