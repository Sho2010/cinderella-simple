package claim

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Sho2010/cinderella-simple/encrypt"
	"github.com/slack-go/slack"
	"k8s.io/apimachinery/pkg/util/validation"
)

type ClaimStatus string

const (
	Accepted ClaimStatus = "accepted"
	Rejected ClaimStatus = "rejected"
	Pending  ClaimStatus = "pending"
	Expired  ClaimStatus = "expired"
)

type ClaimValidationError struct {
	mes   string
	field string
}

func (err *ClaimValidationError) Error() string {
	return fmt.Sprintf("%s, field: [%s]", err.mes, err.field)
}

var (
	ErrorRequireNamespace = &ClaimValidationError{mes: "Require Namespace field", field: "Namespace"}
	ErrorRequireSubject   = &ClaimValidationError{mes: "Require Subject field", field: "Subject"}
)

type ClaimInterface interface {
	GetLabels() map[string]string
	GetAnnotations() map[string]string
	GetServiceAccountName() string
}

type Claim struct {
	ClaimDate        slack.JSONTime      `json:"claim_date"`
	Email            string              `json:"email,omitempty"`
	EncryptType      encrypt.EncryptType `json:"encrypt_type,omitempty"`
	Name             string              `json:"name,omitempty"`
	Namespaces       []string            `json:"namespace"`
	State            ClaimStatus         `json:"status"`
	Subject          string              `json:"subject,omitempty"`
	ZipEncryptOption `json:"zip_option,omitempty"`
	GPGEncryptOption
}

type ZipEncryptOption struct {
	ZipPassword string `json:"-"`
}

type GPGEncryptOption struct {
	PublicKey string `json:"-"`
}

type SlackClaim struct {
	Claim
	SlackUser slack.User
}

func (c *Claim) GetServiceAccountName() (string, error) {
	s, err := NormalizeDNS1123(c.Subject)
	if err != nil {
		return "", fmt.Errorf("Name is not DNS1123: %w", err)
	}

	s = fmt.Sprintf("glass-shoes-%s", s)

	return s, nil
}

func (c *Claim) Validate() error {
	if len(c.Subject) == 0 {
		return ErrorRequireSubject
	}
	if len(c.Namespaces) == 0 {
		return ErrorRequireNamespace
	}

	for _, v := range c.Namespaces {
		errors := validation.IsDNS1123Label(v)

		if len(errors) != 0 {
			return &ClaimValidationError{
				mes:   fmt.Sprintf("Invalid namespace, [%s] is not RFC1123 format, %#v", v, errors),
				field: "Namespaces",
			}
		}
	}
	return nil
}

func NormalizeDNS1123(str string) (string, error) {
	//とりあえずよく使われそうな記号だけでも変換する
	rep := regexp.MustCompile("[@_;:.,=|/]")
	s := strings.ToLower(rep.ReplaceAllString(str, "-"))

	//変換かけてもダメならError
	errs := validation.IsDNS1123Label(s)
	if len(errs) > 0 {
		return "", fmt.Errorf("IsDNS1123Label errors: %#v", errs)
	}

	return s, nil
}
