package service

import (
	"errors"

	"github.com/Sho2010/cinderella-simple/domain/repository"
	"github.com/Sho2010/cinderella-simple/k8s"
)

var (
	ErrInvalidClaimStatus = errors.New("invalid claim status, acceptable claim is only 'pending' state")
)

type AcceptClaimService struct {
	claimRepository repository.ClaimRepository
}

func NewAcceptClaimService(claimRepository repository.ClaimRepository) *AcceptClaimService {
	return &AcceptClaimService{
		claimRepository: claimRepository,
	}
}

func (s *AcceptClaimService) AcceptClaim(subject string) error {
	claim, err := s.claimRepository.FindBySubject(subject)
	if err != nil {
		return err
	}

	if claim.GetState() != "pending" {
		return ErrInvalidClaimStatus
	}

	claim.Accept()

	rc, err := k8s.NewResourceCreator(*claim)
	if err != nil {
		// pendingに戻す
		return err
	}

	if err := rc.Create(); err != nil {
		// pendingに戻す
		return err
	}

	err = s.claimRepository.Save(claim)
	if err != nil {
		return err
	}

	return nil
}
