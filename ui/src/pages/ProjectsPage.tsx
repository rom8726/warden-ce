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
  Menu,
  MenuItem
} from '@mui/material';
import { 
  ArrowForward as ArrowForwardIcon,
  MoreVert as MoreVertIcon,
  Settings as SettingsIcon,
  FolderOutlined as ProjectsIcon
} from '@mui/icons-material';
import { useAuth } from '../auth/AuthContext';
import { Navigate, useNavigate } from 'react-router-dom';
import { useQuery, useQueries } from '@tanstack/react-query';
import apiClient from '../api/apiClient';
import AuthenticatedLayout from "../components/AuthenticatedLayout";
import PageHeader from "../components/PageHeader";
import {GetProjectStatsPeriodEnum} from "../generated/api/client";
import { usePermissions } from '../hooks/usePermissions';

const ProjectsPage: React.FC = () => {
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const { canManageProject } = usePermissions();

  // State for menu
  const [menuAnchorEl, setMenuAnchorEl] = React.useState<null | HTMLElement>(null);
  const [selectedProjectId, setSelectedProjectId] = React.useState<number | null>(null);
  const isMenuOpen = Boolean(menuAnchorEl);

  // Handle menu open
  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>, projectId: number) => {
    event.stopPropagation(); // Prevent triggering the ListItem click
    setMenuAnchorEl(event.currentTarget);
    setSelectedProjectId(projectId);
  };

  // Handle menu close
  const handleMenuClose = () => {
    setMenuAnchorEl(null);
    setSelectedProjectId(null);
  };

  // Handle settings click
  const handleSettingsClick = () => {
    if (selectedProjectId) {
      navigate(`/projects/${selectedProjectId}/settings`);
    }
    handleMenuClose();
  };

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

  const isLoading = projectsLoading || projectStatsQueries.some(query => query.isLoading);
  const error = projectsError || projectStatsQueries.some(query => query.error);

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
        title="Projects"
        subtitle="View and manage all your projects."
        icon={<ProjectsIcon />}
      />

      <Paper sx={{ p: 2, width: '100%', minWidth: '800px' }}>
        <Typography variant="subtitle1" gutterBottom>
          All Projects
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
                        <Tooltip title="More options">
                          <IconButton 
                            size="small"
                            onClick={(e) => handleMenuOpen(e, project.id)}
                          >
                            <MoreVertIcon fontSize="small" />
                          </IconButton>
                        </Tooltip>
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
                            label={`${project.stats.issues_by_level.fatal + project.stats.issues_by_level.exception} critical`} 
                            size="small" 
                            color="error" 
                            variant="outlined"
                          />
                        </>
                      ) : (
                        <Typography variant="caption">No stats available</Typography>
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

      {/* Project options menu */}
      <Menu
        anchorEl={menuAnchorEl}
        open={isMenuOpen}
        onClose={handleMenuClose}
        onClick={handleMenuClose}
      >
        {selectedProjectId && (() => {
          const selectedProject = projects?.find(p => p.id === selectedProjectId);
          return canManageProject(selectedProjectId, selectedProject?.team_id || undefined);
        })() && (
          <MenuItem onClick={handleSettingsClick}>
            <SettingsIcon fontSize="small" sx={{ mr: 1 }} />
            Settings
          </MenuItem>
        )}
      </Menu>
    </AuthenticatedLayout>
  );
};

export default ProjectsPage;
