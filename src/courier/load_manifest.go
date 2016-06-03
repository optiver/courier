package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

func LoadManifest(raw []byte) (Manifest, error) {

	LogDebug(`Loading Manifest %q`, string(raw))

	var manifestMap map[string]map[string]string
	err := json.Unmarshal(raw, &manifestMap)
	if err != nil {
		return Manifest{}, err
	}

	var manifest Manifest = make(map[string]Dependency)

	for dir, dep := range manifestMap {
		if vcs, ok := dep["vcs"]; !ok {
			return Manifest{}, fmt.Errorf("missing required key 'vcs' in dependency '%s'", dir)
		} else if vcs == "git" {
			gitDep, err := LoadGitDependency(dep)
			if err != nil {
				return Manifest{}, err
			}
			manifest[dir] = gitDep
		} else if vcs == "svn" {
			svnDep, err := LoadSVNDependency(dep)
			if err != nil {
				return Manifest{}, err
			}
			manifest[dir] = svnDep
		} else {
			return Manifest{}, fmt.Errorf("unknown dependency vcs '%s'", vcs)
		}
	}

	return manifest, nil
}

func LoadGitDependency(depMap map[string]string) (GitDependency, error) {
	d := GitDependency{VCS: "git"}
	var ok bool
	if d.URL, ok = depMap["url"]; !ok {
		return GitDependency{}, errors.New("missing required key 'url'")
	}
	if d.Ref, ok = depMap["ref"]; !ok {
		return GitDependency{}, errors.New("missing required key 'ref'")
	}
	if d.Dir, ok = depMap["dir"]; !ok {
		return GitDependency{}, errors.New("missing required key 'dir'")
	}
	delete(depMap, "vcs")
	delete(depMap, "url")
	delete(depMap, "ref")
	delete(depMap, "dir")
	for k, v := range depMap {
		LogWarn(`Ignoring unknown key value pair %q:%q`, k, v)
	}
	return d, nil
}

func LoadSVNDependency(depMap map[string]string) (SVNDependency, error) {
	d := SVNDependency{VCS: "svn"}
	var ok bool
	if d.URL, ok = depMap["url"]; !ok {
		return SVNDependency{}, errors.New("missing required key 'url'")
	}
	if _, ok = depMap["rev"]; ok {
		r := depMap["rev"]
		d.Rev = &r
	}
	delete(depMap, "vcs")
	delete(depMap, "url")
	delete(depMap, "rev")
	for k, v := range depMap {
		LogWarn(`Ignoring unknown key value pair %q:%q`, k, v)
	}
	return d, nil
}
