package storage

import (
	"context"
	"github.com/vovk404/course-platform/application-api/internal/entity"
	"github.com/vovk404/course-platform/application-api/pkg/database"
)

type nodeStorage struct {
	*database.PostgreSQL
}

var _ NodeStorage = (*nodeStorage)(nil)

func NewNodeStorage(postgresql *database.PostgreSQL) NodeStorage {
	return &nodeStorage{postgresql}
}

func (n nodeStorage) CreateNode(ctx context.Context, node *entity.Node) (*entity.Node, error) {
	err := n.DB.WithContext(ctx).Create(node).Error
	if err != nil {
		return nil, err
	}

	return node, nil
}
