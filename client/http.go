package client

import (
	"github.com/go-resty/resty/v2"
	"sync"
)

var once sync.Once

var client *resty.Client

var GetClient = func() *resty.Client {
	once.Do(func() {
		client = resty.New()
	})
	return client
}
