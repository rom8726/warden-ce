package dto

import (
	"time"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

// DomainLevelToAPI converts domain.IssueLevel to generatedapi.IssueLevel.
func DomainLevelToAPI(level domain.IssueLevel) generatedapi.IssueLevel {
	return generatedapi.IssueLevel(level)
}

// MakeIssuesListFilter converts generatedapi.ListIssuesParams to domain.ListIssuesFilter.
func MakeIssuesListFilter(params generatedapi.ListIssuesParams) domain.ListIssuesFilter {
	filter := domain.ListIssuesFilter{
		PerPage: params.PerPage,
		PageNum: params.Page,
		// Default sorting
		OrderBy:  domain.OrderByFieldTotalEvents,
		OrderAsc: false, // DESC by default
	}

	if params.ProjectID.Set {
		projectID := domain.ProjectID(params.ProjectID.Value)
		filter.ProjectID = &projectID
	}
	if params.Level.Set {
		level := domain.IssueLevel(params.Level.Value)
		filter.Level = &level
	}
	if params.Status.Set {
		status := domain.IssueStatus(params.Status.Value)
		filter.Status = &status
	}

	// Handle sort_by parameter
	if params.SortBy.Set {
		switch params.SortBy.Value {
		case generatedapi.IssueSortColumnTotalEvents:
			filter.OrderBy = domain.OrderByFieldTotalEvents
		case generatedapi.IssueSortColumnFirstSeen:
			filter.OrderBy = domain.OrderByFieldFirstSeen
		case generatedapi.IssueSortColumnLastSeen:
			filter.OrderBy = domain.OrderByFieldLastSeen
		}
	}

	// Handle sort_order parameter
	if params.SortOrder.Set {
		filter.OrderAsc = params.SortOrder.Value == generatedapi.SortOrderAsc
	}

	return filter
}

// DomainIssueToAPI converts domain.Issue to generatedapi.Issue.
func DomainIssueToAPI(
	issue domain.Issue,
	projectName string,
	resolvedAt *time.Time,
	resolvedBy *domain.UserID,
	resolvedByUsername *string,
) generatedapi.Issue {
	var resolvedAtOpt generatedapi.OptDateTime
	if resolvedAt != nil {
		resolvedAtOpt.Value = *resolvedAt
		resolvedAtOpt.Set = true
	}

	var resolvedByOpt generatedapi.OptString
	if resolvedBy != nil && resolvedByUsername != nil {
		resolvedByOpt.Value = *resolvedByUsername
		resolvedByOpt.Set = true
	}

	return generatedapi.Issue{
		ID:          uint(issue.ID),
		ProjectID:   issue.ProjectID.Uint(),
		Source:      generatedapi.IssueSource(issue.Source),
		Status:      generatedapi.IssueStatus(issue.Status),
		ProjectName: projectName,
		Title:       issue.Title,
		Message:     issue.Title, // TODO: This seems to be a duplicate of Title
		Level:       DomainLevelToAPI(issue.Level),
		Platform:    issue.Platform,
		Count:       issue.TotalEvents,
		FirstSeen:   issue.FirstSeen,
		LastSeen:    issue.LastSeen,
		ResolvedAt:  resolvedAtOpt,
		ResolvedBy:  resolvedByOpt,
	}
}

// MakeIssueResponseWithEvent converts domain.IssueExtendedWithChildren to generatedapi.IssueResponse.
func MakeIssueResponseWithEvent(issue domain.IssueExtendedWithChildren) generatedapi.IssueResponse {
	events := make([]generatedapi.IssueEvent, len(issue.Events))
	for i := range issue.Events {
		events[i] = DomainIssueEventToAPI(issue.Events[i])
	}

	return generatedapi.IssueResponse{
		Source: generatedapi.IssueSource(issue.Source),
		Issue: DomainIssueToAPI(
			issue.Issue,
			issue.ProjectName,
			issue.ResolvedAt,
			issue.ResolvedBy,
			issue.ResolvedByUsername,
		),
		Events: events,
	}
}
