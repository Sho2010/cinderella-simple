package claim

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Sho2010/cinderella-simple/encrypt"
	"k8s.io/apimachinery/pkg/util/validation"
)

type ClaimStatus string

const (
	ClaimStatusAccepted ClaimStatus = "accepted"
	ClaimStatusRejected ClaimStatus = "rejected"
	ClaimStatusPending  ClaimStatus = "pending"
	ClaimStatusExpired  ClaimStatus = "expired"
)

type ClaimValidationError struct {
	mes       string
	field     string
	errorType string
}

func (err *ClaimValidationError) Error() string {
	return fmt.Sprintf("Message:[%s], field:[%s], type:[%s]", err.mes, err.field, err.errorType)
}

func (err *ClaimValidationError) Is(e error) bool {
	if e2, ok := e.(*ClaimValidationError); ok {
		return (err.errorType == e2.errorType && err.field == e2.field)
	}
	return false
}

var (
	ErrorRequireNamespace = &ClaimValidationError{mes: "Require Namespace field", field: "Namespaces", errorType: "empty value"}
	ErrorRequireSubject   = &ClaimValidationError{mes: "Require Subject field", field: "Subject", errorType: "empty value"}
)

//TODO: implement period(制限時間)
type Claim interface {
	GetClaimDate() time.Time
	GetDescription() string
	GetEmail() string
	GetEncryptType() encrypt.EncryptType
	GetName() string
	GetNamespaces() []string
	GetState() ClaimStatus
	GetSubject() string
	GetLabels() map[string]string
	GetAnnotations() map[string]string
	GetServiceAccountName() (string, error)
	Validate() error

	//TODO: 暫定
	GetZipPassword() string
}

type ClaimBase struct {
	ClaimDate        time.Time           `json:"claim_date"`
	Description      string              `json:"description"`
	Email            string              `json:"email,omitempty"`
	EncryptType      encrypt.EncryptType `json:"encrypt_type,omitempty"`
	Name             string              `json:"name,omitempty"`
	Namespaces       []string            `json:"namespace"`
	State            ClaimStatus         `json:"status"`
	Subject          string              `json:"subject,omitempty"`
	ZipEncryptOption `json:"zip_option,omitempty"`
	GPGEncryptOption `json:"gpg_option,omitempty"`
}

type ZipEncryptOption struct {
	ZipPassword string `json:"-"`
}

type GPGEncryptOption struct {
	PublicKey string `json:"-"`
}

func (c *ClaimBase) GetClaimDate() time.Time {
	return c.ClaimDate
}

func (c *ClaimBase) GetDescription() string {
	return c.Description
}

func (c *ClaimBase) GetEncryptType() encrypt.EncryptType {
	return c.EncryptType
}

func (c *ClaimBase) GetNamespaces() []string {
	return c.Namespaces
}

func (c *ClaimBase) GetState() ClaimStatus {
	return c.State
}

func (c *ClaimBase) GetLabels() map[string]string {
	return map[string]string{
		"": "",
	}

	// return make(map[string]string)
}

func (c *ClaimBase) GetAnnotations() map[string]string {
	return map[string]string{
		"": "",
	}
}

func (c *ClaimBase) GetSubject() string {
	return c.Subject
}

func (c *ClaimBase) GetName() string {
	return c.Name
}

func (c *ClaimBase) GetEmail() string {
	return c.Email
}

func (c *ClaimBase) GetZipPassword() string {
	return c.ZipPassword
}

func (c *ClaimBase) GetServiceAccountName() (string, error) {
	s, err := NormalizeDNS1123(c.Subject)
	if err != nil {
		return "", fmt.Errorf("Name is not DNS1123: %w", err)
	}

	s = fmt.Sprintf("glass-shoes-%s", s)

	return s, nil
}

func (c *ClaimBase) Validate() error {
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
				mes:       fmt.Sprintf("Invalid namespace, [%s] is not RFC1123 format, %#v", v, errors),
				field:     "Namespaces",
				errorType: "RFC1123",
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
		return "", fmt.Errorf("%#v", errs)
	}
	return s, nil
}
