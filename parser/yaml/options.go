package yaml

type Option interface {
	apply(*Parser) error
}

type optionFunc func(*Parser) error

func (f optionFunc) apply(r *Parser) error {
	return f(r)
}

func WithAdvancedTypesEnabled() Option {
	return optionFunc(func(r *Parser) error {
		r.advancedTypesEnabled = true
		return nil
	})
}

func WithPkgName(pkgName string) Option {
	return optionFunc(func(r *Parser) error {
		r.pkgName = pkgName
		return nil
	})
}

func WithOutputDir(outputDir string) Option {
	return optionFunc(func(r *Parser) error {
		r.outputDir = outputDir
		return nil
	})
}
