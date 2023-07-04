package dictionar

import (
	"reflect"
	"testing"
)

func TestLoad_xxhsum_file(t *testing.T) {
	type args struct {
		in_file string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{"NO FILE", args{"/Bulba"}, nil, true},
		{"WRONG FILE", args{"/home/lukasz/.profile"}, make(map[string]string), false},
		{"FILE", args{"/home/lukasz/Code/golang/xxhsum/test.xxhsum"}, map[string]string{
			"./golang/goroot/go.mod": "1f809539dbc4e242", "./golang/goroot/sample-app/sample-app": "76b9b81ec1c51248"},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load_xxhsum_file(tt.args.in_file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load_xxhsum_file() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load_xxhsum_file() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDump_xxhsum_dict(t *testing.T) {
	type args struct {
		in_data map[string]string
	}
	tests := []struct {
		name string
		args args
	}{
		{"SINGLE", args{map[string]string{
			"./golang/goroot/go.mod": "1f809539dbc4e242", "./golang/goroot/sample-app/sample-app": "76b9b81ec1c51248"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Dump_xxhsum_dict(tt.args.in_data)
		})
	}
}
