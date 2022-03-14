package wisdom

import (
	"bufio"
	"context"
	"encoding/binary"
	"io"
	"math/big"
	"net"

	"wisdom/pkg/logger"

	"github.com/pkg/errors"
)

// Client is a client for Word of Wisdom server.
type Client struct {
	// host is hostname:port (or just hostname) where to connect.
	host string
	// maxComplexity is byte-length of challenge.
	maxComplexity int

	log logger.Logger
}

// NewClient returns new wisdom client.
func NewClient(host string, maxComplexity int, log logger.Logger) *Client {
	return &Client{
		host:          host,
		maxComplexity: maxComplexity,
		log:           log,
	}
}

// GetQuote gets a wise quote.
func (c *Client) GetQuote(ctx context.Context) (string, error) {
	conn, err := net.Dial("tcp", c.host)
	if err != nil {
		return "", errors.Wrap(err, "can't dial")
	}

	// done chan is used to prevent goroutine leak in goroutine that expects <-ctx.Done().
	done := make(chan struct{}, 1)
	defer func() {
		done <- struct{}{}
	}()

	go func() {
		select {
		case <-ctx.Done():
			err := conn.Close()
			if err != nil {
				c.log.Error(errors.Wrap(err, "can't close conn"))
			}
		case <-done:
		}
	}()

	// Buffer for challenge.
	challenge := make([]byte, c.maxComplexity)

	// Read challenge number into challenge.
	_, err = io.ReadFull(conn, challenge)
	if err != nil {
		return "", errors.Wrap(err, "can't read challenge to buffer")
	}

	challengeNumber := (&big.Int{}).SetBytes(challenge)

	factors := solve(challengeNumber)

	// answer buffer to fill it with data before sending.
	answer := bufio.NewWriter(conn)

	var factorsCount [4]byte
	binary.BigEndian.PutUint32(factorsCount[:], uint32(len(factors)))

	// Write factors count as uint32.
	_, err = answer.Write(factorsCount[:])
	if err != nil {
		return "", errors.Wrap(err, "can't write to buffer")
	}

	for _, factor := range factors {
		encoded := factor.Bytes()

		// factorSize encoded as uint32 so it requires 4 bytes.
		var factorSize [4]byte
		binary.BigEndian.PutUint32(factorSize[:], uint32(len(encoded)))

		// Write uint32 that show how many bytes used to encode factor.
		_, err = answer.Write(factorSize[:])
		if err != nil {
			return "", errors.Wrap(err, "can't write to buffer")
		}

		// Write factor itself.
		_, err = answer.Write(encoded)
		if err != nil {
			return "", errors.Wrap(err, "can't write to buffer")
		}
	}

	// Flush buffer with answer to connection.
	err = answer.Flush()
	if err != nil {
		return "", errors.Wrap(err, "can't flush answer buffer")
	}

	// quoteByteLen is encoded as uint32 value.
	var quoteByteLen [4]byte

	_, err = io.ReadFull(conn, quoteByteLen[:])
	if err != nil {
		return "", errors.Wrap(err, "can't read from conn")
	}

	quote := make([]byte, binary.BigEndian.Uint32(quoteByteLen[:]))

	// Read the quote.
	_, err = io.ReadFull(conn, quote)
	if err != nil {
		return "", errors.Wrap(err, "can't read from conn")
	}

	return string(quote), nil
}
