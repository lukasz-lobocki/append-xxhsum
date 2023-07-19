# append-xxhsum

Recursively adds missing xxhsum (XXH64) hashes from PATH to --xxhsum-filepath.

## Usage

```bash
append-xxhsum [--xxhsum-filepath FILEPATH] \
  [--bsd-style] [--verbose] [--debug] [--help] \
  PATH
```

## Arguments

| arg | description |
| -- | -- |
| PATH | PATH to analyze |

## Parameters

| param | long-param | description |
| -- | -- | -- |
| -x | --xxhsum-filepath | FILEPATH of file to append to. Defaults to PATH\\..\\DIRNAME.xxhsum |
| -b | --bsd-style | BSD-style checksum lines. Defaults to GNU-style |
| -v | --verbose | increase the verbosity |
| -d | --debug | show debug information |
| -h | --help | show this help message and exit |

To verify use `xxhsum --check --quiet FILEPATH`

<details>
<summary>Test run</summary>

```bash
pushd ~/Pictures >/dev/null \
  && time ~/Code/golang/append-xxhsum/bin/append-xxhsum-amd64 ../Code \
  && popd >/dev/null
```

</details>

<details>
<summary>Cross-compilation for ARM</summary>

Use `export GOOS=linux && export GOARCH=arm64` before running `go build`.

Use `lscpu` to find out architecture. Check [this](https://github.com/golang/go/wiki/GoArm) guide for export values.

</details>
