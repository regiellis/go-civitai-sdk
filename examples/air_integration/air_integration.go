//go:build ignore

/*
Copyright (c) 2025 Regi Ellis

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"context"
	"fmt"
	"log"

	civitai "github.com/regiellis/go-civitai-sdk"
)

func main() {
	fmt.Println("=== CivitAI SDK - AIR (AI Resource Identifier) Integration Examples ===")

	// Example 1: Parse AIR identifiers
	fmt.Println("\n--- Example 1: Parsing AIR Identifiers ---")

	// Parse different types of AIR identifiers
	airExamples := []string{
		"urn:air:sdxl:model:civitai:2421",                         // Basic model
		"urn:air:sdxl:model:civitai:2421@43533",                   // Model with version
		"urn:air:sdxl:lora:civitai:328553@368189",                 // LoRA with version
		"urn:air:sd1:model:civitai:2421@43533:layer1.safetensors", // Full AIR
		"urn:air:gpt:model:openai:gpt@4",                          // Non-CivitAI example
	}

	for _, airString := range airExamples {
		air, err := civitai.ParseAIR(airString)
		if err != nil {
			fmt.Printf("❌ Failed to parse %s: %v\n", airString, err)
			continue
		}

		fmt.Printf("✅ Parsed: %s\n", airString)
		fmt.Printf("   Ecosystem: %s, Type: %s, Source: %s, ID: %s\n",
			air.Ecosystem, air.Type, air.Source, air.ID)

		if air.IsVersionSpecific() {
			fmt.Printf("   Version: %s\n", air.Version)
		}
		if air.IsFormatSpecific() {
			fmt.Printf("   Format: %s\n", air.Format)
		}
		if air.HasLayer() {
			fmt.Printf("   Layer: %s\n", air.Layer)
		}
		fmt.Println()
	}

	// Example 2: Create AIR identifiers
	fmt.Println("\n--- Example 2: Creating AIR Identifiers ---")

	// Create AIR for CivitAI model
	modelAIR := civitai.NewCivitAIModelAIR("sdxl", 2421)
	fmt.Printf("Basic CivitAI model AIR: %s\n", modelAIR.String())

	// Create AIR with version
	versionAIR := civitai.NewCivitAIModelAIR("sdxl", 2421, 43533)
	fmt.Printf("CivitAI model with version: %s\n", versionAIR.String())

	// Create AIR with fluent interface
	complexAIR := civitai.NewAIR("sdxl", "lora", "civitai", "328553").
		WithVersion("368189").
		WithFormat("safetensors")
	fmt.Printf("Complex AIR with fluent interface: %s\n", complexAIR.String())

	// Example 3: Using AIR with the CivitAI client
	fmt.Println("\n--- Example 3: Client Integration with AIR ---")

	client := civitai.NewClientWithoutAuth()
	ctx := context.Background()
	// Unused variables for demonstration purposes
	_ = client
	_ = ctx

	// Parse an AIR and use it to fetch the model
	targetAIR, err := civitai.ParseAIR("urn:air:sdxl:model:civitai:2421")
	if err != nil {
		log.Printf("Failed to parse AIR: %v", err)
	} else {
		fmt.Printf("Fetching model using AIR: %s\n", targetAIR.String())

		// This would work with a real API token and model ID
		// model, err := client.GetModelByAIR(ctx, targetAIR)
		// if err != nil {
		//     fmt.Printf("Failed to fetch model: %v\n", err)
		// } else {
		//     fmt.Printf("Successfully fetched model: %s\n", model.Name)
		// }

		fmt.Printf("Would fetch model ID: %d from CivitAI\n", 2421)
	}

	// Example 4: Convert models to AIR
	fmt.Println("\n--- Example 4: Converting Models to AIR ---")

	// Create example model
	exampleModel := &civitai.Model{
		ID:   12345,
		Name: "Example Anime Model",
		Type: civitai.ModelTypeCheckpoint,
		Tags: []string{
			"anime",
			"character",
			"SD 1.5",
		},
	}

	// Convert to AIR with auto-detected ecosystem
	autoAIR := exampleModel.ToAIR("")
	fmt.Printf("Auto-detected ecosystem AIR: %s\n", autoAIR.String())

	// Convert to AIR with specific ecosystem
	sdxlAIR := exampleModel.ToAIR("sdxl", 67890)
	fmt.Printf("SDXL ecosystem with version: %s\n", sdxlAIR.String())

	// Use convenience method
	fluxAIR := exampleModel.GetAIRForEcosystem(civitai.AIREcosystemFlux)
	fmt.Printf("Flux ecosystem AIR: %s\n", fluxAIR.String())

	// Example 5: AIR collections and filtering
	fmt.Println("\n--- Example 5: AIR Collections ---")

	collection := civitai.AIRCollection{
		civitai.NewCivitAIModelAIR("sdxl", 1),
		civitai.NewCivitAIModelAIR("sd1", 2),
		civitai.NewAIR("sdxl", "lora", "civitai", "3"),
		civitai.NewAIR("gpt", "model", "openai", "gpt-4"),
	}

	fmt.Printf("Total AIRs in collection: %d\n", len(collection))

	// Filter by ecosystem
	sdxlAIRs := collection.FilterByEcosystem("sdxl")
	fmt.Printf("SDXL AIRs: %d\n", len(sdxlAIRs))

	// Filter by type
	loraAIRs := collection.FilterByType("lora")
	fmt.Printf("LoRA AIRs: %d\n", len(loraAIRs))

	// Get only CivitAI AIRs
	civitaiAIRs := collection.CivitAIOnly()
	fmt.Printf("CivitAI AIRs: %d\n", len(civitaiAIRs))

	// Convert to strings
	airStrings := collection.Strings()
	fmt.Println("All AIR strings:")
	for i, airStr := range airStrings {
		fmt.Printf("  %d. %s\n", i+1, airStr)
	}

	// Example 6: AIR validation and helpers
	fmt.Println("\n--- Example 6: AIR Validation and Helpers ---")

	validAIR := civitai.NewCivitAIModelAIR("sdxl", 2421, 43533)

	fmt.Printf("AIR: %s\n", validAIR.String())
	fmt.Printf("Is CivitAI: %t\n", validAIR.IsCivitAI())
	fmt.Printf("Is version specific: %t\n", validAIR.IsVersionSpecific())
	fmt.Printf("Has layer: %t\n", validAIR.HasLayer())
	fmt.Printf("Is format specific: %t\n", validAIR.IsFormatSpecific())

	if validAIR.IsCivitAI() {
		modelID, _ := validAIR.GetModelID()
		versionID, _ := validAIR.GetVersionID()
		fmt.Printf("Model ID: %d, Version ID: %d\n", modelID, versionID)
	}

	// Convert to CivitAI model type
	modelType := validAIR.ToModelType()
	fmt.Printf("Converted to CivitAI model type: %s\n", modelType)

	// Clone AIR
	clonedAIR := validAIR.Clone()
	fmt.Printf("Cloned AIR equals original: %t\n", validAIR.Equal(clonedAIR))

	// Modify clone to test independence
	clonedAIR.Version = "99999"
	fmt.Printf("After modification, equals original: %t\n", validAIR.Equal(clonedAIR))

	// Example 7: Future-proofing with AIR
	fmt.Println("\n--- Example 7: Future-Proofing with AIR ---")

	fmt.Println("Benefits of using AIR:")
	fmt.Println("✅ Standardized resource identification across AI platforms")
	fmt.Println("✅ Future compatibility when CivitAI adopts AIR as default")
	fmt.Println("✅ Easy migration between different AI resource platforms")
	fmt.Println("✅ Clear resource versioning and format specification")
	fmt.Println("✅ Type-safe resource handling in applications")

	fmt.Println("\nExample use cases:")
	fmt.Println("• Model sharing: Share 'urn:air:sdxl:model:civitai:2421@43533' instead of URLs")
	fmt.Println("• Automation: Build pipelines that work with any AIR-compatible platform")
	fmt.Println("• Version control: Track model versions across different ecosystems")
	fmt.Println("• Resource discovery: Search and filter by ecosystem, type, or source")

	fmt.Println("\n=== AIR Integration Examples Complete ===")
	fmt.Println("The Go CivitAI SDK is now ready for the future of AI resource identification!")
}
