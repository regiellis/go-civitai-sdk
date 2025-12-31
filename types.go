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

// Package civitai - Core Types and Data Structures
//
// This file defines all the data types used by the CivitAI API, including
// models, images, creators, tags, and various parameter structures.
//
// # Key Types
//
// Model represents an AI model with metadata:
//
//	model := civitai.Model{
//		ID:    12345,
//		Name:  "Realistic Vision",
//		Type:  civitai.ModelTypeCheckpoint,
//		Tags:  []string{"realistic", "photorealistic"},
//	}
//
// SearchParams configures model search requests:
//
//	params := civitai.SearchParams{
//		Tag:    "anime",           // Recommended over Query
//		Types:  []civitai.ModelType{civitai.ModelTypeCheckpoint},
//		Sort:   civitai.SortMostDownload,
//		Limit:  20,
//		Cursor: "cursor-string",   // For pagination
//	}
//
// ImageParams configures image browsing:
//
//	params := civitai.ImageParams{
//		Sort:     "Most Reactions",
//		NSFW:     string(civitai.NSFWLevelNone),
//		Limit:    50,
//		Username: "artist-name",
//	}
//
// # Pagination
//
// The API supports both cursor and page-based pagination. Cursor-based
// is recommended for consistency:
//
//	// Cursor-based (recommended)
//	params.Cursor = metadata.NextCursor
//
//	// Page-based (less reliable)
//	params.Page = 2
//
// # Enumerations
//
// The package provides typed constants for various options:
//
//	// Model types
//	civitai.ModelTypeCheckpoint
//	civitai.ModelTypeLora
//	civitai.ModelTypeTextualInversion
//
//	// Sort options
//	civitai.SortMostDownload
//	civitai.SortHighestRated
//	civitai.SortNewest
//
//	// NSFW levels
//	civitai.NSFWLevelNone
//	civitai.NSFWLevelSoft
//	civitai.NSFWLevelMature

// Package gocivitaisdk provides a Go SDK for the CivitAI API
package civitai

import (
	"encoding/json"
	"time"
)

// FlexibleStringSlice handles API responses that may return either a string or []string
type FlexibleStringSlice []string

// UnmarshalJSON handles both string and []string JSON values
func (f *FlexibleStringSlice) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if str != "" {
			*f = []string{str}
		} else {
			*f = []string{}
		}
		return nil
	}

	// Try to unmarshal as []string
	var slice []string
	if err := json.Unmarshal(data, &slice); err == nil {
		*f = slice
		return nil
	}

	// Default to empty slice
	*f = []string{}
	return nil
}

// MarshalJSON converts back to JSON (as array)
func (f FlexibleStringSlice) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(f))
}

// Common types and structures used across all CivitAI resources

// APIResponse represents the standard API response structure
// type APIResponse struct {
// 	Items    []interface{} `json:"items,omitempty"`
// 	Metadata *Metadata     `json:"metadata,omitempty"`
// 	Success  bool          `json:"success"`
// 	Error    *APIError     `json:"error,omitempty"`
// }

// Metadata contains pagination and response metadata
type Metadata struct {
	TotalItems  int    `json:"totalItems"`
	TotalPages  int    `json:"totalPages"`
	CurrentPage int    `json:"currentPage"`
	PageSize    int    `json:"pageSize"`
	NextCursor  string `json:"nextCursor,omitempty"`
	PrevCursor  string `json:"prevCursor,omitempty"`
	NextPage    string `json:"nextPage,omitempty"`
	PrevPage    string `json:"prevPage,omitempty"`
}

// APIError represents an API error response
// type APIError struct {
// 	Code    string `json:"code"`
// 	Message string `json:"message"`
// 	Details string `json:"details,omitempty"`
// }

// ResourceType represents the type of CivitAI resource
type ResourceType string

const (
	ResourceTypeModel      ResourceType = "Model"
	ResourceTypeCheckpoint ResourceType = "Checkpoint"
	ResourceTypeLORA       ResourceType = "LORA"
	ResourceTypeEmbedding  ResourceType = "TextualInversion"
	ResourceTypeVAE        ResourceType = "VAE"
	ResourceTypeWorkflow   ResourceType = "Workflow"
	ResourceTypeImage      ResourceType = "Image"
	ResourceTypeArticle    ResourceType = "Article"
	ResourceTypeCollection ResourceType = "Collection"
	ResourceTypePost       ResourceType = "Post"
	ResourceTypeWildcard   ResourceType = "Wildcard"
)

// ModelType represents specific model subtypes
type ModelType string

const (
	ModelTypeCheckpoint       ModelType = "Checkpoint"
	ModelTypeLORA             ModelType = "LORA"
	ModelTypeTextualInversion ModelType = "TextualInversion"
	ModelTypeEmbedding        ModelType = "TextualInversion" // Alias for TextualInversion
	ModelTypeHypernetwork     ModelType = "Hypernetwork"
	ModelTypeAestheticGrad    ModelType = "AestheticGradient"
	ModelTypeControlNet       ModelType = "ControlNet"
	ModelTypePose             ModelType = "Pose"
	ModelTypeVAE              ModelType = "VAE"
)

// BaseModel represents the base model architecture
type BaseModel string

const (
	BaseModelSD1_5 BaseModel = "SD 1.5"
	BaseModelSDXL  BaseModel = "SDXL 1.0"
	BaseModelSD2_0 BaseModel = "SD 2.0"
	BaseModelSD2_1 BaseModel = "SD 2.1"
	BaseModelOther BaseModel = "Other"
)

// SortType represents sorting options
type SortType string

const (
	SortHighestRated SortType = "Highest Rated"
	SortMostLiked    SortType = "Most Liked"
	SortMostDownload SortType = "Most Downloaded"
	SortNewest       SortType = "Newest"
	SortOldest       SortType = "Oldest"
)

// Period represents time period filters
type Period string

const (
	PeriodAllTime Period = "AllTime"
	PeriodYear    Period = "Year"
	PeriodMonth   Period = "Month"
	PeriodWeek    Period = "Week"
	PeriodDay     Period = "Day"
)

// User represents a CivitAI user
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Image    string `json:"image,omitempty"`
}

// Stats represents statistics for a resource
type Stats struct {
	DownloadCount   int     `json:"downloadCount"`
	FavoriteCount   int     `json:"favoriteCount"`
	CommentCount    int     `json:"commentCount"`
	RatingCount     int     `json:"ratingCount"`
	Rating          float64 `json:"rating"`
	ThumbsUpCount   int     `json:"thumbsUpCount"`
	ThumbsDownCount int     `json:"thumbsDownCount"`
}

// FileFormat represents supported file formats
type FileFormat string

const (
	FileFormatSafeTensors  FileFormat = "SafeTensor"
	FileFormatPickleTensor FileFormat = "PickleTensor"
	FileFormatCKPT         FileFormat = "CKPT"
	FileFormatOther        FileFormat = "Other"
)

// FileMetadata represents metadata for downloadable files
type FileMetadata struct {
	FP     string     `json:"fp,omitempty"`
	Size   string     `json:"size,omitempty"`
	Format FileFormat `json:"format,omitempty"`
}

// File represents a downloadable file
type File struct {
	ID                int          `json:"id"`
	URL               string       `json:"url"`
	SizeKB            float64      `json:"sizeKB"`
	Name              string       `json:"name"`
	Type              string       `json:"type"`
	Metadata          FileMetadata `json:"metadata"`
	PickleScanResult  string       `json:"pickleScanResult,omitempty"`
	PickleScanMessage string       `json:"pickleScanMessage,omitempty"`
	VirusScanResult   string       `json:"virusScanResult,omitempty"`
	VirusScanMessage  string       `json:"virusScanMessage,omitempty"`
	ScannedAt         *time.Time   `json:"scannedAt,omitempty"`
	Hashes            Hashes       `json:"hashes,omitempty"`
	Primary           bool         `json:"primary,omitempty"`
}

// Hashes represents file hash checksums
type Hashes struct {
	AutoV1 string `json:"AutoV1,omitempty"`
	AutoV2 string `json:"AutoV2,omitempty"`
	SHA256 string `json:"SHA256,omitempty"`
	CRC32  string `json:"CRC32,omitempty"`
	BLAKE3 string `json:"BLAKE3,omitempty"`
}

// Image represents an image associated with a resource
type Image struct {
	ID           int                    `json:"id"`
	URL          string                 `json:"url"`
	NSFW         string                 `json:"nsfw,omitempty"`
	Width        int                    `json:"width,omitempty"`
	Height       int                    `json:"height,omitempty"`
	Hash         string                 `json:"hash,omitempty"`
	Type         string                 `json:"type,omitempty"`
	Metadata     map[string]interface{} `json:"meta,omitempty"`
	Availability string                 `json:"availability,omitempty"`
}

// Tag represents a tag associated with a resource
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type,omitempty"`
}

// SearchParams represents common search parameters
type SearchParams struct {
	Query                 string      `json:"query,omitempty"`
	Types                 []ModelType `json:"types,omitempty"`
	Sort                  SortType    `json:"sort,omitempty"`
	Period                Period      `json:"period,omitempty"`
	Rating                int         `json:"rating,omitempty"`
	Page                  int         `json:"page,omitempty"`
	Limit                 int         `json:"limit,omitempty"`
	Cursor                string      `json:"cursor,omitempty"` // Added cursor support for pagination
	Tag                   string      `json:"tag,omitempty"`
	Username              string      `json:"username,omitempty"`
	Favorites             bool        `json:"favorites,omitempty"`
	Hidden                bool        `json:"hidden,omitempty"`
	PrimaryFileOnly       bool        `json:"primaryFileOnly,omitempty"`
	AllowNoCredit         bool        `json:"allowNoCredit,omitempty"`
	AllowDerivatives      bool        `json:"allowDerivatives,omitempty"`
	AllowDifferentLicense bool        `json:"allowDifferentLicense,omitempty"`
	AllowCommercialUse    []string    `json:"allowCommercialUse,omitempty"`
	NSFW                  *bool       `json:"nsfw,omitempty"`
	SupportsGeneration    *bool       `json:"supportsGeneration,omitempty"`
}

// ModelVersion represents a version of a model
type ModelVersion struct {
	ID                   int        `json:"id"`
	ModelID              int        `json:"modelId,omitempty"`
	Name                 string     `json:"name"`
	Description          string     `json:"description,omitempty"`
	BaseModel            BaseModel  `json:"baseModel,omitempty"`
	BaseModelType        string     `json:"baseModelType,omitempty"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            time.Time  `json:"updatedAt"`
	PublishedAt          *time.Time `json:"publishedAt,omitempty"`
	TrainedWords         []string   `json:"trainedWords,omitempty"`
	Files                []File     `json:"files,omitempty"`
	Images               []Image    `json:"images,omitempty"`
	DownloadURL          string     `json:"downloadUrl,omitempty"`
	EarlyAccessTimeFrame int        `json:"earlyAccessTimeFrame,omitempty"`
	Stats                Stats      `json:"stats,omitempty"`
	Availability         string     `json:"availability,omitempty"`
}

// ToAIR converts the model version to an AIR identifier
func (mv *ModelVersion) ToAIR(ecosystem string) *AIR {
	return ConvertVersionToAIR(mv, ecosystem)
}

// GetAIRForEcosystem returns an AIR for a specific ecosystem
func (mv *ModelVersion) GetAIRForEcosystem(ecosystem AIREcosystem) *AIR {
	return mv.ToAIR(string(ecosystem))
}

// Model represents a CivitAI model
type Model struct {
	ID                    int                 `json:"id"`
	Name                  string              `json:"name"`
	Description           string              `json:"description,omitempty"`
	Type                  ModelType           `json:"type"`
	POI                   bool                `json:"poi,omitempty"`
	NSFW                  bool                `json:"nsfw,omitempty"`
	AllowNoCredit         bool                `json:"allowNoCredit,omitempty"`
	AllowCommercialUse    FlexibleStringSlice `json:"allowCommercialUse,omitempty"`
	AllowDerivatives      bool                `json:"allowDerivatives,omitempty"`
	AllowDifferentLicense bool                `json:"allowDifferentLicense,omitempty"`
	Stats                 Stats               `json:"stats,omitempty"`
	Creator               User                `json:"creator,omitempty"`
	Tags                  []string            `json:"tags,omitempty"`
	ModelVersions         []ModelVersion      `json:"modelVersions,omitempty"`
	Images                []Image             `json:"images,omitempty"`
	CreatedAt             time.Time           `json:"createdAt"`
	UpdatedAt             time.Time           `json:"updatedAt"`
	PublishedAt           *time.Time          `json:"publishedAt,omitempty"`
}

// ToAIR converts the model to an AIR identifier
func (m *Model) ToAIR(ecosystem string, versionID ...int) *AIR {
	return ConvertModelToAIR(m, ecosystem, versionID...)
}

// GetAIRForEcosystem returns an AIR for a specific ecosystem
func (m *Model) GetAIRForEcosystem(ecosystem AIREcosystem, versionID ...int) *AIR {
	return m.ToAIR(string(ecosystem), versionID...)
}

// Article represents a CivitAI article
type Article struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content,omitempty"`
	CoverImage  Image     `json:"coverImage,omitempty"`
	PublishedAt time.Time `json:"publishedAt"`
	User        User      `json:"user"`
	Stats       Stats     `json:"stats,omitempty"`
	Tags        []Tag     `json:"tags,omitempty"`
}

// Collection represents a CivitAI collection
type Collection struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Image       Image   `json:"image,omitempty"`
	User        User    `json:"user"`
	Tags        []Tag   `json:"tags,omitempty"`
	Stats       Stats   `json:"stats,omitempty"`
	Items       []Model `json:"items,omitempty"`
}

// Post represents a CivitAI post/image post
type Post struct {
	ID            int            `json:"id"`
	Title         string         `json:"title,omitempty"`
	URL           string         `json:"url,omitempty"`
	Images        []Image        `json:"images,omitempty"`
	User          User           `json:"user"`
	Stats         Stats          `json:"stats,omitempty"`
	Tags          []Tag          `json:"tags,omitempty"`
	NSFW          bool           `json:"nsfw,omitempty"`
	ModelVersions []ModelVersion `json:"modelVersions,omitempty"`
	PublishedAt   time.Time      `json:"publishedAt"`
}

// DetailedImage represents a detailed image with generation info
type DetailedImage struct {
	Image
	GenerationProcess string                   `json:"generationProcess,omitempty"`
	Prompt            string                   `json:"prompt,omitempty"`
	NegativePrompt    string                   `json:"negativePrompt,omitempty"`
	Steps             int                      `json:"steps,omitempty"`
	Sampler           string                   `json:"sampler,omitempty"`
	CFGScale          float64                  `json:"cfgScale,omitempty"`
	Seed              int64                    `json:"seed,omitempty"`
	Size              string                   `json:"size,omitempty"`
	Model             string                   `json:"model,omitempty"`
	ModelHash         string                   `json:"modelHash,omitempty"`
	Resources         []map[string]interface{} `json:"resources,omitempty"`
	Techniques        []string                 `json:"techniques,omitempty"`
	Tools             []string                 `json:"tools,omitempty"`
}

// Workflow represents a ComfyUI or A1111 workflow
type Workflow struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Type        string                 `json:"type"` // ComfyUI, A1111, etc.
	Definition  map[string]interface{} `json:"definition,omitempty"`
	Nodes       []WorkflowNode         `json:"nodes,omitempty"`
	User        User                   `json:"user"`
	Images      []Image                `json:"images,omitempty"`
	Tags        []Tag                  `json:"tags,omitempty"`
	Stats       Stats                  `json:"stats,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// WorkflowNode represents a node in a workflow
type WorkflowNode struct {
	ID     string                 `json:"id"`
	Type   string                 `json:"type"`
	Data   map[string]interface{} `json:"data,omitempty"`
	Inputs map[string]interface{} `json:"inputs,omitempty"`
}

// VAE represents a Variational Auto-Encoder
type VAE struct {
	Model                    // Inherits from Model
	Architecture string      `json:"architecture,omitempty"`
	TargetModels []BaseModel `json:"targetModels,omitempty"`
}

// Wildcard represents a text file for prompt automation
type Wildcard struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Content     string    `json:"content"`
	Category    string    `json:"category,omitempty"`
	Description string    `json:"description,omitempty"`
	User        User      `json:"user"`
	Tags        []Tag     `json:"tags,omitempty"`
	Stats       Stats     `json:"stats,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Creator represents a CivitAI creator/user from the /creators endpoint
type Creator struct {
	Username   string `json:"username"`
	ModelCount int    `json:"modelCount"`
	Link       string `json:"link"`
}

// ImageParams represents parameters for searching images
type ImageParams struct {
	Limit          int    `json:"limit,omitempty"`
	PostID         int    `json:"postId,omitempty"`
	ModelID        int    `json:"modelId,omitempty"`
	ModelVersionID int    `json:"modelVersionId,omitempty"`
	Username       string `json:"username,omitempty"`
	NSFW           string `json:"nsfw,omitempty"` // None, Soft, Mature, X
	Sort           string `json:"sort,omitempty"` // Most Reactions, Most Comments, Newest
	Period         Period `json:"period,omitempty"`
	Page           int    `json:"page,omitempty"`
}

// CreatorParams represents parameters for searching creators
type CreatorParams struct {
	Limit int    `json:"limit,omitempty"`
	Page  int    `json:"page,omitempty"`
	Query string `json:"query,omitempty"`
}

// TagParams represents parameters for searching tags
type TagParams struct {
	Limit int    `json:"limit,omitempty"`
	Page  int    `json:"page,omitempty"`
	Query string `json:"query,omitempty"`
}

// ImageStats represents statistics for an image
type ImageStats struct {
	CryCount     int `json:"cryCount"`
	LaughCount   int `json:"laughCount"`
	LikeCount    int `json:"likeCount"`
	HeartCount   int `json:"heartCount"`
	CommentCount int `json:"commentCount"`
}

// DetailedImageResponse represents a complete image response from the API
type DetailedImageResponse struct {
	ID        int                    `json:"id"`
	URL       string                 `json:"url"`
	Hash      string                 `json:"hash"`
	Width     int                    `json:"width"`
	Height    int                    `json:"height"`
	NSFW      bool                   `json:"nsfw"`
	NSFWLevel string                 `json:"nsfwLevel"` // None, Soft, Mature, X
	CreatedAt time.Time              `json:"createdAt"`
	PostID    int                    `json:"postId"`
	Stats     ImageStats             `json:"stats"`
	Meta      map[string]interface{} `json:"meta,omitempty"`
	Username  string                 `json:"username"`
}

// ModelVersionByHashResponse represents a model version response when searched by hash
type ModelVersionByHashResponse struct {
	ModelVersion
	Model struct {
		Name string    `json:"name"`
		Type ModelType `json:"type"`
		NSFW bool      `json:"nsfw"`
		POI  bool      `json:"poi"`
		Mode string    `json:"mode,omitempty"` // Archived, TakenDown
	} `json:"model"`
	ModelID int `json:"modelId"`
}

// TagResponse represents a tag from the /tags endpoint
type TagResponse struct {
	Name       string `json:"name"`
	ModelCount int    `json:"modelCount"`
	Link       string `json:"link"`
}

// NSFWLevel represents NSFW content levels
type NSFWLevel string

const (
	NSFWLevelNone   NSFWLevel = "None"
	NSFWLevelSoft   NSFWLevel = "Soft"
	NSFWLevelMature NSFWLevel = "Mature"
	NSFWLevelX      NSFWLevel = "X"
)

// ImageSort represents image sorting options
type ImageSort string

const (
	ImageSortMostReactions ImageSort = "Most Reactions"
	ImageSortMostComments  ImageSort = "Most Comments"
	ImageSortNewest        ImageSort = "Newest"
)

// CommercialUse represents commercial use permissions
type CommercialUse string

const (
	CommercialUseNone  CommercialUse = "None"
	CommercialUseImage CommercialUse = "Image"
	CommercialUseRent  CommercialUse = "Rent"
	CommercialUseSell  CommercialUse = "Sell"
)
