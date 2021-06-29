package httphandle

import (
	"animar/v1/pkg/interfaces/fires"
	"animar/v1/pkg/tools/tools"
	"context"
	"fmt"
)

type BaseRepository struct {
	fires.Firebase
}

func (repo *BaseRepository) GetUserId(idToken string) (userId string, err error) {
	ctx := context.Background()
	client, err := repo.Firebase.Auth(ctx)
	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		tools.ErrorLog(err)
	}
	claims := token.Claims
	userId = claims["user_id"].(string)
	return
}

func (repo *BaseRepository) GetClaims(idToken string) (claims map[string]interface{}, err error) {
	fmt.Print(idToken)
	ctx := context.Background()
	client, err := repo.Firebase.Auth(ctx)
	token, err := client.VerifyIDToken(ctx, idToken)
	claims = token.Claims
	return
}
