import React, { useState, useEffect, useMemo, useCallback } from 'react';
import { Box, Paper, Tabs, Tab, Alert } from '@mui/material';
import { Navigate, useParams } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext';
import Layout from '../components/Layout';
import apiClient from '../api/apiClient';
import { IssueStatus, IssueSource, type IssueResponse, type TimeseriesData } from '../generated/api/client/api';

// Import components
import IssueHeader from '../components/issues/IssueHeader';
import IssueStats from '../components/issues/IssueStats';
import IssueChart from '../components/issues/IssueChart';
import TabPanel from '../components/issues/TabPanel';
import EventsTab from '../components/issues/EventsTab';
import StackTraceTab from '../components/issues/StackTraceTab';
import LoadingState from '../components/issues/LoadingState';
import ErrorState from '../components/issues/ErrorState';

// Import utilities
import { getGranularity, transformTimeseriesData } from '../utils/issues/issueUtils';
import type {SelectChangeEvent} from "@mui/material/Select";

// Main component
const IssuePage: React.FC = () => {
  const { isAuthenticated } = useAuth();
  const { projectId, issueId } = useParams<{ projectId: string, issueId: string }>();
  const [tabValue, setTabValue] = useState(0);
  const [timeRange, setTimeRange] = useState<string>(() => {
    return localStorage.getItem('timeRange') || '6h';
  });
  const [loading, setLoading] = useState(true);
  const [statusChangeLoading, setStatusChangeLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [statusChangeError, setStatusChangeError] = useState<string | null>(null);
  const [issueData, setIssueData] = useState<IssueResponse | null>(null);
  const [timeseriesData, setTimeseriesData] = useState<TimeseriesData[] | null>(null);

  // Fetch issue data
  useEffect(() => {
    const fetchIssue = async () => {
      setLoading(true);
      setError(null);
      try {
        const response = await apiClient.getIssue(Number(projectId), Number(issueId));
        setIssueData(response.data);
      } catch (err) {
        console.error('Error fetching issue:', err);
        setError('Failed to fetch issue data. Please try again later.');
      } finally {
        setLoading(false);
      }
    };

    fetchIssue();
  }, [projectId, issueId]);

  // Fetch timeseries data
  useEffect(() => {
    const fetchTimeseries = async () => {
      if (!projectId || !issueId || !issueData?.issue?.source) return;

      try {
        // Get appropriate granularity based on timeRange
        const granularity = getGranularity(timeRange);

        // Use the same endpoint for both event and exception sources
        const response = await apiClient.getProjectIssueEventsTimeseries(
          Number(projectId), 
          Number(issueId), 
          timeRange, 
          granularity
        );
        setTimeseriesData(response.data);
      } catch (err) {
        console.error('Error fetching timeseries data:', err);
        // We don't set the main error state here to avoid blocking the whole page
      }
    };

    fetchTimeseries();
  }, [projectId, issueId, timeRange, issueData?.issue?.source]);

  // Derived state with useMemo
  const issue = useMemo(() => issueData?.issue, [issueData]);
  const issueSource = useMemo(() => issue?.source, [issue]);

  // Transform timeseries data for the chart
  const chartData = useMemo(() => 
    transformTimeseriesData(timeseriesData, issueSource), 
    [timeseriesData, issueSource]);

  // Get tags based on issue source
  const tags = useMemo(() => {
    if (issueData?.events && issueData.events.length > 0) {
      return issueData.events[0].tags;
    }
    return undefined;
  }, [issueData]);

  // Get stack trace for event source (if available)
  const stackTrace = useMemo(() => {
    if (issueData?.events && issueData.events.length > 0) {
      const firstEvent = issueData.events[0];
      // Only use exception_stacktrace field
      return firstEvent.exception_stacktrace;
    }
    return undefined;
  }, [issueData]);

  // Function to handle status changes
  const handleStatusChange = useCallback(async (newStatus: IssueStatus) => {
    if (!projectId || !issueId || !issue) return;

    setStatusChangeLoading(true);
    setStatusChangeError(null);

    try {
      // Call the API to change the issue status
      await apiClient.changeIssueStatus(
        Number(projectId),
        Number(issueId),
        { status: newStatus }
      );

      // Refresh the issue data
      const response = await apiClient.getIssue(Number(projectId), Number(issueId));
      setIssueData(response.data);
    } catch (err) {
      console.error('Error changing issue status:', err);
      setStatusChangeError('Failed to change issue status. Please try again later.');
    } finally {
      setStatusChangeLoading(false);
    }
  }, [projectId, issueId, issue]);

  // Event handlers with useCallback
  const handleTabChange = useCallback((_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  }, []);

  const handleTimeRangeChange = useCallback((event: SelectChangeEvent) => {
    const newTimeRange = event.target.value as string;
    setTimeRange(newTimeRange);
    localStorage.setItem('timeRange', newTimeRange);
  }, []);

  // Show loading state
  if (loading) {
    return (
      <Layout showBackButton={true} backTo={`/projects/${projectId}`}>
        <LoadingState />
      </Layout>
    );
  }

  // Show error state
  if (error) {
    return (
      <Layout showBackButton={true} backTo={`/projects/${projectId}`}>
        <ErrorState message={error} />
      </Layout>
    );
  }

  // If no issue data is available, show a message
  if (!issue) {
    return (
      <Layout showBackButton={true} backTo={`/projects/${projectId}`}>
        <ErrorState message="Issue data not available. Please try again later." />
      </Layout>
    );
  }

  // If not authenticated, redirect to login
  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  return (
    <Layout showBackButton={true} backTo={`/projects/${projectId}`}>
      {statusChangeError && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {statusChangeError}
        </Alert>
      )}

      <IssueHeader 
        issue={issue} 
        tags={tags} 
        onStatusChange={handleStatusChange}
        statusChangeLoading={statusChangeLoading}
      />

      <IssueStats issue={issue} issueData={issueData} />

      <IssueChart 
        timeRange={timeRange} 
        onTimeRangeChange={handleTimeRangeChange}
        chartData={chartData}
        issueSource={issue.source}
      />

      <Paper>
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs value={tabValue} onChange={handleTabChange} aria-label="issue tabs">
            <Tab label="Occurrences" />
            <Tab label={issue.source === IssueSource.Exception ? "Stack Trace" : "Details"} />
          </Tabs>
        </Box>

        <TabPanel value={tabValue} index={0}>
          <EventsTab 
            issueData={issueData} 
            issue={issue} 
          />
        </TabPanel>

        <TabPanel value={tabValue} index={1}>
          <StackTraceTab 
            issueData={issueData} 
            issue={issue}
            stackTrace={stackTrace || undefined}
          />
        </TabPanel>
      </Paper>
    </Layout>
  );
};

export default IssuePage;
