// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package tool

import (
	genkitai "github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// RegisterTools registers all Ctx AI tools and returns their refs.
func RegisterTools(g *genkit.Genkit) []genkitai.ToolRef {
	return []genkitai.ToolRef{
		RegisterUpdateUserName(g),
		RegisterSummarizePlatformSettings(g),
		RegisterGetMaintenanceModeEnableInstructions(g),
		RegisterCreateWorkspace(g),
	}
}
