import React, { useState, useEffect } from 'react';
import {Paper, useTheme, useMediaQuery, Divider} from '@mui/material';
import type { SelectChangeEvent } from '@mui/material/Select';
import { useAuth } from '../auth/AuthContext';
import { Navigate, useNavigate } from 'react-router-dom';
import Layout from '../components/Layout';
import PageHeader from '../components/PageHeader';
import apiClient from '../api/apiClient';
import { IssueLevel, IssueStatus, IssueSortColumn, SortOrder } from '../generated/api/client';
import type { Issue, TimeseriesData } from '../generated/api/client';
import IssuesTrendsChart from '../components/issues/IssuesTrendsChart';
import IssuesFilterPanel from '../components/issues/IssuesFilterPanel';
import IssuesTabsList from '../components/issues/IssuesTabsList';
import {
  BugReportOutlined as IssuesIcon,
} from '@mui/icons-material';

// Interface for our issue data with additional fields
interface ExtendedIssue extends Issue {
  projectName?: string;
}

// Tab Panel component
const IssuesPage: React.FC = () => {
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  useMediaQuery(theme.breakpoints.down('lg'));
  const [levelFilter, setLevelFilter] = useState<string>('all');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [tabValue, setTabValue] = useState(0);
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
  const [timeRange, setTimeRange] = useState<string>(() => {
    return localStorage.getItem('timeRange') || '6h';
  });
  const [timeseriesData, setTimeseriesData] = useState<TimeseriesData[]>([]);
  const [isChartExpanded, setIsChartExpanded] = useState<boolean>(!isMobile);
  const [sortBy, setSortBy] = useState<IssueSortColumn | undefined>(undefined);
  const [sortOrder, setSortOrder] = useState<SortOrder>(SortOrder.Desc);

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
        return { interval: '3h', granularity: '15m' };
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
        const response = await apiClient.getIssuesTimeseries(interval, granularity);
        setTimeseriesData(response.data);
      } catch (err) {
        console.error('Error fetching timeseries data:', err);
      }
    };

    fetchTimeseriesData();
  }, [timeRange]); // Re-fetch when time range changes

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

        // Call the listIssues method from the API client
        const response = await apiClient.listIssues(perPage, page, apiLevel, apiStatus, undefined, sortBy, sortOrder);

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

    fetchIssues();
  }, [levelFilter, statusFilter, page, perPage, sortBy, sortOrder]); // Re-fetch when level filter, status filter, page, perPage, sortBy, or sortOrder changes

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
      const dataPoint: Record<string, string | number> = { 
        index: i,
        timeLabel: generateTimeLabel(i)
      };

      // Add counts for each level
      Object.keys(levelGroups).forEach(level => {
        dataPoint[level] = levelGroups[level].occurrences[i];
      });

      chartData.push(dataPoint);
    }

    return chartData;
  };

  // Filter issues based on a search query
  const filteredIssues = issues
    .filter(issue => 
      searchQuery === '' || 
      (issue.title && issue.title.toLowerCase().includes(searchQuery.toLowerCase())) ||
      (issue.message && issue.message.toLowerCase().includes(searchQuery.toLowerCase())) ||
      (issue.projectName && issue.projectName.toLowerCase().includes(searchQuery.toLowerCase()))
    );

  // If not authenticated, redirect to login
  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  const handleLevelFilterChange = (event: SelectChangeEvent) => {
    setLevelFilter(event.target.value);
  };

  const handleStatusFilterChange = (event: SelectChangeEvent) => {
    setStatusFilter(event.target.value);
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

  const handleViewModeChange = (_event: React.MouseEvent<HTMLElement>, newViewMode: string) => {
    if (newViewMode !== null) {
      setViewMode(newViewMode);
      // Save to localStorage
      localStorage.setItem('issueViewMode', newViewMode);
    }
  };

  const handleTimeRangeChange = (event: SelectChangeEvent) => {
    const newTimeRange = event.target.value;
    setTimeRange(newTimeRange);
    localStorage.setItem('timeRange', newTimeRange);
  };

  const toggleChartExpanded = () => {
    setIsChartExpanded(!isChartExpanded);
  };
  
  const handleSortByChange = (event: SelectChangeEvent) => {
    setSortBy(event.target.value as IssueSortColumn);
  };
  
  const handleSortOrderChange = (event: SelectChangeEvent) => {
    setSortOrder(event.target.value as SortOrder);
  };

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

        // Call the listIssues method from the API client
        const response = await apiClient.listIssues(perPage, page, apiLevel, apiStatus, undefined, sortBy, sortOrder);

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

    // Re-fetch timeseries data
    const fetchTimeseriesData = async () => {
      try {
        const { interval, granularity } = getIntervalAndGranularity(timeRange);
        const response = await apiClient.getIssuesTimeseries(interval, granularity);
        setTimeseriesData(response.data);
      } catch (err) {
        console.error('Error fetching timeseries data:', err);
      }
    };

    fetchIssues();
    fetchTimeseriesData();
  };

  // Handle pagination change
  const handlePageChange = (_event: React.ChangeEvent<unknown>, value: number) => {
    setPage(value);
  };

  const handleIssueClick = (projectId: number, issueId: number) => {
    navigate(`/projects/${projectId}/issues/${issueId}`);
  };

  return (
    <Layout>
      <PageHeader
        title="Issues"
        subtitle="View and manage issues across all your projects. Track errors, warnings, and information messages to improve your application's stability."
        icon={<IssuesIcon />}
      />

      <Paper sx={{ mb: 2 }}>
        <IssuesTrendsChart
          data={transformTimeseriesData()}
          timeRange={timeRange}
          onTimeRangeChange={handleTimeRangeChange}
          isChartExpanded={isChartExpanded}
          onToggleChartExpanded={toggleChartExpanded}
          chartTitle="Issues Over Time"
        />

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

        <Divider />

        <IssuesTabsList
          issues={issues}
          filteredIssues={filteredIssues}
          isLoading={isLoading}
          error={error}
          viewMode={viewMode as 'compact' | 'medium' | 'large'}
          onIssueClick={handleIssueClick}
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

export default IssuesPage;
