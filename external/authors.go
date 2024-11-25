package external

import (
	"context"
	"fmt"

	"github.com/hilmiikhsan/library-book-service/constants"
	"github.com/hilmiikhsan/library-book-service/external/proto/author"
	"github.com/hilmiikhsan/library-book-service/helpers"
	"github.com/hilmiikhsan/library-book-service/internal/models"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func (e *External) GetDetailAuthor(ctx context.Context, id string) (models.AuthorModel, error) {
	var (
		res models.AuthorModel
	)

	conn, err := grpc.Dial(helpers.GetEnv("AUTHOR_GRPC_HOST", ""), grpc.WithInsecure())
	if err != nil {
		e.Logger.Error("external::GetDetailAuthor - failed to dial grpc server: ", err)
		return res, errors.Wrap(err, "failed to dial Author grpc")
	}
	defer conn.Close()

	client := author.NewAuthorServiceClient(conn)
	request := &author.AuthorRequest{
		Id: id,
	}

	resp, err := client.GetDetailAuthor(ctx, request)
	if err != nil {
		e.Logger.Error("external::GetDetailAuthor - failed to get detail Author: ", err)
		return res, errors.Wrap(err, "failed to get detail Author")
	}

	if resp.Message != constants.SuccessMessage {
		e.Logger.Error("external::GetDetailAuthor - invalid Author: ", resp.Message)
		return res, fmt.Errorf("got response error from authors: %s", resp.Message)
	}

	res.ID = resp.Data.Id
	res.Name = resp.Data.Name

	return res, nil
}
