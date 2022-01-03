package database

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"

	"github.com/arangodb/go-driver"
	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"

	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/model"
	"github.com/SecurityBrewery/catalyst/pointer"
	"github.com/SecurityBrewery/catalyst/role"
	"github.com/SecurityBrewery/catalyst/time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateKey() string {
	b := make([]rune, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func toUser(user *model.UserForm, sha256 *string) *model.User {
	roles := []string{}
	roles = append(roles, role.Strings(role.Explodes(user.Roles))...)
	u := &model.User{
		Blocked: user.Blocked,
		Roles:   roles,
		Sha256:  sha256,
		Apikey:  user.Apikey,
	}

	// log.Println(u)
	// b, _ := json.Marshal(u)
	// loader := gojsonschema.NewBytesLoader(b)
	// res, err := model.UserSchema.Validate(loader)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(res.Errors())

	return u
}

func toUserResponse(key string, user *model.User) *model.UserResponse {
	return &model.UserResponse{
		ID:      key,
		Roles:   user.Roles,
		Blocked: user.Blocked,
		Apikey:  user.Apikey,
	}
}

func toNewUserResponse(key string, user *model.User, secret *string) *model.NewUserResponse {
	return &model.NewUserResponse{
		ID:      key,
		Roles:   user.Roles,
		Secret:  secret,
		Blocked: user.Blocked,
	}
}

func (db *Database) UserGetOrCreate(ctx *gin.Context, newUser *model.UserForm) (*model.UserResponse, error) {
	user, err := db.UserGet(ctx, newUser.ID)
	if err != nil {
		newUser, err := db.UserCreate(ctx, newUser)
		if err != nil {
			return nil, err
		}
		return &model.UserResponse{ID: newUser.ID, Roles: newUser.Roles, Blocked: newUser.Blocked}, nil
	}
	return user, nil
}

func (db *Database) UserCreate(ctx context.Context, newUser *model.UserForm) (*model.NewUserResponse, error) {
	var key string
	var hash *string
	if newUser.Apikey {
		key = generateKey()
		hash = pointer.String(fmt.Sprintf("%x", sha256.Sum256([]byte(key))))
	}

	var doc model.User
	newctx := driver.WithReturnNew(ctx, &doc)
	meta, err := db.userCollection.CreateDocument(ctx, newctx, strcase.ToKebab(newUser.ID), toUser(newUser, hash))
	if err != nil {
		return nil, err
	}

	return toNewUserResponse(meta.Key, &doc, pointer.String(key)), nil
}

func (db *Database) UserCreateSetupAPIKey(ctx context.Context, key string) (*model.UserResponse, error) {
	newUser := &model.UserForm{
		ID:      "setup",
		Roles:   []string{role.Admin},
		Apikey:  true,
		Blocked: false,
	}
	hash := pointer.String(fmt.Sprintf("%x", sha256.Sum256([]byte(key))))

	var doc model.User
	newctx := driver.WithReturnNew(ctx, &doc)
	meta, err := db.userCollection.CreateDocument(ctx, newctx, strcase.ToKebab(newUser.ID), toUser(newUser, hash))
	if err != nil {
		return nil, err
	}

	return toUserResponse(meta.Key, &doc), nil
}

func (db *Database) UserGet(ctx context.Context, id string) (*model.UserResponse, error) {
	var doc model.User
	meta, err := db.userCollection.ReadDocument(ctx, id, &doc)
	if err != nil {
		return nil, err
	}

	return toUserResponse(meta.Key, &doc), nil
}

func (db *Database) UserDelete(ctx context.Context, id string) error {
	_, err := db.userCollection.RemoveDocument(ctx, id)
	return err
}

func (db *Database) UserList(ctx context.Context) ([]*model.UserResponse, error) {
	query := "FOR d IN @@collection RETURN d"
	cursor, _, err := db.Query(ctx, query, map[string]interface{}{"@collection": UserCollectionName}, busdb.ReadOperation)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	var docs []*model.UserResponse
	for {
		var doc model.User
		meta, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		doc.Sha256 = nil
		docs = append(docs, toUserResponse(meta.Key, &doc))
	}

	return docs, err
}

func (db *Database) UserByHash(ctx context.Context, sha256 string) (*model.UserResponse, error) {
	query := `FOR d in @@collection
	FILTER d.sha256 == @sha256
	RETURN d`

	cursor, _, err := db.Query(ctx, query, map[string]interface{}{"@collection": UserCollectionName, "sha256": sha256}, busdb.ReadOperation)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var doc model.User
	meta, err := cursor.ReadDocument(ctx, &doc)
	if err != nil {
		return nil, err
	}

	return toUserResponse(meta.Key, &doc), err
}

func (db *Database) UserUpdate(ctx context.Context, id string, user *model.UserForm) (*model.UserResponse, error) {
	var doc model.User
	_, err := db.userCollection.ReadDocument(ctx, id, &doc)
	if err != nil {
		return nil, err
	}

	if doc.Sha256 != nil {
		return nil, errors.New("cannot update an API key")
	}

	ctx = driver.WithReturnNew(ctx, &doc)

	meta, err := db.userCollection.ReplaceDocument(ctx, id, toUser(user, nil))
	if err != nil {
		return nil, err
	}

	return toUserResponse(meta.Key, &doc), nil
}
