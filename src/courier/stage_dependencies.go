package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"
	"sync"
)

type StagedDependency struct {
	StagingDir string
	Pinned     Dependency
}

func StageDependencies(sources Manifest) (map[string]StagedDependency, error) {

	var mu sync.Mutex
	var fail bool
	stagedDeps := make(map[string]StagedDependency)

	var wg sync.WaitGroup
	wg.Add(len(sources))

	for dir, dep := range sources {
		go func(dir string, dep Dependency) {

			defer wg.Done()

			var err error
			var staged StagedDependency
			LogInfo(`Staging dependency %q`, dir)
			switch dep := dep.(type) {
			case GitDependency:
				staged, err = StageGitDependency(dep)
			case SVNDependency:
				staged, err = StageSVNDependency(dep)
			default:
				staged, err = StagedDependency{}, fmt.Errorf("Unknown dependency type '%v'", reflect.TypeOf(dep))
			}

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				LogWarn(`Error while staging dependency %q: %v`, dir, err)
				fail = true
			}

			LogInfo(`Finished staging dependency %q`, dir)
			stagedDeps[dir] = staged

		}(dir, dep)
	}

	wg.Wait()

	if fail {
		for _, stagedDep := range stagedDeps {
			LogDebug("Removing dir %q", stagedDep.StagingDir)
			if err := os.RemoveAll(stagedDep.StagingDir); err != nil {
				LogWarn("Could not clean up dir %q", stagedDep.StagingDir)
			}
		}
		return nil, errors.New("failed to stage all dependencies")
	}
	return stagedDeps, nil
}

func StageGitDependency(dep GitDependency) (staged StagedDependency, err error) {

	// Get a temp dir.
	staged.StagingDir, err = MakeTmpDir()
	if err != nil {
		return
	}

	// Clean up on exit if something went wrong.
	defer func() {
		if err != nil {
			_ = os.RemoveAll(staged.StagingDir) // If this errors out, there's not much we can do.
			staged = StagedDependency{}         // Clear staged to hide what happen from the user.
		}
	}()

	// Clone into it.
	err = GitClone(staged.StagingDir, dep.URL)
	if err != nil {
		return
	}

	// Check out the right commit.
	err = GitCheckout(staged.StagingDir, dep.Ref)
	if err != nil {
		return
	}

	// Make sure the subdirectory we want actually exits in the repo.
	var fi os.FileInfo
	fi, err = os.Stat(path.Join(staged.StagingDir, dep.Dir))
	if err != nil {
		return
	}
	if !fi.IsDir() {
		err = fmt.Errorf("%q is not a dir", path.Join(staged.StagingDir, dep.Dir))
		return
	}

	// Get the SHA1 so we can reproduce the exact version of the external.
	var pin GitDependency = dep
	pin.Ref, err = GitGetSHA1(staged.StagingDir)
	if err != nil {
		return
	}
	staged.Pinned = pin

	return
}

func StageSVNDependency(dep SVNDependency) (staged StagedDependency, err error) {

	// Get a temp dir.
	staged.StagingDir, err = MakeTmpDir()
	if err != nil {
		return
	}

	// Clean up on exit if something went wrong.
	defer func() {
		if err != nil {
			_ = os.RemoveAll(staged.StagingDir) // If this errors out, there's not much we can do.
			staged = StagedDependency{}         // Clear staged to hide what happen from the user.
		}
	}()

	// Checkout the repo.
	if dep.Rev == nil {
		err = SVNCheckoutLatest(staged.StagingDir, dep.URL)
	} else {
		err = SVNCheckoutAtRev(staged.StagingDir, dep.URL, *dep.Rev)
	}
	if err != nil {
		return
	}

	// Get checked out revision.
	var pin SVNDependency = dep
	var rev string
	rev, err = SVNVersion(staged.StagingDir)
	pin.Rev = &rev
	if err != nil {
		return
	}
	staged.Pinned = pin

	return
}
