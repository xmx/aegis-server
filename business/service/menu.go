package service

import (
	"context"
	"log/slog"

	"github.com/xgfone/ship/v5"
	"github.com/xmx/aegis-server/datalayer/model"
	"github.com/xmx/aegis-server/datalayer/query"
)

func NewMenu(qry *query.Query, log *slog.Logger) *Menu {
	return &Menu{
		qry: qry,
		log: log,
	}
}

type Menu struct {
	qry *query.Query
	log *slog.Logger
}

func (svc *Menu) Tree(ctx context.Context, parentID ...int64) (model.MenuNodes, error) {
	tbl := svc.qry.Menu
	menus, err := tbl.WithContext(ctx).
		Order(tbl.Folder.Desc(), tbl.ID).
		Find()
	if err != nil {
		return nil, err
	}
	nodeMap := make(map[int64]*model.MenuNode, 128) // map[ID]Node
	roots := make(model.MenuNodes, 0, 30)
	misses := make(model.MenuNodes, 0, 100)
	for _, m := range menus {
		node := m.Node()
		id, pid := node.ID, node.ParentID
		nodeMap[id] = node
		if m.IsRoot() {
			roots = append(roots, node)
		} else if parent := nodeMap[pid]; parent != nil {
			parent.Children = append(parent.Children, node)
		} else {
			misses = append(misses, node)
		}
	}
	for _, m := range misses {
		pid := m.ParentID
		if parent := nodeMap[pid]; parent != nil {
			parent.Children = append(parent.Children, m)
		} else {
			roots = append(roots, m)
		}
	}

	roots.Sort()
	if len(parentID) == 0 || parentID[0] == 0 {
		return roots, nil
	}

	pid := parentID[0]
	node := nodeMap[pid]
	ret := make(model.MenuNodes, 0, 1)
	if node != nil {
		ret = append(ret, node)
	}

	return ret, nil
}

func (svc *Menu) Migrate(ctx context.Context, routes []ship.Route) error {
	for _, route := range routes {
		data := route.Data
		if data == nil {
			continue
		}
	}

	return nil
}
