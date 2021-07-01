package httphandle

// type BaseRepository struct {
// 	fires.Firebase
// }

// var (
// 	baseInteractor = usecase.NewBaseInteractor(
// 		&BaseRepository{
// 			Firebase: infrastructure.NewFireBaseClient(),
// 		},
// 	)
// )

// func (repo *BaseRepository) GetUserId(idToken string) (userId string, err error) {
// 	ctx := context.Background()
// 	client, err := repo.Firebase.Auth(ctx)
// 	token, err := client.VerifyIDToken(ctx, idToken)
// 	if err != nil {
// 		tools.ErrorLog(err)
// 	}
// 	claims := token.Claims
// 	fmt.Print(claims)
// 	userId = claims["user_id"].(string)
// 	return
// }

// func (repo *BaseRepository) GetClaims(idToken string) (claims map[string]interface{}, err error) {
// 	ctx := context.Background()
// 	client, err := repo.Firebase.Auth(ctx)
// 	token, err := client.VerifyIDToken(ctx, idToken)
// 	claims = token.Claims
// 	return
// }
