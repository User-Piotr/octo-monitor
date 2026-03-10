package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

func RunMonitoring(ctx context.Context, url string) string {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "ERROR"
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		var netErr *net.OpError
		switch {
		case ctx.Err() == context.DeadlineExceeded:
			return "TIMEOUT"
		case ctx.Err() == context.Canceled:
			return "CANCELED"
		case errors.As(err, &netErr):
			return "NETWORK_ERROR"
		default:
			return "DOWN"
		}
	}
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return "UP"
	}
	return fmt.Sprintf("DOWN(%d)", resp.StatusCode)
}
