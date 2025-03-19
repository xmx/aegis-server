package jsmod

import (
	"context"

	"github.com/dop251/goja"
	"github.com/xmx/aegis-server/jsenv/jsvm"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewMongoDB(db *mongo.Database) jsvm.GlobalRegister {
	return &mongoDB{db: db}
}

type mongoDB struct {
	vm *goja.Runtime
	db *mongo.Database
}

func (db *mongoDB) RegisterGlobal(vm *goja.Runtime) error {
	db.vm = vm
	px := vm.NewProxy(vm.NewObject(), db.proxyTrapConfig())
	return vm.Set("db", px)
}

func (db *mongoDB) proxyTrapConfig() *goja.ProxyTrapConfig {
	props := map[string]goja.Value{
		"getName":            db.vm.ToValue(db.getName),
		"getCollectionNames": db.vm.ToValue(db.getCollectionNames),
	}

	return &goja.ProxyTrapConfig{
		Get: func(target *goja.Object, property string, receiver goja.Value) goja.Value {
			if val, exists := props[property]; exists {
				return val
			}

			coll := &mongoColl{
				vm:   db.vm,
				coll: db.db.Collection(property),
			}

			return db.vm.ToValue(coll)
		},
	}
}

func (db *mongoDB) getName() string {
	return db.db.Name()
}

func (db *mongoDB) getCollectionNames() ([]string, error) {
	return db.db.ListCollectionNames(context.Background(), bson.D{})
}

type mongoColl struct {
	vm   *goja.Runtime
	coll *mongo.Collection
}
