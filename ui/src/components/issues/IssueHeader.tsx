import React, { memo, useState, useCallback } from 'react';
import { 
  Box, 
  Typography, 
  Chip, 
  Menu, 
  MenuItem, 
  useTheme 
} from '@mui/material';
import { 
  Error as ErrorIcon, 
  Warning as WarningIcon, 
  Info as InfoIcon, 
  KeyboardArrowDown as KeyboardArrowDownIcon 
} from '@mui/icons-material';
import { IssueStatus, type Issue } from '../../generated/api/client';
import { useIssuePageStyles } from './IssuePageStyles';
import { getLevelColor, getLevelHexColor, getLevelBadgeStyles, getStatusColor } from '../../utils/issues/issueUtils';

// Get level icon based on level
const getLevelIcon = (level: string) => {
  switch (level) {
    case 'fatal':
      return <ErrorIcon sx={{ color: getLevelHexColor('fatal') }} />;
    case 'error':
      return <ErrorIcon sx={{ color: getLevelHexColor('error') }} />;
    case 'exception':
      return <ErrorIcon sx={{ color: getLevelHexColor('exception') }} />;
    case 'warning':
      return <WarningIcon sx={{ color: getLevelHexColor('warning') }} />;
    case 'info':
      return <InfoIcon sx={{ color: getLevelHexColor('info') }} />;
    case 'debug':
      return <InfoIcon sx={{ color: getLevelHexColor('debug') }} />;
    default:
      return <ErrorIcon sx={{ color: getLevelHexColor('error') }} />;
  }
};

interface IssueHeaderProps {
  issue: Issue;
  tags?: Record<string, string>;
  onStatusChange?: (newStatus: IssueStatus) => void;
  statusChangeLoading?: boolean;
}

const IssueHeader = memo(({ issue, tags, onStatusChange, statusChangeLoading = false }: IssueHeaderProps) => {
  const styles = useIssuePageStyles();
  const theme = useTheme();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);

  const handleClick = useCallback((event: React.MouseEvent<HTMLElement>) => {
    if (statusChangeLoading) return;
    setAnchorEl(event.currentTarget);
  }, [statusChangeLoading]);

  const handleClose = useCallback(() => {
    setAnchorEl(null);
  }, []);

  const handleStatusChange = useCallback((newStatus: IssueStatus) => {
    if (onStatusChange) {
      onStatusChange(newStatus);
    }
    handleClose();
  }, [onStatusChange, handleClose]);

  // Determine available status options based on current status
  const getAvailableStatusOptions = useCallback(() => {
    switch (issue.status) {
      case IssueStatus.Unresolved:
        return [IssueStatus.Resolved, IssueStatus.Ignored];
      case IssueStatus.Resolved:
        return [IssueStatus.Unresolved];
      case IssueStatus.Ignored:
        return [IssueStatus.Resolved, IssueStatus.Unresolved];
      default:
        return [];
    }
  }, [issue.status]);

  const availableStatusOptions = getAvailableStatusOptions();

  return (
    <Box sx={styles.headerContainer}>
      <Box sx={styles.headerTitleContainer}>
        {getLevelIcon(issue.level)}
        <Typography variant="h4" component="h1" sx={styles.headerTitle}>
          {issue.title}
        </Typography>
      </Box>
      <Typography variant="body1" color="text.secondary" sx={styles.headerMessage}>
        {issue.message}
      </Typography>
      <Box sx={styles.tagsContainer}>
        <Chip 
          label={issue.level} 
          color={getLevelColor(issue.level)}
          sx={{
            ...getLevelBadgeStyles(issue.level),
            '& .MuiChip-label': {
              fontSize: '0.875rem', // Larger font size for Issue page
              fontWeight: 500,
              px: 1.2,
              letterSpacing: '0.02em',
              textAlign: 'center',
              display: 'block',
            }
          }}
        />
        {onStatusChange && availableStatusOptions.length > 0 ? (
          <>
            <Chip 
              label={issue.status} 
              color={getStatusColor(issue.status)}
              onClick={handleClick}
              deleteIcon={<KeyboardArrowDownIcon />}
              onDelete={handleClick}
              disabled={statusChangeLoading}
            />
            <Menu
              anchorEl={anchorEl}
              open={open}
              onClose={handleClose}
              anchorOrigin={{
                vertical: 'bottom',
                horizontal: 'left',
              }}
              transformOrigin={{
                vertical: 'top',
                horizontal: 'left',
              }}
            >
              {availableStatusOptions.map((status) => (
                <MenuItem 
                  key={status} 
                  onClick={() => handleStatusChange(status)}
                  sx={{ 
                    color: theme => {
                      switch (status) {
                        case IssueStatus.Resolved:
                          return theme.palette.success.main;
                        case IssueStatus.Unresolved:
                          return theme.palette.error.main;
                        case IssueStatus.Ignored:
                          return theme.palette.text.secondary;
                        default:
                          return theme.palette.text.primary;
                      }
                    }
                  }}
                >
                  Change to {status}
                </MenuItem>
              ))}
            </Menu>
          </>
        ) : (
          <Chip 
            label={issue.status} 
            color={getStatusColor(issue.status)} 
          />
        )}
        {tags && Object.entries(tags).map(([key, value], _index) => (
          <Chip 
            key={`${key}-${value}`}
            label={`${key}: ${value}`} 
            variant="outlined"
          />
        ))}
      </Box>
    </Box>
  );
});

IssueHeader.displayName = 'IssueHeader';

export default IssueHeader;
