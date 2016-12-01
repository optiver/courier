[![Build Status](https://travis-ci.org/optiver/courier.svg?branch=master)](https://travis-ci.org/optiver/courier)

# Courier

Courier is a simple way of organising and fetching external source code into a
project.

It lets you define a manifest of paths to dependencies. It will fetch
dependencies and organise them in your project. Dependencies can be pinned to a
revision, and subdirectories of the project can be referenced.

This is useful for projects that have source-level dependencies, build 

Dependencies are defined in a `deps.json` file, e.g.

```json
{
    "src/github.com/optiver/killerdinosaurs": {
        "vcs": "git",
        "url": "https://github.com/optiver/killerdinosaurs.git",
        "ref": "master"
    },
    ...
}
```

## How to use

1. Create a `deps.json`.
2. Run `courier`.
3. Check-in `deps.json` and `pins.json` (created by courier)
4. To obtain the exact same dependencies later, run `courier --reproduce`.

## Installing

### From Binaries

TODO - do releases from github.

### From Source

```bash
$ go get github.com/optiver/courier
```

Copy the binary (named `bin/courier` in Linux or `bin/courier.exe` in Windows)
to somewhere in your path.

## Design Goals

* Do one thing and do it well.
* Language and VCS agnostic.
* Obvious behaviour by default.
* Simple manifest.
* Obey the principal of [least astonishment](http://en.wikipedia.org/wiki/Principle_of_least_astonishment).

## License

Courier is:

Copyright \[2016\] \[Optiver\]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
