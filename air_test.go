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

package civitai

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAIRParsing(t *testing.T) {
	testCases := []struct {
		name        string
		airString   string
		expected    *AIR
		expectError bool
	}{
		{
			name:      "Basic CivitAI model",
			airString: "urn:air:sdxl:model:civitai:2421",
			expected: &AIR{
				Ecosystem: "sdxl",
				Type:      "model",
				Source:    "civitai",
				ID:        "2421",
			},
			expectError: false,
		},
		{
			name:      "CivitAI model with version",
			airString: "urn:air:sd1:model:civitai:2421@43533",
			expected: &AIR{
				Ecosystem: "sd1",
				Type:      "model",
				Source:    "civitai",
				ID:        "2421",
				Version:   "43533",
			},
			expectError: false,
		},
		{
			name:      "CivitAI LoRA with version",
			airString: "urn:air:sdxl:lora:civitai:328553@368189",
			expected: &AIR{
				Ecosystem: "sdxl",
				Type:      "lora",
				Source:    "civitai",
				ID:        "328553",
				Version:   "368189",
			},
			expectError: false,
		},
		{
			name:      "Full AIR with layer and format",
			airString: "urn:air:sdxl:model:civitai:2421@43533:layer1.safetensors",
			expected: &AIR{
				Ecosystem: "sdxl",
				Type:      "model",
				Source:    "civitai",
				ID:        "2421",
				Version:   "43533",
				Layer:     "layer1",
				Format:    "safetensors",
			},
			expectError: false,
		},
		{
			name:      "GPT model",
			airString: "urn:air:gpt:model:openai:gpt@4",
			expected: &AIR{
				Ecosystem: "gpt",
				Type:      "model",
				Source:    "openai",
				ID:        "gpt",
				Version:   "4",
			},
			expectError: false,
		},
		{
			name:        "Invalid format - missing prefix",
			airString:   "air:sdxl:model:civitai:2421",
			expectError: true,
		},
		{
			name:        "Invalid format - missing components",
			airString:   "urn:air:sdxl:model",
			expectError: true,
		},
		{
			name:        "Empty string",
			airString:   "",
			expectError: true,
		},
		{
			name:        "Invalid ecosystem",
			airString:   "urn:air:invalid:model:civitai:2421",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			air, err := ParseAIR(tc.airString)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tc.airString)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tc.airString, err)
				return
			}

			if air.Ecosystem != tc.expected.Ecosystem {
				t.Errorf("Expected ecosystem %s, got %s", tc.expected.Ecosystem, air.Ecosystem)
			}
			if air.Type != tc.expected.Type {
				t.Errorf("Expected type %s, got %s", tc.expected.Type, air.Type)
			}
			if air.Source != tc.expected.Source {
				t.Errorf("Expected source %s, got %s", tc.expected.Source, air.Source)
			}
			if air.ID != tc.expected.ID {
				t.Errorf("Expected ID %s, got %s", tc.expected.ID, air.ID)
			}
			if air.Version != tc.expected.Version {
				t.Errorf("Expected version %s, got %s", tc.expected.Version, air.Version)
			}
			if air.Layer != tc.expected.Layer {
				t.Errorf("Expected layer %s, got %s", tc.expected.Layer, air.Layer)
			}
			if air.Format != tc.expected.Format {
				t.Errorf("Expected format %s, got %s", tc.expected.Format, air.Format)
			}
		})
	}
}

func TestAIRConstruction(t *testing.T) {
	t.Run("NewAIR", func(t *testing.T) {
		air := NewAIR("sdxl", "model", "civitai", "2421")

		if air.Ecosystem != "sdxl" {
			t.Errorf("Expected ecosystem sdxl, got %s", air.Ecosystem)
		}
		if air.Type != "model" {
			t.Errorf("Expected type model, got %s", air.Type)
		}
		if air.Source != "civitai" {
			t.Errorf("Expected source civitai, got %s", air.Source)
		}
		if air.ID != "2421" {
			t.Errorf("Expected ID 2421, got %s", air.ID)
		}
	})

	t.Run("NewCivitAIModelAIR", func(t *testing.T) {
		air := NewCivitAIModelAIR("sdxl", 2421)

		if air.Ecosystem != "sdxl" {
			t.Errorf("Expected ecosystem sdxl, got %s", air.Ecosystem)
		}
		if air.Type != "model" {
			t.Errorf("Expected type model, got %s", air.Type)
		}
		if air.Source != "civitai" {
			t.Errorf("Expected source civitai, got %s", air.Source)
		}
		if air.ID != "2421" {
			t.Errorf("Expected ID 2421, got %s", air.ID)
		}
	})

	t.Run("NewCivitAIModelAIR with version", func(t *testing.T) {
		air := NewCivitAIModelAIR("sdxl", 2421, 43533)

		if air.Version != "43533" {
			t.Errorf("Expected version 43533, got %s", air.Version)
		}
	})

	t.Run("Fluent interface", func(t *testing.T) {
		air := NewAIR("sdxl", "model", "civitai", "2421").
			WithVersion("43533").
			WithLayer("layer1").
			WithFormat("safetensors")

		if air.Version != "43533" {
			t.Errorf("Expected version 43533, got %s", air.Version)
		}
		if air.Layer != "layer1" {
			t.Errorf("Expected layer layer1, got %s", air.Layer)
		}
		if air.Format != "safetensors" {
			t.Errorf("Expected format safetensors, got %s", air.Format)
		}
	})
}

func TestAIRString(t *testing.T) {
	testCases := []struct {
		name     string
		air      *AIR
		expected string
	}{
		{
			name: "Basic AIR",
			air: &AIR{
				Ecosystem: "sdxl",
				Type:      "model",
				Source:    "civitai",
				ID:        "2421",
			},
			expected: "urn:air:sdxl:model:civitai:2421",
		},
		{
			name: "AIR with version",
			air: &AIR{
				Ecosystem: "sdxl",
				Type:      "model",
				Source:    "civitai",
				ID:        "2421",
				Version:   "43533",
			},
			expected: "urn:air:sdxl:model:civitai:2421@43533",
		},
		{
			name: "Full AIR",
			air: &AIR{
				Ecosystem: "sdxl",
				Type:      "model",
				Source:    "civitai",
				ID:        "2421",
				Version:   "43533",
				Layer:     "layer1",
				Format:    "safetensors",
			},
			expected: "urn:air:sdxl:model:civitai:2421@43533:layer1.safetensors",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.air.String()
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestAIRValidation(t *testing.T) {
	t.Run("Valid AIR", func(t *testing.T) {
		air := NewCivitAIModelAIR("sdxl", 2421)
		err := air.Validate()
		if err != nil {
			t.Errorf("Expected valid AIR, got error: %v", err)
		}
	})

	t.Run("Missing ecosystem", func(t *testing.T) {
		air := &AIR{Type: "model", Source: "civitai", ID: "2421"}
		err := air.Validate()
		if err == nil {
			t.Error("Expected error for missing ecosystem")
		}
	})

	t.Run("Invalid ecosystem", func(t *testing.T) {
		air := &AIR{Ecosystem: "invalid", Type: "model", Source: "civitai", ID: "2421"}
		err := air.Validate()
		if err == nil {
			t.Error("Expected error for invalid ecosystem")
		}
	})

	t.Run("Invalid type", func(t *testing.T) {
		air := &AIR{Ecosystem: "sdxl", Type: "invalid", Source: "civitai", ID: "2421"}
		err := air.Validate()
		if err == nil {
			t.Error("Expected error for invalid type")
		}
	})

	t.Run("Invalid source", func(t *testing.T) {
		air := &AIR{Ecosystem: "sdxl", Type: "model", Source: "invalid", ID: "2421"}
		err := air.Validate()
		if err == nil {
			t.Error("Expected error for invalid source")
		}
	})
}

func TestAIRHelperMethods(t *testing.T) {
	air := NewCivitAIModelAIR("sdxl", 2421, 43533).WithFormat("safetensors")

	t.Run("IsCivitAI", func(t *testing.T) {
		if !air.IsCivitAI() {
			t.Error("Expected true for CivitAI AIR")
		}

		nonCivitAI := NewAIR("gpt", "model", "openai", "gpt-4")
		if nonCivitAI.IsCivitAI() {
			t.Error("Expected false for non-CivitAI AIR")
		}
	})

	t.Run("GetModelID", func(t *testing.T) {
		modelID, err := air.GetModelID()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if modelID != 2421 {
			t.Errorf("Expected model ID 2421, got %d", modelID)
		}
	})

	t.Run("GetVersionID", func(t *testing.T) {
		versionID, err := air.GetVersionID()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if versionID != 43533 {
			t.Errorf("Expected version ID 43533, got %d", versionID)
		}
	})

	t.Run("IsVersionSpecific", func(t *testing.T) {
		if !air.IsVersionSpecific() {
			t.Error("Expected true for version-specific AIR")
		}

		noVersion := NewCivitAIModelAIR("sdxl", 2421)
		if noVersion.IsVersionSpecific() {
			t.Error("Expected false for non-version-specific AIR")
		}
	})

	t.Run("IsFormatSpecific", func(t *testing.T) {
		if !air.IsFormatSpecific() {
			t.Error("Expected true for format-specific AIR")
		}
	})

	t.Run("Equal", func(t *testing.T) {
		air1 := NewCivitAIModelAIR("sdxl", 2421, 43533)
		air2 := NewCivitAIModelAIR("sdxl", 2421, 43533)
		air3 := NewCivitAIModelAIR("sdxl", 2421, 43534)

		if !air1.Equal(air2) {
			t.Error("Expected equal AIRs to be equal")
		}

		if air1.Equal(air3) {
			t.Error("Expected different AIRs to not be equal")
		}

		if air1.Equal(nil) {
			t.Error("Expected AIR to not equal nil")
		}
	})

	t.Run("Clone", func(t *testing.T) {
		clone := air.Clone()
		if !air.Equal(clone) {
			t.Error("Expected clone to be equal to original")
		}

		// Modify clone to ensure independence
		clone.Version = "99999"
		if air.Equal(clone) {
			t.Error("Expected modified clone to not equal original")
		}
	})
}

func TestAIRTypeConversion(t *testing.T) {
	testCases := []struct {
		airType  string
		expected ModelType
	}{
		{"model", ModelTypeCheckpoint},
		{"lora", ModelTypeLORA},
		{"embedding", ModelTypeTextualInversion},
		{"vae", ModelTypeVAE},
		{"control", ModelTypeControlNet},
	}

	for _, tc := range testCases {
		t.Run(tc.airType, func(t *testing.T) {
			air := &AIR{Type: tc.airType}
			result := air.ToModelType()
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestAIRCollection(t *testing.T) {
	collection := AIRCollection{
		NewCivitAIModelAIR("sdxl", 1),
		NewCivitAIModelAIR("sd1", 2),
		NewAIR("sdxl", "lora", "civitai", "3"),
		NewAIR("gpt", "model", "openai", "gpt-4"),
	}

	t.Run("FilterByEcosystem", func(t *testing.T) {
		sdxl := collection.FilterByEcosystem("sdxl")
		if len(sdxl) != 2 {
			t.Errorf("Expected 2 SDXL AIRs, got %d", len(sdxl))
		}
	})

	t.Run("FilterByType", func(t *testing.T) {
		models := collection.FilterByType("model")
		if len(models) != 3 {
			t.Errorf("Expected 3 model AIRs, got %d", len(models))
		}
	})

	t.Run("FilterBySource", func(t *testing.T) {
		civitai := collection.FilterBySource("civitai")
		if len(civitai) != 3 {
			t.Errorf("Expected 3 CivitAI AIRs, got %d", len(civitai))
		}
	})

	t.Run("CivitAIOnly", func(t *testing.T) {
		civitai := collection.CivitAIOnly()
		if len(civitai) != 3 {
			t.Errorf("Expected 3 CivitAI AIRs, got %d", len(civitai))
		}
	})

	t.Run("Strings", func(t *testing.T) {
		strings := collection.Strings()
		if len(strings) != 4 {
			t.Errorf("Expected 4 string representations, got %d", len(strings))
		}
	})
}

func TestClientAIRIntegration(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if r.URL.Path == "/models/2421" {
			w.Write([]byte(`{"id": 2421, "name": "Test Model", "type": "Checkpoint", "createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z"}`))
		} else if r.URL.Path == "/model-versions/43533" {
			w.Write([]byte(`{"id": 43533, "name": "Test Version", "modelId": 2421, "createdAt": "2024-01-01T00:00:00Z", "updatedAt": "2024-01-01T00:00:00Z"}`))
		}
	}))
	defer server.Close()

	client := NewClientWithoutAuth(WithBaseURL(server.URL))
	ctx := context.Background()

	t.Run("GetModelByAIR", func(t *testing.T) {
		air := NewCivitAIModelAIR("sdxl", 2421)
		model, err := client.GetModelByAIR(ctx, air)

		if err != nil {
			t.Fatalf("GetModelByAIR failed: %v", err)
		}

		if model.ID != 2421 {
			t.Errorf("Expected model ID 2421, got %d", model.ID)
		}
	})

	t.Run("GetModelVersionByAIR", func(t *testing.T) {
		air := NewCivitAIModelAIR("sdxl", 2421, 43533)
		version, err := client.GetModelVersionByAIR(ctx, air)

		if err != nil {
			t.Fatalf("GetModelVersionByAIR failed: %v", err)
		}

		if version.ID != 43533 {
			t.Errorf("Expected version ID 43533, got %d", version.ID)
		}
	})

	t.Run("GetModelByAIR with non-CivitAI source", func(t *testing.T) {
		air := NewAIR("gpt", "model", "openai", "gpt-4")
		_, err := client.GetModelByAIR(ctx, air)

		if err == nil {
			t.Error("Expected error for non-CivitAI AIR")
		}
	})

	t.Run("GetModelVersionByAIR without version", func(t *testing.T) {
		air := NewCivitAIModelAIR("sdxl", 2421) // No version
		_, err := client.GetModelVersionByAIR(ctx, air)

		if err == nil {
			t.Error("Expected error for AIR without version")
		}
	})
}

func TestConvertModelToAIR(t *testing.T) {
	model := &Model{
		ID:   2421,
		Name: "Test Model",
		Type: ModelTypeCheckpoint,
		Tags: []string{"SD 1.5"},
	}

	t.Run("Convert with explicit ecosystem", func(t *testing.T) {
		air := ConvertModelToAIR(model, "sdxl")

		if air.Ecosystem != "sdxl" {
			t.Errorf("Expected ecosystem sdxl, got %s", air.Ecosystem)
		}
		if air.Type != "model" {
			t.Errorf("Expected type model, got %s", air.Type)
		}
		if air.ID != "2421" {
			t.Errorf("Expected ID 2421, got %s", air.ID)
		}
	})

	t.Run("Convert with inferred ecosystem", func(t *testing.T) {
		air := ConvertModelToAIR(model, "")

		if air.Ecosystem != "sd1" {
			t.Errorf("Expected inferred ecosystem sd1, got %s", air.Ecosystem)
		}
	})

	t.Run("Convert with version", func(t *testing.T) {
		air := ConvertModelToAIR(model, "sdxl", 43533)

		if air.Version != "43533" {
			t.Errorf("Expected version 43533, got %s", air.Version)
		}
	})

	t.Run("Convert nil model", func(t *testing.T) {
		air := ConvertModelToAIR(nil, "sdxl")
		if air != nil {
			t.Error("Expected nil for nil model")
		}
	})
}

func TestConvertVersionToAIR(t *testing.T) {
	version := &ModelVersion{
		ID:      43533,
		ModelID: 2421,
		Name:    "Test Version",
		Files: []File{
			{Name: "model.safetensors"},
		},
	}

	t.Run("Convert version", func(t *testing.T) {
		air := ConvertVersionToAIR(version, "sdxl")

		if air.Ecosystem != "sdxl" {
			t.Errorf("Expected ecosystem sdxl, got %s", air.Ecosystem)
		}
		if air.ID != "2421" {
			t.Errorf("Expected model ID 2421, got %s", air.ID)
		}
		if air.Version != "43533" {
			t.Errorf("Expected version 43533, got %s", air.Version)
		}
		if air.Format != "safetensors" {
			t.Errorf("Expected format safetensors, got %s", air.Format)
		}
	})

	t.Run("Convert LoRA version", func(t *testing.T) {
		loraVersion := &ModelVersion{
			ID:      43533,
			ModelID: 2421,
			Files: []File{
				{Name: "lora_model.safetensors"},
			},
		}

		air := ConvertVersionToAIR(loraVersion, "sdxl")

		if air.Type != "lora" {
			t.Errorf("Expected type lora, got %s", air.Type)
		}
	})

	t.Run("Convert nil version", func(t *testing.T) {
		air := ConvertVersionToAIR(nil, "sdxl")
		if air != nil {
			t.Error("Expected nil for nil version")
		}
	})
}
