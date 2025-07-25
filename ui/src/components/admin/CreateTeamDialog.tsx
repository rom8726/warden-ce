import React, { useState } from 'react';
import {
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  TextField,
  Button
} from '@mui/material';

interface CreateTeamDialogProps {
  open: boolean;
  onClose: () => void;
  onCreateTeam: (name: string) => void;
}

const CreateTeamDialog: React.FC<CreateTeamDialogProps> = ({
  open,
  onClose,
  onCreateTeam
}) => {
  const [teamName, setTeamName] = useState('');

  const handleCreate = () => {
    onCreateTeam(teamName);
    // Reset form
    setTeamName('');
  };

  const handleCancel = () => {
    // Reset form
    setTeamName('');
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleCancel}>
      <DialogTitle className="gradient-text-purple">Create New Team</DialogTitle>
      <DialogContent>
        <DialogContentText>
          Enter the name for the new team.
        </DialogContentText>
        <TextField
          autoFocus
          margin="dense"
          id="teamName"
          label="Team Name"
          type="text"
          fullWidth
          variant="outlined"
          value={teamName}
          onChange={(e) => setTeamName(e.target.value)}
          sx={{ mb: 2 }}
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
          disabled={!teamName.trim()}
        >
          Create
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default CreateTeamDialog;