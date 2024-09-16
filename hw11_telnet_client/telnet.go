package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	// Place your code here.
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		addr:    address,
		timeout: timeout,
		in:      in,
		out:     out,
		ctx:     ctx,
		cancel:  cancel,
	}
}

type Client struct {
	addr    string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn

	ctx    context.Context
	cancel context.CancelFunc
}

func (c *Client) Connect() error {
	_, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	conn, err := net.DialTimeout("tcp", c.addr, c.timeout)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}

	c.conn = conn
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", c.addr)
	return nil
}

func (c *Client) Close() error {
	fmt.Fprintln(os.Stderr, "...Connection closed")
	c.cancel()
	return c.conn.Close()
}

func (c *Client) Send() error {
	scanner := bufio.NewScanner(c.in)

	for scanner.Scan() {
		select {
		case <-c.ctx.Done():
			return c.ctx.Err()

		default:
			_, err := c.conn.Write([]byte(scanner.Text() + "\n"))
			if err != nil {
				return fmt.Errorf("send message error: %w", err)
			}
		}
	}

	return scanner.Err()
}

func (c *Client) Receive() error {
	scanner := bufio.NewScanner(c.conn)

	for scanner.Scan() {
		select {
		case <-c.ctx.Done():
			return c.ctx.Err()

		default:
			_, err := fmt.Fprintln(c.out, scanner.Text())
			if err != nil {
				return fmt.Errorf("receive message error: %w", err)
			}
		}
	}

	return scanner.Err()
}
