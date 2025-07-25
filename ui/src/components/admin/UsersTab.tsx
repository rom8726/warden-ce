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
  FormControlLabel,
  Switch
} from '@mui/material';
import {
  Add as AddIcon,
  Delete as DeleteIcon
} from '@mui/icons-material';

interface User {
  id: number;
  username: string;
  email: string;
  is_active: boolean;
  is_superuser: boolean;
}

interface UsersTabProps {
  users: User[] | undefined;
  isLoading: boolean;
  error: unknown;
  onCreateUser: () => void;
  onToggleUserStatus: (userId: number, isActive: boolean) => void;
  onToggleSuperuserStatus: (userId: number, isSuperuser: boolean) => void;
  onDeleteUser: (userId: number) => void;
}

const UsersTab: React.FC<UsersTabProps> = ({
  users,
  isLoading,
  error,
  onCreateUser,
  onToggleUserStatus,
  onToggleSuperuserStatus,
  onDeleteUser,
}) => {
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
            className="gradient-subtitle-purple"
          >
            Manage Users
          </Typography>
          <Typography 
            variant="body2" 
            color="text.secondary"
            sx={{ maxWidth: '600px' }}
          >
            Create and manage user accounts and permissions.
          </Typography>
        </Box>
        <Button 
          variant="contained" 
          startIcon={<AddIcon />}
          onClick={onCreateUser}
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
          Create User
        </Button>
      </Box>

      {isLoading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
          <CircularProgress />
        </Box>
      ) : error ? (
        <Typography color="error">
          Error loading users. Please try again.
        </Typography>
      ) : users && users.length > 0 ? (
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>Username</TableCell>
                <TableCell>Email</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Role</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {users.map((user) => (
                <TableRow key={user.id}>
                  <TableCell>{user.id}</TableCell>
                  <TableCell>{user.username}</TableCell>
                  <TableCell>{user.email}</TableCell>
                  <TableCell>
                    <FormControlLabel
                      control={
                        <Switch
                          checked={user.is_active}
                          onChange={() => onToggleUserStatus(user.id, !user.is_active)}
                          color="primary"
                          size="small"
                        />
                      }
                      label={user.is_active ? "Active" : "Inactive"}
                    />
                  </TableCell>
                  <TableCell>
                    <FormControlLabel
                      control={
                        <Switch
                          checked={user.is_superuser}
                          onChange={() => onToggleSuperuserStatus(user.id, !user.is_superuser)}
                          color="primary"
                          size="small"
                        />
                      }
                      label={user.is_superuser ? "Superuser" : "Regular"}
                    />
                  </TableCell>
                  <TableCell>
                    <IconButton 
                      size="small" 
                      color="error"
                      onClick={() => onDeleteUser(user.id)}
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
          No users to display. Create a new user to get started.
        </Typography>
      )}
    </>
  );
};

export default UsersTab;