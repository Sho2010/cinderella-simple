package slack

import (
	"fmt"

	"github.com/Sho2010/cinderella-simple/claim"
)

type ClaimManager struct {
	claims []claim.Claim
}

func (m *ClaimManager) AddClaim(c claim.Claim) {
	m.claims = append(m.claims, c)
}

func (m *ClaimManager) FindClaim(userId string) (*claim.Claim, error) {
	for _, claim := range m.claims {
		if claim.GetSubject() == userId {
			return &claim, nil
		}

	}
	return nil, fmt.Errorf("Could not find claim")
}
