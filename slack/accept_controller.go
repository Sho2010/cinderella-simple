package slack

import (
	"github.com/Sho2010/cinderella-simple/domain/repository"
	"github.com/Sho2010/cinderella-simple/domain/service"
)

type AcceptController struct {
}

func NewAcceptController() AcceptController {
	return AcceptController{}
}

func (c *AcceptController) Accept(userId string) error {
	repo := repository.DefaultClaimRepository()
	s := service.NewAcceptClaimService(repo)
	if err := s.AcceptClaim(userId); err != nil {
		return err
	}
	return nil
}
