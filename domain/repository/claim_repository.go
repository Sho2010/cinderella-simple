package repository

import (
	"fmt"

	"github.com/Sho2010/cinderella-simple/domain/model"
)

type ClaimRepository interface {
	FindBySubject(subject string) (*model.Claim, error)
	List() ([]model.Claim, error)
	Add(claim *model.Claim) error
	Save(claim *model.Claim) error
}

// 今の所永続化を考えてないのでrepositoryの実装もここにおいちゃう
type MemoryClaimRepository struct {
	claims []model.Claim
}

var _memoryClaimRepository MemoryClaimRepository

func init() {
	_memoryClaimRepository = MemoryClaimRepository{}
}

func GetMemoryClaimRepository() MemoryClaimRepository {
	return _memoryClaimRepository
}

func DefaultClaimRepository() ClaimRepository {
	// 他のrepositoryを使う場合はregistryに処理を移す
	return &_memoryClaimRepository
}

// implement ClaimRepository
func (repo *MemoryClaimRepository) FindBySubject(subject string) (*model.Claim, error) {
	for i, claim := range repo.claims {
		if claim.Subject == subject {
			return &repo.claims[i], nil
		}
	}
	return nil, fmt.Errorf("claim not found")
}

func (repo *MemoryClaimRepository) List() ([]model.Claim, error) {
	cp := make([]model.Claim, len(repo.claims))
	copy(cp, repo.claims)

	return cp, nil
}

func (repo *MemoryClaimRepository) Add(claim *model.Claim) error {
	repo.claims = append(repo.claims, *claim)
	return nil
}

func (repo *MemoryClaimRepository) Save(claim *model.Claim) error {
	return nil
	// for i, c := range repo.claims {
	// 	if c.Subject == claim.Subject {
	// 		repo.claims[i] = *claim
	// 		return nil
	// 	}
	// }
	// return fmt.Errorf("claim not found")
}
