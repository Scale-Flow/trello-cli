package trello

import (
	"fmt"
	"net/http"

	"github.com/brettmcdowell/trello-cli/internal/contract"
)

// mapHTTPError converts an HTTP error response to a ContractError.
func mapHTTPError(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return contract.NewError(contract.AuthInvalid, "Trello rejected the credentials")
	case http.StatusNotFound:
		return contract.NewError(contract.NotFound, "resource not found")
	case http.StatusTooManyRequests:
		return contract.NewError(contract.RateLimited, "rate limited by Trello API")
	default:
		return contract.NewError(contract.HTTPError, fmt.Sprintf("Trello API returned status %d", resp.StatusCode))
	}
}
