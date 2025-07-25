import React from 'react';
import {
  Paper,
  Typography,
  Grid,
  Card,
  CardContent,
  Chip,
  Box,
  CircularProgress,
  Avatar,
  LinearProgress,
  useTheme,
} from '@mui/material';
import {
  TrendingUp as TrendingUpIcon,
  TrendingDown as TrendingDownIcon,
  Remove as RemoveIcon,
  CompareArrows as CompareArrowsIcon,
  Speed as SpeedIcon,
  BugReport as BugReportIcon,
  Group as GroupIcon,
  Timeline as TimelineIcon,
} from '@mui/icons-material';
import { useTheme as useAppTheme } from '../../../../theme/ThemeContext';

interface ComparisonData {
  delta: {
    known_issues_total: number;
    new_issues_total: number;
    regressions_total: number;
    resolved_in_version_total: number;
    users_affected: number;
  };
  base_version: string;
  target_version: string;
}

interface ReleasesDeltaProps {
  comparisonData: ComparisonData | null | undefined;
  loading: boolean;
  error: string | null;
  compareMode: boolean;
  compareRelease: string | null;
  selectedRelease: string | null;
}

const ReleasesDelta: React.FC<ReleasesDeltaProps> = ({
  comparisonData,
  loading,
  error,
  compareMode,
  compareRelease,
  selectedRelease,
}) => {
  const theme = useTheme();
  const { mode } = useAppTheme();

  const getDeltaCardBackground = () => {
    switch (mode) {
      case 'dark':
        return 'linear-gradient(135deg, rgba(45, 48, 56, 0.9) 0%, rgba(35, 38, 46, 0.9) 100%)';
      case 'blue':
        return 'linear-gradient(135deg, rgba(25, 55, 84, 0.9) 0%, rgba(16, 42, 66, 0.9) 100%)';
      case 'green':
        return 'linear-gradient(135deg, rgba(50, 95, 65, 0.9) 0%, rgba(40, 85, 50, 0.9) 70%, rgba(60, 105, 75, 0.9) 100%)';
      case 'light':
      default:
        return 'linear-gradient(135deg, rgba(255, 193, 7, 0.1) 0%, rgba(255, 152, 0, 0.1) 100%)';
    }
  };

  const getDeltaCardTextColor = () => {
    switch (mode) {
      case 'light':
        return theme.palette.text.primary;
      default:
        return 'white';
    }
  };

  const getDeltaCardBorderColor = () => {
    switch (mode) {
      case 'light':
        return 'rgba(255, 193, 7, 0.2)';
      default:
        return 'rgba(255, 255, 255, 0.2)';
    }
  };

  const getMetricCardBackground = () => {
    switch (mode) {
      case 'light':
        return 'rgba(255, 193, 7, 0.05)';
      default:
        return 'rgba(255, 255, 255, 0.1)';
    }
  };

  const getMetricCardBorderColor = () => {
    switch (mode) {
      case 'light':
        return 'rgba(255, 193, 7, 0.15)';
      default:
        return 'rgba(255, 255, 255, 0.2)';
    }
  };

  const getMetricCardTextColor = () => {
    switch (mode) {
      case 'light':
        return theme.palette.text.secondary;
      default:
        return 'rgba(255, 255, 255, 0.8)';
    }
  };

  const renderDelta = (value: number | undefined, label: string, icon: React.ReactNode) => {
    if (value === undefined || value === 0) {
      return (
        <Box sx={{ 
          textAlign: 'center',
          p: 2,
          borderRadius: 2,
          bgcolor: getMetricCardBackground(),
          border: `1px solid ${getMetricCardBorderColor()}`
        }}>
          <Avatar sx={{ 
            bgcolor: 'text.secondary',
            width: 48,
            height: 48,
            mx: 'auto',
            mb: 1
          }}>
            {icon}
          </Avatar>
          <Typography variant="h5" sx={{ 
            fontWeight: 'bold',
            color: 'text.secondary',
            mb: 0.5
          }}>
            0
          </Typography>
          <Typography variant="body2" sx={{ color: getMetricCardTextColor(), mb: 1 }}>
            {label}
          </Typography>
          <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 0.5 }}>
            <RemoveIcon sx={{ fontSize: 16, color: 'text.secondary' }} />
            <Typography variant="caption" sx={{ color: 'text.secondary', fontWeight: 600 }}>
              No Change
            </Typography>
          </Box>
        </Box>
      );
    }
    
    const isGoodChange = (val: number, metricLabel: string): boolean => {
      switch (metricLabel) {
        case 'Known Issues':
        case 'New Issues':
        case 'Regressions':
        case 'Users Affected':
          return val < 0;
        case 'Resolved':
          return val > 0;
        default:
          return val < 0;
      }
    };
    
    const isGood = isGoodChange(value, label);
    const Icon = value > 0 ? TrendingUpIcon : TrendingDownIcon;
    const color = isGood ? 'success.main' : 'error.main';
    const progress = Math.min(Math.abs(value) * 10, 100);
    
    return (
      <Box sx={{ 
        textAlign: 'center',
        p: 2,
        borderRadius: 2,
        bgcolor: getMetricCardBackground(),
        border: `1px solid ${getMetricCardBorderColor()}`,
        transition: 'all 0.2s ease-in-out',
        '&:hover': {
          bgcolor: mode === 'light' ? 'rgba(255, 193, 7, 0.08)' : 'rgba(255, 255, 255, 0.15)',
          transform: 'translateY(-2px)',
          boxShadow: mode === 'light' 
            ? '0 6px 16px rgba(255, 193, 7, 0.15)' 
            : '0 6px 16px rgba(255, 255, 255, 0.1)'
        }
      }}>
        <Avatar sx={{ 
          bgcolor: color,
          width: 48,
          height: 48,
          mx: 'auto',
          mb: 1
        }}>
          {icon}
        </Avatar>
        <Typography variant="h5" sx={{ 
          fontWeight: 'bold',
          color: color,
          mb: 0.5
        }}>
          {value > 0 ? '+' : ''}{value}
        </Typography>
        <Typography variant="body2" sx={{ color: getMetricCardTextColor(), mb: 1 }}>
          {label}
        </Typography>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 0.5, mb: 1 }}>
          <Icon sx={{ fontSize: 16, color }} />
          <Typography variant="caption" sx={{ color, fontWeight: 600 }}>
            {isGood ? 'Improved' : 'Worsened'}
          </Typography>
        </Box>
        <LinearProgress 
          variant="determinate" 
          value={progress} 
          sx={{ 
            height: 4,
            borderRadius: 2,
            bgcolor: getMetricCardBackground(),
            '& .MuiLinearProgress-bar': {
              bgcolor: color
            }
          }} 
        />
      </Box>
    );
  };

  if (!compareMode || !comparisonData) {
    return null;
  }

  if (loading) {
    return (
      <Paper sx={{ p: 3, mb: 3 }}>
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
          <CircularProgress />
        </Box>
      </Paper>
    );
  }

  if (error) {
    return (
      <Paper sx={{ p: 3, mb: 3 }}>
        <Typography color="error">
          Error loading comparison data. Please try again.
        </Typography>
      </Paper>
    );
  }

  return (
    <Paper sx={{ p: 3, mb: 3 }}>
      <Typography variant="h6" gutterBottom className="gradient-subtitle-green">
        Release Comparison: {compareRelease} → {selectedRelease}
      </Typography>

      <Card sx={{ 
        background: getDeltaCardBackground(),
        color: getDeltaCardTextColor(),
        position: 'relative',
        overflow: 'hidden'
      }}>
        <Box sx={{
          position: 'absolute',
          top: -20,
          right: -20,
          width: 100,
          height: 100,
          borderRadius: '50%',
          background: mode === 'light' ? 'rgba(255, 193, 7, 0.1)' : 'rgba(255,255,255,0.1)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center'
        }}>
          <CompareArrowsIcon sx={{ 
            fontSize: 40, 
            opacity: mode === 'light' ? 0.4 : 0.3,
            color: mode === 'light' ? '#ffc107' : 'white'
          }} />
        </Box>
        <CardContent sx={{ position: 'relative', zIndex: 1 }}>
          <Typography variant="h6" gutterBottom sx={{}} className="gradient-subtitle-green">
            <CompareArrowsIcon sx={{ mr: 1 }} />
            Releases Delta Analysis
          </Typography>
          
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6} md={2.4}>
              {renderDelta(comparisonData.delta.known_issues_total, 'Known Issues', <BugReportIcon />)}
            </Grid>
            <Grid item xs={12} sm={6} md={2.4}>
              {renderDelta(comparisonData.delta.new_issues_total, 'New Issues', <BugReportIcon />)}
            </Grid>
            <Grid item xs={12} sm={6} md={2.4}>
              {renderDelta(comparisonData.delta.regressions_total, 'Regressions', <TrendingDownIcon />)}
            </Grid>
            <Grid item xs={12} sm={6} md={2.4}>
              {renderDelta(comparisonData.delta.resolved_in_version_total, 'Resolved', <SpeedIcon />)}
            </Grid>
            <Grid item xs={12} sm={6} md={2.4}>
              {renderDelta(comparisonData.delta.users_affected, 'Users Affected', <GroupIcon />)}
            </Grid>
          </Grid>

          {/* Summary */}
          <Box sx={{ 
            mt: 3, 
            p: 2, 
            borderRadius: 2, 
            bgcolor: getMetricCardBackground(),
            border: `1px solid ${getMetricCardBorderColor()}`
          }}>
            <Typography variant="subtitle2" sx={{ color: getMetricCardTextColor(), mb: 1, fontWeight: 600 }}>
              📊 Comparison Summary
            </Typography>
            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
              <Chip 
                label={`${compareRelease} → ${selectedRelease}`}
                size="small"
                sx={{ 
                  bgcolor: mode === 'light' ? 'rgba(255, 193, 7, 0.2)' : 'rgba(255, 255, 255, 0.2)',
                  color: getDeltaCardTextColor(),
                  fontWeight: 600
                }}
              />
              <Chip 
                label={`Total Changes: ${Object.values(comparisonData.delta).reduce((sum, val) => sum + Math.abs(val || 0), 0)}`}
                size="small"
                variant="outlined"
                sx={{ 
                  borderColor: getDeltaCardBorderColor(),
                  color: getDeltaCardTextColor()
                }}
              />
            </Box>
          </Box>
        </CardContent>
      </Card>
    </Paper>
  );
};

export default ReleasesDelta; 