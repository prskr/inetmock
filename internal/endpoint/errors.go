package endpoint

import (
	"errors"
	"net"
	"net/http"

	"github.com/soheilhy/cmux"
)

func IgnoreShutdownError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, http.ErrServerClosed):
		return nil
	case errors.Is(err, cmux.ErrServerClosed):
		return nil
	case errors.Is(err, net.ErrClosed):
		return nil
	default:
		return err
	}
}
