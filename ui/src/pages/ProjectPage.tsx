import React, {useEffect, useMemo, useState} from 'react';
import {Box, Card, CardActionArea, CardContent, Chip, Grid, Paper, Tooltip, Typography, useTheme} from '@mui/material';
import type {SelectChangeEvent} from '@mui/material/Select';
import {useAuth} from '../auth/AuthContext';
import {Navigate, useNavigate, useParams} from 'react-router-dom';
import Layout from '../components/Layout';
import apiClient from '../api/apiClient';
import type {Issue, TimeseriesData} from '../generated/api/client';
import {IssueLevel, IssueStatus, IssueSortColumn, SortOrder} from '../generated/api/client';
import IssuesTrendsChart from '../components/issues/IssuesTrendsChart';
import IssuesFilterPanel from '../components/issues/IssuesFilterPanel';
import IssuesTabsList from '../components/issues/IssuesTabsList';
import {
  CheckCircle as CheckCircleIcon,
  Error as ErrorIcon,
  Info as InfoIcon,
  SettingsOutlined as SettingsIcon,
  Warning as WarningIcon,
} from '@mui/icons-material';
import {getLevelBadgeStyles, getLevelColor, getLevelHexColor} from "../utils/issues/issueUtils.ts";
import StatusBadge from "../components/StatusBadge.tsx";
import { usePermissions } from '../hooks/usePermissions';

// Interface for our project data
interface Project {
  id: number;
  name: string;
  description: string;
  team_id?: number | null;
  team_name?: string | null;
  created_at: string;
  issues_count: number;
  critical_count: number;
}

// Interface for our issue data with additional fields
interface ExtendedIssue extends Issue {
  projectName?: string;
  browser?: string;
  os?: string;
}

const ProjectPage: React.FC = () => {
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const { projectId } = useParams<{ projectId: string }>();
  const theme = useTheme();
  const { canManageProject } = usePermissions();
  // We only have one tab now, so we don't need to track the tab value
  const [timeRange, setTimeRange] = useState<string>(() => {
    return localStorage.getItem('timeRange') || '6h';
  });
  const [levelFilter, setLevelFilter] = useState<string>('all');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [timeseriesData, setTimeseriesData] = useState<TimeseriesData[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [issues, setIssues] = useState<ExtendedIssue[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState<number>(1);
  const [perPage] = useState<number>(20);
  const [totalCount, setTotalCount] = useState<number>(0);
  const [viewMode, setViewMode] = useState<string>(() => {
    // Get saved view mode from localStorage or default to 'medium'
    return localStorage.getItem('issueViewMode') || 'medium';
  });
  const [isChartExpanded, setIsChartExpanded] = useState<boolean>(true);
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [sortBy, setSortBy] = useState<IssueSortColumn | undefined>(undefined);
  const [sortOrder, setSortOrder] = useState<SortOrder>(SortOrder.Desc);
  const [tabValue, setTabValue] = useState<number>(0);

  // Initialize project data
  const [project, setProject] = useState<Project>({
    id: Number(projectId),
    name: `Project ${projectId}`,
    description: '',
    team_id: null,
    team_name: null,
    created_at: new Date().toISOString(),
    issues_count: 0,
    critical_count: 0
  });

  // State for most frequent issues
  const [mostFrequentIssues, setMostFrequentIssues] = useState<Array<any>>([]);

  // Fetch project details
  useEffect(() => {
    const fetchProjectDetails = async () => {
      try {
        // Use the getProject API to get project details
        const response = await apiClient.getProject(Number(projectId));

        // Update project details based on the API response
        setProject(prev => ({
          ...prev,
          name: response.data.project.name,
          description: response.data.project.description,
          team_id: response.data.project.team_id,
          team_name: response.data.project.team_name,
          created_at: response.data.project.created_at
        }));
      } catch (err) {
        console.error('Error fetching project details:', err);
      }
    };

    fetchProjectDetails();
  }, [projectId]);

  // Fetch project stats
  useEffect(() => {
    const fetchProjectStats = async () => {
      try {
        // Use the getProjectStats API to get project statistics
        const response = await apiClient.getProjectStats(Number(projectId), '7d');

        // Update project stats based on the API response
        setProject(prev => ({
          ...prev,
          issues_count: response.data.total_issues,
          critical_count: response.data.issues_by_level.fatal + response.data.issues_by_level.exception
        }));

        // Store most frequent issues
        setMostFrequentIssues(response.data.most_frequent_issues || []);
      } catch (err) {
        console.error('Error fetching project stats:', err);
      }
    };

    fetchProjectStats();
  }, [projectId]);

  // Helper function to get interval and granularity based on time range
  const getIntervalAndGranularity = (period: string): { interval: string, granularity: string } => {
    switch (period) {
      case '10m':
        return { interval: '10m', granularity: '1m' };
      case '30m':
        return { interval: '30m', granularity: '1m' };
      case '1h':
        return { interval: '1h', granularity: '5m' };
      case '3h':
        return { interval: '3h', granularity: '10m' };
      case '6h':
        return { interval: '6h', granularity: '30m' };
      case '12h':
        return { interval: '12h', granularity: '1h' };
      case '24h':
        return { interval: '24h', granularity: '1h' };
      case '3d':
        return { interval: '3d', granularity: '6h' };
      case '7d':
        return { interval: '7d', granularity: '12h' };
      case '14d':
        return { interval: '14d', granularity: '1d' };
      case '30d':
        return { interval: '30d', granularity: '1d' };
      default:
        return { interval: '6h', granularity: '30m' };
    }
  };

  // Fetch timeseries data
  useEffect(() => {
    const fetchTimeseriesData = async () => {
      try {
        const { interval, granularity } = getIntervalAndGranularity(timeRange);
        const response = await apiClient.getEventsTimeseries(interval, granularity, Number(projectId));
        setTimeseriesData(response.data);
      } catch (err) {
        console.error('Error fetching timeseries data:', err);
      }
    };

    fetchTimeseriesData();
  }, [timeRange, projectId]); // Re-fetch when time range or project ID changes

  // Fetch issues from API
  useEffect(() => {
    const fetchIssues = async () => {
      setIsLoading(true);
      setError(null);
      try {
        // Convert levelFilter to the enum type expected by the API
        const apiLevel = levelFilter !== 'all' ? levelFilter as IssueLevel : undefined;

        // Convert statusFilter to the enum type expected by the API
        const apiStatus = statusFilter !== 'all' ? statusFilter as IssueStatus : undefined;

        // Call the listIssues method from the API client with projectId
        const response = await apiClient.listIssues(perPage, page, apiLevel, apiStatus, Number(projectId), sortBy, sortOrder);

        // Transform the response data to match our ExtendedIssue interface
        const fetchedIssues = response.data.issues.map(issue => ({
          ...issue,
          // Add any additional fields needed
          projectName: issue.project_name || `Project ${issue.project_id}`
        }));

        setIssues(fetchedIssues);
        setTotalCount(response.data.total);

        // Update project stats based on issues
        setProject(prev => ({
          ...prev,
          issues_count: response.data.total,
          critical_count: fetchedIssues.filter(issue => issue.level === 'fatal' || issue.level === 'exception').length
        }));
      } catch (err) {
        console.error('Error fetching issues:', err);
        setError('Failed to fetch issues. Please try again later.');
      } finally {
        setIsLoading(false);
      }
    };

    fetchIssues();
  }, [projectId, levelFilter, statusFilter, page, perPage, sortBy, sortOrder]); // Re-fetch when projectId, level filter, status filter, page, perPage, sortBy, or sortOrder changes

  // Function to refresh issues data
  const refreshIssues = () => {
    // Re-fetch issues with current filter
    const fetchIssues = async () => {
      setIsLoading(true);
      setError(null);
      try {
        // Convert levelFilter to the enum type expected by the API
        const apiLevel = levelFilter !== 'all' ? levelFilter as IssueLevel : undefined;

        // Convert statusFilter to the enum type expected by the API
        const apiStatus = statusFilter !== 'all' ? statusFilter as IssueStatus : undefined;

        // Call the listIssues method from the API client with projectId
        const response = await apiClient.listIssues(perPage, page, apiLevel, apiStatus, Number(projectId), sortBy, sortOrder);

        // Transform the response data to match our ExtendedIssue interface
        const fetchedIssues = response.data.issues.map(issue => ({
          ...issue,
          // Add any additional fields needed
          projectName: issue.project_name || `Project ${issue.project_id}`
        }));

        setIssues(fetchedIssues);
        setTotalCount(response.data.total);
      } catch (err) {
        console.error('Error fetching issues:', err);
        setError('Failed to fetch issues. Please try again later.');
      } finally {
        setIsLoading(false);
      }
    };

    // Re-fetch project stats
    const fetchProjectStats = async () => {
      try {
        // Use the getProjectStats API to get project statistics
        const response = await apiClient.getProjectStats(Number(projectId), '6h');

        // Update project stats based on the API response
        setProject(prev => ({
          ...prev,
          issues_count: response.data.total_issues,
          critical_count: response.data.issues_by_level.fatal + response.data.issues_by_level.exception
        }));

        // Store most frequent issues
        setMostFrequentIssues(response.data.most_frequent_issues || []);
      } catch (err) {
        console.error('Error fetching project stats:', err);
      }
    };

    // Re-fetch timeseries data
    const fetchTimeseriesData = async () => {
      try {
        const { interval, granularity } = getIntervalAndGranularity(timeRange);
        const response = await apiClient.getEventsTimeseries(interval, granularity, Number(projectId));
        setTimeseriesData(response.data);
      } catch (err) {
        console.error('Error fetching timeseries data:', err);
      }
    };

    fetchIssues();
    fetchProjectStats();
    fetchTimeseriesData();
  };

  // Handle pagination change
  const handlePageChange = (_event: React.ChangeEvent<unknown>, value: number) => {
    setPage(value);
  };

  // Transform timeseries data for the chart
  const transformTimeseriesData = () => {
    if (!timeseriesData.length) return [];

    // Group by level
    const levelGroups: Record<string, TimeseriesData> = {};

    timeseriesData.forEach(series => {
      levelGroups[series.name] = series;
    });

    // Create chart data
    const chartData = [];
    const occurrencesLength = timeseriesData.length > 0 ? timeseriesData[0].occurrences.length : 0;

    // Generate appropriate time labels based on the selected time range
    const generateTimeLabel = (index: number) => {
      // The server returns data with appropriate granularity based on the period
      // For short periods (minutes), we show time labels
      // For medium periods (hours), we show hour labels
      // For long periods (days), we show date labels

      if (timeRange === '10m' || timeRange === '30m') {
        // For 10m and 30m, show minute labels (e.g., "10:05")
        const now = new Date();
        const minutesAgo = (occurrencesLength - 1 - index) * (timeRange === '10m' ? 1 : 3);
        const date = new Date(now.getTime() - minutesAgo * 60000);
        return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
      } else if (timeRange === '1h' || timeRange === '3h' || 
                timeRange === '6h' || timeRange === '12h' || 
                timeRange === '24h') {
        // For 1h to 24h, show hour labels (e.g., "10:00")
        const now = new Date();
        const hoursAgo = (occurrencesLength - 1 - index);
        const date = new Date(now.getTime() - hoursAgo * 3600000);
        return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
      } else {
        // For 3d, 7d, 14d, 30d, show date labels (e.g., "Jun 15")
        const now = new Date();
        const daysAgo = (occurrencesLength - 1 - index);
        const date = new Date(now.getTime() - daysAgo * 86400000);
        return date.toLocaleDateString([], { month: 'short', day: 'numeric' });
      }
    };

    for (let i = 0; i < occurrencesLength; i++) {
      const dataPoint: Record<string, any> = { 
        index: i,
        timeLabel: generateTimeLabel(i)
      };

      // Add counts for each level
      Object.keys(levelGroups).forEach(level => {
        // Use 0 for zero values to ensure lines touch the x-axis
        dataPoint[level] = levelGroups[level].occurrences[i];
      });

      chartData.push(dataPoint);
    }

    return chartData;
  };
  transformTimeseriesData();

  // If not authenticated, redirect to log in
  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  // We only have one tab now, so we don't need the handleTabChange function

  const handleTimeRangeChange = (event: SelectChangeEvent) => {
    const newTimeRange = event.target.value;
    setTimeRange(newTimeRange);
    localStorage.setItem('timeRange', newTimeRange);
  };

  const handleLevelFilterChange = (event: SelectChangeEvent) => {
    setLevelFilter(event.target.value);
  };

  const handleStatusFilterChange = (event: SelectChangeEvent) => {
    setStatusFilter(event.target.value);
  };

  const toggleChartExpanded = () => {
    setIsChartExpanded(!isChartExpanded);
  };

  const handleViewModeChange = (_event: React.MouseEvent<HTMLElement>, newViewMode: string) => {
    if (newViewMode !== null) {
      setViewMode(newViewMode);
      // Save to localStorage
      localStorage.setItem('issueViewMode', newViewMode);
    }
  };
  
  const handleSortByChange = (event: SelectChangeEvent) => {
    setSortBy(event.target.value as IssueSortColumn);
  };
  
  const handleSortOrderChange = (event: SelectChangeEvent) => {
    setSortOrder(event.target.value as SortOrder);
  };
  
  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
    // Set level filter based on tab
    switch (newValue) {
      case 0: // All
        setLevelFilter('all');
        break;
      case 1: // Fatal
        setLevelFilter('fatal');
        break;
      case 2: // Errors
        setLevelFilter('error');
        break;
      case 3: // Warnings
        setLevelFilter('warning');
        break;
      case 4: // Info
        setLevelFilter('info');
        break;
      case 5: // Exceptions
        setLevelFilter('exception');
        break;
      case 6: // Debug
        setLevelFilter('debug');
        break;
    }
  };

  const handleIssueClick = (issueId: number) => {
    navigate(`/projects/${projectId}/issues/${issueId}`);
  };

  const handleProjectIssueClick = (_projectId: number, issueId: number) => {
    handleIssueClick(issueId);
  }

  const getLevelIcon = (level: string) => {
    switch (level) {
      case 'fatal':
        return <ErrorIcon fontSize="small" sx={{ mr: 1, color: getLevelHexColor('fatal') }} />;
      case 'error':
        return <ErrorIcon fontSize="small" sx={{ mr: 1, color: getLevelHexColor('error') }} />;
      case 'exception':
        return <ErrorIcon fontSize="small" sx={{ mr: 1, color: getLevelHexColor('exception') }} />;
      case 'warning':
        return <WarningIcon fontSize="small" sx={{ mr: 1, color: getLevelHexColor('warning') }} />;
      case 'info':
        return <InfoIcon fontSize="small" sx={{ mr: 1, color: getLevelHexColor('info') }} />;
      case 'debug':
        return <InfoIcon fontSize="small" sx={{ mr: 1, color: getLevelHexColor('debug') }} />;
      default:
        return <ErrorIcon fontSize="small" sx={{ mr: 1, color: getLevelHexColor('error') }} />;
    }
  };

  // Memoized filtered issues - moved to top level to avoid conditional hook calls
  const filteredIssues = useMemo(() => {
    return issues.filter(issue =>
      searchQuery === '' ||
      (issue.title && issue.title.toLowerCase().includes(searchQuery.toLowerCase())) ||
      (issue.message && issue.message.toLowerCase().includes(searchQuery.toLowerCase())) ||
      (issue.projectName && issue.projectName.toLowerCase().includes(searchQuery.toLowerCase()))
    );
  }, [issues, searchQuery]);

  return (
    <Layout showBackButton={true} backTo="/dashboard">
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Box>
          <Typography variant="h6" component="h1" gutterBottom className="gradient-text">
            {project.name}
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5, whiteSpace: 'pre-line' }}>
            {project.description}
          </Typography>
        </Box>
        <Box sx={{ display: 'flex', gap: 0.5, alignItems: 'center' }}>
          {project.team_name && (
            <Chip 
              label={`Team: ${project.team_name}`} 
              color="secondary" 
              variant="outlined"
              size="small"
            />
          )}
          <Chip 
            label={`${project.issues_count} issues`} 
            color="primary" 
            variant="outlined"
            size="small"
          />
          <Chip 
            label={`${project.critical_count} critical`} 
            color="error" 
            variant="outlined"
            size="small"
          />
          <Tooltip title={canManageProject(Number(projectId), project.team_id || undefined) ? "Project Settings" : "You don't have permission to manage this project"}>
            {canManageProject(Number(projectId), project.team_id || undefined) ? (
              <Box
                component="button"
                onClick={() => navigate(`/projects/${projectId}/settings`)}
                sx={{ 
                  ml: 1,
                  display: 'flex',
                  alignItems: 'center',
                  border: 'none',
                  borderRadius: 1,
                  padding: '4px 8px',
                  cursor: 'pointer',
                  fontSize: '0.875rem',
                  bgcolor: theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.05)' : 'rgba(0, 0, 0, 0.04)',
                  color: 'text.primary',
                  '&:hover': {
                    bgcolor: theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.08)',
                  }
                }}
              >
                <Typography variant="body2" sx={{ display: 'flex', alignItems: 'center' }}>
                  Settings <SettingsIcon fontSize="small" sx={{ ml: 0.5 }} />
                </Typography>
              </Box>
            ) : <div/>}
          </Tooltip>
        </Box>
      </Box>

      <Paper sx={{ mb: 2 }}>
        <IssuesTrendsChart
          data={transformTimeseriesData()}
          timeRange={timeRange}
          onTimeRangeChange={handleTimeRangeChange}
          isChartExpanded={isChartExpanded}
          onToggleChartExpanded={toggleChartExpanded}
          chartTitle="Project Events Over Time"
        />
      </Paper>

      {mostFrequentIssues.length > 0 && (
        <Paper sx={{ mb: 3 }}>
          <Box sx={{ p: 3 }}>
            <Typography variant="h6" sx={{ mb: 2 }} className="gradient-subtitle">
              Most Frequent Issues
            </Typography>
            <Grid container spacing={2}>
              {mostFrequentIssues.map((issue) => (
                <Grid item xs={12} md={6} key={issue.id}>
                  <Card 
                    sx={{ 
                      '&:hover': { 
                        boxShadow: 3,
                        '& .MuiCardActionArea-focusHighlight': {
                          opacity: 0.1
                        }
                      } 
                    }}
                  >
                    <CardActionArea onClick={() => handleIssueClick(issue.id)}>
                      <CardContent>
                        <Box sx={{ display: 'flex', alignItems: 'flex-start', gap: 1 }}>
                          {getLevelIcon(issue.level)}
                          <Box sx={{ flexGrow: 1 }}>
                            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                              <Typography variant="subtitle1" component="div" sx={{ fontWeight: 'bold', mb: 0.5 }}>
                                {issue.title}
                              </Typography>
                              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                <StatusBadge status={IssueStatus.Unresolved} size="small" />
                              </Box>
                            </Box>
                            <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                              {issue.message}
                            </Typography>
                            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                              <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
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
                                />
                              </Box>
                              {issue.status === IssueStatus.Resolved && issue.resolved_at && (
                                <Typography 
                                  variant="caption" 
                                  color="success.main" 
                                  sx={{ 
                                    display: 'flex',
                                    alignItems: 'center',
                                    gap: 0.5
                                  }}
                                >
                                  <CheckCircleIcon fontSize="small" />
                                  {new Date(issue.resolved_at).toLocaleString()}
                                  {issue.resolved_by && ` by ${issue.resolved_by}`}
                                </Typography>
                              )}
                            </Box>
                          </Box>
                        </Box>
                      </CardContent>
                    </CardActionArea>
                  </Card>
                </Grid>
              ))}
            </Grid>
          </Box>
        </Paper>
      )}

      <Paper>
        <IssuesFilterPanel
          searchQuery={searchQuery}
          onSearchChange={setSearchQuery}
          levelFilter={levelFilter}
          onLevelFilterChange={handleLevelFilterChange}
          statusFilter={statusFilter}
          onStatusFilterChange={handleStatusFilterChange}
          viewMode={viewMode}
          onViewModeChange={handleViewModeChange}
          onRefresh={refreshIssues}
          isLoading={isLoading}
          sortBy={sortBy}
          onSortByChange={handleSortByChange}
          sortOrder={sortOrder}
          onSortOrderChange={handleSortOrderChange}
        />
        <IssuesTabsList
          issues={issues}
          filteredIssues={filteredIssues}
          isLoading={isLoading}
          error={error}
          viewMode={viewMode as 'compact' | 'medium' | 'large'}
          onIssueClick={handleProjectIssueClick}
          onRefresh={refreshIssues}
          totalCount={totalCount}
          page={page}
          perPage={perPage}
          onPageChange={handlePageChange}
          tabValue={tabValue}
          onTabChange={handleTabChange}
        />
      </Paper>
    </Layout>
  );
};

export default ProjectPage;
