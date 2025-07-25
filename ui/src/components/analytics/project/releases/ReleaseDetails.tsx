import React from 'react';
import { Link } from 'react-router-dom';
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
  BugReport as BugReportIcon,
  PieChart as PieChartIcon,
  Speed as SpeedIcon,
  Group as GroupIcon,
  TrendingUp as TrendingUpIcon,
  TrendingDown as TrendingDownIcon,
  Remove as RemoveIcon,
  Timer as TimerIcon,
  Schedule as ScheduleIcon,
  AccessTime as AccessTimeIcon,
  Speed as SpeedometerIcon,
  Timeline as TimelineIcon,
  Link as LinkIcon,
} from '@mui/icons-material';
import {
  PieChart,
  Pie,
  Cell,
  ResponsiveContainer,
} from 'recharts';
import { IssueLevel } from '../../../../generated/api/client';
import { getLevelHexColor } from '../../../../utils/issues/issueUtils';
import { formatDuration } from '../../../../utils/timeUtils';
import { useTheme as useAppTheme } from '../../../../theme/ThemeContext';

// Local interface matching the API structure
interface ReleaseDetailsData {
  version: string;
  created_at: string;
  stats: {
    version: string;
    created_at: string;
    known_issues_total: number;
    new_issues_total: number;
    regressions_total: number;
    resolved_in_version_total: number;
    users_affected: number;
  };
  top_issues: Array<{
    id: number;
    project_id: number;
    title: string;
    level: IssueLevel;
    count: number;
    last_seen: string;
  }>;
  severity_distribution: { [key: string]: number };
  fix_time: {
    avg?: number;
    median?: number;
    p95?: number;
  };
  segments: {
    platform?: { [key: string]: number };
    browser_name?: { [key: string]: number };
    os_name?: { [key: string]: number };
    device_arch?: { [key: string]: number };
    runtime_name?: { [key: string]: number };
  };
}

interface ReleaseDetailsProps {
  releaseDetails: ReleaseDetailsData | null | undefined;
  loading: boolean;
  error: string | null;
  compareMode: boolean;
  comparisonData: any;
  projectId: string | undefined;
}

const ReleaseDetails: React.FC<ReleaseDetailsProps> = ({
  releaseDetails,
  loading,
  error,
  compareMode,
  comparisonData,
  projectId,
}) => {
  const theme = useTheme();
  const { mode } = useAppTheme();

  const getFixTimeCardBackground = () => {
    switch (mode) {
      case 'dark':
        return 'linear-gradient(135deg, rgba(45, 48, 56, 0.9) 0%, rgba(35, 38, 46, 0.9) 100%)';
      case 'blue':
        return 'linear-gradient(135deg, rgba(25, 55, 84, 0.9) 0%, rgba(16, 42, 66, 0.9) 100%)';
      case 'green':
        return 'linear-gradient(135deg, rgba(50, 95, 65, 0.9) 0%, rgba(40, 85, 50, 0.9) 70%, rgba(60, 105, 75, 0.9) 100%)';
      case 'light':
      default:
        return 'linear-gradient(135deg, rgba(130, 82, 255, 0.1) 0%, rgba(150, 110, 255, 0.1) 100%)';
    }
  };

  const getFixTimeCardTextColor = () => {
    switch (mode) {
      case 'light':
        return theme.palette.text.primary;
      default:
        return 'white';
    }
  };

  const getMetricCardBackground = () => {
    switch (mode) {
      case 'light':
        return 'rgba(130, 82, 255, 0.05)';
      default:
        return 'rgba(255, 255, 255, 0.1)';
    }
  };

  const getMetricCardBorderColor = () => {
    switch (mode) {
      case 'light':
        return 'rgba(130, 82, 255, 0.15)';
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

  const getTopIssuesCardBackground = () => {
    switch (mode) {
      case 'dark':
        return 'linear-gradient(135deg, rgba(45, 48, 56, 0.9) 0%, rgba(35, 38, 46, 0.9) 100%)';
      case 'blue':
        return 'linear-gradient(135deg, rgba(25, 55, 84, 0.9) 0%, rgba(16, 42, 66, 0.9) 100%)';
      case 'green':
        return 'linear-gradient(135deg, rgba(50, 95, 65, 0.9) 0%, rgba(40, 85, 50, 0.9) 70%, rgba(60, 105, 75, 0.9) 100%)';
      case 'light':
      default:
        return 'linear-gradient(135deg, rgba(255, 107, 107, 0.1) 0%, rgba(255, 159, 67, 0.1) 100%)';
    }
  };

  const getSeverityDistributionCardBackground = () => {
    switch (mode) {
      case 'dark':
        return 'linear-gradient(135deg, rgba(45, 48, 56, 0.9) 0%, rgba(35, 38, 46, 0.9) 100%)';
      case 'blue':
        return 'linear-gradient(135deg, rgba(25, 55, 84, 0.9) 0%, rgba(16, 42, 66, 0.9) 100%)';
      case 'green':
        return 'linear-gradient(135deg, rgba(50, 95, 65, 0.9) 0%, rgba(40, 85, 50, 0.9) 70%, rgba(60, 105, 75, 0.9) 100%)';
      case 'light':
      default:
        return 'linear-gradient(135deg, rgba(72, 149, 239, 0.1) 0%, rgba(100, 181, 246, 0.1) 100%)';
    }
  };

  const getUserSegmentsCardBackground = () => {
    switch (mode) {
      case 'dark':
        return 'linear-gradient(135deg, rgba(45, 48, 56, 0.9) 0%, rgba(35, 38, 46, 0.9) 100%)';
      case 'blue':
        return 'linear-gradient(135deg, rgba(25, 55, 84, 0.9) 0%, rgba(16, 42, 66, 0.9) 100%)';
      case 'green':
        return 'linear-gradient(135deg, rgba(50, 95, 65, 0.9) 0%, rgba(40, 85, 50, 0.9) 70%, rgba(60, 105, 75, 0.9) 100%)';
      case 'light':
      default:
        return 'linear-gradient(135deg, rgba(46, 204, 113, 0.1) 0%, rgba(52, 211, 153, 0.1) 100%)';
    }
  };

  const getIssueCardBackground = () => {
    switch (mode) {
      case 'light':
        return 'rgba(255, 107, 107, 0.05)';
      default:
        return 'rgba(255, 255, 255, 0.1)';
    }
  };

  const getIssueCardBorderColor = () => {
    switch (mode) {
      case 'light':
        return 'rgba(255, 107, 107, 0.2)';
      default:
        return 'rgba(255, 255, 255, 0.2)';
    }
  };

  const getLevelIcon = (level: IssueLevel) => {
    switch (level) {
      case 'fatal': return <BugReportIcon />;
      case 'exception': return <BugReportIcon />;
      case 'error': return <BugReportIcon />;
      case 'warning': return <BugReportIcon />;
      case 'info': return <BugReportIcon />;
      case 'debug': return <BugReportIcon />;
      default: return <BugReportIcon />;
    }
  };

  const renderDelta = (value: number | undefined, label: string) => {
    if (value === undefined || value === 0) {
      return (
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
          <RemoveIcon sx={{ fontSize: 16, color: 'text.secondary' }} />
          <Typography variant="body2" color="text.secondary">
            {label}: 0
          </Typography>
        </Box>
      );
    }
    
    // Определяем, является ли изменение "хорошим" или "плохим"
    const isGoodChange = (val: number, metricLabel: string): boolean => {
      switch (metricLabel) {
        case 'Known Issues':
        case 'New Issues':
        case 'Regressions':
        case 'Users':
          // Для этих метрик уменьшение - это хорошо
          return val < 0;
        case 'Resolved':
          // Для resolved issues увеличение - это хорошо
          return val > 0;
        default:
          // По умолчанию считаем, что уменьшение - это хорошо
          return val < 0;
      }
    };
    
    const isGood = isGoodChange(value, label);
    const Icon = value > 0 ? TrendingUpIcon : TrendingDownIcon;
    const color = isGood ? 'success.main' : 'error.main';
    
    return (
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
        <Icon sx={{ fontSize: 16, color }} />
        <Typography variant="body2" color={color}>
          {label}: {value > 0 ? '+' : ''}{value}
        </Typography>
      </Box>
    );
  };

  const prepareSeverityData = (severityDistribution: { [key: string]: number }) => {
    return Object.entries(severityDistribution).map(([level, count]) => ({
      name: level,
      value: count,
      color: getLevelHexColor(level as IssueLevel)
    }));
  };

  const getFixTimeColor = (seconds: number | undefined) => {
    if (!seconds || seconds <= 0) return 'text.secondary';
    if (seconds < 3600) return 'success.main'; // < 1 hour - green
    if (seconds < 86400) return 'warning.main'; // < 1 day - orange
    if (seconds < 604800) return 'error.main'; // < 1 week - red
    return 'error.dark'; // > 1 week - dark red
  };

  const getFixTimeIcon = (type: 'avg' | 'median' | 'p95') => {
    switch (type) {
      case 'avg': return <SpeedometerIcon />;
      case 'median': return <TimelineIcon />;
      case 'p95': return <TimerIcon />;
      default: return <AccessTimeIcon />;
    }
  };

  const getFixTimeLabel = (type: 'avg' | 'median' | 'p95') => {
    switch (type) {
      case 'avg': return 'Average';
      case 'median': return 'Median';
      case 'p95': return '95th Percentile';
      default: return '';
    }
  };

  const getSpeedRating = (seconds: number | undefined) => {
    if (!seconds || seconds <= 0) return { rating: 'N/A', color: 'text.secondary', progress: 0 };
    if (seconds < 3600) return { rating: 'Fast', color: 'success.main', progress: 100 };
    if (seconds < 86400) return { rating: 'Good', color: 'success.main', progress: 75 };
    if (seconds < 604800) return { rating: 'Slow', color: 'warning.main', progress: 50 };
    return { rating: 'Very Slow', color: 'error.main', progress: 25 };
  };

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
          Error loading release details. Please try again.
        </Typography>
      </Paper>
    );
  }

  if (!releaseDetails) {
    return (
      <Paper sx={{ p: 3, mb: 3 }}>
        <Typography color="text.secondary">
          Select a release to view details.
        </Typography>
      </Paper>
    );
  }

  const severityData = prepareSeverityData(releaseDetails.severity_distribution);

  return (
    <Paper sx={{ p: 3, mb: 3 }}>
      <Typography variant="h6" gutterBottom>
        Release Details: {releaseDetails.version}
      </Typography>

      {/* Comparison Summary */}
      {compareMode && comparisonData && (
        <Card sx={{ mb: 3, bgcolor: 'background.paper' }}>
          <CardContent>
            <Typography variant="subtitle2" gutterBottom color="text.secondary">
              📊 Comparison Summary:
            </Typography>
            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 2 }}>
              {renderDelta(comparisonData.delta.known_issues_total, 'Known Issues')}
              {renderDelta(comparisonData.delta.new_issues_total, 'New Issues')}
              {renderDelta(comparisonData.delta.regressions_total, 'Regressions')}
              {renderDelta(comparisonData.delta.resolved_in_version_total, 'Resolved')}
              {renderDelta(comparisonData.delta.users_affected, 'Users')}
            </Box>
          </CardContent>
        </Card>
      )}

      <Grid container spacing={3}>
        {/* Top Issues */}
        <Grid item xs={12} md={6}>
          <Card sx={{ 
            height: '100%', 
            background: getTopIssuesCardBackground(),
            color: getFixTimeCardTextColor(),
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
              background: mode === 'light' ? 'rgba(255, 107, 107, 0.1)' : 'rgba(255,255,255,0.1)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center'
            }}>
              <BugReportIcon sx={{ 
                fontSize: 40, 
                opacity: mode === 'light' ? 0.4 : 0.3,
                color: mode === 'light' ? '#ff6b6b' : 'white'
              }} />
            </Box>
            <CardContent sx={{ position: 'relative', zIndex: 1 }}>
              <Typography variant="h6" gutterBottom sx={{ 
                display: 'flex', 
                alignItems: 'center',
                color: getFixTimeCardTextColor(),
                fontWeight: 600
              }}>
                <BugReportIcon sx={{ mr: 1 }} />
                Top Issues (5)
              </Typography>
              {releaseDetails.top_issues.slice(0, 5).map((issue) => (
                <Link 
                  key={issue.id}
                  to={projectId ? `/projects/${projectId}/issues/${issue.id}` : `/issues/${issue.id}`}
                  style={{ 
                    textDecoration: 'none',
                    color: 'inherit'
                  }}
                >
                  <Box sx={{ 
                    mb: 2, 
                    p: 1.5, 
                    borderRadius: 2,
                    bgcolor: getIssueCardBackground(),
                    border: `1px solid ${getIssueCardBorderColor()}`,
                    transition: 'all 0.2s ease-in-out',
                    cursor: 'pointer',
                    '&:hover': {
                      bgcolor: mode === 'light' ? 'rgba(255, 107, 107, 0.08)' : 'rgba(255, 255, 255, 0.15)',
                      transform: 'translateY(-1px)',
                      boxShadow: mode === 'light' 
                        ? '0 4px 12px rgba(255, 107, 107, 0.15)' 
                        : '0 4px 12px rgba(255, 255, 255, 0.1)'
                    }
                  }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                      <Box sx={{ flex: 1, display: 'flex', alignItems: 'center', gap: 1 }}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                          <LinkIcon sx={{ 
                            fontSize: 14, 
                            color: mode === 'light' ? '#ff6b6b' : '#ff8a8a',
                            opacity: 0.7
                          }} />
                          <Typography 
                            variant="body2" 
                            sx={{ 
                              color: mode === 'light' ? '#ff6b6b' : '#ff8a8a',
                              fontWeight: 600,
                              transition: 'all 0.2s ease-in-out',
                              '&:hover': {
                                color: mode === 'light' ? '#e55a5a' : '#ff6b6b',
                                textDecoration: 'underline'
                              }
                            }}
                          >
                            #{issue.id}
                          </Typography>
                        </Box>
                        <Typography variant="body2" sx={{ 
                          color: getFixTimeCardTextColor(),
                          fontWeight: 500
                        }}>
                          {issue.title}
                        </Typography>
                      </Box>
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                        <Chip 
                          icon={getLevelIcon(issue.level)}
                          label={issue.level}
                          size="small"
                          sx={{ 
                            bgcolor: getLevelHexColor(issue.level),
                            color: 'white',
                            fontWeight: 600,
                            boxShadow: '0 2px 4px rgba(0,0,0,0.1)'
                          }}
                        />
                        <Typography variant="caption" sx={{ 
                          color: getMetricCardTextColor(),
                          fontWeight: 500
                        }}>
                          [{issue.count} events]
                        </Typography>
                      </Box>
                    </Box>
                  </Box>
                </Link>
              ))}
            </CardContent>
          </Card>
        </Grid>

        {/* Severity Distribution */}
        <Grid item xs={12} md={6}>
          <Card sx={{ 
            height: '100%', 
            background: getSeverityDistributionCardBackground(),
            color: getFixTimeCardTextColor(),
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
              background: mode === 'light' ? 'rgba(72, 149, 239, 0.1)' : 'rgba(255,255,255,0.1)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center'
            }}>
              <PieChartIcon sx={{ 
                fontSize: 40, 
                opacity: mode === 'light' ? 0.4 : 0.3,
                color: mode === 'light' ? '#4895ef' : 'white'
              }} />
            </Box>
            <CardContent sx={{ height: '100%', display: 'flex', flexDirection: 'column', position: 'relative', zIndex: 1 }}>
              <Typography variant="h6" gutterBottom sx={{ 
                display: 'flex', 
                alignItems: 'center',
                color: getFixTimeCardTextColor(),
                fontWeight: 600
              }}>
                <PieChartIcon sx={{ mr: 1 }} />
                Severity Distribution
              </Typography>
              <Box sx={{ flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={severityData}
                      cx="50%"
                      cy="50%"
                      labelLine={false}
                      label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                      outerRadius={80}
                      fill="#8884d8"
                      dataKey="value"
                    >
                      {severityData.map((entry, index) => (
                        <Cell key={`severity-cell-${entry.name}-${index}`} fill={entry.color} />
                      ))}
                    </Pie>
                  </PieChart>
                </ResponsiveContainer>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* Fix Time */}
        <Grid item xs={12} md={6}>
          <Card sx={{ 
            background: getFixTimeCardBackground(),
            color: getFixTimeCardTextColor(),
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
              background: mode === 'light' ? 'rgba(130, 82, 255, 0.1)' : 'rgba(255,255,255,0.1)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center'
            }}>
              <SpeedIcon sx={{ 
                fontSize: 40, 
                opacity: mode === 'light' ? 0.4 : 0.3,
                color: mode === 'light' ? theme.palette.primary.main : 'white'
              }} />
            </Box>
            <CardContent sx={{ position: 'relative', zIndex: 1 }}>
              <Typography variant="h6" gutterBottom sx={{ 
                display: 'flex', 
                alignItems: 'center',
                color: getFixTimeCardTextColor(),
                fontWeight: 600
              }}>
                <ScheduleIcon sx={{ mr: 1 }} />
                Fix Time Performance
              </Typography>
              
              <Grid container spacing={2}>
                {/* Average Fix Time */}
                <Grid item xs={12} sm={4}>
                  <Box sx={{ 
                    textAlign: 'center',
                    p: 2,
                    borderRadius: 2,
                    bgcolor: getMetricCardBackground(),
                    border: `1px solid ${getMetricCardBorderColor()}`
                  }}>
                    <Avatar sx={{ 
                      bgcolor: getFixTimeColor(releaseDetails.fix_time.avg),
                      width: 48,
                      height: 48,
                      mx: 'auto',
                      mb: 1
                    }}>
                      {getFixTimeIcon('avg')}
                    </Avatar>
                    <Typography variant="h5" sx={{ 
                      fontWeight: 'bold',
                      color: getFixTimeColor(releaseDetails.fix_time.avg),
                      mb: 0.5
                    }}>
                      {formatDuration(releaseDetails.fix_time.avg)}
                    </Typography>
                    <Typography variant="body2" sx={{ color: getMetricCardTextColor(), mb: 1 }}>
                      {getFixTimeLabel('avg')}
                    </Typography>
                    {(() => {
                      const rating = getSpeedRating(releaseDetails.fix_time.avg);
                      return (
                        <Box>
                          <Typography variant="caption" sx={{ color: rating.color, fontWeight: 600 }}>
                            {rating.rating}
                          </Typography>
                          <LinearProgress 
                            variant="determinate" 
                            value={rating.progress} 
                            sx={{ 
                              mt: 0.5,
                              height: 4,
                              borderRadius: 2,
                              bgcolor: getMetricCardBackground(),
                              '& .MuiLinearProgress-bar': {
                                bgcolor: rating.color
                              }
                            }} 
                          />
                        </Box>
                      );
                    })()}
                  </Box>
                </Grid>

                {/* Median Fix Time */}
                <Grid item xs={12} sm={4}>
                  <Box sx={{ 
                    textAlign: 'center',
                    p: 2,
                    borderRadius: 2,
                    bgcolor: getMetricCardBackground(),
                    border: `1px solid ${getMetricCardBorderColor()}`
                  }}>
                    <Avatar sx={{ 
                      bgcolor: getFixTimeColor(releaseDetails.fix_time.median),
                      width: 48,
                      height: 48,
                      mx: 'auto',
                      mb: 1
                    }}>
                      {getFixTimeIcon('median')}
                    </Avatar>
                    <Typography variant="h5" sx={{ 
                      fontWeight: 'bold',
                      color: getFixTimeColor(releaseDetails.fix_time.median),
                      mb: 0.5
                    }}>
                      {formatDuration(releaseDetails.fix_time.median)}
                    </Typography>
                    <Typography variant="body2" sx={{ color: getMetricCardTextColor(), mb: 1 }}>
                      {getFixTimeLabel('median')}
                    </Typography>
                    {(() => {
                      const rating = getSpeedRating(releaseDetails.fix_time.median);
                      return (
                        <Box>
                          <Typography variant="caption" sx={{ color: rating.color, fontWeight: 600 }}>
                            {rating.rating}
                          </Typography>
                          <LinearProgress 
                            variant="determinate" 
                            value={rating.progress} 
                            sx={{ 
                              mt: 0.5,
                              height: 4,
                              borderRadius: 2,
                              bgcolor: getMetricCardBackground(),
                              '& .MuiLinearProgress-bar': {
                                bgcolor: rating.color
                              }
                            }} 
                          />
                        </Box>
                      );
                    })()}
                  </Box>
                </Grid>

                {/* P95 Fix Time */}
                <Grid item xs={12} sm={4}>
                  <Box sx={{ 
                    textAlign: 'center',
                    p: 2,
                    borderRadius: 2,
                    bgcolor: getMetricCardBackground(),
                    border: `1px solid ${getMetricCardBorderColor()}`
                  }}>
                    <Avatar sx={{ 
                      bgcolor: getFixTimeColor(releaseDetails.fix_time.p95),
                      width: 48,
                      height: 48,
                      mx: 'auto',
                      mb: 1
                    }}>
                      {getFixTimeIcon('p95')}
                    </Avatar>
                    <Typography variant="h5" sx={{ 
                      fontWeight: 'bold',
                      color: getFixTimeColor(releaseDetails.fix_time.p95),
                      mb: 0.5
                    }}>
                      {formatDuration(releaseDetails.fix_time.p95)}
                    </Typography>
                    <Typography variant="body2" sx={{ color: getMetricCardTextColor(), mb: 1 }}>
                      {getFixTimeLabel('p95')}
                    </Typography>
                    {(() => {
                      const rating = getSpeedRating(releaseDetails.fix_time.p95);
                      return (
                        <Box>
                          <Typography variant="caption" sx={{ color: rating.color, fontWeight: 600 }}>
                            {rating.rating}
                          </Typography>
                          <LinearProgress 
                            variant="determinate" 
                            value={rating.progress} 
                            sx={{ 
                              mt: 0.5,
                              height: 4,
                              borderRadius: 2,
                              bgcolor: getMetricCardBackground(),
                              '& .MuiLinearProgress-bar': {
                                bgcolor: rating.color
                              }
                            }} 
                          />
                        </Box>
                      );
                    })()}
                  </Box>
                </Grid>
              </Grid>

              {/* Performance Summary */}
              <Box sx={{ 
                mt: 2, 
                p: 2, 
                borderRadius: 2, 
                bgcolor: getMetricCardBackground(),
                border: `1px solid ${getMetricCardBorderColor()}`
              }}>
                <Typography variant="subtitle2" sx={{ color: getMetricCardTextColor(), mb: 1 }}>
                  ⚡ Performance Summary
                </Typography>
                <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
                  {(() => {
                    const avgRating = getSpeedRating(releaseDetails.fix_time.avg);
                    const medianRating = getSpeedRating(releaseDetails.fix_time.median);
                    const p95Rating = getSpeedRating(releaseDetails.fix_time.p95);
                    
                    const overallRating = Math.round((avgRating.progress + medianRating.progress + p95Rating.progress) / 3);
                    let overallLabel = 'Excellent';
                    let overallColor = 'success.main';
                    
                    if (overallRating < 25) {
                      overallLabel = 'Critical';
                      overallColor = 'error.dark';
                    } else if (overallRating < 50) {
                      overallLabel = 'Poor';
                      overallColor = 'error.main';
                    } else if (overallRating < 75) {
                      overallLabel = 'Fair';
                      overallColor = 'warning.main';
                    } else if (overallRating < 90) {
                      overallLabel = 'Good';
                      overallColor = 'success.main';
                    }
                    
                    return (
                      <>
                        <Chip 
                          label={`Overall: ${overallLabel}`}
                          size="small"
                          sx={{ 
                            bgcolor: overallColor,
                            color: 'white',
                            fontWeight: 600
                          }}
                        />
                        <Chip 
                          label={`Avg: ${avgRating.rating}`}
                          size="small"
                          variant="outlined"
                          sx={{ 
                            borderColor: avgRating.color,
                            color: avgRating.color
                          }}
                        />
                        <Chip 
                          label={`Median: ${medianRating.rating}`}
                          size="small"
                          variant="outlined"
                          sx={{ 
                            borderColor: medianRating.color,
                            color: medianRating.color
                          }}
                        />
                        <Chip 
                          label={`P95: ${p95Rating.rating}`}
                          size="small"
                          variant="outlined"
                          sx={{ 
                            borderColor: p95Rating.color,
                            color: p95Rating.color
                          }}
                        />
                      </>
                    );
                  })()}
                </Box>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        {/* User Segments */}
        <Grid item xs={12} md={6}>
          <Card sx={{ 
            height: '100%', 
            background: getUserSegmentsCardBackground(),
            color: getFixTimeCardTextColor(),
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
              background: mode === 'light' ? 'rgba(46, 204, 113, 0.1)' : 'rgba(255,255,255,0.1)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center'
            }}>
              <GroupIcon sx={{ 
                fontSize: 40, 
                opacity: mode === 'light' ? 0.4 : 0.3,
                color: mode === 'light' ? '#2ecc71' : 'white'
              }} />
            </Box>
            <CardContent sx={{ height: '100%', display: 'flex', flexDirection: 'column', position: 'relative', zIndex: 1 }}>
              <Typography variant="h6" gutterBottom sx={{ 
                display: 'flex', 
                alignItems: 'center',
                color: getFixTimeCardTextColor(),
                fontWeight: 600
              }}>
                <GroupIcon sx={{ mr: 1 }} />
                User Segments
              </Typography>
              
              <Box sx={{ flex: 1, display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}>
                {/* Segments Grid Layout */}
                <Grid container spacing={2}>
                  {releaseDetails.segments.platform && (
                    <Grid item xs={12} sm={6}>
                      <Box>
                        <Typography variant="subtitle2" gutterBottom sx={{ 
                          color: getMetricCardTextColor(),
                          fontWeight: 600,
                          fontSize: '0.8rem'
                        }}>
                          Platforms:
                        </Typography>
                        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                          {Object.entries(releaseDetails.segments.platform).map(([platform, count]) => (
                            <Chip 
                              key={platform}
                              label={`${platform}: ${count}`}
                              size="small"
                              variant="outlined"
                              sx={{ 
                                borderColor: mode === 'light' ? 'rgba(46, 204, 113, 0.3)' : 'rgba(255, 255, 255, 0.3)',
                                color: mode === 'light' ? '#2ecc71' : 'white',
                                fontWeight: 500,
                                fontSize: '0.7rem',
                                height: 24
                              }}
                            />
                          ))}
                        </Box>
                      </Box>
                    </Grid>
                  )}

                  {releaseDetails.segments.os_name && (
                    <Grid item xs={12} sm={6}>
                      <Box>
                        <Typography variant="subtitle2" gutterBottom sx={{ 
                          color: getMetricCardTextColor(),
                          fontWeight: 600,
                          fontSize: '0.8rem'
                        }}>
                          OS:
                        </Typography>
                        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                          {Object.entries(releaseDetails.segments.os_name).map(([os, count]) => (
                            <Chip 
                              key={os}
                              label={`${os}: ${count}`}
                              size="small"
                              variant="outlined"
                              sx={{ 
                                borderColor: mode === 'light' ? 'rgba(46, 204, 113, 0.3)' : 'rgba(255, 255, 255, 0.3)',
                                color: mode === 'light' ? '#2ecc71' : 'white',
                                fontWeight: 500,
                                fontSize: '0.7rem',
                                height: 24
                              }}
                            />
                          ))}
                        </Box>
                      </Box>
                    </Grid>
                  )}

                  {releaseDetails.segments.browser_name && (
                    <Grid item xs={12} sm={6}>
                      <Box>
                        <Typography variant="subtitle2" gutterBottom sx={{ 
                          color: getMetricCardTextColor(),
                          fontWeight: 600,
                          fontSize: '0.8rem'
                        }}>
                          Browsers:
                        </Typography>
                        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                          {Object.entries(releaseDetails.segments.browser_name).map(([browser, count]) => (
                            <Chip 
                              key={browser}
                              label={`${browser}: ${count}`}
                              size="small"
                              variant="outlined"
                              sx={{ 
                                borderColor: mode === 'light' ? 'rgba(46, 204, 113, 0.3)' : 'rgba(255, 255, 255, 0.3)',
                                color: mode === 'light' ? '#2ecc71' : 'white',
                                fontWeight: 500,
                                fontSize: '0.7rem',
                                height: 24
                              }}
                            />
                          ))}
                        </Box>
                      </Box>
                    </Grid>
                  )}

                  {releaseDetails.segments.device_arch && (
                    <Grid item xs={12} sm={6}>
                      <Box>
                        <Typography variant="subtitle2" gutterBottom sx={{ 
                          color: getMetricCardTextColor(),
                          fontWeight: 600,
                          fontSize: '0.8rem'
                        }}>
                          Device Architecture:
                        </Typography>
                        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                          {Object.entries(releaseDetails.segments.device_arch).map(([arch, count]) => (
                            <Chip 
                              key={arch}
                              label={`${arch}: ${count}`}
                              size="small"
                              variant="outlined"
                              sx={{ 
                                borderColor: mode === 'light' ? 'rgba(46, 204, 113, 0.3)' : 'rgba(255, 255, 255, 0.3)',
                                color: mode === 'light' ? '#2ecc71' : 'white',
                                fontWeight: 500,
                                fontSize: '0.7rem',
                                height: 24
                              }}
                            />
                          ))}
                        </Box>
                      </Box>
                    </Grid>
                  )}

                  {releaseDetails.segments.runtime_name && (
                    <Grid item xs={12} sm={6}>
                      <Box>
                        <Typography variant="subtitle2" gutterBottom sx={{ 
                          color: getMetricCardTextColor(),
                          fontWeight: 600,
                          fontSize: '0.8rem'
                        }}>
                          Runtime:
                        </Typography>
                        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                          {Object.entries(releaseDetails.segments.runtime_name).map(([runtime, count]) => (
                            <Chip 
                              key={runtime}
                              label={`${runtime}: ${count}`}
                              size="small"
                              variant="outlined"
                              sx={{ 
                                borderColor: mode === 'light' ? 'rgba(46, 204, 113, 0.3)' : 'rgba(255, 255, 255, 0.3)',
                                color: mode === 'light' ? '#2ecc71' : 'white',
                                fontWeight: 500,
                                fontSize: '0.7rem',
                                height: 24
                              }}
                            />
                          ))}
                        </Box>
                      </Box>
                    </Grid>
                  )}
                </Grid>
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Paper>
  );
};

export default ReleaseDetails; 