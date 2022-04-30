package base

import (
	"reflect"
	"testing"
)

// @Author: Derek
// @Description:
// @Date: 2022/4/30 11:09
// @Version 1.0

func TestNewError(t *testing.T) {
	type args struct {
		code    int
		message string
	}
	tests := []struct {
		name string
		args args
		want Error
	}{
		{"test1", args{code: -1, message: "test1"}, NewError(-1, "test1")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewError(tt.args.code, tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseError_Equal(t *testing.T) {
	type fields struct {
		ErrNo  int
		ErrMsg string
	}
	type args struct {
		e Error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"test1", fields{ErrNo: -1, ErrMsg: "test1"}, args{NewError(-1, "test1")}, true},
		{"test2", fields{ErrNo: -1, ErrMsg: "test1"}, args{NewError(-1, "test2")}, true},
		{"test3", fields{ErrNo: -1, ErrMsg: "test1"}, args{NewError(-2, "test2")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &baseError{
				ErrNo:  tt.fields.ErrNo,
				ErrMsg: tt.fields.ErrMsg,
			}
			if got := err.Equal(tt.args.e); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseError_Error(t *testing.T) {
	type fields struct {
		ErrNo  int
		ErrMsg string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"test1", fields{ErrNo: -1, ErrMsg: "test1"}, "test1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := baseError{
				ErrNo:  tt.fields.ErrNo,
				ErrMsg: tt.fields.ErrMsg,
			}
			if got := err.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseError_SetErrMsg(t *testing.T) {
	type fields struct {
		ErrNo  int
		ErrMsg string
	}
	type args struct {
		format string
		v      []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"test1", fields{ErrNo: -1, ErrMsg: "test1"}, args{"test1_%d", []interface{}{1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &baseError{
				ErrNo:  tt.fields.ErrNo,
				ErrMsg: tt.fields.ErrMsg,
			}

			err.SetErrMsg(tt.args.format, tt.args.v...)

			if err.ErrMsg != "test1_1" {
				t.Error("Fail.", err.ErrMsg)
			}
		})
	}
}

func Test_baseError_Sprintf(t *testing.T) {
	type fields struct {
		ErrNo  int
		ErrMsg string
	}
	type args struct {
		v []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Error
	}{
		{"test1", fields{ErrNo: -1, ErrMsg: "test1_%d"}, args{[]interface{}{1}}, NewError(-1, "test1_1")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &baseError{
				ErrNo:  tt.fields.ErrNo,
				ErrMsg: tt.fields.ErrMsg,
			}
			if got := err.Sprintf(tt.args.v...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sprintf() = %v, want %v", got, tt.want)
			}
		})
	}
}
