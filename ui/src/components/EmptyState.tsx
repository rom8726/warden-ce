import React from 'react';
import { 
  Box, 
  Typography, 
  Button,
  useTheme
} from '@mui/material';
import { 
  Inbox as InboxIcon,
  BugReportOutlined as IssuesIcon,
  FolderOutlined as ProjectsIcon,
  Dashboard as DashboardIcon
} from '@mui/icons-material';

interface EmptyStateProps {
  title: string;
  description: string;
  icon?: React.ReactNode;
  action?: {
    label: string;
    onClick: () => void;
    variant?: 'contained' | 'outlined' | 'text';
  };
  type?: 'default' | 'issues' | 'projects' | 'dashboard';
}

const EmptyState: React.FC<EmptyStateProps> = ({ 
  title, 
  description, 
  icon, 
  action,
  type = 'default'
}) => {
  const theme = useTheme();

  const getDefaultIcon = () => {
    switch (type) {
      case 'issues':
        return <IssuesIcon />;
      case 'projects':
        return <ProjectsIcon />;
      case 'dashboard':
        return <DashboardIcon />;
      default:
        return <InboxIcon />;
    }
  };

  const displayIcon = icon || getDefaultIcon();

  return (
    <Box sx={{ 
      display: 'flex', 
      flexDirection: 'column', 
      alignItems: 'center', 
      justifyContent: 'center',
      py: 8,
      px: 2,
      textAlign: 'center',
      minHeight: 300,
    }}>
      <Box sx={{ 
        display: 'flex', 
        alignItems: 'center', 
        justifyContent: 'center',
        width: 80,
        height: 80,
        borderRadius: '50%',
        backgroundColor: theme.palette.mode === 'dark' 
          ? 'rgba(130, 82, 255, 0.1)' 
          : 'rgba(130, 82, 255, 0.08)',
        color: theme.palette.primary.main,
        mb: 3,
      }}>
        {React.cloneElement(displayIcon as React.ReactElement, { 
          sx: { fontSize: 40 } 
        } as any)}
      </Box>
      
      <Typography 
        variant="h5" 
        component="h2" 
        gutterBottom
        sx={{ 
          fontWeight: 600,
          color: theme.palette.text.primary,
          mb: 1,
        }}
      >
        {title}
      </Typography>
      
      <Typography 
        variant="body1" 
        color="text.secondary"
        sx={{ 
          maxWidth: 400,
          mb: action ? 3 : 0,
          lineHeight: 1.6,
        }}
      >
        {description}
      </Typography>
      
      {action && (
        <Button
          variant={action.variant || 'contained'}
          onClick={action.onClick}
          sx={{
            px: 3,
            py: 1.5,
            borderRadius: 2,
            textTransform: 'none',
            fontWeight: 500,
          }}
        >
          {action.label}
        </Button>
      )}
    </Box>
  );
};

export default EmptyState; 