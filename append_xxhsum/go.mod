module xxhsum/append_xxhsum

go 1.20

replace xxhsum/arg_handling => ../arg_handling

replace xxhsum/dictionar => ../dictionar

replace xxhsum/globals => ../globals

require (
	github.com/briandowns/spinner v1.23.0
	github.com/cespare/xxhash/v2 v2.2.0
	xxhsum/arg_handling v0.0.0-00010101000000-000000000000
	xxhsum/dictionar v0.0.0-00010101000000-000000000000
	xxhsum/globals v0.0.0-00010101000000-000000000000
)

require (
	github.com/fatih/color v1.7.0 // indirect
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mattn/go-isatty v0.0.8 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/term v0.1.0 // indirect
)
