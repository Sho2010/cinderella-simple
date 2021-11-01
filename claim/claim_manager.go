package claim

import (
	"fmt"
	"time"

	"github.com/Sho2010/cinderella-simple/audit"
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

	e := ClaimRegisterEvent{
		message: fmt.Sprintf("Slack User[%s]によって権限が請求されました, %s ", c.GetName(), c.GetEmail()),
		eventAt: time.Now(),
	}
	audit.PublishEvent(&e)
}

func (m *ClaimManager) findClaim(userId string) *Claim {
	for _, claim := range m.claims {
		if claim.GetSubject() == userId {
			return &claim
		}
	}
	return nil
}

func AddClaim(c Claim) {
	cmInstance.addClaim(c)
}

func FindClaim(userId string) *Claim {
	return cmInstance.findClaim(userId)
}
