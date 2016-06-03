package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"strings"
)

const tagged = true
const version = "1.1.0"

var cmdLineArgs struct {
	help            bool
	quiet           bool
	verbose         bool
	colour          bool
	reproduce       bool
	forceCopy       bool
	primaryManifest string
	pinnedManifest  string
}

func main() {
	err := body()
	if err != nil {
		LogError("%v", err)
		os.Exit(1)
	}
}

func body() error {

	// Set up command line args.
	flag.BoolVar(&cmdLineArgs.help, "help", false, "show usage")
	flag.BoolVar(&cmdLineArgs.quiet, "quiet", false, "don't output to stdout")
	flag.BoolVar(&cmdLineArgs.verbose, "verbose", false, "output debug information")
	flag.BoolVar(&cmdLineArgs.colour, "colour", runtime.GOOS != "windows", "use ANSI colour escape codes")
	flag.BoolVar(&cmdLineArgs.reproduce, "reproduce", false, "read from the pinned manifest instead of the primary manifest")
	flag.BoolVar(&cmdLineArgs.forceCopy, "force-copy", false, "force copying dependency even if unchanged/identical")
	flag.StringVar(&cmdLineArgs.primaryManifest, "primary-manifest", "deps.json", "location of the primary manifest")
	flag.StringVar(&cmdLineArgs.pinnedManifest, "pinned-manifest", "pins.json", "location of the pinned manifest")
	flag.Parse()

	// Show help.
	if cmdLineArgs.help {
		fmt.Printf("Courier %s\nUsage:\n", version)
		flag.PrintDefaults()
		return nil
	}

	// Set up logging.
	if err := SetupLogging(cmdLineArgs.quiet, cmdLineArgs.verbose); err != nil {
		return err
	}

	// Show version.
	if tagged {
		LogInfo("Courier version %s", version)
	} else {
		LogWarn("\u001b[41mCourier version %s\u001b[0m", version)
	}

	listenForCtrlC()

	// Determine which manifest to read from.
	var manifestFile string
	if cmdLineArgs.reproduce {
		manifestFile = cmdLineArgs.pinnedManifest
	} else {
		manifestFile = cmdLineArgs.primaryManifest
	}

	// Get the manifest.
	LogInfo("Using manifest %q", manifestFile)
	buf, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		return err
	}
	m, err := LoadManifest(buf)
	if err != nil {
		enrErr := enrichJSONError(err, string(buf))
		return enrErr
	}

	// Stage the dependencies.
	stagedDeps, err := StageDependencies(m)
	if err != nil {
		return err
	}

	// Clean up the staged dependencies when we return.
	defer func() {
		for _, stagedDep := range stagedDeps {
			_ = os.RemoveAll(stagedDep.StagingDir) // If we can't remove... then there's not much we can do.
		}
	}()

	// Copy the dependencies.
	for dir, stagedDep := range stagedDeps {
		src := path.Join(stagedDep.StagingDir, stagedDep.Pinned.DirToCopy())

		if cmdLineArgs.forceCopy {
			LogInfo(`Copying dependency %q (forced)`, dir)
			if err := CopyDirContents(src, dir, stagedDep.Pinned.IgnoreDir()); err != nil {
				return err
			}
		} else {
			// Calculate source and destination hashes; skip copying if equal
			srcHash, err := CreateDirHash(src, stagedDep.Pinned.IgnoreDir())
			if err != nil {
				return err
			}
			dstHash, err := CreateDirHash(dir, stagedDep.Pinned.IgnoreDir())
			if err != nil && !os.IsNotExist(err) {
				return err
			}
			equalDirs := bytes.Equal(srcHash, dstHash)
			if equalDirs {
				LogInfo(`Skipping copying dependency %q (unchanged)`, dir)
			} else {
				LogDebug(`Src hash: %x, Dst hash: %x`, srcHash, dstHash)
				LogInfo(`Copying dependency %q`, dir)
				if err := CopyDirContents(src, dir, stagedDep.Pinned.IgnoreDir()); err != nil {
					return err
				}
			}
		}
	}

	// Save the pinned manifest to file.
	if !cmdLineArgs.reproduce {
		LogInfo("Saving pinned manifest to %q", cmdLineArgs.pinnedManifest)
		var pinned Manifest = make(map[string]Dependency)
		for dir, stagedDep := range stagedDeps {
			pinned[dir] = stagedDep.Pinned
		}
		if raw, err := json.MarshalIndent(pinned, "", "\t"); err != nil {
			return err
		} else if err := ioutil.WriteFile(cmdLineArgs.pinnedManifest, append(raw, '\n'), 0644); err != nil {
			return err
		}
	}

	LogInfo("Finished!")

	return nil
}

func listenForCtrlC() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill)
	go func() {
		sig := <-sigCh
		LogInfo("Received signal %s, panicking now!", sig)
		panic(fmt.Sprintf("%s", sig))
	}()
}

func enrichJSONError(err error, js string) error {
	if s, ok := err.(*json.SyntaxError); ok {
		line, pos, desc := findLineAndPos(s, js)
		return fmt.Errorf("%v\nOccurred on line %v at pos %v: %v", err, line, pos, desc)
	} else {
		return err
	}
}

func findLineAndPos(s *json.SyntaxError, js string) (line int, pos int, desc string) {
	if s.Offset < 0 || s.Offset > int64(len(js)) {
		return -1, -1, "Offset " + strconv.FormatInt(s.Offset, 10) + " is out of bounds."
	}
	if s.Offset == 0 {
		return 0, 0, "Offset 0: Empty file"
	}
	// Syntax error: calculate line number
	start, end := strings.LastIndex(js[:s.Offset], "\n")+1, len(js)
	if idx := strings.Index(js[start:], "\n"); idx >= 0 {
		end = start + idx
	}
	line, pos = strings.Count(js[:start], "\n")+1, int(s.Offset)-start-1
	if pos < 0 {
		pos = 0
	}
	// Return line that has error
	return line, pos, strings.Trim(js[start:end], " \t")
}
