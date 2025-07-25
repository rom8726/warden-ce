import React, { useState } from 'react';
import {
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  TextField,
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
  CircularProgress
} from '@mui/material';

interface Team {
  id: number;
  name: string;
  members: { user_id: number; role: string }[];
}

interface CreateProjectDialogProps {
  open: boolean;
  onClose: () => void;
  onCreateProject: (name: string, description: string, teamId: number | null) => void;
  teams: Team[] | undefined;
  isLoadingTeams: boolean;
}

const CreateProjectDialog: React.FC<CreateProjectDialogProps> = ({
  open,
  onClose,
  onCreateProject,
  teams,
  isLoadingTeams
}) => {
  const [projectName, setProjectName] = useState('');
  const [projectDescription, setProjectDescription] = useState('');
  const [selectedTeamId, setSelectedTeamId] = useState<number | null>(null);
  const [descriptionError, setDescriptionError] = useState<string>('');

  // Validate description length
  const validateDescription = (description: string) => {
    if (description.trim().length < 10) {
      setDescriptionError('Description must be at least 10 characters long');
      return false;
    } else {
      setDescriptionError('');
      return true;
    }
  };

  const handleDescriptionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setProjectDescription(value);
    validateDescription(value);
  };

  // Reset error when dialog opens
  React.useEffect(() => {
    if (open) {
      setDescriptionError('');
    }
  }, [open]);

  const handleCreate = () => {
    // Validate description before creating
    if (!validateDescription(projectDescription)) {
      return;
    }
    
    onCreateProject(projectName, projectDescription, selectedTeamId);
    // Reset form
    setProjectName('');
    setProjectDescription('');
    setSelectedTeamId(null);
    setDescriptionError('');
  };

  const handleCancel = () => {
    // Reset form
    setProjectName('');
    setProjectDescription('');
    setSelectedTeamId(null);
    setDescriptionError('');
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleCancel}>
      <DialogTitle>Create New Project</DialogTitle>
      <DialogContent>
        <DialogContentText>
          Enter the name for the new project and select a team (optional).
        </DialogContentText>
        <TextField
          autoFocus
          margin="dense"
          id="name"
          label="Project Name"
          type="text"
          fullWidth
          variant="outlined"
          value={projectName}
          onChange={(e) => setProjectName(e.target.value)}
        />
        <TextField
          margin="dense"
          id="description"
          label="Project Description"
          type="text"
          fullWidth
          variant="outlined"
          value={projectDescription}
          onChange={handleDescriptionChange}
          multiline
          minRows={2}
          required
          error={!!descriptionError}
          helperText={descriptionError || 'Enter a detailed description of the project (minimum 10 characters)'}
        />
        <Box sx={{ mt: 2 }}>
          <Typography variant="subtitle2" gutterBottom>
            Select Team (Optional)
          </Typography>
          {isLoadingTeams ? (
            <CircularProgress size={24} />
          ) : teams && teams.length > 0 ? (
            <TableContainer component={Paper} sx={{ maxHeight: 200, mb: 2 }}>
              <Table size="small">
                <TableHead>
                  <TableRow>
                    <TableCell>Team Name</TableCell>
                    <TableCell>Members</TableCell>
                    <TableCell>Select</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  <TableRow data-testid="no-team-row">
                    <TableCell colSpan={2}>No team (default)</TableCell>
                    <TableCell>
                      <IconButton 
                        size="small" 
                        color={selectedTeamId === null ? "primary" : "default"}
                        onClick={() => setSelectedTeamId(null)}
                        data-testid="no-team-select-button"
                      >
                        {selectedTeamId === null ? (
                          <span>✓</span>
                        ) : (
                          <span>○</span>
                        )}
                      </IconButton>
                    </TableCell>
                  </TableRow>
                  {teams.map((team) => (
                    <TableRow key={team.id}>
                      <TableCell>{team.name}</TableCell>
                      <TableCell>{team.members.length}</TableCell>
                      <TableCell>
                        <IconButton 
                          size="small" 
                          color={selectedTeamId === team.id ? "primary" : "default"}
                          onClick={() => setSelectedTeamId(team.id)}
                        >
                          {selectedTeamId === team.id ? (
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
          ) : (
            <Typography variant="body2">
              No teams available. You can create a project without a team.
            </Typography>
          )}
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
          onClick={handleCreate} 
          variant="contained"
          color="primary"
          disabled={!projectName.trim() || !projectDescription.trim() || !!descriptionError}
        >
          Create
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default CreateProjectDialog;
