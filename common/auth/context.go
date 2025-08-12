package auth

import "context"

type contextKey string

const userContextKey contextKey = "user"

func GetUserInfo(ctx context.Context) (UserInfo, bool) {
	userInfo, ok := ctx.Value(userContextKey).(UserInfo)
	return userInfo, ok
}

func GetUserID(ctx context.Context) (int64, bool) {
	userInfo, ok := GetUserInfo(ctx)
	return userInfo.UserID, ok
}

func GetUserEmail(ctx context.Context) (string, bool) {
	userInfo, ok := GetUserInfo(ctx)
	return userInfo.Email, ok
}

func GetTenantID(ctx context.Context) (*string, bool) {
	userInfo, ok := GetUserInfo(ctx)
	return userInfo.TenantID, ok
}
