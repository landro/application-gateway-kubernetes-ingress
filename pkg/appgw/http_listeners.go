// -------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// --------------------------------------------------------------------------------------------

package appgw

import "k8s.io/api/extensions/v1beta1"

func (builder *appGwConfigBuilder) HTTPListeners(ingressList []*v1beta1.Ingress) (ConfigBuilder, error) {
	builder.appGwConfig.SslCertificates = builder.getSslCertificates(ingressList)
	builder.appGwConfig.FrontendPorts = builder.getFrontendPorts(ingressList)
	builder.appGwConfig.HTTPListeners, _ = builder.getFrontendListeners(ingressList)

	// App Gateway Rules can be configured to Redirect traffic to another listener (or external site)
	// The controller here re-configures App Gateway automatically by reading specific Ingress annotations.
	builder.appGwConfig.RedirectConfigurations = builder.getRedirectConfigurations(ingressList)

	return builder, nil
}
