package net

type Resolver interface {
	Resolve() (string, error)
}

type UrlResolver struct {
	url string
}

func (u *UrlResolver) Resolve() (string, error) {
	return u.url, nil
}

func NewUrlResolver(url string) *UrlResolver {
	return &UrlResolver{url: url}
}
