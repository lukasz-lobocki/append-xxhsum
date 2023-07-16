package utils

import (
	"reflect"
	"testing"
)

func TestDumpXXHSumDict(t *testing.T) {
	type args struct {
		inputData map[string]string
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
			DumpXXHSumDict(tt.args.inputData)
		})
	}
}

func TestLoadXXHSumFile(t *testing.T) {
	type args struct {
		inputFile string
		bsdStyle  bool
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{"NO_FILE1", args{"/Bulba", false}, nil, true},
		{"NO_FILE2", args{"/Bulba", true}, nil, true},
		{"WRONG_FILE1", args{"/home/lukasz/.profile", true}, make(map[string]string), false},
		{"WRONG_FILE2", args{"/home/lukasz/.gitcommitmessage.txt", false}, make(map[string]string), false},
		{"FILEA", args{"/home/lukasz/Code/golang/append-xxhsum/tst/test1.xxhsum", false}, map[string]string{
			"./golang/goroot/go.mod": "1f809539dbc4e242", "./golang/goroot/sample-app/sample-app": "76b9b81ec1c51248"},
			false},
		{"FILEB", args{"/home/lukasz/Code/golang/append-xxhsum/tst/test2.xxhsum", true}, map[string]string{
			"./golang/goroot/go.mod": "1f809539dbc4e242", "./golang/goroot/sample-app/sample-app": "76b9b81ec1c51248"},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadXXHSumFile(tt.args.inputFile, tt.args.bsdStyle)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadXXHSumFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadXXHSumFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
