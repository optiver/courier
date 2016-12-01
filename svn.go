package main

import (
	"fmt"
	"os/exec"
)

func SVNCheckoutLatest(dir, url string) error {
	LogDebug(`Performing SVN Checkout Latest from %q to %q`, url, dir)
	cmd := exec.Command("svn", "checkout", "--non-interactive", url, dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), trim(out))
	}
	return nil
}

func SVNCheckoutAtRev(dir, url, rev string) error {
	LogDebug(`Performing SVN Checkout from %q at %q to %q`, url, rev, dir)
	cmd := exec.Command("svn", "checkout", "--non-interactive", "--revision", rev, url, dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), trim(out))
	}
	return nil
}

func SVNVersion(dir string) (string, error) {
	LogDebug(`Performing SVN version in %q`, dir)
	cmd := exec.Command("svnversion")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err.Error(), trim(out))
	}
	return trim(out), nil
}
