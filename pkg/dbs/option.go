package dbs

type FindOption interface {
	apply(*Option)
}

type Option struct {
	query    []Query
	order    any
	offset   int
	limit    int
	preloads []string
}

type optionFn func(*Option)

func (f optionFn) apply(opt *Option) {
	f(opt)
}

func WithQuery(query ...Query) FindOption {
	return optionFn(func(opt *Option) {
		opt.query = query
	})
}

func WithOffset(offset int) FindOption {
	return optionFn(func(opt *Option) {
		opt.offset = offset
	})
}

func WithLimit(limit int) FindOption {
	return optionFn(func(opt *Option) {
		opt.limit = limit
	})
}

func WithOrder(order interface{}) FindOption {
	return optionFn(func(opt *Option) {
		opt.order = order
	})
}

func WithPreload(preloads []string) FindOption {
	return optionFn(func(opt *Option) {
		opt.preloads = preloads
	})
}

func getOption(opts ...FindOption) Option {
	opt := Option{
		query:  []Query{},
		offset: 0,
		limit:  1000,
		order:  "id",
	}

	for _, o := range opts {
		o.apply(&opt)
	}

	return opt
}
