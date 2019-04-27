package goforce

import (
	"reflect"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	type args struct {
		salesforceTime string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name:    "No Time",
			args:    args{},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name: "Invalid Format",
			args: args{
				salesforceTime: "Invalid",
			},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name: "RFC 3339",
			args: args{
				salesforceTime: "2019-04-08T00:05:30Z",
			},
			want:    time.Date(2019, 4, 8, 0, 5, 30, 0, time.UTC),
			wantErr: false,
		},
		{
			name: "Salesforce DateTime",
			args: args{
				salesforceTime: "2013-05-08T21:20:00.000+0000",
			},
			want:    time.Date(2013, 5, 8, 21, 20, 00, 00, time.UTC),
			wantErr: false,
		},
		{
			name: "Salesforce Date",
			args: args{
				salesforceTime: "2018-07-26",
			},
			want:    time.Date(2018, 7, 26, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTime(tt.args.salesforceTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
