package daenerys

type Option func(*Daenerys)

func Namespace(namespace string) Option {
	return func(o *Daenerys) {
		o.Namespace = namespace
	}
}
func ConfigPath(path string) Option {
	return func(o *Daenerys) {
		o.ConfigPath = path
	}
}
