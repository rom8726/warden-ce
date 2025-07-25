import React, { useState } from 'react';
import { 
  Box, 
  Typography,
  Paper,
  Card,
  CardContent,
  Chip,
  CircularProgress,
  List,
  ListItem,
  ListItemText,
  Divider,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from '@mui/material';
import { 
  Group as GroupIcon,
  Person as PersonIcon
} from '@mui/icons-material';
import { useAuth } from '../auth/AuthContext';
import { Navigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../api/apiClient';
import AuthenticatedLayout from "../components/AuthenticatedLayout";
import TwoFactorAuthSection from '../components/TwoFactorAuthSection';
import Button from '@mui/material/Button';
import Alert from '@mui/material/Alert';
import VersionInfo from '../components/VersionInfo';
import SystemVersions from '../components/SystemVersions';
import ChangePasswordForm from '../components/ChangePasswordForm';

const SettingsPage: React.FC = () => {
  const { isAuthenticated } = useAuth();
  const [leaveLoading, setLeaveLoading] = useState<number | null>(null);
  const [leaveError, setLeaveError] = useState<string | null>(null);
  const [leaveSuccess, setLeaveSuccess] = useState<string | null>(null);
  const [confirmDialogOpen, setConfirmDialogOpen] = useState(false);
  const [selectedTeam, setSelectedTeam] = useState<{ id: number; name: string } | null>(null);
  const [showPasswordDialog, setShowPasswordDialog] = useState(false);

  // Fetch current user data using React Query
  const { 
    data: userData, 
    isLoading: userLoading, 
    error: userError,
    refetch
  } = useQuery({
    queryKey: ['currentUser'],
    queryFn: async () => {
      const response = await apiClient.getCurrentUser();
      return response.data;
    }
  });


  // If not authenticated, redirect to login
  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  const handleLeaveTeam = async () => {
    if (!selectedTeam) return;
    
    setLeaveError(null);
    setLeaveSuccess(null);
    setLeaveLoading(selectedTeam.id);
    setConfirmDialogOpen(false);
    
    try {
      await apiClient.removeTeamMember(selectedTeam.id, userData?.id || 0);
      setLeaveSuccess('Successfully left the team.');
      await refetch();
    } catch (err: unknown) {
      if (typeof err === 'object' && err !== null && 'response' in err && 
          typeof err.response === 'object' && err.response !== null && 
          'status' in err.response && err.response.status === 403 &&
          'data' in err.response && typeof err.response.data === 'object' && 
          err.response.data !== null && 'error' in err.response.data) {
        const errorData = err.response.data.error;
        setLeaveError(typeof errorData === 'object' && errorData !== null && 'message' in errorData ? 
          String(errorData.message) : 'Cannot leave team: you are the last owner.');
      } else {
        setLeaveError('Error while trying to leave the team.');
      }
    } finally {
      setLeaveLoading(null);
      setSelectedTeam(null);
    }
  };

  return (
    <AuthenticatedLayout>
      <Typography variant="h4" component="h1" gutterBottom className="gradient-text">
        Settings
      </Typography>

      <Typography variant="body1" paragraph>
        Manage your account settings and view your team memberships.
      </Typography>

      <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
        {/* 1. Account Information */}
        <Paper sx={{ p: 3 }}>
          <Typography variant="h6" gutterBottom>
            Account Information
          </Typography>
          {userLoading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
              <CircularProgress />
            </Box>
          ) : userError ? (
            <Typography color="error">
              Error loading user data. Please try again.
            </Typography>
          ) : userData ? (
            <Card>
              <CardContent>
                <Box sx={{ mb: 2 }}>
                  <Typography variant="subtitle2" color="text.secondary">
                    Username
                  </Typography>
                  <Typography variant="body1">
                    {userData.username}
                  </Typography>
                </Box>
                <Box sx={{ mb: 2 }}>
                  <Typography variant="subtitle2" color="text.secondary">
                    Email
                  </Typography>
                  <Typography variant="body1">
                    {userData.email}
                  </Typography>
                </Box>
                <Box sx={{ mb: 2 }}>
                  <Typography variant="subtitle2" color="text.secondary">
                    Account Type
                  </Typography>
                  <Chip 
                    label={userData.is_superuser ? "Administrator" : "Regular User"} 
                    color={userData.is_superuser ? "primary" : "default"}
                    size="small"
                  />
                </Box>
                <Box>
                  <Typography variant="subtitle2" color="text.secondary">
                    Account Created
                  </Typography>
                  <Typography variant="body1">
                    {new Date(userData.created_at).toLocaleDateString()}
                  </Typography>
                </Box>
                <Box sx={{ mt: 2 }}>
                  <Button 
                    variant="contained" 
                    onClick={() => setShowPasswordDialog(true)}
                  >
                    Change Password
                  </Button>
                </Box>
              </CardContent>
            </Card>
          ) : (
            <Typography variant="body2">
              Unable to load account information.
            </Typography>
          )}
        </Paper>

        {/* 2. Two-Factor Authentication */}
        <TwoFactorAuthSection userData={userData} userLoading={userLoading} userError={userError} />

        {/* 3. Your Teams */}
        <Paper sx={{ p: 3 }}>
          <Typography variant="h6" gutterBottom>
            Your Teams
          </Typography>
          {leaveError && <Alert severity="error" sx={{ mb: 2 }}>{leaveError}</Alert>}
          {leaveSuccess && <Alert severity="success" sx={{ mb: 2 }}>{leaveSuccess}</Alert>}
          {userLoading ? (
            <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
              <CircularProgress />
            </Box>
          ) : userError ? (
            <Typography color="error">
              Error loading user data. Please try again.
            </Typography>
          ) : (() => {
              const teams = userData?.teams;
            if (teams && teams.length > 0) {
                return (
                  <List>
                    {teams.map((team, index) => (
                      <React.Fragment key={team.id}>
                        <ListItem>
                          <Box sx={{ display: 'flex', alignItems: 'center', width: '100%' }}>
                            <GroupIcon sx={{ mr: 2, color: 'primary.main' }} />
                            <ListItemText 
                              primary={
                                <Box sx={{ display: 'flex', alignItems: 'center' }}>
                                  <Typography variant="body1" sx={{ mr: 1 }}>
                                    {team.name}
                                  </Typography>
                                </Box>
                              }
                              secondary={
                                <div style={{ display: 'flex', alignItems: 'center', marginTop: '4px' }}>
                                  <PersonIcon sx={{ fontSize: 16, mr: 0.5, color: 'text.secondary' }} />
                                  <Typography variant="body2" color="text.secondary" component="div">
                                    Your role: 
                                  </Typography>
                                  <Chip 
                                    label={team.role} 
                                    size="small" 
                                    color={
                                      team.role === 'owner' 
                                        ? 'primary' 
                                        : team.role === 'admin' 
                                          ? 'secondary' 
                                          : 'default'
                                    }
                                    sx={{ ml: 1, height: 20 }}
                                  />
                                </div>
                              }
                            />
                            {/* Leave Team button */}
                            <Box sx={{ ml: 2 }}>
                              <Button
                                variant="outlined"
                                color="error"
                                size="small"
                                disabled={leaveLoading === team.id}
                                onClick={() => {
                                  setSelectedTeam({ id: team.id, name: team.name });
                                  setConfirmDialogOpen(true);
                                }}
                              >
                                Leave Team
                              </Button>
                            </Box>
                          </Box>
                        </ListItem>
                        {index < teams.length - 1 && <Divider />}
                      </React.Fragment>
                    ))}
                  </List>
                );
              }
              return (
                <Typography variant="body2">
                  You are not a member of any teams.
                </Typography>
              );
            })()
          }
        </Paper>


        {/* 5. Version Information */}
        <Paper sx={{ p: 3 }}>
          <Typography variant="h6" gutterBottom>
            Version Information
          </Typography>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            <VersionInfo showBuildTime={true} variant="body2" />
          </Box>
        </Paper>

        {/* 6. System Components Versions */}
        <Paper sx={{ p: 3 }}>
          <SystemVersions showBuildTime={true} showStatus={true} />
        </Paper>
      </Box>

      {/* Confirmation Dialog */}
      <Dialog
        open={confirmDialogOpen}
        onClose={() => {
          setConfirmDialogOpen(false);
          setSelectedTeam(null);
        }}
        aria-labelledby="leave-team-dialog-title"
      >
        <DialogTitle id="leave-team-dialog-title">
          Leave Team
        </DialogTitle>
        <DialogContent>
          <Typography>
            Are you sure you want to leave the team "{selectedTeam?.name}"? 
            This action cannot be undone.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button 
            onClick={() => {
              setConfirmDialogOpen(false);
              setSelectedTeam(null);
            }}
            disabled={leaveLoading !== null}
          >
            Cancel
          </Button>
          <Button 
            onClick={handleLeaveTeam}
            color="error"
            variant="contained"
            disabled={leaveLoading !== null}
          >
            {leaveLoading === selectedTeam?.id ? 'Leaving...' : 'Leave Team'}
          </Button>
        </DialogActions>
      </Dialog>

      <ChangePasswordForm open={showPasswordDialog} onClose={() => setShowPasswordDialog(false)} />
    </AuthenticatedLayout>
  );
};

export default SettingsPage;