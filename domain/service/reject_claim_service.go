package service

import (
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
		return ErrInvalidClaimStatus
	}

	claim.Reject()
	err = s.claimRepository.Save(claim)
	if err != nil {
		return err
	}
	return nil
}
