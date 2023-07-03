package arg_handling

import "testing"

func Test_expand_tilde(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"no change", args{"./Bulba"}, "./Bulba"},
		{"just tilde", args{"~"}, "/home/lukasz"},
		{"tilde", args{"~/Documents"}, "/home/lukasz/Documents"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := expand_tilde(tt.args.path); got != tt.want {
				t.Errorf("expand_tilde() = %v, want %v", got, tt.want)
			}
		})
	}
}
