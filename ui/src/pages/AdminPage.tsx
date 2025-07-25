import React, { useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  Tabs,
  Tab,
} from '@mui/material';
import { Navigate, useNavigate } from 'react-router-dom';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '../auth/AuthContext';
import { apiClient } from '../api/apiClient';
import { AddTeamMemberRequestRoleEnum, ChangeTeamMemberRoleRequestRoleEnum } from '../generated/api/client';
import AuthenticatedLayout from "../components/AuthenticatedLayout";
import TabPanel from '../components/admin/TabPanel';
import ProjectsTab from '../components/admin/ProjectsTab';
import UsersTab from '../components/admin/UsersTab';
import TeamsTab from '../components/admin/TeamsTab';

import LicenseTab from '../components/admin/LicenseTab';
import CreateProjectDialog from '../components/admin/CreateProjectDialog';
import CreateUserDialog from '../components/admin/CreateUserDialog';
import CreateTeamDialog from '../components/admin/CreateTeamDialog';
import AddTeamMemberDialog from '../components/admin/AddTeamMemberDialog';
import ConfirmationDialog from '../components/admin/ConfirmationDialog';
import { useNotification } from '../App';

const AdminPage: React.FC = () => {
  const { isAuthenticated, user } = useAuth();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { showNotification } = useNotification();
  const [tabValue, setTabValue] = useState(0);
  const [openProjectDialog, setOpenProjectDialog] = useState(false);
  const [openUserDialog, setOpenUserDialog] = useState(false);
  const [openTeamDialog, setOpenTeamDialog] = useState(false);
  const [openAddTeamMemberDialog, setOpenAddTeamMemberDialog] = useState(false);
  const [projectName, setProjectName] = useState('');
  const [selectedTeamId, setSelectedTeamId] = useState<number | null>(null);
  const [teamName, setTeamName] = useState('');
  const [selectedTeam, setSelectedTeam] = useState<number | null>(null);
  const [selectedUser, setSelectedUser] = useState<number | null>(null);
  const [selectedRole, setSelectedRole] = useState<string>('member');
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isSuperuser, setIsSuperuser] = useState(false);

  // State for confirmation dialogs
  const [deleteUserConfirmation, setDeleteUserConfirmation] = useState<{open: boolean, userId: number | null}>({
    open: false,
    userId: null
  });
  const [deleteTeamConfirmation, setDeleteTeamConfirmation] = useState<{open: boolean, teamId: number | null}>({
    open: false,
    teamId: null
  });
  const [removeTeamMemberConfirmation, setRemoveTeamMemberConfirmation] = useState<{
    open: boolean,
    teamId: number | null,
    userId: number | null
  }>({
    open: false,
    teamId: null,
    userId: null
  });
  const [archiveProjectConfirmation, setArchiveProjectConfirmation] = useState<{open: boolean, projectId: number | null}>({
    open: false,
    projectId: null
  });

  // Fetch projects using React Query
  const {
    data: projects,
    isLoading: isLoadingProjects,
    error: projectsError
  } = useQuery({
    queryKey: ['projects'],
    queryFn: async () => {
      const response = await apiClient.listProjects();
      return response.data;
    }
  });

  // Fetch users using React Query
  const {
    data: users,
    isLoading: isLoadingUsers,
    error: usersError
  } = useQuery({
    queryKey: ['users'],
    queryFn: async () => {
      const response = await apiClient.listUsers();
      return response.data;
    }
  });

  // Fetch teams using React Query
  const {
    data: teams,
    isLoading: isLoadingTeams,
    error: teamsError
  } = useQuery({
    queryKey: ['teams'],
    queryFn: async () => {
      const response = await apiClient.listTeams();
      return response.data;
    }
  });

  // If not authenticated or not a superuser, redirect to dashboard
  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }

  if (user && !user.is_superuser) {
    return <Navigate to="/dashboard" />;
  }

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  const handleOpenProjectDialog = () => {
    setProjectName('');
    setSelectedTeamId(null);
    setOpenProjectDialog(true);
  };

  const handleCloseProjectDialog = () => {
    setOpenProjectDialog(false);
  };

  const handleOpenUserDialog = () => {
    setUsername('');
    setEmail('');
    setPassword('');
    setIsSuperuser(false);
    setOpenUserDialog(true);
  };

  const handleCloseUserDialog = () => {
    setOpenUserDialog(false);
  };

  const handleCreateProject = async (name: string, description: string, teamId: number | null) => {
    try {
      // Create the project using the addProject API method
      await apiClient.addProject({
        name,
        description,
        team_id: teamId
      });

      // Show success message
      showNotification(`Project "${name}" создан успешно!`, 'success');
      handleCloseProjectDialog();

      // Refetch projects after creating a new one
      queryClient.invalidateQueries({ queryKey: ['projects'] });
    } catch (error) {
      console.error('Ошибка при создании проекта:', error);
      showNotification('Не удалось создать проект. Попробуйте еще раз.', 'error');
    }
  };

  const handleCreateUser = async (username: string, email: string, password: string, isSuperuser: boolean) => {
    try {
      // Create the user using the createUser API method
      await apiClient.createUser({
        username,
        email,
        password,
        is_superuser: isSuperuser
      });

      // Show success message
      showNotification(`User "${username}" created successfully!`, 'success');
      handleCloseUserDialog();

      // Refetch users after creating a new one
      queryClient.invalidateQueries({ queryKey: ['users'] });
    } catch (error) {
      console.error('Error creating user:', error);
      showNotification('Failed to create user. Please try again.', 'error');
    }
  };

  const handleToggleUserStatus = async (userId: number, isActive: boolean) => {
    try {
      // Use the new API endpoint for toggling user active status
      await apiClient.setUserActiveStatus(userId, {
        is_active: isActive
      });

      showNotification(`User ${isActive ? 'activated' : 'deactivated'} successfully!`, 'success');

      // Refetch users after updating status
      queryClient.invalidateQueries({ queryKey: ['users'] });
    } catch (error) {
      console.error('Error toggling user status:', error);
      showNotification('Failed to update user status. Please try again.', 'error');
    }
  };

  const handleToggleSuperuserStatus = async (userId: number, isSuperuser: boolean) => {
    try {
      // Use the new API endpoint for toggling superuser status
      await apiClient.setSuperuserStatus(userId, {
        is_superuser: isSuperuser
      });

      showNotification(`User ${isSuperuser ? 'granted' : 'revoked'} superuser privileges successfully!`, 'success');

      // Refetch users after updating status
      queryClient.invalidateQueries({ queryKey: ['users'] });
    } catch (error) {
      console.error('Error toggling superuser status:', error);
      showNotification('Failed to update superuser status. Please try again.', 'error');
    }
  };

  const handleOpenTeamDialog = () => {
    setTeamName('');
    setOpenTeamDialog(true);
  };

  const handleCloseTeamDialog = () => {
    setOpenTeamDialog(false);
  };

  const handleCreateTeam = async (name: string) => {
    try {
      // Create the team using the createTeam API method
      await apiClient.createTeam({
        name
      });

      // Show success message
      showNotification(`Team "${name}" created successfully!`, 'success');
      handleCloseTeamDialog();

      // Refetch teams after creating a new one
      queryClient.invalidateQueries({ queryKey: ['teams'] });
    } catch (error) {
      console.error('Error creating team:', error);
      showNotification('Failed to create team. Please try again.', 'error');
    }
  };

  const handleConfirmDeleteTeam = (teamId: number) => {
    setDeleteTeamConfirmation({
      open: true,
      teamId
    });
  };

  const handleCloseDeleteTeamConfirmation = () => {
    setDeleteTeamConfirmation({
      open: false,
      teamId: null
    });
  };

  const handleDeleteTeam = async () => {
    if (!deleteTeamConfirmation.teamId) return;

    try {
      // Delete the team using the deleteTeam API method
      await apiClient.deleteTeam(deleteTeamConfirmation.teamId);

      // Show success message
      showNotification('Team deleted successfully!', 'success');

      // Close the confirmation dialog
      handleCloseDeleteTeamConfirmation();

      // Refetch teams after deleting one
      queryClient.invalidateQueries({ queryKey: ['teams'] });
    } catch (error: any) {
      console.error('Error deleting team:', error);
      let message = 'Failed to delete team. Please try again.';
      if (error?.response?.data?.error?.message) {
        message = error.response.data.error.message;
      }
      showNotification(message, 'error');
    }
  };

  const handleOpenAddTeamMemberDialog = (teamId: number) => {
    setSelectedTeam(teamId);
    setSelectedUser(null);
    setSelectedRole('member');
    setOpenAddTeamMemberDialog(true);
  };

  const handleCloseAddTeamMemberDialog = () => {
    setOpenAddTeamMemberDialog(false);
  };

  const handleAddTeamMember = async (userId: number, role: AddTeamMemberRequestRoleEnum) => {
    if (!selectedTeam) return;

    try {
      // Add the team member using the addTeamMember API method
      await apiClient.addTeamMember(selectedTeam, {
        user_id: userId,
        role
      });

      // Show success message (only for single user addition)
      showNotification('Team member added successfully!', 'success');
      
      // Refetch teams after adding a member
      queryClient.invalidateQueries({ queryKey: ['teams'] });
    } catch (error) {
      console.error('Error adding team member:', error);
      showNotification('Failed to add team member. Please try again.', 'error');
      throw error; // Re-throw for bulk operations
    }
  };

  const handleBulkAddTeamMembers = async (userIds: number[], role: AddTeamMemberRequestRoleEnum) => {
    if (!selectedTeam) return;

    let successCount = 0;
    let errorCount = 0;

    for (const userId of userIds) {
      try {
        await apiClient.addTeamMember(selectedTeam, {
          user_id: userId,
          role
        });
        successCount++;
      } catch (error) {
        errorCount++;
        console.error(`Failed to add user ${userId}:`, error);
      }
    }

    // Show results
    if (successCount > 0) {
      const message = errorCount > 0 
        ? `Added ${successCount}/${userIds.length} members from group (${errorCount} failed)`
        : `Successfully added ${successCount} members from group`;
      showNotification(message, errorCount > 0 ? 'warning' : 'success');
    } else {
      showNotification('Failed to add any members from group', 'error');
    }

    // Refetch teams after adding members
    queryClient.invalidateQueries({ queryKey: ['teams'] });
  };

  const handleConfirmRemoveTeamMember = (teamId: number, userId: number) => {
    setRemoveTeamMemberConfirmation({
      open: true,
      teamId,
      userId
    });
  };

  const handleCloseRemoveTeamMemberConfirmation = () => {
    setRemoveTeamMemberConfirmation({
      open: false,
      teamId: null,
      userId: null
    });
  };

  const handleRemoveTeamMember = async () => {
    if (!removeTeamMemberConfirmation.teamId || !removeTeamMemberConfirmation.userId) return;

    try {
      // Remove the team member using the removeTeamMember API method
      await apiClient.removeTeamMember(
        removeTeamMemberConfirmation.teamId,
        removeTeamMemberConfirmation.userId
      );

      // Show success message
      showNotification('Team member removed successfully!', 'success');

      // Close the confirmation dialog
      handleCloseRemoveTeamMemberConfirmation();

      // Refetch teams after removing a member
      queryClient.invalidateQueries({ queryKey: ['teams'] });
    } catch (error) {
      console.error('Error removing team member:', error);
      showNotification('Failed to remove team member. Please try again.', 'error');
    }
  };

  const handleChangeMemberRole = async (teamId: number, userId: number, role: ChangeTeamMemberRoleRequestRoleEnum) => {
    try {
      // Change the team member role using the changeTeamMemberRole API method
      await apiClient.changeTeamMemberRole(teamId, userId, { role });

      // Show success message
      showNotification('Team member role changed successfully!', 'success');

      // Refetch teams after changing role
      queryClient.invalidateQueries({ queryKey: ['teams'] });
    } catch (error) {
      console.error('Error changing team member role:', error);
      showNotification('Failed to change team member role. Please try again.', 'error');
      throw error; // Re-throw to let the dialog handle the error
    }
  };

  const handleConfirmDeleteUser = (userId: number) => {
    setDeleteUserConfirmation({
      open: true,
      userId
    });
  };

  const handleCloseDeleteUserConfirmation = () => {
    setDeleteUserConfirmation({
      open: false,
      userId: null
    });
  };

  const handleDeleteUser = async () => {
    if (!deleteUserConfirmation.userId) return;

    try {
      // Use the new API endpoint for deleting a user
      await apiClient.deleteUser(deleteUserConfirmation.userId);

      showNotification('User deleted successfully!', 'success');

      // Close the confirmation dialog
      handleCloseDeleteUserConfirmation();

      // Refetch users after deleting
      queryClient.invalidateQueries({ queryKey: ['users'] });
    } catch (error) {
      console.error('Error deleting user:', error);
      showNotification('Failed to delete user. Please try again.', 'error');
    }
  };
  
  const handleConfirmArchiveProject = (projectId: number) => {
    setArchiveProjectConfirmation({
      open: true,
      projectId
    });
  };

  const handleCloseArchiveProjectConfirmation = () => {
    setArchiveProjectConfirmation({
      open: false,
      projectId: null
    });
  };

  const handleArchiveProject = async () => {
    if (!archiveProjectConfirmation.projectId) return;

    try {
      // Use the archiveProject API method
      await apiClient.archiveProject(archiveProjectConfirmation.projectId);

      showNotification('Project archived successfully!', 'success');

      // Close the confirmation dialog
      handleCloseArchiveProjectConfirmation();

      // Refetch projects after archiving
      queryClient.invalidateQueries({ queryKey: ['projects'] });
    } catch (error) {
      console.error('Error archiving project:', error);
      showNotification('Failed to archive project. Please try again.', 'error');
    }
  };

  return (
    <AuthenticatedLayout>
      <Box sx={{ mb: 4, width: '100%' }}>
        <Typography
          variant="h4"
          component="h1"
          gutterBottom
          className="gradient-text-purple"
          sx={{
            fontWeight: 600,
            mb: 1
          }}
        >
          Admin Panel
        </Typography>
        <Typography
          variant="body1"
          paragraph
          sx={{
            color: 'text.secondary',
            maxWidth: '800px',
            fontSize: '1.05rem'
          }}
        >
          Manage projects, users, and system settings. Control access and permissions across the platform.
        </Typography>
      </Box>

      <Paper
        sx={{
          width: '100%',
          p: 2,
          background: (theme) => theme.palette.mode === 'dark'
            ? 'linear-gradient(to bottom, rgba(65, 68, 74, 0.5), rgba(55, 58, 64, 0.5))'
            : 'linear-gradient(to bottom, rgba(255, 255, 255, 0.9), rgba(245, 245, 245, 0.9))',
          backdropFilter: 'blur(10px)',
          boxShadow: '0 4px 20px 0 rgba(0, 0, 0, 0.05)',
          borderRadius: 2
        }}>
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs
            value={tabValue}
            onChange={handleTabChange}
            aria-label="admin tabs"
            sx={{
              '& .MuiTab-root': {
                fontWeight: 500,
                transition: 'all 0.2s ease-in-out',
                '&:hover': {
                  color: 'primary.main',
                  opacity: 0.8
                }
              },
              '& .Mui-selected': {
                fontWeight: 600
              }
            }}
          >
            <Tab label="Projects" />
            <Tab label="Users" />
            <Tab label="Teams" />
            <Tab label="License" />
          </Tabs>
        </Box>

        {/* Projects Tab */}
        <TabPanel value={tabValue} index={0}>
          <ProjectsTab
            projects={projects}
            isLoading={isLoadingProjects}
            error={projectsError}
            onCreateProject={handleOpenProjectDialog}
            onArchiveProject={handleConfirmArchiveProject}
            setSnackbar={({ message, severity }) => showNotification(message, severity)}
          />
        </TabPanel>

        {/* Users Tab */}
        <TabPanel value={tabValue} index={1}>
          <UsersTab
            users={users}
            isLoading={isLoadingUsers}
            error={usersError}
            onCreateUser={handleOpenUserDialog}
            onToggleUserStatus={handleToggleUserStatus}
            onToggleSuperuserStatus={handleToggleSuperuserStatus}
            onDeleteUser={handleConfirmDeleteUser}
          />
        </TabPanel>

        {/* Teams Tab */}
        <TabPanel value={tabValue} index={2}>
          <TeamsTab
            teams={teams}
            users={users}
            isLoading={isLoadingTeams}
            error={teamsError}
            onCreateTeam={handleOpenTeamDialog}
            onDeleteTeam={handleConfirmDeleteTeam}
            onAddTeamMember={handleOpenAddTeamMemberDialog}
            onRemoveTeamMember={handleConfirmRemoveTeamMember}
            onChangeMemberRole={handleChangeMemberRole}
          />
        </TabPanel>



        {/* License Tab */}
        <TabPanel 
          value={tabValue} 
          index={3}
        >
          <LicenseTab />
        </TabPanel>
      </Paper>

      {/* Create Project Dialog */}
      <CreateProjectDialog
        open={openProjectDialog}
        onClose={handleCloseProjectDialog}
        onCreateProject={handleCreateProject}
        teams={teams}
        isLoadingTeams={isLoadingTeams}
      />

      {/* Create User Dialog */}
      <CreateUserDialog
        open={openUserDialog}
        onClose={handleCloseUserDialog}
        onCreateUser={handleCreateUser}
      />

      {/* Create Team Dialog */}
      <CreateTeamDialog
        open={openTeamDialog}
        onClose={handleCloseTeamDialog}
        onCreateTeam={handleCreateTeam}
      />

      {/* Add Team Member Dialog */}
      <AddTeamMemberDialog
        open={openAddTeamMemberDialog}
        onClose={handleCloseAddTeamMemberDialog}
        onAddTeamMember={handleAddTeamMember}
        onBulkAddTeamMembers={handleBulkAddTeamMembers}
        users={users}
        teamId={selectedTeam}
        teamMembers={teams?.find(t => t.id === selectedTeam)?.members || []}
      />

      {/* Delete User Confirmation Dialog */}
      <ConfirmationDialog
        open={deleteUserConfirmation.open}
        title="Confirm Delete User"
        message="Are you sure you want to delete this user? This action cannot be undone."
        onConfirm={handleDeleteUser}
        onCancel={handleCloseDeleteUserConfirmation}
        confirmButtonText="Delete"
      />

      {/* Delete Team Confirmation Dialog */}
      <ConfirmationDialog
        open={deleteTeamConfirmation.open}
        title="Confirm Delete Team"
        message="Are you sure you want to delete this team? This action cannot be undone."
        onConfirm={handleDeleteTeam}
        onCancel={handleCloseDeleteTeamConfirmation}
        confirmButtonText="Delete"
      />

      {/* Remove Team Member Confirmation Dialog */}
      <ConfirmationDialog
        open={removeTeamMemberConfirmation.open}
        title="Confirm Remove Team Member"
        message="Are you sure you want to remove this user from the team? This action cannot be undone."
        onConfirm={handleRemoveTeamMember}
        onCancel={handleCloseRemoveTeamMemberConfirmation}
        confirmButtonText="Remove"
      />

      {/* Archive Project Confirmation Dialog */}
      <ConfirmationDialog
        open={archiveProjectConfirmation.open}
        title="Confirm Archive Project"
        message="Are you sure you want to archive this project? This action cannot be undone."
        onConfirm={handleArchiveProject}
        onCancel={handleCloseArchiveProjectConfirmation}
        confirmButtonText="Archive"
      />
    </AuthenticatedLayout>
  );
};

export default AdminPage;
