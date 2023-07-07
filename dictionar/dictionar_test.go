package dictionar

import (
	"reflect"
	"testing"
)

func TestDumpXXHSumDict(t *testing.T) {
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
			DumpXXHSumDict(&tt.args.in_data)
		})
	}
}

func TestLoadXXHSumFile(t *testing.T) {
	type args struct {
		in_file   string
		bsd_style bool
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{"NO FILE", args{"/Bulba", false}, nil, true},
		{"NO FILE", args{"/Bulba", true}, nil, true},
		{"WRONG FILE", args{"/home/lukasz/.profile", false}, make(map[string]string), false},
		{"WRONG FILE", args{"/home/lukasz/.profile", true}, make(map[string]string), false},
		{"FILEA", args{"/home/lukasz/Code/golang/xxhsum/test.xxhsum", false}, map[string]string{
			"./golang/goroot/go.mod": "1f809539dbc4e242", "./golang/goroot/sample-app/sample-app": "76b9b81ec1c51248"},
			false},
		{"FILEB", args{"/home/lukasz/Code/golang/xxhsum/test-BSD.xxhsum", true}, map[string]string{
			"./golang/goroot/go.mod": "1f809539dbc4e242", "./golang/goroot/sample-app/sample-app": "76b9b81ec1c51248"},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadXXHSumFile(tt.args.in_file, tt.args.bsd_style)
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
