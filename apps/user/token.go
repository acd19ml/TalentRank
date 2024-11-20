package user

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func VerifyToken(token string) bool {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// 检查 Token 是否有效
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		fmt.Println("Invalid Token:", err)
		return false
	}
	fmt.Printf("Token is valid. User: %s\n", user.GetLogin())
	return true
}
