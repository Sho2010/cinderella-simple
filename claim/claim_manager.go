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
var _cmInstance *ClaimManager

func init() {
	_cmInstance = &ClaimManager{}
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

func (m *ClaimManager) findClaim(userId string) Claim {
	for _, claim := range m.claims {
		if claim.GetSubject() == userId {
			return claim
		}
	}
	return nil
}

func AddClaim(c Claim) {
	_cmInstance.addClaim(c)
}

func FindClaim(userId string) Claim {
	return _cmInstance.findClaim(userId)
}

func ListClaims() []Claim {
	return _cmInstance.claims
}
