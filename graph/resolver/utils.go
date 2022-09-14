package resolver

import (
	"context"

	"github.com/garlicgarrison/chessvars-backend/middleware"
	"github.com/garlicgarrison/chessvars-backend/pkg/format"
)

func GetAuthUserID(ctx context.Context) (format.UserID, bool) {
	userID, ok := ctx.Value(middleware.AUTH_USER_CONTEXT_KEY).(format.UserID)
	return userID, ok
}
