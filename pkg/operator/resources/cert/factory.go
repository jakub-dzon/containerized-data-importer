/*
Copyright 2020 The CDI Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cert

import (
	"time"

	sdkcertapi "github.com/jakub-dzon/operator-cert-rotation-sdk/pkg/sdk/certrotation/api"

	"kubevirt.io/containerized-data-importer/pkg/operator/resources/utils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const day = 24 * time.Hour

// FactoryArgs contains the required parameters to generate certs
type FactoryArgs struct {
	Namespace string
}

// CreateCertificateDefinitions creates certificate definitions
func CreateCertificateDefinitions(args *FactoryArgs) []sdkcertapi.CertificateDefinition {
	defs := createCertificateDefinitions()
	for _, def := range defs {
		if def.SignerSecret != nil {
			addNamespace(args.Namespace, def.SignerSecret)
		}

		if def.CertBundleConfigmap != nil {
			addNamespace(args.Namespace, def.CertBundleConfigmap)
		}

		if def.TargetSecret != nil {
			addNamespace(args.Namespace, def.TargetSecret)
		}
	}

	return defs
}

func addNamespace(namespace string, obj metav1.Object) {
	if obj.GetNamespace() == "" {
		obj.SetNamespace(namespace)
	}
}

func createCertificateDefinitions() []sdkcertapi.CertificateDefinition {
	return []sdkcertapi.CertificateDefinition{
		{
			SignerSecret:        utils.ResourcesBuiler.CreateSecret("cdi-apiserver-signer"),
			SignerValidity:      30 * day,
			SignerRefresh:       15 * day,
			CertBundleConfigmap: utils.ResourcesBuiler.CreateConfigMap("cdi-apiserver-signer-bundle"),
			TargetSecret:        utils.ResourcesBuiler.CreateSecret("cdi-apiserver-server-cert"),
			TargetValidity:      48 * time.Hour,
			TargetRefresh:       24 * time.Hour,
			TargetService:       &[]string{"cdi-api"}[0],
		},
		{
			SignerSecret:        utils.ResourcesBuiler.CreateSecret("cdi-uploadproxy-signer"),
			SignerValidity:      30 * day,
			SignerRefresh:       15 * day,
			CertBundleConfigmap: utils.ResourcesBuiler.CreateConfigMap("cdi-uploadproxy-signer-bundle"),
			TargetSecret:        utils.ResourcesBuiler.CreateSecret("cdi-uploadproxy-server-cert"),
			TargetValidity:      48 * time.Hour,
			TargetRefresh:       24 * time.Hour,
			TargetService:       &[]string{"cdi-uploadproxy"}[0],
		},
		{
			SignerSecret:        utils.ResourcesBuiler.CreateSecret("cdi-uploadserver-signer"),
			SignerValidity:      10 * 365 * day,
			SignerRefresh:       8 * 365 * day,
			CertBundleConfigmap: utils.ResourcesBuiler.CreateConfigMap("cdi-uploadserver-signer-bundle"),
		},
		{
			SignerSecret:        utils.ResourcesBuiler.CreateSecret("cdi-uploadserver-client-signer"),
			SignerValidity:      10 * 365 * day,
			SignerRefresh:       8 * 365 * day,
			CertBundleConfigmap: utils.ResourcesBuiler.CreateConfigMap("cdi-uploadserver-client-signer-bundle"),
			TargetSecret:        utils.ResourcesBuiler.CreateSecret("cdi-uploadserver-client-cert"),
			TargetValidity:      48 * time.Hour,
			TargetRefresh:       24 * time.Hour,
			TargetUser:          &[]string{"client.upload-server.cdi.kubevirt.io"}[0],
		},
	}
}
