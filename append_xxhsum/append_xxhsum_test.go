package main

import "testing"

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
		{"FILE", args{"/home/lukasz/Code/golang/xxhsum/test.xxhsum"}, "7832b0cd476c893e", false},
		{"FILE", args{"/home/lukasz/Code/golang/xxhsum/tes.xxhm"}, "", true},
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
