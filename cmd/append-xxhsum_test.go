package main

import (
	"testing"
)

func Test_appendToFile(t *testing.T) {
	type args struct {
		filename string
		content  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"FILE", args{"/home/lukasz/Code/golang/append-xxhsum/tst/test3.xx_append", "Lorem ipsum\n"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := appendToFile(tt.args.filename, tt.args.content); (err != nil) != tt.wantErr {
				t.Errorf("appendToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_calculateXXHash(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"EXISTS", args{"/home/lukasz/Code/golang/append-xxhsum/tst/test1.xxhsum"}, "1b4378db293122d8", false},
		{"DOES_NOT_EXIST", args{"/home/lukasz/Code/golang/append-xxhsum/tst/tes.xxhm"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateXXHash(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateXXHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calculateXXHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_debugVariables(t *testing.T) {
	type args struct {
		verbose          bool
		givenPath        string
		xxhsumFilepath   string
		xxhsumFileExists bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"DUMMY", args{true, "/home/lukasz", "/home/lukasz/test1.xxhsum", true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debugVariables(tt.args.verbose, tt.args.givenPath, tt.args.xxhsumFilepath, tt.args.xxhsumFileExists)
		})
	}
}

func Test_calculateLine(t *testing.T) {
	type args struct {
		bsdStyle bool
		relPath  string
		checksum string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"BSD", args{true, "/home/lukasz", "123567890123456"}, "XXH64 (/home/lukasz) = 123567890123456\n"},
		{"GNU", args{false, "/home/lukasz", "123567890123456"}, "123567890123456 */home/lukasz\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateLine(tt.args.bsdStyle, tt.args.relPath, tt.args.checksum); got != tt.want {
				t.Errorf("calculateLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
