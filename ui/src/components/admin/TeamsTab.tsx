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
} from '@mui/material';
import {
  Add as AddIcon,
  Delete as DeleteIcon,
  PersonAdd as PersonAddIcon,
  Edit as EditIcon
} from '@mui/icons-material';
import ChangeTeamMemberRoleDialog from './ChangeTeamMemberRoleDialog';
import { ChangeTeamMemberRoleRequestRoleEnum } from '../../generated/api/client';

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

interface User {
  id: number;
  username: string;
  email: string;
}

interface TeamsTabProps {
  teams: Team[] | undefined;
  users: User[] | undefined;
  isLoading: boolean;
  error: unknown;
  onCreateTeam: () => void;
  onDeleteTeam: (teamId: number) => void;
  onAddTeamMember: (teamId: number) => void;
  onRemoveTeamMember: (teamId: number, userId: number) => void;
  onChangeMemberRole: (teamId: number, userId: number, role: ChangeTeamMemberRoleRequestRoleEnum) => Promise<void>;
}

const TeamsTab: React.FC<TeamsTabProps> = ({
  teams,
  users,
  isLoading,
  error,
  onCreateTeam,
  onDeleteTeam,
  onAddTeamMember,
  onRemoveTeamMember,
  onChangeMemberRole
}) => {
  const [changeRoleDialog, setChangeRoleDialog] = useState<{
    open: boolean;
    teamId: number;
    member: TeamMember;
  } | null>(null);
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
            Manage Teams
          </Typography>
          <Typography 
            variant="body2" 
            color="text.secondary"
            sx={{ maxWidth: '600px' }}
          >
            Create and manage teams and team memberships.
          </Typography>
        </Box>
        <Button 
          variant="contained" 
          startIcon={<AddIcon />}
          onClick={onCreateTeam}
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
          Create Team
        </Button>
      </Box>

      {isLoading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
          <CircularProgress />
        </Box>
      ) : error ? (
        <Typography color="error">
          Error loading teams. Please try again.
        </Typography>
      ) : teams && teams.length > 0 ? (
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>Name</TableCell>
                <TableCell>Created At</TableCell>
                <TableCell>Members</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {teams.map((team) => (
                <TableRow 
                  key={team.id}
                  sx={{
                    backgroundColor: 'inherit'
                  }}
                >
                  <TableCell>{team.id}</TableCell>
                  <TableCell>
                    <Box sx={{ display: 'flex', alignItems: 'center' }}>
                      <Typography variant="body2" sx={{ mr: 1 }}>
                        {team.name}
                      </Typography>
                    </Box>
                  </TableCell>
                  <TableCell>{new Date(team.created_at).toLocaleString()}</TableCell>
                  <TableCell>{team.members.length}</TableCell>
                  <TableCell>
                    <IconButton 
                      size="small" 
                      color="primary"
                      onClick={() => onAddTeamMember(team.id)}
                    >
                      <PersonAddIcon fontSize="small" />
                    </IconButton>
                    <IconButton 
                      size="small" 
                      color="error"
                      onClick={() => onDeleteTeam(team.id)}
                    >
                      <DeleteIcon fontSize="small" />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      ) : (
        <Typography variant="body2" sx={{ p: 3 }}>
          No teams to display. Create a new team to get started.
        </Typography>
      )}

      {/* Team Members Section */}
      {teams && teams.length > 0 && (
        <Box sx={{ mt: 5, pt: 3, borderTop: (theme) => `1px solid ${theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.05)'}` }}>
          <Box sx={{ mb: 3 }}>
            <Typography 
              variant="h6" 
              sx={{ 
                fontWeight: 600,
                mb: 0.5
              }}
            >
              Team Members
            </Typography>
            <Typography 
              variant="body2" 
              color="text.secondary"
              sx={{ maxWidth: '600px' }}
            >
              View and manage team memberships across all teams.
            </Typography>
          </Box>

          {teams.map((team) => (
            <Box 
              key={team.id} 
              sx={{ 
                mb: 4,
                background: (theme) => theme.palette.mode === 'dark'
                  ? 'linear-gradient(135deg, rgba(60, 63, 70, 0.3) 0%, rgba(55, 58, 64, 0.3) 100%)'
                  : 'linear-gradient(135deg, rgba(255, 255, 255, 0.7) 0%, rgba(245, 245, 245, 0.7) 100%)',
                borderRadius: 2,
                overflow: 'hidden',
                boxShadow: '0 2px 8px 0 rgba(0, 0, 0, 0.05)',
              }}
            >
              <Box 
                sx={{ 
                  p: 2, 
                  borderBottom: (theme) => `1px solid ${theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.05)'}`,
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  background: (theme) => theme.palette.mode === 'dark'
                    ? 'rgba(50, 53, 58, 0.4)'
                    : 'rgba(235, 235, 235, 0.4)',
                }}
              >
                <Box sx={{ display: 'flex', alignItems: 'center' }}>
                  <Box 
                    component="span" 
                    sx={{ 
                      width: 8, 
                      height: 8, 
                      borderRadius: '50%', 
                      bgcolor: 'primary.main',
                      mr: 1.5,
                      display: 'inline-block'
                    }} 
                  />
                  <Typography 
                    variant="subtitle1" 
                    sx={{ 
                      fontWeight: 600,
                      display: 'flex',
                      alignItems: 'center',
                      mr: 1
                    }}
                  >
                    {team.name}
                  </Typography>
                </Box>
                <Chip 
                  label={`${team.members.length} ${team.members.length === 1 ? 'member' : 'members'}`} 
                  size="small" 
                  color="primary" 
                  variant="outlined"
                />
              </Box>

              {team.members.length > 0 ? (
                <TableContainer sx={{ mb: 0 }}>
                  <Table size="small">
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
                      {team.members.map((member) => (
                        <TableRow 
                          key={`${team.id}-${member.user_id}`}
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
                                teamId: team.id,
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
                            <IconButton 
                              size="small" 
                              color="error"
                              onClick={() => onRemoveTeamMember(team.id, member.user_id)}
                              sx={{ 
                                transition: 'transform 0.2s ease-in-out',
                                '&:hover': { transform: 'scale(1.1)' }
                              }}
                            >
                              <DeleteIcon fontSize="small" />
                            </IconButton>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </TableContainer>
              ) : (
                <Box sx={{ p: 3, textAlign: 'center' }}>
                  <Typography variant="body2" color="text.secondary">
                    No members in this team. Add members using the "Add User" button above.
                  </Typography>
                </Box>
              )}
            </Box>
          ))}
        </Box>
      )}

      {/* Change Role Dialog */}
      {changeRoleDialog && (
        <ChangeTeamMemberRoleDialog
          open={changeRoleDialog.open}
          onClose={() => setChangeRoleDialog(null)}
          onConfirm={onChangeMemberRole}
          teamId={changeRoleDialog.teamId}
          member={changeRoleDialog.member}
          user={users?.find(u => u.id === changeRoleDialog.member.user_id)}
          teamName={teams?.find(t => t.id === changeRoleDialog.teamId)?.name || ''}
        />
      )}
    </>
  );
};

export default TeamsTab;