package auth

import "github.com/quic-go/quic-go"

type Authenticated struct {
	store quic.TokenStore
	gets  chan<- string
	puts  chan<- string
}

var _ quic.TokenStore = (*Authenticated)(nil)

func NewAuthenticated(store quic.TokenStore, gets, puts chan<- string) *Authenticated {
	return &Authenticated{
		store: store,
		gets:  gets,
		puts:  puts,
	}
}

func (a *Authenticated) Pop(key string) (token *quic.ClientToken) {
	a.gets <- key
	return a.store.Pop(key)
}

func (a *Authenticated) Put(key string, token *quic.ClientToken) {
	a.puts <- key
	a.store.Put(key, token)
}
