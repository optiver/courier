package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func trim(p []byte) string {
	return strings.TrimSpace(string(p))
}

func GitClone(dir, url string) error {
	LogDebug(`Performing Git Clone from %q to %q`, url, dir)
	cmd := exec.Command("git", "clone", url, dir)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s: %s", err.Error(), trim(out))
	}
	return nil
}

func GitCheckout(dir, ref string) error {
	LogDebug(`Performing Git Checkout in %q to %q`, dir, ref)
	cmd := exec.Command("git", "checkout", ref)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s: %s", err.Error(), trim(out))
	}
	return nil
}

func GitGetSHA1(dir string) (string, error) {
	LogDebug(`Performing Git Rev-Parse in %q`, dir)
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err.Error(), trim(out))
	}
	return trim(out), nil
}
