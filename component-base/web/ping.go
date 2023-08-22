package web

import (
	"context"
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"time"
)

// Ping test a url status. t is timeout and url just like "http://127.0.0.1:8081/healthz".
func Ping(c context.Context, url string, t time.Duration) error {
	if !govalidator.IsURL(url) {
		return fmt.Errorf("wrong url request, got: %s", url)
	}

	ctx, cancel := context.WithTimeout(c, t)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	for {
		res, err := http.DefaultClient.Do(req)
		if err == nil && res.StatusCode == http.StatusOK {
			_ = res.Body.Close()
			return nil
		}

		time.Sleep(time.Second)

		select {
		case <-c.Done():
			return fmt.Errorf("ping timeout: %s", url)
		default:
		}
	}
}
