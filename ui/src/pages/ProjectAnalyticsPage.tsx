import React, { useState } from 'react';
import { Box } from '@mui/material';
import { InsightsOutlined as AnalyticsIcon } from '@mui/icons-material';
import { useAuth } from '../auth/AuthContext';
import { Navigate, useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import apiClient from '../api/apiClient';
import Layout from '../components/Layout';
import PageHeader from '../components/PageHeader';
import { IssueLevel } from '../generated/api/client';
import {
  ReleasesTable,
  ReleaseDetails,
  ErrorsChart,
  ComparisonAlert,
  ReleasesDelta,
} from '../components/analytics/project/releases';

const ProjectAnalyticsPage: React.FC = () => {
  const { isAuthenticated } = useAuth();
  const { projectId } = useParams<{ projectId: string }>();
  const [selectedRelease, setSelectedRelease] = useState<string | null>(null);
  const [selectedLevel, setSelectedLevel] = useState<IssueLevel | undefined>(undefined);
  const [compareMode, setCompareMode] = useState<boolean>(false);
  const [compareRelease, setCompareRelease] = useState<string | null>(null);
  const [selectedTimeWindow, setSelectedTimeWindow] = useState('7d');

  // If not authenticated, redirect to login
  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  const projectIdNum = parseInt(projectId || '0');

  // Fetch releases analytics
  const { 
    data: releases, 
    isLoading: releasesLoading, 
    error: releasesError 
  } = useQuery({
    queryKey: ['projectReleasesAnalytics', projectIdNum],
    queryFn: async () => {
      const response = await apiClient.getProjectReleasesAnalytics(projectIdNum);
      return response.data;
    },
    enabled: !!projectIdNum
  });

  // Fetch comparison data when in compare mode
  const { 
    data: comparisonData, 
    isLoading: comparisonLoading, 
    error: comparisonError 
  } = useQuery({
    queryKey: ['compareProjectReleasesAnalytics', projectIdNum, selectedRelease, compareRelease],
    queryFn: async () => {
      if (!selectedRelease || !compareRelease) return null;
      const response = await apiClient.compareProjectReleasesAnalytics(projectIdNum, {
        base_version: compareRelease,
        target_version: selectedRelease
      });
      return response.data;
    },
    enabled: !!selectedRelease && !!compareRelease && !!projectIdNum && compareMode
  });

  // Fetch selected release details
  const { 
    data: releaseDetails, 
    isLoading: detailsLoading, 
    error: detailsError 
  } = useQuery({
    queryKey: ['projectReleaseAnalyticsDetails', projectIdNum, selectedRelease],
    queryFn: async () => {
      if (!selectedRelease) return null;
      const response = await apiClient.getProjectReleaseAnalyticsDetails(projectIdNum, selectedRelease);
      return response.data;
    },
    enabled: !!selectedRelease && !!projectIdNum
  });

  // Fetch errors timeseries for selected release
  const { 
    data: errorsTimeseries, 
    isLoading: timeseriesLoading 
  } = useQuery({
    queryKey: ['projectReleaseErrorsTimeseries', projectIdNum, selectedRelease, selectedLevel, selectedTimeWindow],
    queryFn: async () => {
      if (!selectedRelease) return [];
      console.log('Fetching timeseries for release:', selectedRelease, 'level:', selectedLevel, 'timeWindow:', selectedTimeWindow);
      
      // Determine granularity based on selected time window
      let granularity = '1d';
      if (selectedTimeWindow === '1d' || selectedTimeWindow === '3d' || selectedTimeWindow === '5d' || selectedTimeWindow === '7d' || selectedTimeWindow === '10d') {
        granularity = '1h';
      } else if (selectedTimeWindow === '14d' || selectedTimeWindow === '20d') {
        granularity = '4h';
      } else if (selectedTimeWindow === '25d' || selectedTimeWindow === '30d') {
        granularity = '6h';
      } else if (selectedTimeWindow === '45d' || selectedTimeWindow === '60d') {
        granularity = '12h';
      } else {
        granularity = '1d';
      }
      
      const response = await apiClient.getProjectReleaseErrorsTimeseries(
        projectIdNum, 
        selectedRelease, 
        selectedTimeWindow, 
        granularity,
        selectedLevel,
        'level'
      );
      console.log('Timeseries response:', response.data);
      return response.data;
    },
    enabled: !!selectedRelease && !!projectIdNum
  });

  // Fetch errors timeseries for comparison release
  const { 
    data: compareTimeseries, 
    isLoading: compareTimeseriesLoading 
  } = useQuery({
    queryKey: ['projectReleaseErrorsTimeseries', projectIdNum, compareRelease, selectedLevel, selectedTimeWindow],
    queryFn: async () => {
      if (!compareRelease) return [];
      console.log('Fetching timeseries for comparison release:', compareRelease, 'level:', selectedLevel, 'timeWindow:', selectedTimeWindow);
      
      // Determine granularity based on selected time window
      let granularity = '1d';
      if (selectedTimeWindow === '1d' || selectedTimeWindow === '3d' || selectedTimeWindow === '5d' || selectedTimeWindow === '7d' || selectedTimeWindow === '10d') {
        granularity = '1h';
      } else if (selectedTimeWindow === '14d' || selectedTimeWindow === '20d') {
        granularity = '4h';
      } else if (selectedTimeWindow === '25d' || selectedTimeWindow === '30d') {
        granularity = '6h';
      } else if (selectedTimeWindow === '45d' || selectedTimeWindow === '60d') {
        granularity = '12h';
      } else {
        granularity = '1d';
      }
      
      const response = await apiClient.getProjectReleaseErrorsTimeseries(
        projectIdNum, 
        compareRelease, 
        selectedTimeWindow, 
        granularity,
        selectedLevel,
        'level'
      );
      console.log('Comparison timeseries response:', response.data);
      return response.data;
    },
    enabled: !!compareRelease && !!projectIdNum && compareMode
  });

  const handleReleaseSelect = (version: string) => {
    setSelectedRelease(version);
  };

  const handleCompareRelease = (releaseVersion: string) => {
    if (compareMode && compareRelease === releaseVersion) {
      // Exit compare mode
      setCompareMode(false);
      setCompareRelease(null);
    } else {
      // Enter compare mode with selected release
      setCompareMode(true);
      setCompareRelease(releaseVersion);
      
      // If no release is selected, select the comparison release
      if (!selectedRelease) {
        setSelectedRelease(releaseVersion);
      }
    }
  };

  const handleLevelChange = (level: IssueLevel | undefined) => {
    setSelectedLevel(level);
  };

  const handleTimeWindowChange = (timeWindow: string) => {
    setSelectedTimeWindow(timeWindow);
  };

  const handleCompareModeToggle = () => {
    setCompareMode(false);
    setCompareRelease(null);
  };

  const handleSwitchComparison = () => {
    if (selectedRelease && compareRelease) {
      // Swap the releases
      setSelectedRelease(compareRelease);
      setCompareRelease(selectedRelease);
    }
  };

  return (
    <Layout>
      <PageHeader
        title="Project Analytics"
        subtitle="Project release analytics"
        icon={<AnalyticsIcon />}
        gradientVariant="green"
      />
      
      <Box sx={{ mt: 2 }}>
        {/* Comparison Alert */}
        <ComparisonAlert
          compareMode={compareMode}
          loading={comparisonLoading}
          error={comparisonError?.message || null}
          comparisonData={comparisonData}
          onSwitchComparison={handleSwitchComparison}
        />

        {/* Release Summary Table */}
        <ReleasesTable
          releases={releases || []}
          loading={releasesLoading}
          error={releasesError?.message || null}
          selectedRelease={selectedRelease}
          compareMode={compareMode}
          comparisonData={comparisonData}
          onReleaseSelect={handleReleaseSelect}
          onCompareRelease={handleCompareRelease}
        />

        {/* Release Details - only show when not in compare mode */}
        {!compareMode && (
          <ReleaseDetails
            releaseDetails={releaseDetails}
            loading={detailsLoading}
            error={detailsError?.message || null}
            compareMode={compareMode}
            comparisonData={comparisonData}
            projectId={projectId}
          />
        )}

        {/* Releases Delta - only show when in compare mode */}
        <ReleasesDelta
          comparisonData={comparisonData as any}
          loading={comparisonLoading}
          error={comparisonError?.message || null}
          compareMode={compareMode}
          compareRelease={compareRelease}
          selectedRelease={selectedRelease}
        />

        {/* Errors Chart */}
        <ErrorsChart
          selectedRelease={selectedRelease}
          selectedLevel={selectedLevel}
          selectedTimeWindow={selectedTimeWindow}
          errorsTimeseries={errorsTimeseries || []}
          compareTimeseries={compareTimeseries || []}
          compareMode={compareMode}
          compareRelease={compareRelease}
          comparisonData={comparisonData}
          loading={timeseriesLoading || compareTimeseriesLoading}
          onTimeWindowChange={handleTimeWindowChange}
          onLevelChange={handleLevelChange}
          onCompareModeToggle={handleCompareModeToggle}
        />
      </Box>
    </Layout>
  );
};

export default ProjectAnalyticsPage; 