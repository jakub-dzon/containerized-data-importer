package api

import (
	"time"

	corev1 "k8s.io/api/core/v1"
)

// CertificateDefinition contains the data required to create/manage certtificate chains
type CertificateDefinition struct {
	// current CA key/cert
	SignerSecret   *corev1.Secret
	SignerValidity time.Duration
	SignerRefresh  time.Duration

	// all valid CA certs
	CertBundleConfigmap *corev1.ConfigMap

	// current key/cert for target
	TargetSecret   *corev1.Secret
	TargetValidity time.Duration
	TargetRefresh  time.Duration

	// only one of the following should be set
	// contains target key/cert for server
	TargetService *string
	// contains target user name
	TargetUser *string
}
