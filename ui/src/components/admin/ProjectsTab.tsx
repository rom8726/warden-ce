import React from 'react';
import {
  Box,
  Typography,
  Button,
  TableContainer,
  Table,
  TableHead,
  TableBody,
  TableRow,
  TableCell,
  IconButton,
  CircularProgress,
  Chip
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  PersonAdd as PersonAddIcon
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';

interface Project {
  id: number;
  name: string;
  description: string;
  team_id?: number | null;
  team_name?: string | null;
  created_at: string;
}

interface ProjectsTabProps {
  projects: Project[] | undefined;
  isLoading: boolean;
  error: unknown;
  onCreateProject: () => void;
  onArchiveProject?: (projectId: number) => void;
  setSnackbar: (snackbar: { open: boolean; message: string; severity: 'success' | 'error' | 'info' | 'warning' }) => void;
}

const ProjectsTab: React.FC<ProjectsTabProps> = ({
  projects,
  isLoading,
  error,
  onCreateProject,
  onArchiveProject,
  setSnackbar
}) => {
  const navigate = useNavigate();

  return (
    <>
      <Box 
        sx={{ 
          display: 'flex', 
          justifyContent: 'space-between', 
          alignItems: 'center', 
          mb: 4,
          pb: 2,
          borderBottom: (theme) => `1px solid ${theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.05)'}`
        }}
      >
        <Box>
          <Typography 
            variant="h6"
            sx={{ 
              fontWeight: 600,
              mb: 0.5
            }}
          >
            Manage Projects
          </Typography>
          <Typography 
            variant="body2" 
            color="text.secondary"
            sx={{ maxWidth: '600px' }}
          >
            Create and manage projects across your organization.
          </Typography>
        </Box>
        <Button 
          variant="contained" 
          startIcon={<AddIcon />}
          onClick={onCreateProject}
          sx={{
            px: 2,
            py: 1,
            boxShadow: (theme) => theme.palette.mode === 'dark' 
              ? '0 4px 12px rgba(0, 0, 0, 0.3)' 
              : '0 4px 12px rgba(94, 114, 228, 0.2)',
            '&:hover': {
              transform: 'translateY(-2px)',
              boxShadow: (theme) => theme.palette.mode === 'dark' 
                ? '0 6px 16px rgba(0, 0, 0, 0.4)' 
                : '0 6px 16px rgba(94, 114, 228, 0.3)',
            },
            transition: 'all 0.2s ease-in-out'
          }}
        >
          Create Project
        </Button>
      </Box>

      {isLoading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
          <CircularProgress />
        </Box>
      ) : error ? (
        <Typography color="error">
          Error loading projects. Please try again.
        </Typography>
      ) : projects && projects.length > 0 ? (
        <TableContainer 
          sx={{ 
            borderRadius: 2,
            boxShadow: '0 2px 10px 0 rgba(0, 0, 0, 0.05)',
            overflow: 'hidden'
          }}
        >
          <Table sx={{ minWidth: 650 }}>
            <TableHead>
              <TableRow sx={{ 
                backgroundColor: (theme) => theme.palette.mode === 'dark' 
                  ? 'rgba(60, 63, 70, 0.5)' 
                  : 'rgba(245, 245, 245, 0.8)'
              }}>
                <TableCell sx={{ fontWeight: 600, py: 2 }}>ID</TableCell>
                <TableCell sx={{ fontWeight: 600, py: 2 }}>Name</TableCell>
                <TableCell sx={{ fontWeight: 600, py: 2 }}>Team</TableCell>
                <TableCell sx={{ fontWeight: 600, py: 2 }}>Created At</TableCell>
                <TableCell sx={{ fontWeight: 600, py: 2 }}>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {projects.map((project) => (
                <TableRow 
                  key={project.id}
                  sx={{ 
                    '&:hover': { 
                      backgroundColor: (theme) => theme.palette.mode === 'dark' 
                        ? 'rgba(60, 63, 70, 0.3)' 
                        : 'rgba(245, 245, 245, 0.5)'
                    },
                    transition: 'background-color 0.2s ease-in-out'
                  }}
                >
                  <TableCell sx={{ py: 1.5 }}>{project.id}</TableCell>
                  <TableCell sx={{ py: 1.5, fontWeight: 500 }}>
                    {project.name}
                    <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5, whiteSpace: 'pre-line' }}>
                      {project.description}
                    </Typography>
                  </TableCell>
                  <TableCell sx={{ py: 1.5 }}>
                    {project.team_name ? (
                      <Chip 
                        label={project.team_name} 
                        size="small" 
                        variant="outlined" 
                        color="secondary"
                      />
                    ) : (
                      "common"
                    )}
                  </TableCell>
                  <TableCell sx={{ py: 1.5 }}>{new Date(project.created_at).toLocaleString()}</TableCell>
                  <TableCell sx={{ py: 1.5 }}>
                    <Box sx={{ display: 'flex', gap: 1 }}>
                      <IconButton 
                        size="small" 
                        color="primary"
                        onClick={() => navigate(`/projects/${project.id}/settings`)}
                        sx={{ 
                          transition: 'transform 0.2s ease-in-out',
                          '&:hover': { transform: 'scale(1.1)' }
                        }}
                      >
                        <EditIcon fontSize="small" />
                      </IconButton>
                      <IconButton 
                        size="small" 
                        color="error"
                        onClick={() => {
                          if (onArchiveProject) {
                            onArchiveProject(project.id);
                          } else {
                            setSnackbar({
                              open: true,
                              message: 'Archive functionality not implemented yet',
                              severity: 'info'
                            });
                          }
                        }}
                        title="Archive project"
                        sx={{ 
                          transition: 'transform 0.2s ease-in-out',
                          '&:hover': { transform: 'scale(1.1)' }
                        }}
                      >
                        <DeleteIcon fontSize="small" />
                      </IconButton>
                    </Box>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      ) : (
        <Typography variant="body2" sx={{ p: 3 }}>
          No projects to display. Create a new project to get started.
        </Typography>
      )}
    </>
  );
};

export default ProjectsTab;