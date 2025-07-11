package node

import (
	"context"
	"errors"
	"github.com/HuckOps/forge/model"
	"github.com/HuckOps/forge/server/repository/pagination"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetNodesByPagination(ctx context.Context, skip, limit int, filter bson.M) (pagination.PaginationResult[model.Node], error) {
	query := pagination.PaginationQuery[model.Node]{
		Repository: (&model.Node{}).Repository(),
		Skip:       skip,
		Limit:      limit,
		Filter:     filter,
		Sort:       bson.D{{"created_at", -1}},
	}
	return pagination.GetByPagination[model.Node](ctx, query)
}

func SetNodeLabel(ctx context.Context, nodeIDs []bson.ObjectID, labelIDs []bson.ObjectID) error {
	nodes, err := (&model.Node{}).Repository().FindByIDs(ctx, nodeIDs)
	if err != nil {
		return err
	}

	if len(nodes) != len(nodeIDs) {
		return errors.New("node ids not match")
	}

	labels, err := (&model.Label{}).Repository().FindByIDs(ctx, labelIDs)

	if err != nil {
		return err
	}

	if len(labels) != len(labelIDs) {
		return errors.New("label ids not match")
	}

	nodeLabels := make([]model.NodeLabel, 0)
	for _, nodeID := range nodeIDs {
		for _, labelID := range labelIDs {
			nodeLabels = append(nodeLabels, model.NodeLabel{
				NodeID:  nodeID,
				LabelID: labelID,
			})
		}
	}

	_, err = (&model.NodeLabel{}).Repository().CreateMany(ctx, nodeLabels)
	return err
}
