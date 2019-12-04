package bulk

import (
	"encoding/csv"
	"reflect"
	"strings"
	"testing"
)

func TestNewFormatter(t *testing.T) {
	type args struct {
		job    *Job
		fields []string
	}
	sb := &strings.Builder{}
	csvWriter := csv.NewWriter(sb)

	tests := []struct {
		name    string
		args    args
		want    *Formatter
		wantErr bool
	}{
		{
			name: "Passing",
			args: args{
				job: &Job{
					info: Response{
						ColumnDelimiter: string(Pipe),
						LineEnding:      string(Linefeed),
					},
				},
				fields: []string{
					"Name",
					`"Site"`,
				},
			},
			want: &Formatter{
				job: &Job{
					info: Response{
						ColumnDelimiter: string(Pipe),
						LineEnding:      string(Linefeed),
					},
				},
				fields: []string{
					"Name",
					`"Site"`,
				},
				sb:        sb,
				csvWriter: csvWriter,
			},
			wantErr: false,
		},
		{
			name: "No Job",
			args: args{
				job: nil,
				fields: []string{
					"Name",
					"Site",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "No Fields",
			args: args{
				job: &Job{
					info: Response{
						ColumnDelimiter: string(Pipe),
						LineEnding:      string(Linefeed),
					},
				},
				fields: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFormatter(tt.args.job, tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFormatter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				tt.want.csvWriter.Comma = rune(tt.args.job.delimiter()[0])
				if err := tt.want.csvWriter.Write(tt.want.fields); err != nil {
					t.Errorf("error writting: %v", err)
				}
				tt.want.csvWriter.Flush()
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Logf("got: %q want: %q", got.sb.String(), tt.want.sb.String())
				t.Errorf("NewFormatter() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

type testRecord struct {
	fields     map[string]interface{}
	insertNull bool
}

func (t *testRecord) Fields() map[string]interface{} {
	return t.fields
}
func (t *testRecord) InsertNull() bool {
	return t.insertNull
}
func TestFormatter_Add(t *testing.T) {
	type fields struct {
		job        *Job
		fields     []string
		insertNull bool
		sb         *strings.Builder
	}
	type args struct {
		records []Record
	}

	sb := &strings.Builder{}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Adding",
			fields: fields{
				job: &Job{
					info: Response{
						ColumnDelimiter: string(Pipe),
						LineEnding:      string(Linefeed),
					},
				},
				fields: []string{
					"Name",
					"Site",
				},
				insertNull: true,
				sb:         sb,
			},
			args: args{
				records: []Record{
					&testRecord{
						fields: map[string]interface{}{
							"Name": "name 1",
							"Site": "good site",
						},
					},
					&testRecord{
						fields: map[string]interface{}{
							"Name": "name 2",
							"Site": "great site",
						},
					},
				},
			},
			want:    "name 1|good site\nname 2|great site\n",
			wantErr: false,
		},
		{
			name: "Adding",
			fields: fields{
				job: &Job{
					info: Response{
						ColumnDelimiter: string(Pipe),
						LineEnding:      string(Linefeed),
					},
				},
				fields: []string{
					"Name",
					"Site",
				},
				insertNull: true,
				sb:         sb,
			},
			args: args{
				records: nil,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Formatter{
				job:       tt.fields.job,
				fields:    tt.fields.fields,
				sb:        tt.fields.sb,
				csvWriter: csv.NewWriter(sb),
			}
			f.csvWriter.Comma = rune(f.job.delimiter()[0])

			var err error
			if err = f.Add(tt.args.records...); (err != nil) != tt.wantErr {
				t.Errorf("Formatter.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				if tt.want != f.sb.String() {
					t.Errorf("Formatter.Add() want = %v, got %v", tt.want, f.sb.String())
				}
			}
		})
	}
}
