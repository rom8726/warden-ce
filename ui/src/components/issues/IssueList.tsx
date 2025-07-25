import React from 'react';
import { Box, Grid, Pagination, useMediaQuery, useTheme } from '@mui/material';
import LoadingSpinner from '../LoadingSpinner';
import EmptyState from '../EmptyState';
import IssueCard from './IssueCard';
import type { Issue } from '../../generated/api/client';

// Interface for our issue data with additional fields
interface ExtendedIssue extends Issue {
  projectName?: string;
}

interface IssueListProps {
  issues: ExtendedIssue[];
  isLoading: boolean;
  error: string | null;
  viewMode: 'compact' | 'medium' | 'large';
  onIssueClick: (projectId: number, issueId: number) => void;
  onRefresh: () => void;
  // Pagination props
  totalCount?: number;
  page?: number;
  perPage?: number;
  onPageChange?: (event: React.ChangeEvent<unknown>, value: number) => void;
}

const IssueList: React.FC<IssueListProps> = ({
  issues,
  isLoading,
  error,
  viewMode,
  onIssueClick,
  onRefresh,
  totalCount,
  page,
  perPage,
  onPageChange,
}) => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  if (isLoading) {
    return <LoadingSpinner message="Loading issues..." fullHeight />;
  }

  if (error) {
    return (
      <EmptyState
        title="Error Loading Issues"
        description={error}
        type="issues"
        action={{
          label: 'Try Again',
          onClick: onRefresh,
          variant: 'contained',
        }}
      />
    );
  }

  if (issues.length === 0) {
    return (
      <EmptyState
        title="No Issues Found"
        description="No issues found matching the selected filters. Try adjusting your search criteria or filters."
        type="issues"
      />
    );
  }

  return (
    <Box sx={{ px: { xs: 1, sm: 2 }, pb: 2, width: '100%' }}>
      <Grid container spacing={1.5} sx={{ width: '100%' }}>
        {issues.map((issue) => (
          <Grid item xs={12} key={`${issue.project_id}-${issue.id}`}>
            <IssueCard 
              issue={issue} 
              viewMode={viewMode}
              onIssueClick={onIssueClick}
            />
          </Grid>
        ))}
      </Grid>

      {/* Pagination */}
      {totalCount && page && perPage && onPageChange && totalCount > perPage && (
        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 3 }}>
          <Pagination
            count={Math.ceil(totalCount / perPage)}
            page={page}
            onChange={onPageChange}
            color="primary"
            size={isMobile ? 'small' : 'medium'}
            showFirstButton
            showLastButton
          />
        </Box>
      )}
    </Box>
  );
};

export default IssueList; 