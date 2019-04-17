package collections

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/g8rswimmer/goforce/sobject"

	"github.com/g8rswimmer/goforce"

	"github.com/g8rswimmer/goforce/session"
)

func TestDelete_values(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		records []string
	}
	type args struct {
		allOrNone bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *url.Values
	}{
		{
			name: "values",
			fields: fields{
				records: []string{"id1", "id2", "id3"},
			},
			args: args{
				allOrNone: true,
			},
			want: &url.Values{
				"ids":       []string{"id1,id2,id3"},
				"allOrNone": []string{"true"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Delete{
				session: tt.fields.session,
				records: tt.fields.records,
			}
			if got := d.values(tt.args.allOrNone); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete.values() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelete_Records(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		records []string
	}
	type args struct {
		records []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Delete
	}{
		{
			name:   "records",
			fields: fields{},
			args: args{
				records: []string{"id1", "id2", "id3"},
			},
			want: &Delete{
				records: []string{"id1", "id2", "id3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Delete{
				session: tt.fields.session,
				records: tt.fields.records,
			}
			d.Records(tt.args.records...)
			if !reflect.DeepEqual(d, tt.want) {
				t.Errorf("Delete.Records() = %v, want %v", d, tt.want)
			}
		})
	}
}

func TestDelete_Callout(t *testing.T) {
	type fields struct {
		session session.ServiceFormatter
		records []string
	}
	type args struct {
		allOrNone bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []DeleteValue
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				records: []string{"id1", "id2", "id3"},
				session: &mockSessionFormatter{
					url: "something.com",
					client: mockHTTPClient(func(req *http.Request) *http.Response {

						if strings.HasPrefix(req.URL.String(), "something.com/composite/sobjects") == false {
							return &http.Response{
								StatusCode: 500,
								Status:     "Bad URL: " + req.URL.String(),
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}

						if req.Method != http.MethodDelete {
							return &http.Response{
								StatusCode: 500,
								Status:     "Bad Method",
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}

						values := req.URL.Query()
						if _, ok := values["allOrNone"]; ok == false {

							return &http.Response{
								StatusCode: 500,
								Status:     "allOrNone",
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}
						if _, ok := values["ids"]; ok == false {

							return &http.Response{
								StatusCode: 500,
								Status:     "ids",
								Body:       ioutil.NopCloser(strings.NewReader("resp")),
								Header:     make(http.Header),
							}
						}
						resp := `
						[
							{
								"id" : "001RM000003oLrfYAE",
								"success" : true,
								"errors" : [ ]
							 },
							 {
								"success" : false,
								"errors" : [
								   {
									  "statusCode" : "MALFORMED_ID",
									  "message" : "malformed id 001RM000003oLrB000",
									  "fields" : [ ]
								   }
								]
							 }
						  ]`

						return &http.Response{
							StatusCode: http.StatusOK,
							Status:     "Some Status",
							Body:       ioutil.NopCloser(strings.NewReader(resp)),
							Header:     make(http.Header),
						}
					}),
				},
			},
			want: []DeleteValue{
				{
					sobject.InsertValue{
						Success: true,
						ID:      "001RM000003oLrfYAE",
						Errors:  make([]goforce.Error, 0),
					},
				},
				{
					sobject.InsertValue{
						Success: false,
						Errors: []goforce.Error{
							{
								ErrorCode: "MALFORMED_ID",
								Message:   "malformed id 001RM000003oLrB000",
								Fields:    make([]string, 0),
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Delete{
				session: tt.fields.session,
				records: tt.fields.records,
			}
			got, err := d.Callout(tt.args.allOrNone)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete.Callout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete.Callout() = %v, want %v", got, tt.want)
			}
		})
	}
}
