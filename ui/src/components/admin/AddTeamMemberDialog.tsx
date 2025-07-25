import React, { useState, useMemo } from 'react';
import {
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  Button,
  Box,
  Typography,
  TableContainer,
  Table,
  TableHead,
  TableBody,
  TableRow,
  TableCell,
  IconButton,
  Paper,
  Radio,
  RadioGroup,
  FormControlLabel,
  FormControl
} from '@mui/material';
import { AddTeamMemberRequestRoleEnum } from '../../generated/api/client';
import type { TeamMember } from '../../generated/api/client/api';

interface User {
  id: number;
  username: string;
  email: string;
}

interface AddTeamMemberDialogProps {
  open: boolean;
  onClose: () => void;
  onAddTeamMember: (userId: number, role: AddTeamMemberRequestRoleEnum) => void;
  onBulkAddTeamMembers?: (userIds: number[], role: AddTeamMemberRequestRoleEnum) => void;
  users: User[] | undefined;
  teamId: number | null;
  teamMembers?: TeamMember[];
}

const AddTeamMemberDialog: React.FC<AddTeamMemberDialogProps> = ({
  open,
  onClose,
  onAddTeamMember,
  users,
  teamMembers = []
}) => {
  const [selectedUser, setSelectedUser] = useState<number | null>(null);
  const [selectedRole, setSelectedRole] = useState<AddTeamMemberRequestRoleEnum>(AddTeamMemberRequestRoleEnum.Member);

  // Filter out users who are already team members
  const availableUsers = useMemo(() => {
    if (!users) return [];
    
    const existingMemberIds = new Set(teamMembers.map(member => member.user_id));
    return users.filter(user => !existingMemberIds.has(user.id));
  }, [users, teamMembers]);

  const handleAdd = async () => {
    if (selectedUser) {
      // Errors are handled in the parent component
      await onAddTeamMember(selectedUser, selectedRole);
      // Reset form and close dialog for single user addition
      setSelectedUser(null);
      setSelectedRole(AddTeamMemberRequestRoleEnum.Member);
      onClose();
    }
  };

  const handleCancel = () => {
    // Reset form
    setSelectedUser(null);
    setSelectedRole(AddTeamMemberRequestRoleEnum.Member);
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleCancel} maxWidth="md" fullWidth>
      <DialogTitle>Add User to Team</DialogTitle>
      <DialogContent>
        <DialogContentText>
          Select users and role to add to the team.
        </DialogContentText>
        
        <Box sx={{ mt: 2 }}>
          <Typography variant="subtitle2" gutterBottom>
            Select User
          </Typography>
          {availableUsers.length === 0 ? (
            <Typography variant="body2" color="text.secondary" sx={{ p: 2, textAlign: 'center' }}>
              All users are already members of this team.
            </Typography>
          ) : (
            <TableContainer component={Paper} sx={{ maxHeight: 200, mb: 2 }}>
              <Table size="small">
                <TableHead>
                  <TableRow>
                    <TableCell>Username</TableCell>
                    <TableCell>Email</TableCell>
                    <TableCell>Select</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {availableUsers.map((user) => (
                    <TableRow key={user.id}>
                      <TableCell>{user.username}</TableCell>
                      <TableCell>{user.email}</TableCell>
                      <TableCell>
                        <IconButton 
                          size="small" 
                          color={selectedUser === user.id ? "primary" : "default"}
                          onClick={() => setSelectedUser(user.id)}
                        >
                          {selectedUser === user.id ? (
                            <span>✓</span>
                          ) : (
                            <span>○</span>
                          )}
                        </IconButton>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          )}
        </Box>

        {/* Role Selection */}
        <Box sx={{ mt: 3 }}>
          <Typography variant="subtitle2" gutterBottom>
            Select Role
          </Typography>
          <FormControl component="fieldset">
            <RadioGroup
              value={selectedRole}
              onChange={(e) => setSelectedRole(e.target.value as AddTeamMemberRequestRoleEnum)}
              row
            >
              <FormControlLabel
                value={AddTeamMemberRequestRoleEnum.Owner}
                control={<Radio />}
                label="Owner"
              />
              <FormControlLabel
                value={AddTeamMemberRequestRoleEnum.Admin}
                control={<Radio />}
                label="Admin"
              />
              <FormControlLabel
                value={AddTeamMemberRequestRoleEnum.Member}
                control={<Radio />}
                label="Member"
              />
            </RadioGroup>
          </FormControl>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button 
          onClick={handleCancel}
          color="primary"
        >
          Cancel
        </Button>
        <Button 
          onClick={handleAdd} 
          variant="contained"
          color="primary"
          disabled={!selectedUser || availableUsers.length === 0}
        >
          Add to Team
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default AddTeamMemberDialog;