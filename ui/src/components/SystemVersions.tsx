import React from 'react';
import {
  Box,
  Typography,
  Card,
  CardContent,
  Grid,
  Chip,
  CircularProgress,
  Alert,
  Tooltip,
} from '@mui/material';
import { useVersions } from '../hooks/useVersions';
import { ComponentVersionStatusEnum } from '../generated/api/client';

interface SystemVersionsProps {
  showBuildTime?: boolean;
  showStatus?: boolean;
}

export const SystemVersions: React.FC<SystemVersionsProps> = ({
  showBuildTime = true,
  showStatus = true,
}) => {
  const { data: versions, isLoading, error } = useVersions();

  const formatBuildTime = (buildTime: string) => {
    if (!buildTime || buildTime === 'unknown') {
      return 'Unknown';
    }
    
    try {
      // Try to parse the date directly
      const date = new Date(buildTime);
      
      // Check if the date is valid
      if (!isNaN(date.getTime())) {
        return date.toLocaleString();
      }
      
      // If direct parsing failed, try to clean up the string
      // Remove any non-standard characters but keep the ISO format
      const cleanedTime = buildTime.replace(/[^\d\-T:Z]/g, '');
      const cleanedDate = new Date(cleanedTime);
      
      if (!isNaN(cleanedDate.getTime())) {
        return cleanedDate.toLocaleString();
      }
      
      // If still invalid, return the original string
      return buildTime;
    } catch {
      return buildTime;
    }
  };

  const getStatusColor = (status: ComponentVersionStatusEnum) => {
    switch (status) {
      case ComponentVersionStatusEnum.Available:
        return 'success';
      case ComponentVersionStatusEnum.Unavailable:
        return 'error';
      default:
        return 'default';
    }
  };

  const getStatusText = (status: ComponentVersionStatusEnum) => {
    switch (status) {
      case ComponentVersionStatusEnum.Available:
        return 'Available';
      case ComponentVersionStatusEnum.Unavailable:
        return 'Unavailable';
      default:
        return 'Unknown';
    }
  };

  if (isLoading) {
    return (
      <Box display="flex" justifyContent="center" p={2}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Alert severity="error">
        Failed to load system versions. Please try again later.
      </Alert>
    );
  }

  if (!versions?.components || versions.components.length === 0) {
    return (
      <Alert severity="info">
        No version information available.
      </Alert>
    );
  }

  return (
    <Box>
      <Typography variant="h6" gutterBottom>
        System Components Versions
      </Typography>
      
      {versions.collected_at && (
        <Typography variant="caption" color="text.secondary" gutterBottom>
          Last updated: {formatBuildTime(versions.collected_at)}
        </Typography>
      )}

      <Grid container spacing={2} sx={{ mt: 1 }}>
        {versions.components.map((component) => (
          <Grid item xs={12} sm={6} md={4} key={component.name}>
            <Card variant="outlined">
              <CardContent>
                <Box display="flex" justifyContent="space-between" alignItems="flex-start" mb={1}>
                  <Typography variant="subtitle2" fontWeight="bold">
                    {component.name}
                  </Typography>
                  {showStatus && (
                    <Chip
                      label={getStatusText(component.status)}
                      color={getStatusColor(component.status)}
                      size="small"
                    />
                  )}
                </Box>
                
                <Typography variant="body2" color="primary" gutterBottom>
                  {component.version}
                </Typography>
                
                {showBuildTime && component.build_time && (
                  <Tooltip title="Build time">
                    <Typography variant="caption" color="text.secondary">
                      Built: {formatBuildTime(component.build_time)}
                    </Typography>
                  </Tooltip>
                )}
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Box>
  );
};

export default SystemVersions; 