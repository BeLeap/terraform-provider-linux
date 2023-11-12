package user

import (
	"context"
	"terraform-provider-linux/internal/lib"
)

func GetUser(ctx context.Context, session *lib.CustomSsh, username string) {
	stdout, sshErr := session.RunCommand(ctx, "getent passwd"+" "+username)
}
