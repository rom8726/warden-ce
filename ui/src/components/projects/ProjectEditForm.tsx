import React, { useState } from 'react';
import {
  Box,
  Typography,
  TextField,
  Button,
  Paper,
  Alert,
  CircularProgress,
  Divider
} from '@mui/material';
import { Save as SaveIcon, Cancel as CancelIcon } from '@mui/icons-material';
import apiClient from '../../api/apiClient';
import type { UpdateProjectRequest } from '../../generated/api/client';

interface ProjectEditFormProps {
  projectId: number;
  initialName: string;
  initialDescription: string;
  onSave: (name: string, description: string) => void;
  onCancel: () => void;
}

const ProjectEditForm: React.FC<ProjectEditFormProps> = ({
  projectId,
  initialName,
  initialDescription,
  onSave,
  onCancel
}) => {
  const [name, setName] = useState(initialName);
  const [description, setDescription] = useState(initialDescription);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [nameError, setNameError] = useState<string | null>(null);

  const validateForm = () => {
    let isValid = true;
    
    if (!name.trim()) {
      setNameError('Project name is required');
      isValid = false;
    } else if (name.trim().length < 2) {
      setNameError('Project name must be at least 2 characters long');
      isValid = false;
    } else {
      setNameError(null);
    }

    return isValid;
  };

  const handleSave = async () => {
    if (!validateForm()) {
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const updateRequest: UpdateProjectRequest = {
        name: name.trim(),
        description: description.trim()
      };

      const response = await apiClient.updateProject(projectId, updateRequest);
      
      // Call the parent callback with updated values
      onSave(response.data.project.name, response.data.project.description);
    } catch (err: unknown) {
      console.error('Error updating project:', err);
      const errorMessage = err instanceof Error 
        ? err.message 
        : 'Failed to update project. Please try again.';
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCancel = () => {
    // Reset form to initial values
    setName(initialName);
    setDescription(initialDescription);
    setError(null);
    setNameError(null);
    onCancel();
  };

  return (
    <Paper sx={{ p: 3, mb: 3 }}>
      <Typography variant="h6" gutterBottom>
        Edit Project Details
      </Typography>
      
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}

      <Box sx={{ mb: 3 }}>
        <TextField
          fullWidth
          label="Project Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          error={!!nameError}
          helperText={nameError}
          disabled={isLoading}
          sx={{ mb: 2 }}
        />
        
        <TextField
          fullWidth
          label="Description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          multiline
          rows={4}
          disabled={isLoading}
          placeholder="Enter project description..."
        />
      </Box>

      <Divider sx={{ mb: 2 }} />

      <Box sx={{ display: 'flex', gap: 2, justifyContent: 'flex-end' }}>
        <Button
          variant="outlined"
          onClick={handleCancel}
          disabled={isLoading}
          startIcon={<CancelIcon />}
        >
          Cancel
        </Button>
        <Button
          variant="contained"
          onClick={handleSave}
          disabled={isLoading}
          startIcon={isLoading ? <CircularProgress size={20} /> : <SaveIcon />}
        >
          {isLoading ? 'Saving...' : 'Save Changes'}
        </Button>
      </Box>
    </Paper>
  );
};

export default ProjectEditForm; 