package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	siteRoot, err := locateSiteRoot()
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}

	port := ln.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://127.0.0.1:%d/", port)

	handler := http.FileServer(http.Dir(siteRoot))
	server := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- server.Serve(ln)
	}()

	if os.Getenv("NO_OPEN_BROWSER") == "" {
		if err := openURL(url); err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Serving %s at %s", siteRoot, url)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	select {
	case <-ctx.Done():
	case err := <-srvErr:
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}

func locateSiteRoot() (string, error) {
	candidates := []string{}

	if exe, err := os.Executable(); err == nil {
		exe, _ = filepath.EvalSymlinks(exe)
		exeDir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Clean(filepath.Join(exeDir, "..", "..", "Resources", "site")),
			filepath.Clean(filepath.Join(exeDir, "..", "..", "..", "..")),
		)
	}

	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, wd)
	}

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if looksLikeSite(candidate) {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("could not find site files")
}

func looksLikeSite(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, "index.html"))
	if err != nil || info.IsDir() {
		return false
	}
	return true
}

func openURL(url string) error {
	cmd := exec.Command("open", url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("Stephen Bird Site: ")
}
