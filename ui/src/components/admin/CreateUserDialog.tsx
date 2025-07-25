import React, { useState } from 'react';
import {
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  TextField,
  Button,
  FormControlLabel,
  Switch
} from '@mui/material';

interface CreateUserDialogProps {
  open: boolean;
  onClose: () => void;
  onCreateUser: (username: string, email: string, password: string, isSuperuser: boolean) => void;
}

const CreateUserDialog: React.FC<CreateUserDialogProps> = ({
  open,
  onClose,
  onCreateUser
}) => {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isSuperuser, setIsSuperuser] = useState(false);

  const handleCreate = () => {
    onCreateUser(username, email, password, isSuperuser);
    // Reset form
    setUsername('');
    setEmail('');
    setPassword('');
    setIsSuperuser(false);
  };

  const handleCancel = () => {
    // Reset form
    setUsername('');
    setEmail('');
    setPassword('');
    setIsSuperuser(false);
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleCancel}>
      <DialogTitle className="gradient-text-purple">Create New User</DialogTitle>
      <DialogContent>
        <DialogContentText>
          Enter the details for the new user.
        </DialogContentText>
        <TextField
          autoFocus
          margin="dense"
          id="username"
          label="Username"
          type="text"
          fullWidth
          variant="outlined"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <TextField
          margin="dense"
          id="email"
          label="Email Address"
          type="email"
          fullWidth
          variant="outlined"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
        <TextField
          margin="dense"
          id="password"
          label="Password"
          type="password"
          fullWidth
          variant="outlined"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        <FormControlLabel
          control={
            <Switch
              checked={isSuperuser}
              onChange={(e) => setIsSuperuser(e.target.checked)}
              color="primary"
            />
          }
          label="Superuser"
          sx={{ mt: 2 }}
        />
      </DialogContent>
      <DialogActions>
        <Button 
          onClick={handleCancel}
          color="primary"
        >
          Cancel
        </Button>
        <Button 
          onClick={handleCreate} 
          variant="contained"
          color="primary"
          disabled={!username.trim() || !email.trim() || !password.trim()}
        >
          Create
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default CreateUserDialog;