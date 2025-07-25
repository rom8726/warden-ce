import React, { useState } from 'react';
import { 
  Box, 
  Typography, 
  Card,
  CardContent,
  CardActionArea,
  CardActions,
  Chip,
  IconButton,
  useTheme,
  Collapse
} from '@mui/material';
import { 
  Error as ErrorIcon,
  Warning as WarningIcon,
  Info as InfoIcon,
  MoreVert as MoreVertIcon,
  ContentCopy as ContentCopyIcon,
  ExpandMore as ExpandMoreIcon,
  ExpandLess as ExpandLessIcon
} from '@mui/icons-material';
import StatusBadge from "../StatusBadge";
import { getLevelColor, getLevelHexColor, getLevelBadgeStyles } from '../../utils/issues/issueUtils';
import type { Issue } from '../../generated/api/client';

// Interface for our issue data with additional fields
interface ExtendedIssue extends Issue {
  projectName?: string;
  fingerprint?: string;
}

interface IssueCardProps {
  issue: ExtendedIssue;
  viewMode: 'compact' | 'medium' | 'large';
  onIssueClick: (projectId: number, issueId: number) => void;
}

// Helper functions (could be moved to utils if needed elsewhere)
const getLevelIcon = (level: string) => {
  switch (level) {
    case 'fatal':
      return <ErrorIcon sx={{ color: getLevelHexColor('fatal') }} fontSize="small" />;
    case 'error':
      return <ErrorIcon sx={{ color: getLevelHexColor('error') }} fontSize="small" />;
    case 'exception':
      return <ErrorIcon sx={{ color: getLevelHexColor('exception') }} fontSize="small" />;
    case 'warning':
      return <WarningIcon sx={{ color: getLevelHexColor('warning') }} fontSize="small" />;
    case 'info':
      return <InfoIcon sx={{ color: getLevelHexColor('info') }} fontSize="small" />;
    case 'debug':
      return <InfoIcon sx={{ color: getLevelHexColor('debug') }} fontSize="small" />;
    default:
      return <ErrorIcon sx={{ color: getLevelHexColor('error') }} fontSize="small" />;
  }
};

const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffDays === 0) {
    return `Today at ${date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}`;
  } else if (diffDays === 1) {
    return `Yesterday at ${date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}`;
  } else if (diffDays < 7) {
    return `${diffDays} days ago`;
  } else {
    return date.toLocaleDateString([], { month: 'short', day: 'numeric' });
  }
};

const IssueCard: React.FC<IssueCardProps> = ({ issue, viewMode, onIssueClick }) => {
  const theme = useTheme();
  const [showFingerprint, setShowFingerprint] = useState(false);

  return (
    <Card 
      sx={{ 
        borderRadius: 2,
        background: theme.palette.mode === 'dark'
          ? 'linear-gradient(135deg, rgba(45, 48, 56, 0.7) 0%, rgba(35, 38, 46, 0.7) 100%)'
          : 'linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(248, 249, 252, 0.95) 100%)',
        backdropFilter: 'blur(10px)',
        boxShadow: theme.palette.mode === 'dark'
          ? '0 4px 20px 0 rgba(0, 0, 0, 0.2)'
          : '0 4px 20px 0 rgba(0, 0, 0, 0.05)',
        transition: 'all 0.3s ease-in-out',
        border: `1px solid ${
          theme.palette.mode === 'dark'
            ? 'rgba(255, 255, 255, 0.05)'
            : 'rgba(0, 0, 0, 0.03)'
        }`,
        '&:hover': { 
          transform: 'translateY(-4px)',
          boxShadow: theme.palette.mode === 'dark'
            ? '0 8px 30px 0 rgba(0, 0, 0, 0.3)'
            : '0 8px 30px 0 rgba(0, 0, 0, 0.1)',
          borderColor: theme.palette.mode === 'dark'
            ? 'rgba(255, 255, 255, 0.1)'
            : 'rgba(0, 0, 0, 0.05)',
        }
      }}
    >
      <CardActionArea 
        onClick={() => onIssueClick(issue.project_id, issue.id)}
        sx={{ 
          borderRadius: 2,
          '&:hover .MuiCardActionArea-focusHighlight': {
            opacity: 0.05
          }
        }}
      >
        <CardContent sx={{ p: viewMode === 'compact' ? 1.5 : 2 }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 1 }}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, flex: 1, minWidth: 0 }}>
              {getLevelIcon(issue.level)}
              <Typography 
                variant={viewMode === 'compact' ? 'body2' : 'body1'} 
                sx={{ 
                  fontWeight: 600,
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  whiteSpace: 'nowrap',
                  flex: 1
                }}
              >
                {issue.title}
              </Typography>
            </Box>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, ml: 1 }}>
              <StatusBadge status={issue.status} size="small" tooltip />
              <IconButton 
                size="small" 
                sx={{ opacity: 0.6 }}
                onClick={(e) => {
                  e.stopPropagation();
                  // Handle menu click here
                }}
              >
                <MoreVertIcon fontSize="small" />
              </IconButton>
            </Box>
          </Box>

          <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mb: 1 }}>
            <Chip 
              label={issue.level} 
              size="small" 
              color={getLevelColor(issue.level)} 
              sx={getLevelBadgeStyles(issue.level, 'small')} 
            />
            <Chip 
              label={`${issue.count} occurrences`} 
              size="small" 
              variant="outlined"
              sx={{ 
                height: 20,
                '& .MuiChip-label': {
                  px: 1,
                  fontSize: '0.7rem',
                },
              }} 
            />
            <Chip 
              label={issue.projectName || `Project ${issue.project_id}`} 
              size="small" 
              variant="outlined"
              sx={{ 
                height: 20,
                '& .MuiChip-label': {
                  px: 1,
                  fontSize: '0.7rem',
                },
              }} 
            />
          </Box>

          {viewMode !== 'compact' && (
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <Typography variant="caption" color="text.secondary">
                Last seen: {formatDate(issue.last_seen)}
              </Typography>
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <Typography variant="caption" color="text.secondary">
                  First seen: {formatDate(issue.first_seen)}
                </Typography>
                {viewMode === 'large' && issue.status === 'resolved' && issue.resolved_by && (
                  <Typography
                    variant="caption"
                    sx={{
                      color: theme.palette.success.main,
                      fontWeight: 500,
                      ml: 2,
                    }}
                  >
                    Resolved by: {issue.resolved_by}
                  </Typography>
                )}
              </Box>
            </Box>
          )}
        </CardContent>
      </CardActionArea>
      {viewMode === 'large' && (
        <Collapse in={viewMode === 'large'} timeout="auto" unmountOnExit>
          <CardActions sx={{ pt: 0, pl: 2, pr: 2 }}>
            <Box sx={{ width: '100%' }}>
              {issue.fingerprint && (
                <Box>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Typography variant="caption" sx={{ fontWeight: 'bold' }}>
                      Fingerprint:
                    </Typography>
                    <IconButton size="small" onClick={e => { e.stopPropagation(); setShowFingerprint(v => !v); }}>
                      {showFingerprint ? <ExpandLessIcon fontSize="small" /> : <ExpandMoreIcon fontSize="small" />}
                    </IconButton>
                    <IconButton size="small" onClick={e => { e.stopPropagation(); navigator.clipboard.writeText(issue.fingerprint); }}>
                      <ContentCopyIcon fontSize="small" />
                    </IconButton>
                  </Box>
                  <Collapse in={showFingerprint}>
                    <Typography 
                      variant="caption"
                      sx={{ 
                        wordBreak: 'break-all', 
                        p: 1, 
                        borderRadius: 1, 
                        color: 'text.secondary',
                        fontFamily: 'monospace',
                        display: 'block'
                      }}>
                      {issue.fingerprint}
                    </Typography>
                  </Collapse>
                </Box>
              )}
              {issue.platform && (
                <Chip 
                  label={`Platform: ${issue.platform}`}
                  size="small"
                  sx={{ mt: 1, fontSize: '0.75rem' }}
                />
              )}
            </Box>
          </CardActions>
        </Collapse>
      )}
    </Card>
  );
};

export default IssueCard; 