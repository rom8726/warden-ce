import React, { useState } from 'react';
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
  Chip,
  Paper,
  Alert,
  Tooltip
} from '@mui/material';
import {
  Delete as DeleteIcon,
  PersonAdd as PersonAddIcon,
  Edit as EditIcon
} from '@mui/icons-material';
import { useQuery, useQueryClient, useMutation } from '@tanstack/react-query';
import apiClient from '../../api/apiClient';
import { AddTeamMemberRequestRoleEnum, ChangeTeamMemberRoleRequestRoleEnum } from '../../generated/api/client';
import AddTeamMemberDialog from './AddTeamMemberDialog';
import ConfirmationDialog from '../admin/ConfirmationDialog';
import ChangeTeamMemberRoleDialog from '../admin/ChangeTeamMemberRoleDialog';
import { useNotification } from '../../App';
import { usePermissions } from '../../hooks/usePermissions';

interface TeamMembersManagementProps {
  teamId: number;
  projectId: number;
}

interface User {
  id: number;
  username: string;
  email: string;
}

interface TeamMember {
  user_id: number;
  role: string;
}

interface Team {
  id: number;
  name: string;
  created_at: string;
  members: TeamMember[];
}

const TeamMembersManagement: React.FC<TeamMembersManagementProps> = ({
  teamId,
  projectId
}) => {
  const { showNotification } = useNotification();
  const { isTeamAdmin, isSuperuser } = usePermissions();
  const queryClient = useQueryClient();
  const [openAddDialog, setOpenAddDialog] = useState(false);
  const [removeConfirmation, setRemoveConfirmation] = useState<{
    open: boolean;
    userId: number | null;
  }>({
    open: false,
    userId: null
  });

  const [changeRoleDialog, setChangeRoleDialog] = useState<{
    open: boolean;
    member: TeamMember;
  } | null>(null);

  // Fetch team details
  const {
    data: team,
    isLoading: isLoadingTeam,
    error: teamError
  } = useQuery<{ team: Team }>({
    queryKey: ['project-team', projectId],
    queryFn: async () => {
      const response = await apiClient.getProjectTeam(projectId);
      return response.data;
    },
    enabled: !!projectId
  });

  // Check if user can manage team
  const canManageTeam = isSuperuser() || (team?.team && isTeamAdmin(team.team.id));

  // Fetch users for team
  const {
    data: users,
    isLoading: isLoadingUsers,
    error: usersError
  } = useQuery<User[]>({
    queryKey: ['users', team?.team?.id],
    queryFn: async () => {
      const response = await apiClient.listUsersForTeam(team!.team.id);
      return response.data;
    },
    enabled: !!team?.team?.id
  });

  // Add team member mutation
  const addTeamMemberMutation = useMutation({
    mutationFn: async ({ userId, role }: { userId: number; role: AddTeamMemberRequestRoleEnum }) => {
      await apiClient.addTeamMember(team!.team.id, {
        user_id: userId,
        role
      });
    },
    onSuccess: () => {
      showNotification('Team member added successfully!', 'success');
      setOpenAddDialog(false);
      queryClient.invalidateQueries({ queryKey: ['project-team', projectId] });
      queryClient.invalidateQueries({ queryKey: ['users', team?.team?.id] });
    },
    onError: (error) => {
      console.error('Error adding team member:', error);
      showNotification('Failed to add team member. Please try again.', 'error');
    }
  });

  // Remove team member mutation
  const removeTeamMemberMutation = useMutation({
    mutationFn: async (userId: number) => {
      await apiClient.removeTeamMember(team!.team.id, userId);
    },
    onSuccess: () => {
      showNotification('Team member removed successfully!', 'success');
      setRemoveConfirmation({ open: false, userId: null });
      queryClient.invalidateQueries({ queryKey: ['project-team', projectId] });
      queryClient.invalidateQueries({ queryKey: ['users', team?.team?.id] });
    },
    onError: (error) => {
      console.error('Error removing team member:', error);
      showNotification('Failed to remove team member. Please try again.', 'error');
    }
  });

  // Change team member role mutation
  const changeMemberRoleMutation = useMutation({
    mutationFn: async ({ userId, role }: { userId: number; role: ChangeTeamMemberRoleRequestRoleEnum }) => {
      await apiClient.changeTeamMemberRole(team!.team.id, userId, { role });
    },
    onSuccess: () => {
      showNotification('Team member role changed successfully!', 'success');
      setChangeRoleDialog(null);
      queryClient.invalidateQueries({ queryKey: ['project-team', projectId] });
      queryClient.invalidateQueries({ queryKey: ['users', team?.team?.id] });
    },
    onError: (error) => {
      console.error('Error changing team member role:', error);
      showNotification('Failed to change team member role. Please try again.', 'error');
      throw error; // Re-throw to let the dialog handle the error
    }
  });

  const handleAddTeamMember = (userId: number, role: AddTeamMemberRequestRoleEnum) => {
    addTeamMemberMutation.mutate({ userId, role });
  };

  const handleRemoveTeamMember = (userId: number) => {
    // Find the member to check their role
    const member = team?.team.members.find(m => m.user_id === userId);

    // Prevent removing team owner
    if (member?.role === 'owner') {
      showNotification('Cannot remove team owner. Only superusers can change team ownership.', 'warning');
      return;
    }

    setRemoveConfirmation({ open: true, userId });
  };

  const handleConfirmRemove = () => {
    if (removeConfirmation.userId) {
      removeTeamMemberMutation.mutate(removeConfirmation.userId);
    }
  };

  const handleCloseRemoveConfirmation = () => {
    setRemoveConfirmation({ open: false, userId: null });
  };

  const handleChangeMemberRole = async (teamId: number, userId: number, role: ChangeTeamMemberRoleRequestRoleEnum) => {
    changeMemberRoleMutation.mutate({ userId, role });
  };

  if (isLoadingTeam || isLoadingUsers) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (teamError || usersError) {
    return (
      <Alert severity="error">
        Error loading team data. Please try again.
      </Alert>
    );
  }

  if (!team) {
    return (
      <Alert severity="warning">
        Team not found.
      </Alert>
    );
  }

  if (team && !canManageTeam) {
    return (
      <Alert severity="warning">
        You don't have permission to manage this team. Only team admins and owners can manage team members.
      </Alert>
    );
  }

  return (
    <Box>
      <Box sx={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center', 
        mb: 3,
        pb: 2,
        borderBottom: (theme) => `1px solid ${theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.05)'}`
      }}>
        <Box>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 0.5 }}>
            <Typography 
              variant="h6"
              sx={{ 
                fontWeight: 600,
                mr: 1
              }}
            >
              Team Members - {team.team.name}
            </Typography>
          </Box>
          <Typography 
            variant="body2" 
            color="text.secondary"
            sx={{ maxWidth: '600px' }}
          >
            Manage team members and their roles for this project.
          </Typography>
        </Box>
        <Button 
          variant="contained" 
          startIcon={<PersonAddIcon />}
          onClick={() => setOpenAddDialog(true)}
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
            transition: 'all 0.2s ease-in-out',
            opacity: 1
          }}
        >
          Add Member
        </Button>
      </Box>

      <Paper sx={{ 
        background: (theme) => theme.palette.mode === 'dark'
          ? 'linear-gradient(135deg, rgba(60, 63, 70, 0.3) 0%, rgba(55, 58, 64, 0.3) 100%)'
          : 'linear-gradient(135deg, rgba(255, 255, 255, 0.7) 0%, rgba(245, 245, 245, 0.7) 100%)',
        borderRadius: 2,
        overflow: 'hidden',
        boxShadow: '0 2px 8px 0 rgba(0, 0, 0, 0.05)',
      }}>
        <Box sx={{ p: 3, pb: 2 }}>
          <Typography variant="subtitle1" sx={{ fontWeight: 600, mb: 2 }}>
            {team.team.name} Team
          </Typography>
        </Box>

        {team.team.members.length > 0 ? (
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow sx={{ 
                  backgroundColor: (theme) => theme.palette.mode === 'dark' 
                    ? 'rgba(60, 63, 70, 0.3)' 
                    : 'rgba(245, 245, 245, 0.5)'
                }}>
                  <TableCell sx={{ fontWeight: 600, py: 1.5 }}>User</TableCell>
                  <TableCell sx={{ fontWeight: 600, py: 1.5 }}>Role</TableCell>
                  <TableCell sx={{ fontWeight: 600, py: 1.5 }}>Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {team.team.members.map((member) => (
                  <TableRow 
                    key={member.user_id}
                    sx={{ 
                      '&:hover': { 
                        backgroundColor: (theme) => theme.palette.mode === 'dark' 
                          ? 'rgba(60, 63, 70, 0.2)' 
                          : 'rgba(245, 245, 245, 0.3)'
                      },
                      transition: 'background-color 0.2s ease-in-out'
                    }}
                  >
                    <TableCell sx={{ py: 1.5 }}>
                      <Box sx={{ display: 'flex', alignItems: 'center' }}>
                        {users?.find(u => u.id === member.user_id)?.username || member.user_id}
                      </Box>
                    </TableCell>
                    <TableCell sx={{ py: 1.5 }}>
                      <Chip 
                        label={member.role} 
                        size="small" 
                        color={
                          member.role === 'owner' 
                            ? 'primary' 
                            : member.role === 'admin' 
                              ? 'secondary' 
                              : 'default'
                        }
                        sx={{ height: 20 }}
                      />
                    </TableCell>
                    <TableCell sx={{ py: 1.5 }}>
                      <IconButton 
                        size="small" 
                        color="primary"
                        onClick={() => setChangeRoleDialog({
                          open: true,
                          member
                        })}
                        sx={{ 
                          transition: 'transform 0.2s ease-in-out',
                          '&:hover': { transform: 'scale(1.1)' },
                          mr: 1
                        }}
                      >
                        <EditIcon fontSize="small" />
                      </IconButton>
                      {member.role !== 'owner' ? (
                        <IconButton 
                          size="small" 
                          color="error"
                          onClick={() => handleRemoveTeamMember(member.user_id)}
                          sx={{ 
                            transition: 'transform 0.2s ease-in-out',
                            '&:hover': { transform: 'scale(1.1)' }
                          }}
                        >
                          <DeleteIcon fontSize="small" />
                        </IconButton>
                      ) : (
                        <Tooltip title="Team owner cannot be removed. Only superusers can change team ownership.">
                          <Box sx={{ display: 'inline-block' }}>
                            <IconButton 
                              size="small" 
                              disabled
                              sx={{ opacity: 0.3 }}
                            >
                              <DeleteIcon fontSize="small" />
                            </IconButton>
                          </Box>
                        </Tooltip>
                      )}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        ) : (
          <Box sx={{ p: 3, textAlign: 'center' }}>
            <Typography variant="body2" color="text.secondary">
              No members in this team. Add members using the "Add Member" button above.
            </Typography>
          </Box>
        )}
      </Paper>

      {/* Add Team Member Dialog */}
      <AddTeamMemberDialog
        open={openAddDialog}
        onClose={() => setOpenAddDialog(false)}
        onAddTeamMember={handleAddTeamMember}
        users={users}
        teamId={teamId}
        teamMembers={team.team.members}
      />

      {/* Remove Team Member Confirmation Dialog */}
      <ConfirmationDialog
        open={removeConfirmation.open}
        title="Confirm Remove Team Member"
        message="Are you sure you want to remove this user from the team? This action cannot be undone."
        onConfirm={handleConfirmRemove}
        onCancel={handleCloseRemoveConfirmation}
        confirmButtonText="Remove"
      />

      {/* Change Role Dialog */}
      {changeRoleDialog && (
        <ChangeTeamMemberRoleDialog
          open={changeRoleDialog.open}
          onClose={() => setChangeRoleDialog(null)}
          onConfirm={handleChangeMemberRole}
          teamId={teamId}
          member={changeRoleDialog.member}
          user={users?.find(u => u.id === changeRoleDialog.member.user_id)}
          teamName={team.team.name}
        />
      )}
    </Box>
  );
};

export default TeamMembersManagement; 
