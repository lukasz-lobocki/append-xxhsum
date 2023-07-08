# append-xxhsum

Recursively adds missing xxhsum (XXH64) hashes from PATH to --xxhsum-filepath.

## Usage

```bash
append-xxhsum [--xxhsum-filepath FILEPATH] \
  [--bsd-style] [--verbose] [--help] \
  PATH
```

## Arguments

| arg | description |
| -- | -- |
| PATH | PATH to analyze. |

## Parameters

| param | long-param | description |
| -- | -- | -- |
| -x | --xxhsum-filepath | FILEPATH to file to append to. |
| -b | --bsd-style | BSD-style checksum lines. Defaults to GNU-style. |
| -v | --verbose | increase the verbosity of the bash script. |
| -h | --help | show this help message and exit. |
