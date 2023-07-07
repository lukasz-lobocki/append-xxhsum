package arg_handling

import (
	"testing"
)

func Test_expandTilde(t *testing.T) {
	type args struct {
		in_path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"NO CHANGE", args{"./Bulba"}, "Bulba", false},
		{"NO CHANGE", args{"./Bulba/"}, "Bulba", false},
		{"NO CHANGE", args{"kolo/Domenats"}, "kolo/Domenats", false},
		{"NO CHANGE", args{"/kolo/Dmenats"}, "/kolo/Dmenats", false},
		{"NO CHANGE", args{"kolo/Domenats/"}, "kolo/Domenats", false},
		{"NO CHANGE", args{"/kolo/Domenats/"}, "/kolo/Domenats", false},
		{"EXPAND JUST TILDE", args{"~"}, "/home/lukasz", false},
		{"EXPAND TILDE PREFIX", args{"~/Documents"}, "/home/lukasz/Documents", false},
		{"EXPAND TILDE PREFIX", args{"~/Documents/"}, "/home/lukasz/Documents", false},
		{"EXPAND TILDE PREFIX", args{"~/Documents/Bulba"}, "/home/lukasz/Documents/Bulba", false},
		{"EXPAND TILDE PREFIX", args{"~/kolo/~/Documenats"}, "/home/lukasz/kolo/~/Documenats", false},
		{"GIBBERISH1", args{"~kolo/Documenats"}, "", true},
		{"GIBBERISH1", args{"~kolo/Documenats/"}, "", true},
		{"GIBBERISH1", args{"~.kolo/Documenats/"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expandTilde(tt.args.in_path)
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

func TestArgParse(t *testing.T) {
	type args struct {
		arg     string
		verbose bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"DIR", args{"/home/lukasz", true}, "/home/lukasz", false},
		{"DIR", args{"/home/lukasz/", true}, "/home/lukasz", false},
		{"ROOT PATH", args{"/", true}, "/", false},
		{"NON-EXISTING", args{"./Bulba", true}, "", true},
		{"NODIR", args{"/home/lukasz/.profile", true}, "", true},
		{"WRONG PATH", args{"/home/luksza/.profile", true}, "", true},
		{"ROOT HOMEDIR", args{"/root/snap", true}, "", true},
		{"GIBBERISH", args{"|.~", true}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ArgParse(tt.args.arg, tt.args.verbose)
			if (err != nil) != tt.wantErr {
				t.Errorf("Arg_parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Arg_parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParamParse(t *testing.T) {
	type args struct {
		param   string
		verbose bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   bool
		wantErr bool
	}{
		{"FILE", args{"/home/lukasz/.profile", true}, "/home/lukasz/.profile", true, false},
		{"NON EXISTING DIR", args{"/home/luksiz/.profile", true}, "/home/luksiz/.profile", false, false},
		{"FILE with TILDE", args{"~/.profile", true}, "/home/lukasz/.profile", true, false},
		{"NON EXISTING FILE with TILDE", args{"~/.profilwee", true}, "/home/lukasz/.profilwee", false, false},
		{"DIR", args{"/home/lukasz/Documents", true}, "", true, true},
		{"DIR", args{"/home/lukasz/Documents/", true}, "", true, true},
		{"ROOT HOMEDIR FILE", args{"/root/.profile", true}, "", false, true},
		{"FILE WITH WRONG TILDE", args{"~.profile", true}, "", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParamParse(tt.args.param, tt.args.verbose)
			if (err != nil) != tt.wantErr {
				t.Errorf("Param_parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Param_parse() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Param_parse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
