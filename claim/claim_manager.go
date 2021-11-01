package claim

import (
	"fmt"
)

//NOTE:
// Application全体として一つのClaimManagerで管理する
// あんまいい実装とは言えないけど今の所永続化したり、Datastore, Repository用意したりするまでもない

// ClaimManager唯一のインスタンス
var cmInstance *ClaimManager

func init() {
	cmInstance = &ClaimManager{}
}

type ClaimManager struct {
	claims []Claim
}

func (m *ClaimManager) addClaim(c Claim) {
	m.claims = append(m.claims, c)
}

func (m *ClaimManager) findClaim(userId string) (*Claim, error) {
	for _, claim := range m.claims {
		if claim.GetSubject() == userId {
			return &claim, nil
		}

	}
	return nil, fmt.Errorf("Could not find claim")
}

func AddClaim(c Claim) {
	cmInstance.addClaim(c)
}

func FindClaim(userId string) (*Claim, error) {
	return cmInstance.findClaim(userId)
}
