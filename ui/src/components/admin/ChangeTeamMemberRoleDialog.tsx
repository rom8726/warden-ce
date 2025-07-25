import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Typography,
  Box,
  Chip,
  CircularProgress,
  Alert,
} from '@mui/material';
import { 
  ChangeTeamMemberRoleRequestRoleEnum
} from '../../generated/api/client';
import type { 
  ChangeTeamMemberRoleRequest 
} from '../../generated/api/client';

interface TeamMember {
  user_id: number;
  role: string;
}

interface User {
  id: number;
  username: string;
  email: string;
}

interface ChangeTeamMemberRoleDialogProps {
  open: boolean;
  onClose: () => void;
  onConfirm: (teamId: number, userId: number, role: ChangeTeamMemberRoleRequestRoleEnum) => Promise<void>;
  teamId: number;
  member: TeamMember;
  user: User | undefined;
  teamName: string;
}

const ChangeTeamMemberRoleDialog: React.FC<ChangeTeamMemberRoleDialogProps> = ({
  open,
  onClose,
  onConfirm,
  teamId,
  member,
  user,
  teamName,
}) => {
  const [selectedRole, setSelectedRole] = useState<ChangeTeamMemberRoleRequestRoleEnum>(
    member.role as ChangeTeamMemberRoleRequestRoleEnum
  );
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleConfirm = async () => {
    if (selectedRole === member.role) {
      onClose();
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      await onConfirm(teamId, member.user_id, selectedRole);
      onClose();
    } catch (err: any) {
      setError(err?.response?.data?.error?.message || 'Failed to change member role');
    } finally {
      setIsLoading(false);
    }
  };

  const handleClose = () => {
    if (!isLoading) {
      setSelectedRole(member.role as ChangeTeamMemberRoleRequestRoleEnum);
      setError(null);
      onClose();
    }
  };

  const getRoleColor = (role: string) => {
    switch (role) {
      case 'owner':
        return 'primary';
      case 'admin':
        return 'secondary';
      default:
        return 'default';
    }
  };

  const getRoleDescription = (role: string) => {
    switch (role) {
      case 'owner':
        return 'Full control over team and projects';
      case 'admin':
        return 'Can manage team members and projects';
      case 'member':
        return 'Basic access to team projects';
      default:
        return '';
    }
  };

  return (
    <Dialog 
      open={open} 
      onClose={handleClose}
      maxWidth="sm"
      fullWidth
    >
      <DialogTitle className="gradient-text-purple">
        Change Team Member Role
      </DialogTitle>
      
      <DialogContent>
        <Box sx={{ mb: 3 }}>
          <Typography variant="body2" color="text.secondary" gutterBottom>
            Team: <strong>{teamName}</strong>
          </Typography>
          <Typography variant="body2" color="text.secondary" gutterBottom>
            Member: <strong>{user?.username || `User ${member.user_id}`}</strong>
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Current role: 
            <Chip 
              label={member.role} 
              size="small" 
              color={getRoleColor(member.role)}
              sx={{ ml: 1, height: 20 }}
            />
          </Typography>
        </Box>

        <FormControl fullWidth sx={{ mb: 2 }}>
          <InputLabel>New Role</InputLabel>
          <Select
            value={selectedRole}
            label="New Role"
            onChange={(e) => setSelectedRole(e.target.value as ChangeTeamMemberRoleRequestRoleEnum)}
            disabled={isLoading}
          >
            <MenuItem value={ChangeTeamMemberRoleRequestRoleEnum.Owner}>
              <Box>
                <Typography variant="body2" sx={{ fontWeight: 500 }}>
                  Owner
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  {getRoleDescription('owner')}
                </Typography>
              </Box>
            </MenuItem>
            <MenuItem value={ChangeTeamMemberRoleRequestRoleEnum.Admin}>
              <Box>
                <Typography variant="body2" sx={{ fontWeight: 500 }}>
                  Admin
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  {getRoleDescription('admin')}
                </Typography>
              </Box>
            </MenuItem>
            <MenuItem value={ChangeTeamMemberRoleRequestRoleEnum.Member}>
              <Box>
                <Typography variant="body2" sx={{ fontWeight: 500 }}>
                  Member
                </Typography>
                <Typography variant="caption" color="text.secondary">
                  {getRoleDescription('member')}
                </Typography>
              </Box>
            </MenuItem>
          </Select>
        </FormControl>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        <Box sx={{ 
          p: 2, 
          bgcolor: 'background.paper', 
          borderRadius: 1,
          border: '1px solid',
          borderColor: 'divider'
        }}>
          <Typography variant="body2" color="text.secondary">
            <strong>Note:</strong> Changing a member's role will affect their permissions 
            within this team. Owners have full control, admins can manage members and projects, 
            and members have basic access to team projects.
          </Typography>
        </Box>
      </DialogContent>

      <DialogActions>
        <Button 
          onClick={handleClose} 
          disabled={isLoading}
        >
          Cancel
        </Button>
        <Button 
          onClick={handleConfirm}
          variant="contained"
          disabled={isLoading || selectedRole === member.role}
          startIcon={isLoading ? <CircularProgress size={16} /> : null}
        >
          {isLoading ? 'Changing...' : 'Change Role'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default ChangeTeamMemberRoleDialog; 