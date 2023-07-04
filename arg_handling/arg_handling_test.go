package arg_handling

import "testing"

func Test_expand_tilde(t *testing.T) {
	type args struct {
		in_path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"NO CHANGE", args{"./Bulba"}, "./Bulba", false},
		{"NO CHANGE", args{"/kolo/Dmenats"}, "/kolo/Dmenats", false},
		{"NO CHANGE", args{"kolo/Domenats"}, "kolo/Domenats", false},
		{"EXPAND JUST TILDE", args{"~"}, "/home/lukasz", false},
		{"EXPAND TILDE PREFIX", args{"~/Documents"}, "/home/lukasz/Documents", false},
		{"EXPAND TILDE PREFIX", args{"~/Documents/Bulba"}, "/home/lukasz/Documents/Bulba", false},
		{"EXPAND TILDE PREFIX", args{"~/kolo/~/Documenats"}, "/home/lukasz/kolo/~/Documenats", false},
		{"GIBBERISH1", args{"~kolo/Documenats"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expand_tilde(tt.args.in_path)
			if (err != nil) != tt.wantErr {
				t.Errorf("expand_tilde() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("expand_tilde() = %v, want %v", got, tt.want)
			}
		})
	}
}
