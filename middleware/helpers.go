// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package middleware

import "strings"

// IsAssetPath returns true if path is a static asset
func IsAssetPath(path string) bool {
	return strings.HasPrefix(path, "/assets/")
}
