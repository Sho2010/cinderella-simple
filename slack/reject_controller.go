package slack

import (
	"github.com/Sho2010/cinderella-simple/domain/repository"
	"github.com/Sho2010/cinderella-simple/domain/service"
)

type RejectController struct {
}

func NewRejectController() RejectController {
	return RejectController{}
}

func (c *RejectController) Reject(userId string) error {
	repo := repository.DefaultClaimRepository()
	s := service.NewRejectClaimService(repo)

	if err := s.RejectClaim(userId); err != nil {
		return err
	}
	return nil
}
