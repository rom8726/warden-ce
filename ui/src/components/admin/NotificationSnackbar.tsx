import React from 'react';
import { Snackbar, Alert } from '@mui/material';

interface NotificationSnackbarProps {
  open: boolean;
  message: string;
  severity: 'success' | 'error' | 'info' | 'warning';
  onClose: () => void;
}

const NotificationSnackbar: React.FC<NotificationSnackbarProps> = ({
  open,
  message,
  severity,
  onClose
}) => {
  return (
    <Snackbar
      open={open}
      autoHideDuration={6000}
      onClose={onClose}
      anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      sx={{
        zIndex: 1400,
        '& .MuiAlert-root': {
          minWidth: 300,
          maxWidth: 600,
          boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)',
          borderRadius: 2,
          '& .MuiAlert-message': {
            fontSize: '0.875rem',
            lineHeight: 1.4,
            fontWeight: 500,
            color: 'inherit'
          }
        }
      }}
    >
      <Alert 
        onClose={onClose} 
        severity={severity}
        sx={{ 
          width: '100%',
          '& .MuiAlert-message': {
            color: 'inherit',
            fontWeight: 500
          }
        }}
      >
        {message}
      </Alert>
    </Snackbar>
  );
};

export default NotificationSnackbar;