package app

type Options struct {
	HTTPAddr string
}

type Option func(*Options)

func DefaultOptions() *Options {
	return &Options{}
}

func WithHTTPAddr(addr string) Option {
	return func(o *Options) {
		o.HTTPAddr = addr
	}
}

func NewOptions(opts ...Option) *Options {
	opt := DefaultOptions()
	for _, o := range opts {
		o(opt)
	}
	return opt
}
