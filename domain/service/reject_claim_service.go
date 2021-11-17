package service

import (
	"github.com/Sho2010/cinderella-simple/domain/event"
	"github.com/Sho2010/cinderella-simple/domain/model"
	"github.com/Sho2010/cinderella-simple/domain/repository"
)

type RejectClaimService struct {
	claimRepository repository.ClaimRepository
}

func NewRejectClaimService(claimRepository repository.ClaimRepository) RejectClaimService {
	return RejectClaimService{
		claimRepository: claimRepository,
	}
}

func (s *RejectClaimService) RejectClaim(subject string) error {
	claim, err := s.claimRepository.FindBySubject(subject)
	if err != nil {
		return err
	}

	if claim.GetState() != "pending" {
	if claim.GetState() != model.ClaimStatusPending {
		return ErrInvalidClaimStatus
	}

	claim.Reject()
	err = s.claimRepository.Save(claim)
	if err != nil {
		return err
	}
	event.PublishClaimEvent(event.NewClaimEvent(claim.Subject, event.ClaimEventRejected))
	return nil
}
