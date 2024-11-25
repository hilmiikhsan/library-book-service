package external

import (
	"context"
	"fmt"

	"github.com/hilmiikhsan/library-book-service/constants"
	"github.com/hilmiikhsan/library-book-service/external/proto/tokenvalidation"
	"github.com/hilmiikhsan/library-book-service/helpers"
	"github.com/hilmiikhsan/library-book-service/internal/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type External struct {
	Logger *logrus.Logger
}

func (e *External) ValidateToken(ctx context.Context, token string) (models.TokenData, error) {
	var (
		res models.TokenData
	)

	conn, err := grpc.Dial(helpers.GetEnv("AUTH_GRPC_HOST", ""), grpc.WithInsecure())
	if err != nil {
		e.Logger.Error("external::ValidateToken - failed to dial grpc server: ", err)
		return res, errors.Wrap(err, "failed to dial ums grpc")
	}
	defer conn.Close()

	client := tokenvalidation.NewTokenValidationClient(conn)

	req := &tokenvalidation.TokenRequest{
		Token: token,
	}

	response, err := client.ValidateToken(ctx, req)
	if err != nil {
		e.Logger.Error("external::ValidateToken - failed to validate token: ", err)
		return res, errors.Wrap(err, "failed to validate token")
	}

	if response.Message != constants.SuccessMessage {
		e.Logger.Error("external::ValidateToken - invalid token: ", response.Message)
		return res, fmt.Errorf("got response error from auths: %s", response.Message)
	}

	res.UserID = response.Data.UserId
	res.Username = response.Data.Username
	res.FullName = response.Data.FullName
	res.Role = response.Data.Role

	return res, nil
}
