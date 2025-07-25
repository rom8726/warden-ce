import React from 'react';
import { Box, Tabs, Tab, Badge, Divider } from '@mui/material';
import IssueList from './IssueList';
import { getLevelHexColor } from '../../utils/issues/issueUtils';

interface IssuesTabsListProps {
  issues: any[];
  filteredIssues: any[];
  isLoading: boolean;
  error: string | null;
  viewMode: 'compact' | 'medium' | 'large';
  onIssueClick: (projectId: number, issueId: number) => void;
  onRefresh: () => void;
  totalCount: number;
  page: number;
  perPage: number;
  onPageChange: (event: React.ChangeEvent<unknown>, value: number) => void;
  tabValue: number;
  onTabChange: (event: React.SyntheticEvent, newValue: number) => void;
}

const LEVELS = [
  { key: 'all', label: 'All Issues' },
  { key: 'fatal', label: 'Fatal' },
  { key: 'error', label: 'Errors' },
  { key: 'warning', label: 'Warnings' },
  { key: 'info', label: 'Info' },
  { key: 'exception', label: 'Exceptions' },
  { key: 'debug', label: 'Debug' },
];

const IssuesTabsList: React.FC<IssuesTabsListProps> = ({
  issues,
  filteredIssues,
  isLoading,
  error,
  viewMode,
  onIssueClick,
  onRefresh,
  totalCount,
  page,
  perPage,
  onPageChange,
  tabValue,
  onTabChange,
}) => {
  // Подсчёт количества по каждому уровню
  const counts = {
    all: issues.length,
    fatal: issues.filter(issue => issue.level === 'fatal').length,
    error: issues.filter(issue => issue.level === 'error').length,
    warning: issues.filter(issue => issue.level === 'warning').length,
    info: issues.filter(issue => issue.level === 'info').length,
    exception: issues.filter(issue => issue.level === 'exception').length,
    debug: issues.filter(issue => issue.level === 'debug').length,
  };

  return (
    <>
      <Box sx={{ overflowX: 'auto', '& .MuiTabs-root': { minWidth: 600 } }}>
        <Tabs 
          value={tabValue} 
          onChange={onTabChange} 
          aria-label="issue tabs"
          variant="scrollable"
          scrollButtons="auto"
          sx={{ 
            px: 2,
            '& .MuiTabs-indicator': {
              height: 2,
              borderRadius: '2px 2px 0 0',
            },
            '& .MuiTab-root': {
              textTransform: 'none',
              fontWeight: 500,
              fontSize: '0.85rem',
              minHeight: 40,
              minWidth: 100,
            }
          }}
        >
          {LEVELS.map((level, _idx) => (
            <Tab
              key={level.key}
              label={level.label}
              icon={
                <Badge 
                  badgeContent={counts[level.key as keyof typeof counts]} 
                  color={level.key === 'all' ? 'primary' : undefined}
                  sx={level.key !== 'all' ? { '& .MuiBadge-badge': { fontSize: '0.7rem', bgcolor: getLevelHexColor(level.key), color: '#fff' } } : { '& .MuiBadge-badge': { fontSize: '0.7rem' } }}
                />
              }
              iconPosition="end"
            />
          ))}
        </Tabs>
      </Box>
      <Divider />
      {LEVELS.map((level, idx) => (
        <div
          key={level.key}
          style={{ display: tabValue === idx ? 'block' : 'none' }}
        >
          <IssueList
            issues={filteredIssues}
            isLoading={isLoading}
            error={error}
            viewMode={viewMode}
            onIssueClick={onIssueClick}
            onRefresh={onRefresh}
            totalCount={totalCount}
            page={page}
            perPage={perPage}
            onPageChange={onPageChange}
          />
        </div>
      ))}
    </>
  );
};

export default IssuesTabsList; 