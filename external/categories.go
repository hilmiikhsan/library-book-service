package external

import (
	"context"
	"fmt"

	"github.com/hilmiikhsan/library-book-service/constants"
	"github.com/hilmiikhsan/library-book-service/external/proto/category"
	"github.com/hilmiikhsan/library-book-service/helpers"
	"github.com/hilmiikhsan/library-book-service/internal/models"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func (e *External) GetDetailCategory(ctx context.Context, id string) (models.CategoryModel, error) {
	var (
		res models.CategoryModel
	)

	conn, err := grpc.Dial(helpers.GetEnv("CATEGORY_GRPC_HOST", ""), grpc.WithInsecure())
	if err != nil {
		e.Logger.Error("external::GetDetailCategory - failed to dial grpc server: ", err)
		return res, errors.Wrap(err, "failed to dial category grpc")
	}
	defer conn.Close()

	client := category.NewCategoryServiceClient(conn)
	request := &category.CategoryRequest{
		Id: id,
	}

	resp, err := client.GetDetailCategory(ctx, request)
	if err != nil {
		e.Logger.Error("external::GetDetailCategory - failed to get detail category: ", err)
		return res, errors.Wrap(err, "failed to get detail category")
	}

	if resp.Message != constants.SuccessMessage {
		e.Logger.Error("external::GetDetailCategory - invalid category: ", resp.Message)
		return res, fmt.Errorf("got response error from categories: %s", resp.Message)
	}

	res.ID = resp.Data.Id
	res.Name = resp.Data.Name

	return res, nil
}
