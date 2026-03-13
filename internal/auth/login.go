package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/Scale-Flow/trello-cli/internal/contract"
	"github.com/Scale-Flow/trello-cli/internal/credentials"
)

const trelloAuthorizeBase = "https://trello.com/1/authorize"
const trelloCallbackHost = "localhost"
const trelloCallbackPort = "3007"

// LoginResult is the response shape for auth login.
type LoginResult struct {
	Configured bool    `json:"configured"`
	AuthMode   *string `json:"authMode"`
	Member     *Member `json:"member"`
}

// BrowserOpener opens the Trello authorize URL in a browser.
type BrowserOpener func(string) error

func defaultBrowserOpener(authorizeURL string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", authorizeURL)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", authorizeURL)
	default:
		cmd = exec.Command("xdg-open", authorizeURL)
	}
	return cmd.Start()
}

// BuildAuthorizeURL builds the Trello authorization URL for interactive login.
func BuildAuthorizeURL(apiKey, callbackURL string) string {
	params := url.Values{
		"expiration":      {"never"},
		"name":            {"Trello CLI"},
		"scope":           {"read,write"},
		"response_type":   {"token"},
		"callback_method": {"fragment"},
		"key":             {apiKey},
		"return_url":      {callbackURL},
	}
	return trelloAuthorizeBase + "?" + params.Encode()
}

// CompleteLogin validates a captured token and stores credentials.
func CompleteLogin(ctx context.Context, store credentials.Store, profile, apiKey, token, baseURL, authMode string) (LoginResult, error) {
	member, err := getMember(ctx, baseURL, apiKey, token)
	if err != nil {
		return LoginResult{}, err
	}

	creds := credentials.Credentials{
		APIKey:   apiKey,
		Token:    token,
		AuthMode: authMode,
	}
	if err := store.Set(profile, creds); err != nil {
		return LoginResult{}, fmt.Errorf("failed to store credentials: %w", err)
	}

	return LoginResult{
		Configured: true,
		AuthMode:   &authMode,
		Member:     member,
	}, nil
}

// Login performs the interactive browser authorization flow.
func Login(ctx context.Context, store credentials.Store, profile, baseURL, apiKey string, openBrowser BrowserOpener, stderr io.Writer) (LoginResult, error) {
	resolvedAPIKey, err := resolveLoginAPIKey(store, profile, apiKey)
	if err != nil {
		return LoginResult{}, err
	}

	captureServer, err := newTokenCaptureServer()
	if err != nil {
		return LoginResult{}, contract.NewError(contract.HTTPError, fmt.Sprintf("failed to start local callback server: %v", err))
	}
	defer captureServer.Close()

	authorizeURL := BuildAuthorizeURL(resolvedAPIKey, captureServer.callbackURL())
	if openBrowser == nil {
		openBrowser = defaultBrowserOpener
	}
	if stderr == nil {
		stderr = io.Discard
	}
	if err := openBrowser(authorizeURL); err != nil {
		fmt.Fprintf(stderr, "Open this URL in your browser to continue login: %s\n", authorizeURL)
	}

	loginCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	token, err := captureServer.waitForToken(loginCtx)
	if err != nil {
		return LoginResult{}, err
	}

	return CompleteLogin(ctx, store, profile, resolvedAPIKey, token, baseURL, "interactive")
}

// LoginWithDeviceFlow runs the RFC 8628 device authorization flow.
func LoginWithDeviceFlow(ctx context.Context, store credentials.Store, profile, baseURL, pairingServiceURL string, stderr io.Writer) (LoginResult, error) {
	client := NewDeviceClient(pairingServiceURL)

	codeResp, err := client.RequestCode()
	if err != nil {
		return LoginResult{}, fmt.Errorf("requesting pairing code: %w", err)
	}

	fmt.Fprintf(stderr, "\nEnter this code in your Trello board's CLI Connector Power-Up:\n\n")
	fmt.Fprintf(stderr, "    %s\n\n", FormatDeviceUserCode(codeResp.UserCode))
	fmt.Fprintf(stderr, "Waiting for authorization...\n")

	interval := time.Duration(codeResp.Interval) * time.Second
	deadline := time.Now().Add(time.Duration(codeResp.ExpiresIn) * time.Second)
	networkFailures := 0
	const maxNetworkRetries = 3

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return LoginResult{}, ctx.Err()
		case <-time.After(interval):
		}

		result, err := client.PollToken(codeResp.DeviceCode)
		if err != nil {
			errMsg := err.Error()
			if errMsg == "authorization_pending" {
				networkFailures = 0
				continue
			}
			if errMsg == "slow_down" {
				interval += 5 * time.Second
				networkFailures = 0
				continue
			}
			if errMsg == "expired_token" {
				return LoginResult{}, contract.NewError(contract.AuthRequired, "pairing code expired — run 'trello auth login' to try again")
			}
			networkFailures++
			if networkFailures >= maxNetworkRetries {
				return LoginResult{}, fmt.Errorf("pairing failed after %d retries: %w", maxNetworkRetries, err)
			}
			interval = time.Duration(networkFailures) * 5 * time.Second
			continue
		}

		return CompleteLogin(ctx, store, profile, result.APIKey, result.AccessToken, baseURL, "device")
	}

	return LoginResult{}, contract.NewError(contract.AuthRequired, "pairing code expired — run 'trello auth login' to try again")
}

// FormatDeviceUserCode inserts a dash: WDJBMJHT → WDJB-MJHT.
func FormatDeviceUserCode(code string) string {
	if len(code) == 8 {
		return code[:4] + "-" + code[4:]
	}
	return code
}

func resolveLoginAPIKey(store credentials.Store, profile, apiKey string) (string, error) {
	if apiKey != "" {
		return apiKey, nil
	}

	if creds, err := store.Get(profile); err == nil && creds.APIKey != "" {
		return creds.APIKey, nil
	}

	if envKey := os.Getenv("TRELLO_API_KEY"); envKey != "" {
		return envKey, nil
	}

	return "", contract.NewError(contract.ValidationError, "interactive login requires a Trello API key via TRELLO_API_KEY or existing stored credentials")
}

type tokenCaptureServer struct {
	server   *http.Server
	listener net.Listener
	tokenCh  chan string
	once     sync.Once
}

func newTokenCaptureServer() (*tokenCaptureServer, error) {
	listener, err := net.Listen("tcp", trelloCallbackHost+":"+trelloCallbackPort)
	if err != nil {
		return nil, err
	}

	capture := &tokenCaptureServer{
		listener: listener,
		tokenCh:  make(chan string, 1),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", capture.handleCallback)
	mux.HandleFunc("/callback/token", capture.handleToken)
	capture.server = &http.Server{Handler: mux}

	go func() {
		_ = capture.server.Serve(listener)
	}()

	return capture, nil
}

func (s *tokenCaptureServer) callbackURL() string {
	return "http://" + trelloCallbackHost + ":" + trelloCallbackPort + "/callback"
}

func (s *tokenCaptureServer) waitForToken(ctx context.Context) (string, error) {
	select {
	case token := <-s.tokenCh:
		if token == "" {
			return "", contract.NewError(contract.AuthInvalid, "Trello authorization returned an empty token")
		}
		return token, nil
	case <-ctx.Done():
		return "", contract.NewError(contract.HTTPError, "timed out waiting for Trello authorization callback")
	}
}

func (s *tokenCaptureServer) Close() error {
	return s.server.Close()
}

func (s *tokenCaptureServer) handleCallback(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = io.WriteString(w, `<!doctype html>
<html>
<body>
<p>Completing Trello login...</p>
<script>
const hash = new URLSearchParams(window.location.hash.slice(1));
const token = hash.get("token");
if (!token) {
  document.body.innerHTML = "<p>Authorization token missing.</p>";
} else {
  fetch("/callback/token", {
    method: "POST",
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify({token})
  }).then(() => {
    document.body.innerHTML = "<p>Login complete. You can close this window.</p>";
  }).catch(() => {
    document.body.innerHTML = "<p>Failed to send the token back to the CLI.</p>";
  });
}
</script>
</body>
</html>`)
}

func (s *tokenCaptureServer) handleToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid token payload", http.StatusBadRequest)
		return
	}
	if payload.Token == "" {
		http.Error(w, "missing token", http.StatusBadRequest)
		return
	}

	s.once.Do(func() {
		s.tokenCh <- payload.Token
	})
	w.WriteHeader(http.StatusNoContent)
}

// NewTokenCaptureServerForTest exposes callback-server creation for tests.
func NewTokenCaptureServerForTest() (*tokenCaptureServer, error) {
	return newTokenCaptureServer()
}

// CallbackURLForTest exposes the callback URL for tests.
func (s *tokenCaptureServer) CallbackURLForTest() string {
	return s.callbackURL()
}
