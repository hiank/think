package net

import "github.com/hiank/think/net/pb"

type client struct {
	dialer IDialer
}

func newClient(dialer IDialer) IClient {
	return &client{dialer: dialer}
}

func (cli *client) Send(*pb.Carrier) error {
	return nil
}

func (cli *client) Close() error {
	return nil
}
