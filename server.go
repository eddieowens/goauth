package goauth

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func fetchAuthCode(config *ServerConfig) string {
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	var authCode string

	server := &http.Server{Addr: ":" + strconv.Itoa(config.Port)}

	config.redirectUri = "http://localhost:" + strconv.Itoa(config.Port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authCode = onAuthRequest(w, r, config)
		waitGroup.Done()
	})

	go func() {
		_ = server.ListenAndServe()
	}()

	if err := openBrowser(config.oAuthUrl); err != nil {
		return ""
	}

	waitGroup.Wait()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)

	return authCode
}

func fetchAuthCodeAsync(config *ServerConfig, function onAuthFunction) error {
	var authCode string

	config.redirectUri = "http://localhost:" + strconv.Itoa(config.Port)

	server := &http.Server{Addr: ":" + strconv.Itoa(config.Port)}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authCode = onAuthRequest(w, r, config)
		function(authCode)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			wg.Wait()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_ = server.Shutdown(ctx)
		}()
		wg.Done()
	})

	go func() {
		_ = server.ListenAndServe()
	}()

	return openBrowser(config.oAuthUrl)
}

func onAuthRequest(w http.ResponseWriter, r *http.Request, config *ServerConfig) string {
	authCode := r.URL.Query().Get("code")
	if authCode == "" && config.AuthFailedUrl != "" {
		http.Redirect(w, r, config.AuthFailedUrl, http.StatusSeeOther)
	} else if authCode != "" && config.AuthSuccessUrl != "" {
		http.Redirect(w, r, config.AuthSuccessUrl, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "https://www.google.com/", http.StatusSeeOther)
	}

	return authCode
}
