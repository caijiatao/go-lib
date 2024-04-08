package interface_value

type Route interface {
	Register(func()) Route
}

type RouteImpl struct {
	fs []func()
}

func (r *RouteImpl) Register(f func()) Route {
	r.fs = append(r.fs, f)
	return r
}
