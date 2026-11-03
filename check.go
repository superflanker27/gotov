package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	currentVersion  = "2.4.0"
	manifestURL     = "https://raw.githubusercontent.com/superflanker27/gotov/refs/heads/main/latest.json"
	updateScriptURL = "https://raw.githubusercontent.com/superflanker27/gotov/refs/heads/main/upd.go"
)

type Manifest struct {
	Version string `json:"version"`
}

func tryGet(fn func() error, attempts int) error {
	var err error
	for i := 1; i <= attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		if i == attempts {
			return err
		}
		time.Sleep(time.Duration(500*i) * time.Millisecond)
	}
	return err
}

func fetchManifest() (*Manifest, error) {
	var manifest Manifest
	client := &http.Client{Timeout: 10 * time.Second}

	err := tryGet(func() error {
		resp, err := client.Get(manifestURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("bad status: %s", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return json.Unmarshal(body, &manifest)
	}, 3)

	if err != nil {
		return nil, err
	}
	return &manifest, nil
}

func main() {

	manifest, err := fetchManifest()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch manifest: %v\n", err)
		os.Exit(1)
	}

	if manifest == nil || strings.TrimSpace(manifest.Version) == "" {
		fmt.Fprintln(os.Stderr, "Other manifest error")
		os.Exit(1)
	}

	latestVersion := strings.TrimSpace(manifest.Version)
	trimmedCurrent := strings.TrimSpace(currentVersion)

	if latestVersion == trimmedCurrent {
		fmt.Printf("Already up to date (%s)\n", currentVersion)
		return
	}

	fmt.Printf("%s updating...\n", latestVersion)

	dirName := "updater"
	if err := os.MkdirAll(dirName, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create directory '%s': %v\n", dirName, err)
		os.Exit(1)
	}

	goFile := filepath.Join(dirName, "upd.go")

	if _, err := os.Stat(goFile); err == nil {
		_ = os.Remove(goFile)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	err = tryGet(func() error {
		resp, err := client.Get(updateScriptURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return err
		}

		out, err := os.Create(goFile)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		return err
	}, 3)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to download update script: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("go", "run", goFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}
