import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  Paper,
  CircularProgress,
  Tabs,
  Tab,
  IconButton,
  Button,
  Alert
} from '@mui/material';
import { useAuth } from '../auth/AuthContext';
import { Navigate, useParams } from 'react-router-dom';
import apiClient from '../api/apiClient';
import AuthenticatedLayout from '../components/AuthenticatedLayout';
import NotificationSettings from '../components/NotificationSettings';
import TeamMembersManagement from '../components/projects/TeamMembersManagement';
import ProjectEditForm from '../components/projects/ProjectEditForm';
import { ContentCopy as ContentCopyIcon, Edit as EditIcon } from '@mui/icons-material';
import Tooltip from '@mui/material/Tooltip';
import Snackbar from '@mui/material/Snackbar';
import { usePermissions } from '../hooks/usePermissions';
import { useTheme } from '@mui/material/styles';

// Interface for our project data
interface Project {
  id: number;
  name: string;
  description: string;
  team_id?: number | null;
  team_name?: string | null;
  public_key?: string;
  created_at: string;
  issues_count: number;
  critical_count: number;
}

// Interface for TabPanel props
interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

// TabPanel component
function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`settings-tabpanel-${index}`}
      aria-labelledby={`settings-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ pt: 3 }}>
          {children}
        </Box>
      )}
    </div>
  );
}

const ProjectSettingsPage: React.FC = () => {
  const { isAuthenticated } = useAuth();
  const { projectId } = useParams<{ projectId: string }>();
  const { canManageProject } = usePermissions();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [tabValue, setTabValue] = useState(0);
  const [copySuccess, setCopySuccess] = useState(false);
  const [isDataLoaded, setIsDataLoaded] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const theme = useTheme();

  // Handle tab change
  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  // Initialize project data
  const [project, setProject] = useState<Project>({
    id: Number(projectId),
    name: `Project ${projectId}`,
    description: '',
    team_id: null,
    team_name: null,
    public_key: '',
    created_at: new Date().toISOString(),
    issues_count: 0,
    critical_count: 0
  });

  // Fetch project details and team data
  useEffect(() => {
    const fetchProjectDetails = async () => {
      setIsLoading(true);
      setError(null);
      setIsDataLoaded(false);
      try {
        // Use the getProject API to get project details
        const response = await apiClient.getProject(Number(projectId));

        console.log('API Response:', response.data.project);
        console.log('Project team_id from API:', response.data.project.team_id);

        // Update project details based on the API response
        setProject(prev => ({
          ...prev,
          name: response.data.project.name,
          description: response.data.project.description,
          team_id: response.data.project.team_id,
          team_name: response.data.project.team_name,
          public_key: response.data.project.public_key,
          created_at: response.data.project.created_at
        }));

        // If project has a team, fetch team data
        if (response.data.project.team_id) {
          try {
            await apiClient.getProjectTeam(Number(projectId));
            // Team data is fetched by the TeamMembersManagement component
          } catch (teamErr) {
            console.error('Error fetching team data:', teamErr);
            // Don't fail the whole page if team data fails to load
          }
        }

        setIsDataLoaded(true);
      } catch (err) {
        console.error('Error fetching project details:', err);
        setError('Failed to fetch project details. Please try again later.');
      } finally {
        setIsLoading(false);
      }
    };

    fetchProjectDetails();
  }, [projectId]);

  // If not authenticated, redirect to login
  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  // Show loading while fetching project details
  if (isLoading) {
    return (
      <AuthenticatedLayout showBackButton={true} backTo={`/projects/${projectId}`}>
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
          <CircularProgress />
        </Box>
      </AuthenticatedLayout>
    );
  }

  // If user doesn't have permission to manage this project, redirect to project page
  // Only check after project details are loaded
  if (isDataLoaded) {
    console.log('Checking access:', { projectId, team_id: project.team_id, project });
    if (!canManageProject(Number(projectId), project.team_id || undefined)) {
      console.log('Access denied, redirecting to project page');
      return <Navigate to={`/projects/${projectId}`} />;
    }
  }

  // Handle project update
  const handleProjectUpdate = (newName: string, newDescription: string) => {
    setProject(prev => ({
      ...prev,
      name: newName,
      description: newDescription
    }));
    setIsEditing(false);
    setSuccessMessage('Project details updated successfully!');
    
    // Clear success message after 3 seconds
    setTimeout(() => {
      setSuccessMessage(null);
    }, 3000);
  };

  const handleCancelEdit = () => {
    setIsEditing(false);
  };

  return (
    <AuthenticatedLayout showBackButton={true} backTo={`/projects/${projectId}`}>
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom className="gradient-text-purple">
          Project Settings: {project.name}
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5, whiteSpace: 'pre-line' }}>
          {project.description}
        </Typography>
        <Typography variant="body1" paragraph>
          Configure settings for this project.
        </Typography>
      </Box>

      {successMessage && (
        <Alert severity="success" sx={{ mb: 2 }} onClose={() => setSuccessMessage(null)}>
          {successMessage}
        </Alert>
      )}

      {isLoading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
          <CircularProgress />
        </Box>
      ) : error ? (
        <Paper sx={{ p: 3, bgcolor: 'error.light', color: 'error.contrastText' }}>
          <Typography>{error}</Typography>
        </Paper>
      ) : (
        <Paper>
          <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
            <Tabs value={tabValue} onChange={handleTabChange} aria-label="project settings tabs">
              <Tab label="General" />
              <Tab label="Notifications" />
              <Tab label="API Keys" />
              <Tab label="Team" />
            </Tabs>
          </Box>

          {/* General Settings Tab */}
          <TabPanel value={tabValue} index={0}>
            <Box sx={{ p: 3 }}>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                <Typography variant="h6" gutterBottom>
                  General Settings
                </Typography>
                {!isEditing && (
                  <Button
                    variant="outlined"
                    startIcon={<EditIcon />}
                    onClick={() => setIsEditing(true)}
                  >
                    Edit Project
                  </Button>
                )}
              </Box>
              
              {isEditing ? (
                <ProjectEditForm
                  projectId={Number(projectId)}
                  initialName={project.name}
                  initialDescription={project.description}
                  onSave={handleProjectUpdate}
                  onCancel={handleCancelEdit}
                />
              ) : (
                <Paper
                  sx={{
                    p: 3,
                    bgcolor:
                      theme.palette.mode === 'dark'
                        ? 'rgba(40, 44, 52, 0.85)'
                        : 'grey.50',
                    border: theme.palette.mode === 'dark' ? '1px solid rgba(255,255,255,0.07)' : undefined,
                    boxShadow: theme.palette.mode === 'dark' ? '0 2px 8px 0 rgba(0,0,0,0.25)' : undefined
                  }}
                >
                  <Typography variant="subtitle1" gutterBottom>
                    Project Information
                  </Typography>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="body2" color="text.secondary" gutterBottom>
                      Project Name
                    </Typography>
                    <Typography variant="body1">
                      {project.name}
                    </Typography>
                  </Box>
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="body2" color="text.secondary" gutterBottom>
                      Description
                    </Typography>
                    <Typography variant="body1" sx={{ whiteSpace: 'pre-line' }}>
                      {project.description || 'No description provided'}
                    </Typography>
                  </Box>
                  <Box>
                    <Typography variant="body2" color="text.secondary" gutterBottom>
                      Created
                    </Typography>
                    <Typography variant="body1">
                      {new Date(project.created_at).toLocaleDateString()}
                    </Typography>
                  </Box>
                </Paper>
              )}
            </Box>
          </TabPanel>

          {/* Notification Settings Tab */}
          <TabPanel value={tabValue} index={1}>
            <Box sx={{ p: 3 }}>
              <Typography variant="h6" gutterBottom>
                Notification Settings
              </Typography>
              <NotificationSettings projectId={Number(projectId)} />
            </Box>
          </TabPanel>

          {/* API Keys Tab */}
          <TabPanel value={tabValue} index={2}>
            <Box sx={{ p: 3 }}>
              <Typography variant="h6" gutterBottom>
                API Keys
              </Typography>
              <Typography variant="body1">
                Manage API keys for your project.
              </Typography>
              {/* Add API keys content here */}
              {project.public_key && (
                <Box sx={{ mt: 2, maxWidth: '100%' }}>
                  <Typography variant="subtitle1">Public Key</Typography>
                  <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
                    <Paper sx={{ p: 2, bgcolor: 'grey.100', fontFamily: 'monospace', fontSize: 15, overflowX: 'auto', wordBreak: 'break-all', mr: 1 }}>
                      <Typography variant="body2" sx={{ fontFamily: 'monospace', fontSize: 15, m: 0 }}>
                        {project.public_key}
                      </Typography>
                    </Paper>
                    <Tooltip title="Copy public key to clipboard" placement="top" arrow>
                      <IconButton
                        aria-label="copy public key"
                        onClick={() => {
                          navigator.clipboard.writeText(project.public_key || '');
                          setCopySuccess(true);
                        }}
                        size="small"
                      >
                        <ContentCopyIcon fontSize="small" />
                      </IconButton>
                    </Tooltip>
                  </Box>
                  <Snackbar
                    open={copySuccess}
                    autoHideDuration={2000}
                    onClose={() => setCopySuccess(false)}
                    message="Copied to clipboard!"
                    anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
                  />
                </Box>
              )}
            </Box>
          </TabPanel>

          {/* Team Management Tab */}
          <TabPanel value={tabValue} index={3}>
            <Box sx={{ p: 3 }}>
              <Typography variant="h6" gutterBottom>
                Team Management
              </Typography>
              <Typography variant="body1" sx={{ mb: 3 }}>
                Manage team members and their roles for this project.
              </Typography>
              {project.team_id ? (
                <TeamMembersManagement 
                  teamId={project.team_id} 
                  projectId={Number(projectId)} 
                />
              ) : (
                <Box sx={{ p: 3, textAlign: 'center' }}>
                  <Typography variant="body2" color="text.secondary">
                    This project is not associated with any team.
                  </Typography>
                </Box>
              )}
            </Box>
          </TabPanel>
        </Paper>
      )}
    </AuthenticatedLayout>
  );
};

export default ProjectSettingsPage;
