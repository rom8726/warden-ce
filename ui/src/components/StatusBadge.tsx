import React from 'react';
import { 
  Chip, 
  useTheme,
  Tooltip
} from '@mui/material';
import { 
  CheckCircle as CheckCircleIcon,
  Cancel as CancelIcon,
  PlayArrow as PlayIcon,
  Help as HelpIcon
} from '@mui/icons-material';
import { IssueStatus } from '../generated/api/client/api';

interface StatusBadgeProps {
  status: IssueStatus;
  size?: 'small' | 'medium';
  variant?: 'filled' | 'outlined';
  showIcon?: boolean;
  tooltip?: boolean;
}

const StatusBadge: React.FC<StatusBadgeProps> = ({ 
  status, 
  size = 'medium',
  variant = 'filled',
  showIcon = true,
  tooltip = false
}) => {
  const theme = useTheme();

  const getStatusConfig = () => {
    switch (status) {
      case IssueStatus.Resolved:
        return {
          label: 'Resolved',
          color: 'success' as const,
          icon: <CheckCircleIcon fontSize="small" />,
          description: 'Issue has been resolved'
        };
      case IssueStatus.Ignored:
        return {
          label: 'Ignored',
          color: 'default' as const,
          icon: <CancelIcon fontSize="small" />,
          description: 'Issue has been ignored'
        };
      case IssueStatus.Unresolved:
        return {
          label: 'Unresolved',
          color: 'error' as const,
          icon: <PlayIcon fontSize="small" />,
          description: 'Issue is unresolved'
        };
      default:
        return {
          label: 'Unknown',
          color: 'default' as const,
          icon: <HelpIcon fontSize="small" />,
          description: 'Unknown status'
        };
    }
  };

  const config = getStatusConfig();

  const chip = (
    <Chip
      label={config.label}
      color={config.color}
      variant={variant}
      size={size}
      icon={showIcon ? config.icon : undefined}
      sx={{
        fontWeight: 500,
        '& .MuiChip-icon': {
          fontSize: size === 'small' ? '1rem' : '1.25rem',
        },
        '& .MuiChip-label': {
          fontSize: size === 'small' ? '0.75rem' : '0.875rem',
          fontWeight: 500,
        },
        ...(variant === 'outlined' && {
          borderWidth: 2,
          '&:hover': {
            borderWidth: 2,
          },
        }),
      }}
    />
  );

  if (tooltip) {
    return (
      <Tooltip title={config.description} arrow>
        {chip}
      </Tooltip>
    );
  }

  return chip;
};

export default StatusBadge; 