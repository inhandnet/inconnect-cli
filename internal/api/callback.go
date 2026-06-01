package api

import (
	"fmt"
	"html"
	"net"
	"net/http"
	"time"

	"github.com/inhandnet/ics-cli/internal/debug"
)

type CallbackResult struct {
	Code  string
	State string
}

func WaitForCallback(port int, timeout time.Duration) (*CallbackResult, error) {
	resultCh := make(chan *CallbackResult, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		debug.Log("callback: received request code=%v state=%v", code != "", state)

		if code == "" {
			errMsg := r.URL.Query().Get("error_description")
			if errMsg == "" {
				errMsg = r.URL.Query().Get("error")
			}
			if errMsg == "" {
				errMsg = "no authorization code received"
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, errorHTML, html.EscapeString(errMsg))
			errCh <- fmt.Errorf("OAuth error: %s", errMsg)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, successHTML)

		resultCh <- &CallbackResult{Code: code, State: state}
	})

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to start callback server on port %d: %w", port, err)
	}

	go func() { _ = server.Serve(ln) }()
	defer func() { _ = server.Close() }()

	select {
	case result := <-resultCh:
		time.Sleep(1 * time.Second)
		return result, nil
	case err := <-errCh:
		return nil, err
	case <-time.After(timeout):
		return nil, fmt.Errorf("login timed out after %s", timeout)
	}
}

const successHTML = `<!DOCTYPE html>
<html><head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Login Successful</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;display:flex;align-items:center;justify-content:center;min-height:100vh;background:#f0fdf4;padding-bottom:20vh}
.card{text-align:center;padding:3rem;background:#fff;border-radius:16px;box-shadow:0 4px 24px rgba(0,0,0,.08)}
.icon{width:64px;height:64px;margin:0 auto 1.5rem;background:#22c55e;border-radius:50%;display:flex;align-items:center;justify-content:center}
.icon svg{width:32px;height:32px;stroke:#fff;stroke-width:3;fill:none}
h1{font-size:1.5rem;color:#111;margin-bottom:.5rem}
p{color:#666;font-size:.95rem}
</style>
</head><body>
<div class="card">
<div class="icon"><svg viewBox="0 0 24 24"><path d="M5 13l4 4L19 7"/></svg></div>
<h1>Login Successful</h1>
<p>CLI has captured the token. You can close this tab.</p>
</div>
</body></html>`

const errorHTML = `<!DOCTYPE html>
<html><head>
<meta charset="utf-8">
<title>Login Failed</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;display:flex;align-items:center;justify-content:center;min-height:100vh;background:#fef2f2;padding-bottom:20vh}
.card{text-align:center;padding:3rem;background:#fff;border-radius:16px;box-shadow:0 4px 24px rgba(0,0,0,.08)}
h1{font-size:1.5rem;color:#111;margin-bottom:.5rem}
p{color:#dc2626;font-size:.95rem}
</style>
</head><body>
<div class="card">
<h1>Login Failed</h1>
<p>%s</p>
</div>
</body></html>`
