package impl

import (
	"context"

	"github.com/acd19ml/TalentRank/apps/user"
	"github.com/acd19ml/TalentRank/utils"
)

func (u *UserServiceImpl) CreateUser(ctx context.Context, username string) (*user.User, error) {
	ins := user.NewUser()
	// 注入默认id
	ins.InjectDefault()
	handler := NewHandler()
	owner, err := handler.CreateUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	ins.Username = owner.Username
	ins.Name = owner.Name
	ins.Company = owner.Company
	ins.Location = owner.Location
	ins.Email = owner.Email
	ins.Bio = owner.Bio
	ins.TotalStar = owner.TotalStar
	ins.TotalFork = owner.TotalFork
	ins.Followers = owner.Followers
	ins.Dependents = owner.Dependents
	ins.Organizations = owner.Organizations
	ins.Readme = owner.Readme
	ins.Commits = owner.Commits

	if err := ins.Validate(); err != nil {
		return nil, err
	}

	return ins, nil
}

func NewHandler() *Handler {
	return &Handler{}
}

var handler = &Handler{}

type Handler struct {
	svc utils.Service
}

func (h *Handler) CreateUserByUsername(ctx context.Context, username string) (*user.User, error) {
	ins := user.NewUser()

	ins.Username = username

	name, err := h.svc.GetName(ctx, username)
	if err != nil {
		return nil, err
	}
	ins.Name = name

	company, err := h.svc.GetCompany(ctx, username)
	if err != nil {
		return nil, err
	}
	ins.Company = company

	location, err := h.svc.GetLocation(ctx, username)
	if err != nil {
		return nil, err
	}
	ins.Location = location

	email, err := h.svc.GetEmail(ctx, username)
	if err != nil {
		return nil, err
	}
	ins.Email = email

	bio, err := h.svc.GetBio(ctx, username)
	if err != nil {
		return nil, err
	}
	ins.Bio = bio

	totalStar, err := h.svc.GetTotalStars(ctx, username)
	if err != nil {
		return nil, err
	}
	ins.TotalStar = totalStar

	totalFork, err := h.svc.GetTotalForks(ctx, username)
	if err != nil {
		return nil, err
	}
	ins.TotalFork = totalFork

	followers, err := h.svc.GetFollowers(ctx, username)
	if err != nil {
		return nil, err
	}
	ins.Followers = followers

	dependents, err := h.svc.GetDependentRepositories(ctx, username)
	if err != nil {
		return nil, err
	}
	ins.Dependents = dependents

	organizations, err := h.svc.GetOrganizations(ctx, username)
	if err != nil {
		return nil, err
	}
	ins.Organizations = organizations

	readme, err := h.svc.GetReadme(ctx, username, 100)
	if err != nil {
		return nil, err
	}
	ins.Readme = readme

	commits, err := h.svc.GetCommits(ctx, username, 100)
	if err != nil {
		return nil, err
	}
	ins.Commits = commits
	return ins, nil

}
