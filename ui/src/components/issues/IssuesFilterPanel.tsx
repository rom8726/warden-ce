import React from 'react';
import { Box, FormControl, InputLabel, Select, MenuItem, ToggleButtonGroup, ToggleButton, IconButton, Tooltip } from '@mui/material';
import type { SelectChangeEvent } from '@mui/material/Select';
import ViewCompactIcon from '@mui/icons-material/ViewCompact';
import ViewStreamIcon from '@mui/icons-material/ViewStream';
import ViewAgendaIcon from '@mui/icons-material/ViewAgenda';
import RefreshIcon from '@mui/icons-material/Refresh';
import SearchBar from '../SearchBar';
import { IssueSortColumn, SortOrder } from '../../generated/api/client';

interface IssuesFilterPanelProps {
  searchQuery: string;
  onSearchChange: (value: string) => void;
  levelFilter: string;
  onLevelFilterChange: (event: SelectChangeEvent) => void;
  statusFilter: string;
  onStatusFilterChange: (event: SelectChangeEvent) => void;
  viewMode: string;
  onViewModeChange: (event: React.MouseEvent<HTMLElement>, newViewMode: string) => void;
  onRefresh: () => void;
  isLoading?: boolean;
  sortBy?: IssueSortColumn;
  onSortByChange?: (event: SelectChangeEvent) => void;
  sortOrder?: SortOrder;
  onSortOrderChange?: (event: SelectChangeEvent) => void;
}

const IssuesFilterPanel: React.FC<IssuesFilterPanelProps> = ({
  searchQuery,
  onSearchChange,
  levelFilter,
  onLevelFilterChange,
  statusFilter,
  onStatusFilterChange,
  viewMode,
  onViewModeChange,
  onRefresh,
  isLoading = false,
  sortBy,
  onSortByChange,
  sortOrder,
  onSortOrderChange,
}) => {
  return (
    <Box sx={{ px: 3, pt: 2, pb: 1 }}>
      <Box sx={{ display: 'flex', flexWrap: 'nowrap', alignItems: 'center', minHeight: 48, gap: 1.5, minWidth: 0, overflowX: 'auto' }}>
        <Box sx={{ flexGrow: 1 }} />
        <Box sx={{ width: 220, flex: '0 0 220px', mr: 1.5, display: 'flex', alignItems: 'center' }}>
          <SearchBar
            placeholder="Search issues..."
            value={searchQuery}
            onChange={onSearchChange}
            size="small"
          />
        </Box>
        <FormControl size="small" sx={{ minWidth: 120 }}>
          <InputLabel id="level-filter-label">Level</InputLabel>
          <Select
            labelId="level-filter-label"
            id="level-filter"
            value={levelFilter}
            label="Level"
            onChange={onLevelFilterChange}
          >
            <MenuItem value="all">All Levels</MenuItem>
            <MenuItem value="fatal">Fatal</MenuItem>
            <MenuItem value="error">Error</MenuItem>
            <MenuItem value="exception">Exception</MenuItem>
            <MenuItem value="warning">Warning</MenuItem>
            <MenuItem value="info">Info</MenuItem>
            <MenuItem value="debug">Debug</MenuItem>
          </Select>
        </FormControl>
        <FormControl size="small" sx={{ minWidth: 120 }}>
          <InputLabel id="status-filter-label">Status</InputLabel>
          <Select
            labelId="status-filter-label"
            id="status-filter"
            value={statusFilter}
            label="Status"
            onChange={onStatusFilterChange}
          >
            <MenuItem value="all">All Statuses</MenuItem>
            <MenuItem value="unresolved">Unresolved</MenuItem>
            <MenuItem value="resolved">Resolved</MenuItem>
            <MenuItem value="ignored">Ignored</MenuItem>
          </Select>
        </FormControl>
        {onSortByChange && (
          <FormControl size="small" sx={{ minWidth: 140 }}>
            <InputLabel id="sort-by-label">Sort By</InputLabel>
            <Select
              labelId="sort-by-label"
              id="sort-by"
              value={sortBy || ''}
              label="Sort By"
              onChange={onSortByChange}
            >
              <MenuItem value={IssueSortColumn.TotalEvents}>Total Events</MenuItem>
              <MenuItem value={IssueSortColumn.FirstSeen}>First Seen</MenuItem>
              <MenuItem value={IssueSortColumn.LastSeen}>Last Seen</MenuItem>
            </Select>
          </FormControl>
        )}
        {onSortOrderChange && sortBy && (
          <FormControl size="small" sx={{ minWidth: 120 }}>
            <InputLabel id="sort-order-label">Order</InputLabel>
            <Select
              labelId="sort-order-label"
              id="sort-order"
              value={sortOrder || SortOrder.Desc}
              label="Order"
              onChange={onSortOrderChange}
            >
              <MenuItem value={SortOrder.Asc}>Ascending</MenuItem>
              <MenuItem value={SortOrder.Desc}>Descending</MenuItem>
            </Select>
          </FormControl>
        )}
        <ToggleButtonGroup
          value={viewMode}
          exclusive
          onChange={onViewModeChange}
          aria-label="view mode"
          size="small"
          sx={{ height: 40 }}
        >
          <ToggleButton value="compact" aria-label="compact view">
            <Tooltip title="Compact View">
              <ViewCompactIcon fontSize="small" />
            </Tooltip>
          </ToggleButton>
          <ToggleButton value="medium" aria-label="medium view">
            <Tooltip title="Medium View">
              <ViewStreamIcon fontSize="small" />
            </Tooltip>
          </ToggleButton>
          <ToggleButton value="large" aria-label="large view">
            <Tooltip title="Large View">
              <ViewAgendaIcon fontSize="small" />
            </Tooltip>
          </ToggleButton>
        </ToggleButtonGroup>
        <IconButton 
          size="small" 
          onClick={onRefresh}
          disabled={isLoading}
          sx={{ 
            bgcolor: theme => theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.05)' : 'rgba(0, 0, 0, 0.04)',
            '&:hover': {
              bgcolor: theme => theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.08)',
            },
            height: 40, width: 40
          }}
        >
          <RefreshIcon fontSize="small" />
        </IconButton>
      </Box>
    </Box>
  );
};

export default IssuesFilterPanel; 