package bulk

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	"github.com/g8rswimmer/go-sfdc/session"
)

func TestJob_formatOptions(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		info    Response
	}
	type args struct {
		options *Options
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Options
		wantErr bool
	}{
		{
			name:   "passing",
			fields: fields{},
			args: args{
				options: &Options{
					ColumnDelimiter:     Caret,
					ContentType:         CSV,
					ExternalIDFieldName: "Some External Field",
					LineEnding:          Linefeed,
					Object:              "Account",
					Operation:           Insert,
				},
			},
			want: &Options{
				ColumnDelimiter:     Caret,
				ContentType:         CSV,
				ExternalIDFieldName: "Some External Field",
				LineEnding:          Linefeed,
				Object:              "Account",
				Operation:           Insert,
			},
			wantErr: false,
		},
		{
			name:   "defaults",
			fields: fields{},
			args: args{
				options: &Options{
					ExternalIDFieldName: "Some External Field",
					Object:              "Account",
					Operation:           Insert,
				},
			},
			want: &Options{
				ColumnDelimiter:     Comma,
				ContentType:         CSV,
				ExternalIDFieldName: "Some External Field",
				LineEnding:          Linefeed,
				Object:              "Account",
				Operation:           Insert,
			},
			wantErr: false,
		},
		{
			name:   "no object",
			fields: fields{},
			args: args{
				options: &Options{
					ExternalIDFieldName: "Some External Field",
					Operation:           Insert,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "no operation",
			fields: fields{},
			args: args{
				options: &Options{
					ExternalIDFieldName: "Some External Field",
					Object:              "Account",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "no external fields",
			fields: fields{},
			args: args{
				options: &Options{
					Object:    "Account",
					Operation: Upsert,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				session: tt.fields.session,
				info:    tt.fields.info,
			}
			err := j.formatOptions(tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Job.formatOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(tt.args.options, tt.want) {
				t.Errorf("Job.formatOptions() = %v, want %v", tt.args.options, tt.want)
			}
		})
	}
}

func TestJob_newline(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		info    Response
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Carrage return",
			fields: fields{
				info: Response{
					LineEnding: "CRLF",
				},
			},
			want: "\r\n",
		},
		{
			name: "Line feed",
			fields: fields{
				info: Response{
					LineEnding: "LF",
				},
			},
			want: "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				session: tt.fields.session,
				info:    tt.fields.info,
			}
			if got := j.newline(); got != tt.want {
				t.Errorf("Job.newline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_delimiter(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		info    Response
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "tab",
			fields: fields{
				info: Response{
					ColumnDelimiter: "TAB",
				},
			},
			want: "\t",
		},
		{
			name: "back quote",
			fields: fields{
				info: Response{
					ColumnDelimiter: "BACKQUOTE",
				},
			},
			want: "`",
		},
		{
			name: "caret",
			fields: fields{
				info: Response{
					ColumnDelimiter: "CARET",
				},
			},
			want: "^",
		},
		{
			name: "comma",
			fields: fields{
				info: Response{
					ColumnDelimiter: "COMMA",
				},
			},
			want: ",",
		},
		{
			name: "pipe",
			fields: fields{
				info: Response{
					ColumnDelimiter: "PIPE",
				},
			},
			want: "|",
		},
		{
			name: "semi colon",
			fields: fields{
				info: Response{
					ColumnDelimiter: "SEMICOLON",
				},
			},
			want: ";",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				session: tt.fields.session,
				info:    tt.fields.info,
			}
			if got := j.delimiter(); got != tt.want {
				t.Errorf("Job.delimiter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_record(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		info    Response
	}
	type args struct {
		fields []string
		values []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		{
			name:   "make record",
			fields: fields{},
			args: args{
				fields: []string{
					"first",
					"last",
					"DOB",
				},
				values: []string{
					"john",
					"doe",
					"1/1/1970",
				},
			},
			want: map[string]string{
				"first": "john",
				"last":  "doe",
				"DOB":   "1/1/1970",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				session: tt.fields.session,
				info:    tt.fields.info,
			}
			if got := j.record(tt.args.fields, tt.args.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.record() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJob_fields(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		info    Response
	}
	type args struct {
		scanner   *bufio.Scanner
		delimiter string
		offset    int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:   "passing",
			fields: fields{},
			args: args{
				scanner:   bufio.NewScanner(strings.NewReader("sf_id|first|last|DOB")),
				delimiter: "|",
				offset:    1,
			},
			want: []string{
				"first",
				"last",
				"DOB",
			},
			wantErr: false,
		},
		{
			name:   "error",
			fields: fields{},
			args: args{
				scanner:   bufio.NewScanner(strings.NewReader("")),
				delimiter: "|",
				offset:    1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				session: tt.fields.session,
				info:    tt.fields.info,
			}
			got, err := j.fields(tt.args.scanner, tt.args.delimiter, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("Job.fields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.fields() = %v, want %v", got, tt.want)
			}
		})
	}
}
