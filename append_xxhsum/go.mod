module xxhsum/append_xxhsum

go 1.20

replace xxhsum/arg_handling => ../arg_handling
replace xxhsum/dictionar => ../dictionar

require (
	github.com/cespare/xxhash/v2 v2.2.0
	xxhsum/arg_handling v0.0.0-00010101000000-000000000000
	xxhsum/dictionar v0.0.0-00010101000000-000000000000
)
