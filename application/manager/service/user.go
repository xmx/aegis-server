package service

import (
	"context"
	"log/slog"

	"github.com/xmx/aegis-server/application/manager/request"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	db  *repository.BaseDB
	log *slog.Logger
}

func NewUser(db *repository.BaseDB, log *slog.Logger) *User {
	return &User{
		db:  db,
		log: log,
	}
}

func (usr *User) Page(ctx context.Context, req *request.Pages) (*repository.Pages[model.User], error) {
	coll := usr.db.User()
	return coll.Page(ctx, bson.D{}, req.Page, req.Size)
}
