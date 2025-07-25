import React from 'react';
import { 
  Box, 
  CircularProgress, 
  Typography,
  useTheme
} from '@mui/material';

interface LoadingSpinnerProps {
  message?: string;
  size?: 'small' | 'medium' | 'large';
  fullHeight?: boolean;
  overlay?: boolean;
}

const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({ 
  message = 'Loading...', 
  size = 'medium',
  fullHeight = false,
  overlay = false
}) => {
  const theme = useTheme();

  const getSize = () => {
    switch (size) {
      case 'small':
        return 24;
      case 'large':
        return 48;
      default:
        return 32;
    }
  };

  const content = (
    <Box sx={{ 
      display: 'flex', 
      flexDirection: 'column', 
      alignItems: 'center', 
      justifyContent: 'center',
      gap: 2,
      ...(fullHeight && { height: '100%' }),
      ...(overlay && { 
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundColor: theme.palette.mode === 'dark' 
          ? 'rgba(0, 0, 0, 0.7)' 
          : 'rgba(255, 255, 255, 0.8)',
        backdropFilter: 'blur(4px)',
        zIndex: 1000,
      }),
    }}>
      <CircularProgress 
        size={getSize()} 
        sx={{ 
          color: theme.palette.primary.main,
        }} 
      />
      {message && (
        <Typography 
          variant="body2" 
          color="text.secondary"
          sx={{ 
            textAlign: 'center',
            fontWeight: 500,
          }}
        >
          {message}
        </Typography>
      )}
    </Box>
  );

  if (overlay) {
    return (
      <Box sx={{ position: 'relative' }}>
        {content}
      </Box>
    );
  }

  return content;
};

export default LoadingSpinner; 