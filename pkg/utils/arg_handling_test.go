package utils

import (
	"testing"
)

func Test_expandTilde(t *testing.T) {
	type args struct {
		inputPath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"NO_CHANGE1", args{"./Bulba"}, "Bulba", false},
		{"NO_CHANGE2", args{"./Bulba/"}, "Bulba", false},
		{"NO_CHANGE3", args{"kolo/Domenats"}, "kolo/Domenats", false},
		{"NO_CHANGE4", args{"/kolo/Dmenats"}, "/kolo/Dmenats", false},
		{"NO_CHANGE5", args{"kolo/Domenats/"}, "kolo/Domenats", false},
		{"NO_CHANGE6", args{"/kolo/Domenats/"}, "/kolo/Domenats", false},
		{"EXPAND_JUST TILDE", args{"~"}, "/home/lukasz", false},
		{"EXPAND_TILDE_PREFIX1", args{"~/Documents"}, "/home/lukasz/Documents", false},
		{"EXPAND_TILDE_PREFIX2", args{"~/Documents/"}, "/home/lukasz/Documents", false},
		{"EXPAND_TILDE_PREFIX3", args{"~/Documents/Bulba"}, "/home/lukasz/Documents/Bulba", false},
		{"EXPAND_TILDE_PREFIX4", args{"~/kolo/~/Documenats"}, "/home/lukasz/kolo/~/Documenats", false},
		{"GIBBERISH1", args{"~kolo/Documenats"}, "", true},
		{"GIBBERISH2", args{"~kolo/Documenats/"}, "", true},
		{"GIBBERISH3", args{"~.kolo/Documenats/"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expandTilde(tt.args.inputPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("expandTilde() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("expandTilde() = %v, want %v", got, tt.want)
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
				t.Errorf("ArgParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ArgParse() = %v, want %v", got, tt.want)
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
				t.Errorf("ParamParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParamParse() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParamParse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
