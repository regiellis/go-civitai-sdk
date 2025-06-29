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

// Package civitai - AI Resource (AIR) Specification Implementation
//
// This file implements the AI Resource (AIR) specification, a standardized
// way to identify and reference AI models, datasets, and other resources
// across different platforms and ecosystems.
//
// # AIR Format
//
// AIR follows this structure:
// air://{ecosystem}/{type}/{source}/{identifier}[/{version}][#{layer}][?{format}]
//
// Examples:
//   - air://civitai/model/133005                    (Basic model reference)
//   - air://civitai/model/133005/v1.0               (Specific version)
//   - air://civitai/lora/456789/v2.1#adapter        (LoRA with layer)
//   - air://civitai/model/133005?safetensor         (Format specification)
//
// # Creating AIR Identifiers
//
// Create AIR identifiers for CivitAI models:
//
//	// Basic model AIR
//	air := civitai.NewCivitAIModelAIR(133005)
//	fmt.Println(air.String()) // air://civitai/model/133005
//
//	// Model with version
//	air = civitai.NewCivitAIModelAIR(133005).WithVersion("v1.5")
//	fmt.Println(air.String()) // air://civitai/model/133005/v1.5
//
//	// LoRA with specific format
//	air = civitai.NewAIR("civitai", "lora", "civitai", "456789").
//		WithVersion("v2.0").
//		WithFormat("safetensor")
//	fmt.Println(air.String()) // air://civitai/lora/civitai/456789/v2.0?safetensor
//
// # Parsing AIR Strings
//
// Parse AIR strings back into structured data:
//
//	airStr := "air://civitai/model/133005/v1.0"
//	air, err := civitai.ParseAIR(airStr)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("Ecosystem: %s\n", air.Ecosystem)  // civitai
//	fmt.Printf("Type: %s\n", air.Type)            // model
//	fmt.Printf("ID: %s\n", air.Identifier)        // 133005
//	fmt.Printf("Version: %s\n", air.Version)      // v1.0
//
// # AIR Collections
//
// Work with collections of AIR identifiers:
//
//	collection := civitai.AIRCollection{
//		civitai.NewCivitAIModelAIR(133005),
//		civitai.NewCivitAIModelAIR(456789),
//		civitai.NewAIR("huggingface", "model", "microsoft", "DialoGPT-large"),
//	}
//
//	// Filter by ecosystem
//	civitaiOnly := collection.FilterByEcosystem("civitai")
//
//	// Filter by type
//	modelsOnly := collection.FilterByType("model")
//
//	// Get CivitAI models only
//	civitaiModels := collection.CivitAIOnly()
//
// # Client Integration
//
// Use AIR identifiers with the client:
//
//	client := civitai.NewClientWithoutAuth()
//
//	// Get model by AIR
//	air := civitai.NewCivitAIModelAIR(133005)
//	model, err := client.GetModelByAIR(ctx, air)
//
//	// Get model version by AIR
//	airWithVersion := civitai.NewCivitAIModelAIR(133005).WithVersion("456")
//	version, err := client.GetModelVersionByAIR(ctx, airWithVersion)
//
// # Model Conversion
//
// Convert existing models to AIR format:
//
//	// Convert Model to AIR
//	model := models[0]
//	air := civitai.ConvertModelToAIR(model, "civitai")
//
//	// Convert ModelVersion to AIR
//	version := model.ModelVersions[0]
//	versionAIR := civitai.ConvertVersionToAIR(version, "civitai")
//
// # Validation and Helper Methods
//
//	air := civitai.ParseAIR("air://civitai/model/133005/v1.0")
//
//	// Check if valid
//	if air.IsValid() {
//		fmt.Println("Valid AIR identifier")
//	}
//
//	// Check ecosystem
//	if air.IsCivitAI() {
//		fmt.Println("This is a CivitAI resource")
//	}
//
//	// Extract IDs
//	modelID := air.GetModelID()        // Returns 133005
//	versionID := air.GetVersionID()    // Returns parsed version ID
//
//	// Check specificity
//	if air.IsVersionSpecific() {
//		fmt.Println("Version is specified")
//	}
//
//	if air.IsFormatSpecific() {
//		fmt.Println("Format is specified")
//	}

package civitai

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// AIR represents an AI Resource Identifier
// Format: urn:air:{ecosystem}:{type}:{source}:{id}@{version?}:{layer?}.?{format?}
type AIR struct {
	// Core components (required)
	Ecosystem string // e.g., "sd1", "sd2", "sdxl", "gpt"
	Type      string // e.g., "model", "lora", "embedding"
	Source    string // e.g., "civitai", "huggingface", "openai"
	ID        string // Unique resource identifier

	// Optional components
	Version string // Specific model version
	Layer   string // Specific model layer
	Format  string // Model file format (e.g., "safetensors", "ckpt")

	// Raw AIR string for reference
	Raw string
}

// AIRType represents supported AIR resource types
type AIRType string

const (
	AIRTypeModel     AIRType = "model"
	AIRTypeLora      AIRType = "lora"
	AIRTypeEmbedding AIRType = "embedding"
	AIRTypeVAE       AIRType = "vae"
	AIRTypeControl   AIRType = "control"
)

// AIREcosystem represents supported AI ecosystems
type AIREcosystem string

const (
	AIREcosystemSD1  AIREcosystem = "sd1"
	AIREcosystemSD2  AIREcosystem = "sd2"
	AIREcosystemSDXL AIREcosystem = "sdxl"
	AIREcosystemGPT  AIREcosystem = "gpt"
	AIREcosystemFlux AIREcosystem = "flux"
)

// AIRSource represents supported source platforms
type AIRSource string

const (
	AIRSourceCivitAI     AIRSource = "civitai"
	AIRSourceHuggingFace AIRSource = "huggingface"
	AIRSourceOpenAI      AIRSource = "openai"
)

// Regular expression for parsing AIR identifiers
var airRegex = regexp.MustCompile(`^urn:air:([^:]+):([^:]+):([^:]+):([^@]+)(?:@([^:.]+))?(?::([^.]+))?(?:\.(.+))?$`)

// ParseAIR parses an AIR string into an AIR struct
func ParseAIR(airString string) (*AIR, error) {
	if airString == "" {
		return nil, errors.New("AIR string cannot be empty")
	}

	matches := airRegex.FindStringSubmatch(airString)
	if matches == nil {
		return nil, fmt.Errorf("invalid AIR format: %s", airString)
	}

	air := &AIR{
		Raw:       airString,
		Ecosystem: matches[1],
		Type:      matches[2],
		Source:    matches[3],
		ID:        matches[4],
		Version:   matches[5], // Optional, may be empty
		Layer:     matches[6], // Optional, may be empty
		Format:    matches[7], // Optional, may be empty
	}

	// Validate required components
	if err := air.Validate(); err != nil {
		return nil, fmt.Errorf("invalid AIR: %w", err)
	}

	return air, nil
}

// NewAIR creates a new AIR with required components
func NewAIR(ecosystem, resourceType, source, id string) *AIR {
	return &AIR{
		Ecosystem: ecosystem,
		Type:      resourceType,
		Source:    source,
		ID:        id,
	}
}

// NewCivitAIModelAIR creates an AIR for a CivitAI model
func NewCivitAIModelAIR(ecosystem string, modelID int, versionID ...int) *AIR {
	air := &AIR{
		Ecosystem: ecosystem,
		Type:      string(AIRTypeModel),
		Source:    string(AIRSourceCivitAI),
		ID:        strconv.Itoa(modelID),
	}

	if len(versionID) > 0 && versionID[0] > 0 {
		air.Version = strconv.Itoa(versionID[0])
	}

	return air
}

// WithVersion adds a version to the AIR
func (a *AIR) WithVersion(version string) *AIR {
	a.Version = version
	return a
}

// WithLayer adds a layer to the AIR
func (a *AIR) WithLayer(layer string) *AIR {
	a.Layer = layer
	return a
}

// WithFormat adds a format to the AIR
func (a *AIR) WithFormat(format string) *AIR {
	a.Format = format
	return a
}

// String returns the AIR as a formatted string
func (a *AIR) String() string {
	air := fmt.Sprintf("urn:air:%s:%s:%s:%s", a.Ecosystem, a.Type, a.Source, a.ID)

	if a.Version != "" {
		air += "@" + a.Version
	}

	if a.Layer != "" {
		air += ":" + a.Layer
	}

	if a.Format != "" {
		air += "." + a.Format
	}

	return air
}

// Validate checks if the AIR has valid required components
func (a *AIR) Validate() error {
	if a.Ecosystem == "" {
		return errors.New("ecosystem is required")
	}
	if a.Type == "" {
		return errors.New("type is required")
	}
	if a.Source == "" {
		return errors.New("source is required")
	}
	if a.ID == "" {
		return errors.New("ID is required")
	}

	// Validate ecosystem
	if !a.IsValidEcosystem() {
		return fmt.Errorf("unsupported ecosystem: %s", a.Ecosystem)
	}

	// Validate type
	if !a.IsValidType() {
		return fmt.Errorf("unsupported type: %s", a.Type)
	}

	// Validate source
	if !a.IsValidSource() {
		return fmt.Errorf("unsupported source: %s", a.Source)
	}

	return nil
}

// IsValidEcosystem checks if the ecosystem is supported
func (a *AIR) IsValidEcosystem() bool {
	validEcosystems := []string{
		string(AIREcosystemSD1),
		string(AIREcosystemSD2),
		string(AIREcosystemSDXL),
		string(AIREcosystemGPT),
		string(AIREcosystemFlux),
	}

	for _, valid := range validEcosystems {
		if a.Ecosystem == valid {
			return true
		}
	}
	return false
}

// IsValidType checks if the type is supported
func (a *AIR) IsValidType() bool {
	validTypes := []string{
		string(AIRTypeModel),
		string(AIRTypeLora),
		string(AIRTypeEmbedding),
		string(AIRTypeVAE),
		string(AIRTypeControl),
	}

	for _, valid := range validTypes {
		if a.Type == valid {
			return true
		}
	}
	return false
}

// IsValidSource checks if the source is supported
func (a *AIR) IsValidSource() bool {
	validSources := []string{
		string(AIRSourceCivitAI),
		string(AIRSourceHuggingFace),
		string(AIRSourceOpenAI),
	}

	for _, valid := range validSources {
		if a.Source == valid {
			return true
		}
	}
	return false
}

// IsCivitAI returns true if this AIR refers to a CivitAI resource
func (a *AIR) IsCivitAI() bool {
	return a.Source == string(AIRSourceCivitAI)
}

// GetModelID returns the model ID as an integer (for CivitAI resources)
func (a *AIR) GetModelID() (int, error) {
	if !a.IsCivitAI() {
		return 0, fmt.Errorf("not a CivitAI resource: %s", a.Source)
	}

	modelID, err := strconv.Atoi(a.ID)
	if err != nil {
		return 0, fmt.Errorf("invalid model ID: %s", a.ID)
	}

	return modelID, nil
}

// GetVersionID returns the version ID as an integer (for CivitAI resources)
func (a *AIR) GetVersionID() (int, error) {
	if !a.IsCivitAI() {
		return 0, fmt.Errorf("not a CivitAI resource: %s", a.Source)
	}

	if a.Version == "" {
		return 0, errors.New("no version specified in AIR")
	}

	versionID, err := strconv.Atoi(a.Version)
	if err != nil {
		return 0, fmt.Errorf("invalid version ID: %s", a.Version)
	}

	return versionID, nil
}

// ToModelType converts AIR type to CivitAI ModelType
func (a *AIR) ToModelType() ModelType {
	switch AIRType(a.Type) {
	case AIRTypeModel:
		// Determine specific model type based on ecosystem
		switch AIREcosystem(a.Ecosystem) {
		case AIREcosystemSD1, AIREcosystemSD2, AIREcosystemSDXL:
			return ModelTypeCheckpoint
		default:
			return ModelTypeCheckpoint
		}
	case AIRTypeLora:
		return ModelTypeLORA
	case AIRTypeEmbedding:
		return ModelTypeTextualInversion
	case AIRTypeVAE:
		return ModelTypeVAE
	case AIRTypeControl:
		return ModelTypeControlNet
	default:
		return ModelTypeCheckpoint
	}
}

// Equal compares two AIR identifiers for equality
func (a *AIR) Equal(other *AIR) bool {
	if other == nil {
		return false
	}

	return a.Ecosystem == other.Ecosystem &&
		a.Type == other.Type &&
		a.Source == other.Source &&
		a.ID == other.ID &&
		a.Version == other.Version &&
		a.Layer == other.Layer &&
		a.Format == other.Format
}

// IsVersionSpecific returns true if the AIR includes a specific version
func (a *AIR) IsVersionSpecific() bool {
	return a.Version != ""
}

// IsFormatSpecific returns true if the AIR includes a specific format
func (a *AIR) IsFormatSpecific() bool {
	return a.Format != ""
}

// HasLayer returns true if the AIR includes a layer specification
func (a *AIR) HasLayer() bool {
	return a.Layer != ""
}

// Clone creates a copy of the AIR
func (a *AIR) Clone() *AIR {
	return &AIR{
		Ecosystem: a.Ecosystem,
		Type:      a.Type,
		Source:    a.Source,
		ID:        a.ID,
		Version:   a.Version,
		Layer:     a.Layer,
		Format:    a.Format,
		Raw:       a.Raw,
	}
}

// AIRCollection represents a collection of AIR identifiers
type AIRCollection []*AIR

// FilterByEcosystem filters AIRs by ecosystem
func (ac AIRCollection) FilterByEcosystem(ecosystem string) AIRCollection {
	var result AIRCollection
	for _, air := range ac {
		if air.Ecosystem == ecosystem {
			result = append(result, air)
		}
	}
	return result
}

// FilterByType filters AIRs by type
func (ac AIRCollection) FilterByType(resourceType string) AIRCollection {
	var result AIRCollection
	for _, air := range ac {
		if air.Type == resourceType {
			result = append(result, air)
		}
	}
	return result
}

// FilterBySource filters AIRs by source
func (ac AIRCollection) FilterBySource(source string) AIRCollection {
	var result AIRCollection
	for _, air := range ac {
		if air.Source == source {
			result = append(result, air)
		}
	}
	return result
}

// CivitAIOnly returns only CivitAI AIRs
func (ac AIRCollection) CivitAIOnly() AIRCollection {
	return ac.FilterBySource(string(AIRSourceCivitAI))
}

// Strings returns all AIRs as formatted strings
func (ac AIRCollection) Strings() []string {
	result := make([]string, len(ac))
	for i, air := range ac {
		result[i] = air.String()
	}
	return result
}
