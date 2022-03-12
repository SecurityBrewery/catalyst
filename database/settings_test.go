package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SecurityBrewery/catalyst/generated/model"
	"github.com/SecurityBrewery/catalyst/generated/pointer"
	"github.com/SecurityBrewery/catalyst/test"
)

var bob = &model.UserData{
	Email: pointer.String("bob@example.org"),
	Name:  pointer.String("Bob"),
}

var bobResponse = &model.UserDataResponse{
	ID:    "bob",
	Email: pointer.String("bob@example.org"),
	Name:  pointer.String("Bob"),
}

func TestDatabase_UserDataCreate(t *testing.T) {
	type args struct {
		id      string
		setting *model.UserData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Normal setting", args: args{id: "bob", setting: bob}, wantErr: false},
		{name: "Nil setting", args: args{id: "bob"}, wantErr: true},
		{name: "UserData without settingname", args: args{id: ""}, wantErr: true},
		{name: "Only settingname", args: args{id: "bob"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, _, _, db, cleanup, err := test.DB(t)
			if err != nil {
				t.Fatal(err)
			}
			defer cleanup()

			if err := db.UserDataCreate(test.Context(), tt.args.id, tt.args.setting); (err != nil) != tt.wantErr {
				t.Errorf("settingCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDatabase_UserDataGet(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.UserDataResponse
		wantErr bool
	}{
		{name: "Normal get", args: args{id: "bob"}, want: bobResponse},
		{name: "Not existing", args: args{id: "foo"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, _, _, db, cleanup, err := test.DB(t)
			if err != nil {
				t.Fatal(err)
			}
			defer cleanup()

			if err := db.UserDataCreate(test.Context(), "bob", bob); err != nil {
				t.Errorf("settingCreate() error = %v", err)
			}

			got, err := db.UserDataGet(test.Context(), tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserDataGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDatabase_UserDataList(t *testing.T) {
	tests := []struct {
		name    string
		want    []*model.UserDataResponse
		wantErr bool
	}{
		{name: "Normal list", want: []*model.UserDataResponse{bobResponse}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, _, _, db, cleanup, err := test.DB(t)
			if err != nil {
				t.Fatal(err)
			}
			defer cleanup()

			if err := db.UserDataCreate(test.Context(), "bob", bob); err != nil {
				t.Errorf("settingCreate() error = %v", err)
			}

			got, err := db.UserDataList(test.Context())
			if (err != nil) != tt.wantErr {
				t.Errorf("UserDataList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDatabase_UserDataUpdate(t *testing.T) {
	type args struct {
		id      string
		setting *model.UserData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Normal", args: args{id: "bob", setting: bob}},
		{name: "Not existing", args: args{id: "foo"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, _, _, db, cleanup, err := test.DB(t)
			if err != nil {
				t.Fatal(err)
			}
			defer cleanup()

			if err := db.UserDataCreate(test.Context(), "bob", bob); err != nil {
				t.Errorf("settingCreate() error = %v", err)
			}

			if _, err := db.UserDataUpdate(test.Context(), tt.args.id, tt.args.setting); (err != nil) != tt.wantErr {
				t.Errorf("UserDataUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
