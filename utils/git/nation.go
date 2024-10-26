package git

import "context"

func (g *git) GetName(ctx context.Context, username string) (string, error) {
	user, err := g.GetUser(username)
	if err != nil {
		panic(err)
	}
	return user.GetName(), nil
}

func (g *git) GetCompany(ctx context.Context, username string) (string, error) {
	user, err := g.GetUser(username)
	if err != nil {
		panic(err)
	}
	return user.GetCompany(), nil
}

func (g *git) GetLocation(ctx context.Context, username string) (string, error) {
	user, err := g.GetUser(username)
	if err != nil {
		panic(err)
	}
	return user.GetLocation(), nil
}

func (g *git) GetEmail(ctx context.Context, username string) (string, error) {
	user, err := g.GetUser(username)
	if err != nil {
		panic(err)
	}
	return user.GetEmail(), nil
}

func (g *git) GetBio(ctx context.Context, username string) (string, error) {
	user, err := g.GetUser(username)
	if err != nil {
		panic(err)
	}
	return user.GetBio(), nil
}
