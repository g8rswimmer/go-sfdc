package sobject

import (
	"reflect"
	"testing"
	"time"

	"github.com/TuSKan/go-sfdc"
	"github.com/TuSKan/go-sfdc/session"
)

func TestSalesforceAPI_Metadata(t *testing.T) {
	type fields struct {
		metadata *metadata
		describe *describe
		dml      *dml
		query    *query
	}
	type args struct {
		sobject string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    MetadataValue
		wantErr bool
	}{
		{
			name:    "No Metadata field",
			want:    MetadataValue{},
			wantErr: true,
		},
		{
			name: "Invalid Args",
			fields: fields{
				metadata: &metadata{
					session: &mockSessionFormatter{
						url: "http://wwww.google.com",
					},
				},
			},
			want:    MetadataValue{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Resources{
				metadata: tt.fields.metadata,
				describe: tt.fields.describe,
				dml:      tt.fields.dml,
				query:    tt.fields.query,
			}
			got, err := a.Metadata(tt.args.sobject)
			if (err != nil) != tt.wantErr {
				t.Errorf("SalesforceAPI.Metadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SalesforceAPI.Metadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSalesforceAPI_Describe(t *testing.T) {
	type fields struct {
		metadata *metadata
		describe *describe
		dml      *dml
		query    *query
	}
	type args struct {
		sobject string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    DescribeValue
		wantErr bool
	}{
		{
			name:    "No Describe field",
			want:    DescribeValue{},
			wantErr: true,
		},
		{
			name: "Invalid Args",
			fields: fields{
				describe: &describe{
					session: &mockSessionFormatter{
						url: "http://wwww.google.com",
					},
				},
			},
			want:    DescribeValue{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Resources{
				metadata: tt.fields.metadata,
				describe: tt.fields.describe,
				dml:      tt.fields.dml,
				query:    tt.fields.query,
			}
			got, err := a.Describe(tt.args.sobject)
			if (err != nil) != tt.wantErr {
				t.Errorf("SalesforceAPI.Describe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SalesforceAPI.Describe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSalesforceAPI_Insert(t *testing.T) {
	type fields struct {
		metadata *metadata
		describe *describe
		dml      *dml
		query    *query
	}
	type args struct {
		inserter Inserter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    InsertValue
		wantErr bool
	}{
		{
			name:    "No DML field",
			want:    InsertValue{},
			wantErr: true,
		},
		{
			name: "Invalid Args",
			fields: fields{
				dml: &dml{
					session: &mockSessionFormatter{
						url: "http://wwww.google.com",
					},
				},
			},
			want:    InsertValue{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Resources{
				metadata: tt.fields.metadata,
				describe: tt.fields.describe,
				dml:      tt.fields.dml,
				query:    tt.fields.query,
			}
			got, err := a.Insert(tt.args.inserter)
			if (err != nil) != tt.wantErr {
				t.Errorf("SalesforceAPI.Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SalesforceAPI.Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSalesforceAPI_Update(t *testing.T) {
	type fields struct {
		metadata *metadata
		describe *describe
		dml      *dml
		query    *query
	}
	type args struct {
		updater Updater
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "No DML field",
			wantErr: true,
		},
		{
			name: "Invalid Args",
			fields: fields{
				dml: &dml{
					session: &mockSessionFormatter{
						url: "http://wwww.google.com",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Resources{
				metadata: tt.fields.metadata,
				describe: tt.fields.describe,
				dml:      tt.fields.dml,
				query:    tt.fields.query,
			}
			if err := a.Update(tt.args.updater); (err != nil) != tt.wantErr {
				t.Errorf("SalesforceAPI.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSalesforceAPI_Upsert(t *testing.T) {
	type fields struct {
		metadata *metadata
		describe *describe
		dml      *dml
		query    *query
	}
	type args struct {
		upserter Upserter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    UpsertValue
		wantErr bool
	}{
		{
			name:    "No DML field",
			want:    UpsertValue{},
			wantErr: true,
		},
		{
			name: "Invalid Args",
			fields: fields{
				dml: &dml{
					session: &mockSessionFormatter{
						url: "http://wwww.google.com",
					},
				},
			},
			want:    UpsertValue{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Resources{
				metadata: tt.fields.metadata,
				describe: tt.fields.describe,
				dml:      tt.fields.dml,
				query:    tt.fields.query,
			}
			got, err := a.Upsert(tt.args.upserter)
			if (err != nil) != tt.wantErr {
				t.Errorf("SalesforceAPI.Upsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SalesforceAPI.Upsert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSalesforceAPI_Delete(t *testing.T) {
	type fields struct {
		metadata *metadata
		describe *describe
		dml      *dml
		query    *query
	}
	type args struct {
		deleter Deleter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "No DML field",
			wantErr: true,
		},
		{
			name: "Invalid Args",
			fields: fields{
				dml: &dml{
					session: &mockSessionFormatter{
						url: "http://wwww.google.com",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Resources{
				metadata: tt.fields.metadata,
				describe: tt.fields.describe,
				dml:      tt.fields.dml,
				query:    tt.fields.query,
			}
			if err := a.Delete(tt.args.deleter); (err != nil) != tt.wantErr {
				t.Errorf("SalesforceAPI.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSalesforceAPI_Query(t *testing.T) {
	type fields struct {
		metadata *metadata
		describe *describe
		dml      *dml
		query    *query
	}
	type args struct {
		querier Querier
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *sfdc.Record
		wantErr bool
	}{
		{
			name:    "No Query field",
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid Args",
			fields: fields{
				query: &query{
					session: &mockSessionFormatter{
						url: "http://wwww.google.com",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Resources{
				metadata: tt.fields.metadata,
				describe: tt.fields.describe,
				dml:      tt.fields.dml,
				query:    tt.fields.query,
			}
			got, err := a.Query(tt.args.querier)
			if (err != nil) != tt.wantErr {
				t.Errorf("SalesforceAPI.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SalesforceAPI.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSalesforceAPI_ExternalQuery(t *testing.T) {
	type fields struct {
		metadata *metadata
		describe *describe
		dml      *dml
		query    *query
	}
	type args struct {
		querier ExternalQuerier
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *sfdc.Record
		wantErr bool
	}{
		{
			name:    "No Query field",
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid Args",
			fields: fields{
				query: &query{
					session: &mockSessionFormatter{
						url: "http://wwww.google.com",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Resources{
				metadata: tt.fields.metadata,
				describe: tt.fields.describe,
				dml:      tt.fields.dml,
				query:    tt.fields.query,
			}
			got, err := a.ExternalQuery(tt.args.querier)
			if (err != nil) != tt.wantErr {
				t.Errorf("SalesforceAPI.ExternalQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SalesforceAPI.ExternalQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSalesforceAPI_DeletedRecords(t *testing.T) {
	type fields struct {
		metadata *metadata
		describe *describe
		dml      *dml
		query    *query
	}
	type args struct {
		sobject   string
		startDate time.Time
		endDate   time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    DeletedRecords
		wantErr bool
	}{
		{
			name:    "No Query field",
			want:    DeletedRecords{},
			wantErr: true,
		},
		{
			name: "Invalid Args",
			fields: fields{
				query: &query{
					session: &mockSessionFormatter{
						url: "http://wwww.google.com",
					},
				},
			},
			want:    DeletedRecords{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Resources{
				metadata: tt.fields.metadata,
				describe: tt.fields.describe,
				dml:      tt.fields.dml,
				query:    tt.fields.query,
			}
			got, err := a.DeletedRecords(tt.args.sobject, tt.args.startDate, tt.args.endDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("SalesforceAPI.DeletedRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SalesforceAPI.DeletedRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSalesforceAPI_UpdatedRecords(t *testing.T) {
	type fields struct {
		metadata *metadata
		describe *describe
		dml      *dml
		query    *query
	}
	type args struct {
		sobject   string
		startDate time.Time
		endDate   time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    UpdatedRecords
		wantErr bool
	}{
		{
			name:    "No Query field",
			want:    UpdatedRecords{},
			wantErr: true,
		},
		{
			name: "Invalid Args",
			fields: fields{
				query: &query{
					session: &mockSessionFormatter{
						url: "http://wwww.google.com",
					},
				},
			},
			want:    UpdatedRecords{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Resources{
				metadata: tt.fields.metadata,
				describe: tt.fields.describe,
				dml:      tt.fields.dml,
				query:    tt.fields.query,
			}
			got, err := a.UpdatedRecords(tt.args.sobject, tt.args.startDate, tt.args.endDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("SalesforceAPI.UpdatedRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SalesforceAPI.UpdatedRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSalesforceAPI_GetContent(t *testing.T) {
	type fields struct {
		metadata *metadata
		describe *describe
		dml      *dml
		query    *query
	}
	type args struct {
		id      string
		content ContentType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "No Query field",
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid Args",
			fields: fields{
				query: &query{
					session: &mockSessionFormatter{
						url: "http://wwww.google.com",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid Content",
			fields: fields{
				query: &query{
					session: &mockSessionFormatter{
						url: "http://wwww.google.com",
					},
				},
			},
			args: args{
				id:      "12345",
				content: ContentType("Invalid"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Resources{
				metadata: tt.fields.metadata,
				describe: tt.fields.describe,
				dml:      tt.fields.dml,
				query:    tt.fields.query,
			}
			got, err := a.GetContent(tt.args.id, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("SalesforceAPI.GetContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SalesforceAPI.GetContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewResources(t *testing.T) {
	type args struct {
		session session.ServiceFormatter
	}
	tests := []struct {
		name    string
		args    args
		want    *Resources
		wantErr bool
	}{
		{
			name: "passing",
			args: args{
				session: &mockSessionFormatter{
					url: "https://test.salesforce.com",
				},
			},
			want: &Resources{
				metadata: &metadata{
					session: &mockSessionFormatter{
						url: "https://test.salesforce.com",
					},
				},
				describe: &describe{
					session: &mockSessionFormatter{
						url: "https://test.salesforce.com",
					},
				},
				dml: &dml{
					session: &mockSessionFormatter{
						url: "https://test.salesforce.com",
					},
				},
				query: &query{
					session: &mockSessionFormatter{
						url: "https://test.salesforce.com",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "error",
			args:    args{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewResources(tt.args.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewResources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResources() = %v, want %v", got, tt.want)
			}
		})
	}
}
