package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/masterfabric-go/masterfabric/internal/application/urbantransform/dto"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/event"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/model"
	"github.com/masterfabric-go/masterfabric/internal/domain/urbantransform/repository"
	domainErr "github.com/masterfabric-go/masterfabric/internal/shared/errors"
	"github.com/masterfabric-go/masterfabric/internal/shared/events"
)

// ContractorCommandHandler handles write operations for contractors (CQRS write side).
type ContractorCommandHandler struct {
	repo     repository.ContractorRepository
	eventBus events.EventBus
}

// NewContractorCommandHandler creates a new ContractorCommandHandler.
func NewContractorCommandHandler(repo repository.ContractorRepository, eventBus events.EventBus) *ContractorCommandHandler {
	return &ContractorCommandHandler{repo: repo, eventBus: eventBus}
}

// Create creates a new contractor within the tenant scope.
func (h *ContractorCommandHandler) Create(ctx context.Context, orgID, appID uuid.UUID, req dto.CreateContractorRequest) (*dto.ContractorResponse, error) {
	if existing, _ := h.repo.GetByTaxNumber(ctx, orgID, req.TaxNumber); existing != nil {
		return nil, domainErr.New(domainErr.ErrAlreadyExists, "contractor tax number already exists", nil)
	}

	contractor := &model.Contractor{
		OrganizationID:   orgID,
		AppID:            appID,
		UserID:           req.UserID,
		CompanyName:      req.CompanyName,
		TaxNumber:        req.TaxNumber,
		TradeRegistryNo:  req.TradeRegistryNo,
		AuthorizedPerson: req.AuthorizedPerson,
		Phone:            req.Phone,
		Email:            req.Email,
		Address:          req.Address,
		Status:           model.ContractorStatusActive,
	}

	if err := h.repo.Create(ctx, contractor); err != nil {
		return nil, err
	}

	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypeContractor, event.ActionCreated, contractor.ID, orgID, appID))

	resp := dto.ToContractorResponse(contractor)
	return &resp, nil
}

// Update applies a partial update to a contractor.
func (h *ContractorCommandHandler) Update(ctx context.Context, orgID, appID, id uuid.UUID, req dto.UpdateContractorRequest) (*dto.ContractorResponse, error) {
	contractor, err := h.repo.GetByID(ctx, orgID, appID, id)
	if err != nil {
		return nil, err
	}

	if req.CompanyName != nil {
		contractor.CompanyName = *req.CompanyName
	}
	if req.TradeRegistryNo != nil {
		contractor.TradeRegistryNo = *req.TradeRegistryNo
	}
	if req.AuthorizedPerson != nil {
		contractor.AuthorizedPerson = *req.AuthorizedPerson
	}
	if req.Phone != nil {
		contractor.Phone = *req.Phone
	}
	if req.Email != nil {
		contractor.Email = *req.Email
	}
	if req.Address != nil {
		contractor.Address = *req.Address
	}
	if req.Status != nil {
		status := model.ContractorStatus(*req.Status)
		if !model.IsValidContractorStatus(status) {
			return nil, domainErr.New(domainErr.ErrValidation, "invalid contractor status", nil)
		}
		contractor.Status = status
	}

	if err := h.repo.Update(ctx, contractor); err != nil {
		return nil, err
	}

	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypeContractor, event.ActionUpdated, contractor.ID, orgID, appID))

	resp := dto.ToContractorResponse(contractor)
	return &resp, nil
}

// Delete removes a contractor within the tenant scope.
func (h *ContractorCommandHandler) Delete(ctx context.Context, orgID, appID, id uuid.UUID) error {
	if _, err := h.repo.GetByID(ctx, orgID, appID, id); err != nil {
		return err
	}
	if err := h.repo.Delete(ctx, orgID, appID, id); err != nil {
		return err
	}
	_ = h.eventBus.Publish(ctx, events.TopicUrbanTransform,
		event.NewEntityEvent(event.EntityTypeContractor, event.ActionDeleted, id, orgID, appID))
	return nil
}
