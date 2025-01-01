// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/pkg/util"

	"github.com/samber/lo"
)

var (
	ErrPromptNotFound                   = errors.New("prompt not found")
	ErrPromptVersionNotFound            = errors.New("prompt version not found")
	ErrFailedCreatePromptVersion        = errors.New("failed create prompt version")
	ErrFailedListPrompts                = errors.New("failed list prompts")
	ErrFailedListPromptSummary          = errors.New("failed list prompt summaries")
	ErrFailedListPromptVersions         = errors.New("failed list prompt versions")
	ErrFailedGetPrompt                  = errors.New("failed get prompt")
	ErrFailedGetPromptVersion           = errors.New("failed get prompt version")
	ErrFailedUpdatePromptVersion        = errors.New("failed update prompt version")
	ErrFailedSetPromptVersionProduction = errors.New("failed set prompt version production")
	ErrFailedDeletePromptVersion        = errors.New("failed delete prompt version")
	ErrFailedDeletePrompt               = errors.New("failed delete prompt")
	ErrPromptHandleTaken                = errors.New("prompt handle taken")
	ErrInvalidPromptHandle              = errors.New("invalid prompt handle")
	ErrInvalidPromptVersion             = errors.New("invalid prompt version")

	PromptLabelLatest  = "latest"
	promptHandleMaxLen = 100
)

// Prompt is the module for workspace prompt CRUD.
type Prompt struct {
	PromptRepository    db.PromptRepository
	WorkspaceRepository db.WorkspaceRepository
}

// NewPrompt creates a prompt module with the given repositories.
func NewPrompt(pr db.PromptRepository, wr db.WorkspaceRepository) *Prompt {
	return &Prompt{
		PromptRepository:    pr,
		WorkspaceRepository: wr,
	}
}

// CreatePromptRequest is the body for creating a prompt with its first version.
type CreatePromptRequest struct {
	Name          string  `json:"name" validate:"required,min=1,max=150" label:"Name"`
	Description   string  `json:"description" validate:"omitempty,max=2000" label:"Description"`
	Type          string  `json:"type" validate:"required,oneof=text chat" label:"Type"`
	Content       string  `json:"content" validate:"required" label:"Content"`
	Config        *string `json:"config" validate:"omitempty,json" label:"Config"`
	Labels        *string `json:"labels" validate:"omitempty,json" label:"Labels"`
	CommitMessage *string `json:"commitMessage" validate:"omitempty,max=255" label:"Commit message"`
	Meta          *string `json:"meta" validate:"omitempty,json" label:"Meta"`
	Production    bool    `json:"production" label:"Production"`
}

// CreatePromptVersionRequest is the body for creating a version on an existing prompt.
type CreatePromptVersionRequest struct {
	Type          string  `json:"type" validate:"required,oneof=text chat" label:"Type"`
	Content       string  `json:"content" validate:"required" label:"Content"`
	Config        *string `json:"config" validate:"omitempty,json" label:"Config"`
	Labels        *string `json:"labels" validate:"omitempty,json" label:"Labels"`
	CommitMessage *string `json:"commitMessage" validate:"omitempty,max=255" label:"Commit message"`
	Meta          *string `json:"meta" validate:"omitempty,json" label:"Meta"`
	Production    bool    `json:"production" label:"Production"`
}

// UpdatePromptVersionRequest is the body for updating a prompt version.
type UpdatePromptVersionRequest struct {
	Type          string  `json:"type" validate:"omitempty,oneof=text chat" label:"Type"`
	Content       string  `json:"content" validate:"required" label:"Content"`
	Config        *string `json:"config" validate:"omitempty,json" label:"Config"`
	Labels        *string `json:"labels" validate:"omitempty,json" label:"Labels"`
	CommitMessage *string `json:"commitMessage" validate:"omitempty,max=255" label:"Commit message"`
	Meta          *string `json:"meta" validate:"omitempty,json" label:"Meta"`
	Status        string  `json:"status" validate:"omitempty,oneof=active archived" label:"Status"`
	Production    bool    `json:"production" label:"Production"`
}

// PromptVersionResponse is a prompt version shaped for API responses.
type PromptVersionResponse struct {
	PromptId      db.Id           `json:"promptId"`
	Id            db.Id           `json:"id"`
	WorkspaceId   db.Id           `json:"workspaceId"`
	Name          string          `json:"name"`
	Handle        string          `json:"handle"`
	Description   string          `json:"description"`
	Version       int             `json:"version"`
	Type          string          `json:"type"`
	Content       string          `json:"content"`
	Config        json.RawMessage `json:"config,omitempty"`
	Labels        []string        `json:"labels,omitempty"`
	CommitMessage *string         `json:"commitMessage,omitempty"`
	CommitHash    string          `json:"commitHash"`
	Meta          *string         `json:"meta,omitempty"`
	Status        string          `json:"status"`
	Production    bool            `json:"production"`
	CreatedAt     string          `json:"createdAt"`
	UpdatedAt     string          `json:"updatedAt"`
}

// ListPromptsResponse is returned when listing prompts.
type ListPromptsResponse struct {
	Prompts []*PromptSummaryResponse
	Total   int64
}

// PromptSummaryResponse is a prompt aggregate for list views.
type PromptSummaryResponse struct {
	PromptId        string   `json:"promptId"`
	Name            string   `json:"name"`
	Handle          string   `json:"handle"`
	Description     string   `json:"description"`
	VersionCount    int64    `json:"versionCount"`
	Type            string   `json:"type"`
	LatestVersion   int      `json:"latestVersion"`
	LatestCreatedAt string   `json:"latestCreatedAt"`
	Labels          []string `json:"labels,omitempty"`
	HasProduction   bool     `json:"hasProduction"`
}

// ListPromptVersionsResponse is returned when listing prompt versions.
type ListPromptVersionsResponse struct {
	Versions []*PromptVersionResponse
	Total    int64
}

// CreatePrompt creates a new prompt with its first version.
func (p *Prompt) CreatePrompt(workspaceId db.Id, req *CreatePromptRequest) (*PromptVersionResponse, error) {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	handle := util.HandleFromName(req.Name, promptHandleMaxLen)
	if lo.IsEmpty(handle) {
		return nil, ErrInvalidPromptHandle
	}
	exists, err := p.PromptRepository.ExistsByHandle(workspaceId, handle)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreatePromptVersion, err)
	}
	if exists {
		return nil, ErrPromptHandleTaken
	}

	promptId, err := p.PromptRepository.CreatePrompt(workspaceId, req.Name, handle, req.Description)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreatePromptVersion, err)
	}

	version, err := p.CreatePromptVersion(workspaceId, promptId, &CreatePromptVersionRequest{
		Type:          req.Type,
		Content:       req.Content,
		Config:        req.Config,
		Labels:        req.Labels,
		CommitMessage: req.CommitMessage,
		Meta:          req.Meta,
		Production:    req.Production,
	})
	if err != nil {
		_ = p.PromptRepository.DeletePrompt(workspaceId, promptId)
		return nil, err
	}

	return version, nil
}

// CreatePromptVersion creates a new version for an existing prompt.
func (p *Prompt) CreatePromptVersion(workspaceId, promptId db.Id, req *CreatePromptVersionRequest) (*PromptVersionResponse, error) {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	row, err := p.PromptRepository.GetPrompt(workspaceId, promptId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreatePromptVersion, err)
	}
	if row == nil {
		return nil, ErrPromptNotFound
	}

	var existing []string
	if req.Labels != nil && lo.IsNotEmpty(*req.Labels) {
		if err := json.Unmarshal([]byte(*req.Labels), &existing); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedCreatePromptVersion, err)
		}
	}
	latestLabels := lo.Uniq(append(existing, PromptLabelLatest))
	labelsRaw, err := json.Marshal(latestLabels)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreatePromptVersion, err)
	}
	labels := string(labelsRaw)

	err = p.PromptRepository.ClearLabelFromPromptId(promptId, PromptLabelLatest, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreatePromptVersion, err)
	}

	version := &db.PromptVersion{
		PromptId:      promptId,
		WorkspaceId:   workspaceId,
		Name:          row.Name,
		Handle:        row.Handle,
		Description:   row.Description,
		Type:          req.Type,
		Content:       req.Content,
		Config:        req.Config,
		Labels:        &labels,
		CommitMessage: req.CommitMessage,
		Meta:          req.Meta,
		Status:        "active",
		Production:    req.Production,
	}

	err = p.PromptRepository.Create(version)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreatePromptVersion, err)
	}

	return &PromptVersionResponse{
		PromptId:      version.PromptId,
		Id:            version.Id,
		WorkspaceId:   version.WorkspaceId,
		Name:          version.Name,
		Handle:        version.Handle,
		Description:   version.Description,
		Version:       version.Version,
		Type:          version.Type,
		Content:       version.Content,
		Config:        util.JSONRawFromString(version.Config),
		Labels:        util.JSONSliceFromString[string](version.Labels),
		CommitMessage: version.CommitMessage,
		CommitHash:    version.CommitHash,
		Meta:          version.Meta,
		Status:        version.Status,
		Production:    version.Production,
		CreatedAt:     version.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     version.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// GetPromptVersionByHandle returns a prompt version resolved by handle and version selector.
func (p *Prompt) GetPromptVersionByHandle(workspaceId db.Id, handle, versionSpec string) (*PromptVersionResponse, error) {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}
	if lo.IsEmpty(handle) {
		return nil, ErrInvalidPromptHandle
	}

	versionSpec = strings.TrimSpace(strings.ToLower(versionSpec))

	var version *db.PromptVersion
	switch versionSpec {
	case "", PromptLabelLatest:
		version, err = p.PromptRepository.GetLatestByHandle(workspaceId, handle)
	case "production":
		version, err = p.PromptRepository.GetByHandleAndLabel(workspaceId, handle, "production")
	default:
		vn, err := strconv.Atoi(strings.TrimPrefix(versionSpec, "v"))
		if err != nil || vn < 1 {
			return nil, ErrInvalidPromptVersion
		}
		version, err = p.PromptRepository.GetByHandleAndVersion(workspaceId, handle, vn)
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetPromptVersion, err)
	}
	if version == nil {
		return nil, ErrPromptVersionNotFound
	}

	return &PromptVersionResponse{
		PromptId:      version.PromptId,
		Id:            version.Id,
		WorkspaceId:   version.WorkspaceId,
		Name:          version.Name,
		Handle:        version.Handle,
		Description:   version.Description,
		Version:       version.Version,
		Type:          version.Type,
		Content:       version.Content,
		Config:        util.JSONRawFromString(version.Config),
		Labels:        util.JSONSliceFromString[string](version.Labels),
		CommitMessage: version.CommitMessage,
		CommitHash:    version.CommitHash,
		Meta:          version.Meta,
		Status:        version.Status,
		Production:    version.Production,
		CreatedAt:     version.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     version.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// GetPromptVersion returns one prompt version by Id within a prompt.
func (p *Prompt) GetPromptVersion(workspaceId, promptId, versionId db.Id) (*PromptVersionResponse, error) {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	prompt, err := p.PromptRepository.GetById(versionId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetPromptVersion, err)
	}
	if prompt == nil || prompt.WorkspaceId != workspaceId || prompt.PromptId != promptId {
		return nil, ErrPromptVersionNotFound
	}

	return &PromptVersionResponse{
		PromptId:      prompt.PromptId,
		Id:            prompt.Id,
		WorkspaceId:   prompt.WorkspaceId,
		Name:          prompt.Name,
		Handle:        prompt.Handle,
		Description:   prompt.Description,
		Version:       prompt.Version,
		Type:          prompt.Type,
		Content:       prompt.Content,
		Config:        util.JSONRawFromString(prompt.Config),
		Labels:        util.JSONSliceFromString[string](prompt.Labels),
		CommitMessage: prompt.CommitMessage,
		CommitHash:    prompt.CommitHash,
		Meta:          prompt.Meta,
		Status:        prompt.Status,
		Production:    prompt.Production,
		CreatedAt:     prompt.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     prompt.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// UpdatePromptVersion updates a prompt version.
func (p *Prompt) UpdatePromptVersion(workspaceId, promptId, versionId db.Id, req *UpdatePromptVersionRequest) (*PromptVersionResponse, error) {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}
	prompt, err := p.PromptRepository.GetById(versionId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedUpdatePromptVersion, err)
	}
	if prompt == nil || prompt.WorkspaceId != workspaceId || prompt.PromptId != promptId {
		return nil, ErrPromptVersionNotFound
	}

	if lo.IsNotEmpty(req.Type) {
		prompt.Type = req.Type
	}
	prompt.Content = req.Content
	prompt.Config = req.Config
	prompt.Labels = req.Labels
	prompt.CommitMessage = req.CommitMessage
	prompt.Meta = req.Meta
	if lo.IsNotEmpty(req.Status) {
		prompt.Status = req.Status
	}
	if req.Production {
		err = p.PromptRepository.ClearProductionFromPromptId(prompt.PromptId, &versionId)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedSetPromptVersionProduction, err)
		}
		prompt.Production = true
	}

	err = p.PromptRepository.Update(prompt)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedUpdatePromptVersion, err)
	}

	return &PromptVersionResponse{
		PromptId:      prompt.PromptId,
		Id:            prompt.Id,
		WorkspaceId:   prompt.WorkspaceId,
		Name:          prompt.Name,
		Handle:        prompt.Handle,
		Description:   prompt.Description,
		Version:       prompt.Version,
		Type:          prompt.Type,
		Content:       prompt.Content,
		Config:        util.JSONRawFromString(prompt.Config),
		Labels:        util.JSONSliceFromString[string](prompt.Labels),
		CommitMessage: prompt.CommitMessage,
		CommitHash:    prompt.CommitHash,
		Meta:          prompt.Meta,
		Status:        prompt.Status,
		Production:    prompt.Production,
		CreatedAt:     prompt.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     prompt.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// DeletePromptVersion deletes a prompt version by Id.
func (p *Prompt) DeletePromptVersion(workspaceId, promptId, versionId db.Id) error {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return err
	}
	if workspace == nil {
		return ErrWorkspaceNotFound
	}

	prompt, err := p.PromptRepository.GetById(versionId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeletePromptVersion, err)
	}
	if prompt == nil || prompt.WorkspaceId != workspaceId || prompt.PromptId != promptId {
		return ErrPromptVersionNotFound
	}

	err = p.PromptRepository.DeleteVersion(versionId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeletePromptVersion, err)
	}

	return nil
}

// ListPrompts returns paginated prompt summaries for a workspace.
func (p *Prompt) ListPrompts(workspaceId db.Id, limit, offset int) (*ListPromptsResponse, error) {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	total, err := p.PromptRepository.CountByWorkspaceId(workspaceId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListPromptSummary, err)
	}

	summaries, err := p.PromptRepository.ListSummariesByWorkspace(workspaceId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListPromptSummary, err)
	}

	list := make([]*PromptSummaryResponse, 0, len(summaries))
	for _, item := range summaries {
		var summaryLabels []string
		if item.LatestLabels != nil && lo.IsNotEmpty(*item.LatestLabels) {
			_ = json.Unmarshal([]byte(*item.LatestLabels), &summaryLabels)
		}
		list = append(list, &PromptSummaryResponse{
			PromptId:        item.PromptId.String(),
			Name:            item.Name,
			Handle:          item.Handle,
			Description:     item.Description,
			VersionCount:    item.VersionCount,
			Type:            item.Type,
			LatestVersion:   item.LatestVersion,
			LatestCreatedAt: item.LatestCreatedAt.UTC().Format(time.RFC3339),
			Labels:          summaryLabels,
			HasProduction:   item.HasProduction,
		})
	}

	return &ListPromptsResponse{Prompts: list, Total: total}, nil
}

// GetPrompt returns one prompt summary by Id.
func (p *Prompt) GetPrompt(workspaceId, promptId db.Id) (*PromptSummaryResponse, error) {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	summary, err := p.PromptRepository.GetSummaryByPromptId(workspaceId, promptId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetPrompt, err)
	}
	if summary == nil {
		return nil, ErrPromptNotFound
	}

	var summaryLabels []string
	if summary.LatestLabels != nil && lo.IsNotEmpty(*summary.LatestLabels) {
		_ = json.Unmarshal([]byte(*summary.LatestLabels), &summaryLabels)
	}

	return &PromptSummaryResponse{
		PromptId:        summary.PromptId.String(),
		Name:            summary.Name,
		Handle:          summary.Handle,
		Description:     summary.Description,
		VersionCount:    summary.VersionCount,
		Type:            summary.Type,
		LatestVersion:   summary.LatestVersion,
		LatestCreatedAt: summary.LatestCreatedAt.UTC().Format(time.RFC3339),
		Labels:          summaryLabels,
		HasProduction:   summary.HasProduction,
	}, nil
}

// ListPromptVersions returns paginated versions for a prompt.
func (p *Prompt) ListPromptVersions(workspaceId, promptId db.Id, limit, offset int) (*ListPromptVersionsResponse, error) {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	total, err := p.PromptRepository.CountVersionsByPromptId(workspaceId, promptId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListPromptVersions, err)
	}
	if total == 0 {
		return nil, ErrPromptNotFound
	}

	versions, err := p.PromptRepository.ListVersionsByPromptId(workspaceId, promptId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListPromptVersions, err)
	}

	list := make([]*PromptVersionResponse, 0, len(versions))
	for _, item := range versions {
		list = append(list, &PromptVersionResponse{
			PromptId:      item.PromptId,
			Id:            item.Id,
			WorkspaceId:   item.WorkspaceId,
			Name:          item.Name,
			Handle:        item.Handle,
			Description:   item.Description,
			Version:       item.Version,
			Type:          item.Type,
			Content:       item.Content,
			Config:        util.JSONRawFromString(item.Config),
			Labels:        util.JSONSliceFromString[string](item.Labels),
			CommitMessage: item.CommitMessage,
			CommitHash:    item.CommitHash,
			Meta:          item.Meta,
			Status:        item.Status,
			Production:    item.Production,
			CreatedAt:     item.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:     item.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	return &ListPromptVersionsResponse{Versions: list, Total: total}, nil
}

// DeletePrompt deletes a prompt and all of its versions.
func (p *Prompt) DeletePrompt(workspaceId, promptId db.Id) error {
	workspace, err := p.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return err
	}
	if workspace == nil {
		return ErrWorkspaceNotFound
	}

	err = p.PromptRepository.DeletePrompt(workspaceId, promptId)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrPromptNotFound
		}
		return fmt.Errorf("%w: %v", ErrFailedDeletePrompt, err)
	}

	return nil
}
