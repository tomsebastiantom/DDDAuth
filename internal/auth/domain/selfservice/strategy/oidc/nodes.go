// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package oidc

import (
	"my.com/secrets/internal/auth/domain/text"
	"my.com/secrets/internal/auth/domain/ui/node"
)

func NewLinkNode(provider string) *node.Node {
	return node.NewInputField("link", provider, node.OpenIDConnectGroup, node.InputAttributeTypeSubmit).WithMetaLabel(text.NewInfoSelfServiceSettingsUpdateLinkOIDC(provider))
}

func NewUnlinkNode(provider string) *node.Node {
	return node.NewInputField("unlink", provider, node.OpenIDConnectGroup, node.InputAttributeTypeSubmit).WithMetaLabel(text.NewInfoSelfServiceSettingsUpdateUnlinkOIDC(provider))
}
