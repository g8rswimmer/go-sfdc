package goforce

import "testing"

func TestError_UnmarshalJSON(t *testing.T) {
	type fields struct {
		ErrorCode string
		Message   string
		Fields    []string
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Error
		wantErr bool
	}{
		{
			name:   "Success with Status Code",
			fields: fields{},
			args: args{
				data: []byte(`
				{
					"statusCode" : "MALFORMED_ID",
					"message" : "Contact ID: id value of incorrect type: 001xx000003DGb2999",
					"fields" : [
					   "Id"
					]
				 }`),
			},
			want: &Error{
				ErrorCode: "MALFORMED_ID",
				Message:   "Contact ID: id value of incorrect type: 001xx000003DGb2999",
				Fields:    []string{"Id"},
			},
			wantErr: false,
		},
		{
			name:   "Success with Error Code",
			fields: fields{},
			args: args{
				data: []byte(`
				{
					"fields" : [ "Id" ],
					"message" : "Account ID: id value of incorrect type: 001900K0001pPuOAAU",
					"errorCode" : "MALFORMED_ID"
				  }`),
			},
			want: &Error{
				ErrorCode: "MALFORMED_ID",
				Message:   "Account ID: id value of incorrect type: 001900K0001pPuOAAU",
				Fields:    []string{"Id"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				ErrorCode: tt.fields.ErrorCode,
				Message:   tt.fields.Message,
				Fields:    tt.fields.Fields,
			}
			if err := e.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Error.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
