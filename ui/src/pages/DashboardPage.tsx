import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography,
  Paper,
  Grid,
  Card,
  CardContent,
  CardActionArea,
  Chip,
  CircularProgress,
  Tooltip
} from '@mui/material';
import { 
  Error as ErrorIcon,
  ArrowForward as ArrowForwardIcon,
  TrendingUp as TrendingUpIcon,
  Dashboard as DashboardIcon
} from '@mui/icons-material';
import { getLevelColor, getLevelHexColor, getLevelBadgeStyles } from '../utils/issues/issueUtils';
import ErrorTrends from '../components/ErrorTrends';
import { useAuth } from '../auth/AuthContext';
import { Navigate, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import apiClient from '../api/apiClient';
import AuthenticatedLayout from "../components/AuthenticatedLayout";
import PageHeader from "../components/PageHeader";
import {GetProjectStatsPeriodEnum, type ProjectStatsResponse} from '../generated/api/client';

// Limit for recent issues to fetch
const RECENT_ISSUES_LIMIT = 5;

const DashboardPage: React.FC = () => {
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const [projectStats, setProjectStats] = useState<Record<number, ProjectStatsResponse>>({});

  // Fetch projects using React Query
  const { data: projects, isLoading: projectsLoading, error: projectsError } = useQuery({
    queryKey: ['projects'],
    queryFn: async () => {
      const response = await apiClient.recentProjectsList();
      return response.data;
    }
  });

  // Fetch recent issues using React Query
  const { data: recentIssuesResponse, isLoading: issuesLoading, error: issuesError } = useQuery({
    queryKey: ['recentIssues'],
    queryFn: async () => {
      const response = await apiClient.getRecentIssues(RECENT_ISSUES_LIMIT);
      return response.data;
    }
  });

  // Extract issues from the response
  const recentIssues = recentIssuesResponse?.issues || [];

  // Fetch project stats for each project
  useEffect(() => {
    const fetchProjectStats = async () => {
      if (!projects || projects.length === 0) return;

      const statsPromises = projects.map(async (project) => {
        try {
          const response = await apiClient.getProjectStats(project.id, GetProjectStatsPeriodEnum._7d);
          return { projectId: project.id, stats: response.data };
        } catch (err) {
          console.error(`Error fetching stats for project ${project.id}:`, err);
          return { projectId: project.id, stats: null };
        }
      });

      const results = await Promise.all(statsPromises);
      const statsMap: Record<number, ProjectStatsResponse> = {};

      results.forEach(result => {
        if (result.stats) {
          statsMap[result.projectId] = result.stats;
        }
      });

      setProjectStats(statsMap);
    };

    fetchProjectStats();
  }, [projects]);

  // If not authenticated, redirect to login
  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  const handleProjectClick = (projectId: number) => {
    navigate(`/projects/${projectId}`);
  };

  return (
    <AuthenticatedLayout>
      <PageHeader
        title="Dashboard"
        subtitle="Overview of your projects and issues."
        icon={<DashboardIcon />}
        gradientVariant="default"
      />

      <Grid container spacing={2}>
        <Grid item xs={12} md={4}>
          <Paper 
            sx={{ 
              p: 3,
              background: (theme) => theme.palette.mode === 'dark' 
                ? 'linear-gradient(to bottom, rgba(65, 68, 74, 0.5), rgba(55, 58, 64, 0.5))'
                : 'linear-gradient(to bottom, rgba(255, 255, 255, 0.9), rgba(245, 245, 245, 0.9))',
              backdropFilter: 'blur(10px)',
              boxShadow: '0 4px 20px 0 rgba(0, 0, 0, 0.05)'
            }}>
            <Typography variant="h6" gutterBottom>
              Recent Issues
            </Typography>
            {issuesLoading ? (
              <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
                <CircularProgress />
              </Box>
            ) : issuesError ? (
              <Typography color="error">
                Error loading recent issues. Please try again.
              </Typography>
            ) : recentIssues.length > 0 ? (
              <Box>
                {recentIssues.map((issue) => (
                  <Box 
                    key={issue.id} 
                    sx={{ 
                      mb: 2, 
                      p: 2, 
                      borderRadius: 1, 
                      background: (theme) => theme.palette.mode === 'dark'
                        ? 'linear-gradient(135deg, rgba(55, 58, 64, 0.6) 0%, rgba(50, 53, 58, 0.6) 100%)'
                        : 'linear-gradient(135deg, rgba(250, 250, 250, 0.9) 0%, rgba(240, 240, 240, 0.9) 100%)',
                      backdropFilter: 'blur(5px)',
                      boxShadow: '0 2px 8px 0 rgba(0, 0, 0, 0.05)',
                      cursor: 'pointer',
                      transition: 'all 0.2s ease-in-out',
                      '&:hover': {
                        background: (theme) => theme.palette.mode === 'dark'
                          ? 'linear-gradient(135deg, rgba(60, 63, 70, 0.7) 0%, rgba(55, 58, 64, 0.7) 100%)'
                          : 'linear-gradient(135deg, rgba(255, 255, 255, 1) 0%, rgba(245, 245, 245, 1) 100%)',
                        boxShadow: '0 4px 12px 0 rgba(0, 0, 0, 0.1)',
                        transform: 'translateY(-2px)'
                      }
                    }}
                    onClick={() => navigate(`/projects/${issue.project_id}/issues/${issue.id}`)}
                  >
                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                      <ErrorIcon 
                        sx={{ 
                          mr: 1,
                          color: getLevelHexColor(issue.level)
                        }} 
                        fontSize="small" 
                      />
                      <Typography variant="body2" sx={{ fontWeight: 'bold' }}>
                        {issue.title}
                      </Typography>
                    </Box>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                      <Chip 
                        label={issue.level} 
                        size="small" 
                        color={getLevelColor(issue.level)} 
                        sx={getLevelBadgeStyles(issue.level, 'small')} 
                      />
                      <Typography variant="caption" color="text.secondary">
                        {projects?.find(p => p.id === issue.project_id)?.name || `Project ${issue.project_id}`}
                      </Typography>
                    </Box>
                  </Box>
                ))}
              </Box>
            ) : (
              <Typography variant="body2">
                No recent issues to display.
              </Typography>
            )}
          </Paper>
        </Grid>

        <Grid item xs={12} md={8}>
          <Paper 
            sx={{ 
              p: 3,
              background: (theme) => theme.palette.mode === 'dark' 
                ? 'linear-gradient(to bottom, rgba(65, 68, 74, 0.5), rgba(55, 58, 64, 0.5))'
                : 'linear-gradient(to bottom, rgba(255, 255, 255, 0.9), rgba(245, 245, 245, 0.9))',
              backdropFilter: 'blur(10px)',
              boxShadow: '0 4px 20px 0 rgba(0, 0, 0, 0.05)'
            }}>
            <Typography variant="h6" gutterBottom>
              Most Active Projects
            </Typography>
            {projectsLoading ? (
              <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
                <CircularProgress />
              </Box>
            ) : projectsError ? (
              <Typography color="error">
                Error loading projects. Please try again.
              </Typography>
            ) : projects && projects.length > 0 ? (
              <Grid container spacing={2}>
                {projects.map((project) => (
                  <Grid item xs={12} key={project.id}>
                    <Card 
                      sx={{ 
                        background: (theme) => theme.palette.mode === 'dark'
                          ? 'linear-gradient(135deg, rgba(60, 63, 70, 0.6) 0%, rgba(55, 58, 64, 0.6) 100%)'
                          : 'linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(245, 245, 245, 0.95) 100%)',
                        backdropFilter: 'blur(8px)',
                        boxShadow: '0 2px 10px 0 rgba(0, 0, 0, 0.05)',
                        transition: 'all 0.2s ease-in-out',
                        '&:hover': { 
                          background: (theme) => theme.palette.mode === 'dark'
                            ? 'linear-gradient(135deg, rgba(65, 68, 75, 0.7) 0%, rgba(60, 63, 70, 0.7) 100%)'
                            : 'linear-gradient(135deg, rgba(255, 255, 255, 1) 0%, rgba(250, 250, 250, 1) 100%)',
                          boxShadow: '0 5px 15px 0 rgba(0, 0, 0, 0.1)',
                          transform: 'translateY(-3px)',
                          '& .MuiCardActionArea-focusHighlight': {
                            opacity: 0.1
                          }
                        }
                      }}
                    >
                      <CardActionArea onClick={() => handleProjectClick(project.id)}>
                        <CardContent>
                          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                            <Box sx={{ flexGrow: 1 }}>
                              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                <Typography variant="h6" component="div">
                                  {project.name}
                                </Typography>
                              </Box>
                              <Box sx={{ mt: 2, display: 'flex', gap: 1 }}>
                                <Chip 
                                  label={`${projectStats[project.id]?.total_issues || 0} issues`} 
                                  size="small" 
                                  color="primary" 
                                  variant="outlined"
                                />
                                <Chip 
                                  label={`${(projectStats[project.id]?.issues_by_level.fatal || 0) + (projectStats[project.id]?.issues_by_level.exception || 0)} critical`} 
                                  size="small" 
                                  color="error" 
                                  variant="outlined"
                                />
                              </Box>
                            </Box>

                            <Tooltip title="Error trends (last 3 hours)">
                              <Box 
                                sx={{ 
                                  ml: 2, 
                                  display: 'flex', 
                                  flexDirection: 'column', 
                                  alignItems: 'center',
                                  border: (theme) => `1px solid ${theme.palette.mode === 'dark' 
                                    ? 'rgba(255, 255, 255, 0.1)' 
                                    : 'rgba(0, 0, 0, 0.08)'}`,
                                  borderRadius: 2,
                                  p: 1.2,
                                  width: 'auto',
                                  minWidth: 160,
                                  backgroundColor: (theme) => theme.palette.mode === 'dark' 
                                    ? 'rgba(0, 0, 0, 0.1)' 
                                    : 'rgba(255, 255, 255, 0.5)',
                                  boxShadow: (theme) => theme.palette.mode === 'dark'
                                    ? '0 2px 8px 0 rgba(0, 0, 0, 0.2)'
                                    : '0 2px 8px 0 rgba(0, 0, 0, 0.05)',
                                  transition: 'all 0.2s ease-in-out',
                                  '&:hover': {
                                    transform: 'translateY(-2px)',
                                    boxShadow: (theme) => theme.palette.mode === 'dark'
                                      ? '0 4px 12px 0 rgba(0, 0, 0, 0.3)'
                                      : '0 4px 12px 0 rgba(0, 0, 0, 0.1)',
                                    borderColor: (theme) => theme.palette.mode === 'dark'
                                      ? 'rgba(255, 255, 255, 0.2)'
                                      : 'rgba(0, 0, 0, 0.12)',
                                  }
                                }}
                              >
                                <Box sx={{ display: 'flex', alignItems: 'center', mb: 0.5 }}>
                                  <TrendingUpIcon 
                                    color="error" 
                                    fontSize="small" 
                                    sx={{ mr: 0.5, fontSize: '0.875rem' }} 
                                  />
                                  <Typography 
                                    variant="caption" 
                                    color="text.secondary" 
                                    sx={{ fontSize: '0.7rem', fontWeight: 500 }}
                                  >
                                    Error Trends
                                  </Typography>
                                </Box>
                                <ErrorTrends width={200} height={50} projectId={project.id} />
                              </Box>
                            </Tooltip>
                            <Box sx={{ display: 'flex', alignItems: 'center', marginLeft: 2, marginRight: 1 }}></Box>
                            <ArrowForwardIcon color="primary" />
                          </Box>
                        </CardContent>
                      </CardActionArea>
                    </Card>
                  </Grid>
                ))}
              </Grid>
            ) : (
              <Typography variant="body2">
                No projects to display.
              </Typography>
            )}
          </Paper>
        </Grid>
      </Grid>
    </AuthenticatedLayout>
  );
};

export default DashboardPage;
