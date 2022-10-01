package database

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"log"
	"math/rand"

	"github.com/arangodb/go-driver"
	"github.com/iancoleman/strcase"
	maut "github.com/jonas-plum/maut/auth"

	"github.com/SecurityBrewery/catalyst/database/busdb"
	"github.com/SecurityBrewery/catalyst/generated/model"
	"github.com/SecurityBrewery/catalyst/generated/pointer"
	"github.com/SecurityBrewery/catalyst/generated/time"
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

func toUser(user *model.UserForm, salt, sha256, sha512 *string) *model.User {
	u := &model.User{
		Blocked: user.Blocked,
		Roles:   user.Roles,
		Salt:    salt,
		Sha256:  sha256,
		Sha512:  sha512,
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

func (db *Database) UserGetOrCreate(ctx context.Context, newUser *model.UserForm) (*model.UserResponse, error) {
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
	var key, salt, sha256Hash, sha512Hash *string
	if newUser.Apikey {
		key, sha256Hash = generateAPIKey()
	} else if newUser.Password != nil {
		salt, sha512Hash = hashUserPassword(newUser)
	}

	var doc model.User
	newctx := driver.WithReturnNew(ctx, &doc)
	meta, err := db.userCollection.CreateDocument(ctx, newctx, strcase.ToKebab(newUser.ID), toUser(newUser, salt, sha256Hash, sha512Hash))
	if err != nil {
		return nil, err
	}

	return toNewUserResponse(meta.Key, &doc, key), nil
}

func (db *Database) UserCreateSetupAPIKey(ctx context.Context, key string) (*model.UserResponse, error) {
	newUser := &model.UserForm{
		ID:      "setup",
		Roles:   []string{maut.AdminRole},
		Apikey:  true,
		Blocked: false,
	}
	sha256Hash := pointer.String(fmt.Sprintf("%x", sha256.Sum256([]byte(key))))

	var doc model.User
	newctx := driver.WithReturnNew(ctx, &doc)
	meta, err := db.userCollection.CreateDocument(ctx, newctx, strcase.ToKebab(newUser.ID), toUser(newUser, nil, sha256Hash, nil))
	if err != nil {
		return nil, err
	}

	return toUserResponse(meta.Key, &doc), nil
}

func (db *Database) UserUpdate(ctx context.Context, id string, user *model.UserForm) (*model.UserResponse, error) {
	var doc model.User
	_, err := db.userCollection.ReadDocument(ctx, id, &doc)
	if err != nil {
		return nil, err
	}

	if doc.Apikey {
		return nil, errors.New("cannot update an API key")
	}

	var salt, sha512Hash *string
	if user.Password != nil {
		salt, sha512Hash = hashUserPassword(user)
	} else {
		salt = doc.Salt
		sha512Hash = doc.Sha512
	}

	ctx = driver.WithReturnNew(ctx, &doc)

	user.ID = id

	meta, err := db.userCollection.ReplaceDocument(ctx, id, toUser(user, salt, nil, sha512Hash))
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
	cursor, _, err := db.Query(ctx, query, map[string]any{"@collection": UserCollectionName}, busdb.ReadOperation)
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

func (db *Database) UserAPIKeyByHash(ctx context.Context, sha256 string) (*model.UserResponse, error) {
	query := `FOR d in @@collection
	FILTER d.apikey && d.sha256 == @sha256
	RETURN d`

	vars := map[string]any{"@collection": UserCollectionName, "sha256": sha256}
	cursor, _, err := db.Query(ctx, query, vars, busdb.ReadOperation)
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

func (db *Database) UserByIDAndPassword(ctx context.Context, id, password string) (*model.UserResponse, error) {
	log.Println("UserByIDAndPassword", id, password)
	query := `FOR d in @@collection
	FILTER d._key == @id && !d.apikey && d.sha512 == SHA512(CONCAT(d.salt, @password))
	RETURN d`

	vars := map[string]any{"@collection": UserCollectionName, "id": id, "password": password}
	cursor, _, err := db.Query(ctx, query, vars, busdb.ReadOperation)
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

func generateAPIKey() (key, sha256Hash *string) {
	newKey := generateKey()
	sha256Hash = pointer.String(fmt.Sprintf("%x", sha256.Sum256([]byte(newKey))))

	return &newKey, sha256Hash
}

func hashUserPassword(newUser *model.UserForm) (salt, sha512Hash *string) {
	if newUser.Password != nil {
		saltKey := generateKey()
		salt = &saltKey
		sha512Hash = pointer.String(fmt.Sprintf("%x", sha512.Sum512([]byte(saltKey+*newUser.Password))))
	}

	return salt, sha512Hash
}
