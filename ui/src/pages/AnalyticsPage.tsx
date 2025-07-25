import React from 'react';
import { 
  Box, 
  Typography,
  Paper,
  List,
  ListItem,
  Chip,
  CircularProgress,
  IconButton,
  Tooltip,
  Divider,
  Grid,
  Card,
  CardContent
} from '@mui/material';
import { 
  ArrowForward as ArrowForwardIcon,
  InsightsOutlined as AnalyticsIcon,
  TrendingUp as TrendingUpIcon,
  BugReport as BugReportIcon,
  Warning as WarningIcon,
  CheckCircle as CheckCircleIcon
} from '@mui/icons-material';
import { useAuth } from '../auth/AuthContext';
import { Navigate, useNavigate } from 'react-router-dom';
import { useQuery, useQueries } from '@tanstack/react-query';
import apiClient from '../api/apiClient';
import Layout from '../components/Layout';
import PageHeader from '../components/PageHeader';
import Breadcrumbs from '../components/Breadcrumbs';
import { GetProjectStatsPeriodEnum } from '../generated/api/client';

const AnalyticsPage: React.FC = () => {
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();

  // Fetch projects using React Query
  const { data: projects, isLoading: projectsLoading, error: projectsError } = useQuery({
    queryKey: ['projects'],
    queryFn: async () => {
      const response = await apiClient.listProjects();
      return response.data;
    }
  });

  // Fetch stats for each project
  const projectStatsQueries = useQueries({
    queries: (projects || []).map(project => ({
      queryKey: ['projectStats', project.id],
      queryFn: async () => {
        const response = await apiClient.getProjectStats(project.id, GetProjectStatsPeriodEnum._7d);
        return response.data;
      },
      enabled: !!projects
    }))
  });

  // Combine projects with their stats
  const projectsWithStats = React.useMemo(() => {
    if (!projects) return [];

    return projects.map((project, index) => {
      const statsQuery = projectStatsQueries[index];
      return {
        ...project,
        stats: statsQuery.data,
        statsLoading: statsQuery.isLoading,
        statsError: statsQuery.error
      };
    });
  }, [projects, projectStatsQueries]);

  // Calculate overall analytics summary
  const analyticsSummary = React.useMemo(() => {
    if (!projectsWithStats.length) return null;

    const totalIssues = projectsWithStats.reduce((sum, project) => {
      return sum + (project.stats?.total_issues || 0);
    }, 0);

    const criticalIssues = projectsWithStats.reduce((sum, project) => {
      if (!project.stats?.issues_by_level) return sum;
      return sum + (project.stats.issues_by_level.fatal || 0) + (project.stats.issues_by_level.exception || 0);
    }, 0);

    const resolvedIssues = projectsWithStats.reduce((sum, project) => {
      if (!project.stats?.issues_by_level) return sum;
      return sum + (project.stats.issues_by_level.info || 0);
    }, 0);

    return {
      totalProjects: projectsWithStats.length,
      totalIssues,
      criticalIssues,
      resolvedIssues
    };
  }, [projectsWithStats]);

  const isLoading = projectsLoading || projectStatsQueries.some(query => query.isLoading);
  const error = projectsError || projectStatsQueries.some(query => query.error);

  // If not authenticated, redirect to login
  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  const handleProjectClick = (projectId: number) => {
    navigate(`/analytics/projects/${projectId}`);
  };

  return (
    <Layout>
      <PageHeader
        title="Analytics"
        subtitle="Analytics for project releases and events."
        icon={<AnalyticsIcon />}
        gradientVariant="green"
      />
      <Breadcrumbs />
      
      {/* Analytics Summary Panel */}
      <Paper sx={{ p: 3, mb: 3 }}>
        <Typography variant="h6" gutterBottom>
          Overall Summary (Last 7 Days)
        </Typography>
        {isLoading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 2 }}>
            <CircularProgress />
          </Box>
        ) : analyticsSummary ? (
          <Grid container spacing={3}>
            <Grid item xs={12} sm={6} md={3}>
              <Card sx={{ 
                bgcolor: (theme) => theme.palette.mode === 'dark' ? 'primary.light' : 'primary.50',
                color: (theme) => theme.palette.mode === 'dark' ? 'primary.contrastText' : 'primary.main',
                boxShadow: (theme) => theme.palette.mode === 'light' ? '0 4px 12px rgba(25, 118, 210, 0.15)' : 'none',
                border: (theme) => theme.palette.mode === 'light' ? '1px solid rgba(25, 118, 210, 0.2)' : 'none'
              }}>
                <CardContent>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                    <TrendingUpIcon sx={{ mr: 1, color: (theme) => theme.palette.mode === 'dark' ? 'inherit' : 'primary.main' }} />
                    <Typography variant="h4">{analyticsSummary.totalProjects}</Typography>
                  </Box>
                  <Typography variant="body2">Total Projects</Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <Card sx={{ 
                bgcolor: (theme) => theme.palette.mode === 'dark' ? 'error.light' : 'error.50',
                color: (theme) => theme.palette.mode === 'dark' ? 'error.contrastText' : 'error.main',
                boxShadow: (theme) => theme.palette.mode === 'light' ? '0 4px 12px rgba(211, 47, 47, 0.15)' : 'none',
                border: (theme) => theme.palette.mode === 'light' ? '1px solid rgba(211, 47, 47, 0.2)' : 'none'
              }}>
                <CardContent>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                    <BugReportIcon sx={{ mr: 1, color: (theme) => theme.palette.mode === 'dark' ? 'inherit' : 'error.main' }} />
                    <Typography variant="h4">{analyticsSummary.totalIssues}</Typography>
                  </Box>
                  <Typography variant="body2">Total Issues (7d)</Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <Card sx={{ 
                bgcolor: (theme) => theme.palette.mode === 'dark' ? 'warning.light' : 'warning.50',
                color: (theme) => theme.palette.mode === 'dark' ? 'warning.contrastText' : 'warning.main',
                boxShadow: (theme) => theme.palette.mode === 'light' ? '0 4px 12px rgba(237, 108, 2, 0.15)' : 'none',
                border: (theme) => theme.palette.mode === 'light' ? '1px solid rgba(237, 108, 2, 0.2)' : 'none'
              }}>
                <CardContent>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                    <WarningIcon sx={{ mr: 1, color: (theme) => theme.palette.mode === 'dark' ? 'inherit' : 'warning.main' }} />
                    <Typography variant="h4">{analyticsSummary.criticalIssues}</Typography>
                  </Box>
                  <Typography variant="body2">Critical Issues (7d)</Typography>
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <Card sx={{ 
                bgcolor: (theme) => theme.palette.mode === 'dark' ? 'success.light' : 'success.50',
                color: (theme) => theme.palette.mode === 'dark' ? 'success.contrastText' : 'success.main',
                boxShadow: (theme) => theme.palette.mode === 'light' ? '0 4px 12px rgba(46, 125, 50, 0.15)' : 'none',
                border: (theme) => theme.palette.mode === 'light' ? '1px solid rgba(46, 125, 50, 0.2)' : 'none'
              }}>
                <CardContent>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                    <CheckCircleIcon sx={{ mr: 1, color: (theme) => theme.palette.mode === 'dark' ? 'inherit' : 'success.main' }} />
                    <Typography variant="h4">{analyticsSummary.resolvedIssues}</Typography>
                  </Box>
                  <Typography variant="body2">Resolved Issues (7d)</Typography>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
        ) : (
          <Typography color="text.secondary">
            No data available for summary display.
          </Typography>
        )}
      </Paper>

      {/* Projects List */}
      <Paper sx={{ p: 2, width: '100%', minWidth: '800px' }}>
        <Typography variant="subtitle1" gutterBottom>
          Projects for Analytics
        </Typography>
        {isLoading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
            <CircularProgress />
          </Box>
        ) : error ? (
          <Typography color="error">
            Error loading projects. Please try again.
          </Typography>
        ) : projectsWithStats && projectsWithStats.length > 0 ? (
          <List sx={{ width: '100%', minWidth: '800px' }}>
            {projectsWithStats.map((project, index) => (
              <React.Fragment key={project.id}>
                {index > 0 && <Divider component="li" />}
                <ListItem 
                  sx={{ 
                    py: 1.5,
                    width: '100%',
                    maxWidth: 'none',
                    transition: 'all 0.2s ease-in-out',
                    '&:hover': { 
                      backgroundColor: (theme) => theme.palette.mode === 'dark'
                        ? 'rgba(65, 68, 75, 0.2)'
                        : 'rgba(0, 0, 0, 0.03)',
                    }
                  }}
                  button
                  onClick={() => handleProjectClick(project.id)}
                >
                  <Box sx={{ display: 'flex', flexDirection: 'column', width: '100%', maxWidth: 'none' }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', width: '100%' }}>
                      <Typography 
                        variant="subtitle2" 
                        component="div" 
                        sx={{ 
                          overflow: 'hidden',
                          textOverflow: 'ellipsis',
                          whiteSpace: 'nowrap'
                        }}
                      >
                        {project.name}
                        <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5, whiteSpace: 'pre-line' }}>
                          {project.description}
                        </Typography>
                      </Typography>
                      <Box sx={{ display: 'flex', alignItems: 'center' }}>
                        <ArrowForwardIcon color="primary" sx={{ mr: 1 }} />
                      </Box>
                    </Box>
                    <Box sx={{ mt: 0.5, display: 'flex', flexWrap: 'wrap', gap: 1 }}>
                      {project.statsLoading ? (
                        <CircularProgress size={16} />
                      ) : project.statsError ? (
                        <Typography variant="caption" color="error">Error loading stats</Typography>
                      ) : project.stats ? (
                        <>
                          <Chip 
                            label={`${project.stats.total_issues} issues`} 
                            size="small" 
                            color="primary" 
                            variant="outlined"
                          />
                          <Chip 
                            label={`${(project.stats.issues_by_level.fatal || 0) + (project.stats.issues_by_level.exception || 0)} critical`} 
                            size="small" 
                            color="error" 
                            variant="outlined"
                          />
                        </>
                      ) : (
                        <Typography variant="caption">Stats unavailable</Typography>
                      )}
                    </Box>
                  </Box>
                </ListItem>
              </React.Fragment>
            ))}
          </List>
        ) : (
          <Typography variant="body2">
            No projects to display. Create a new project to get started.
          </Typography>
        )}
      </Paper>
    </Layout>
  );
};

export default AnalyticsPage; 