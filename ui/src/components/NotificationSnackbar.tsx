import React, { useState, useEffect } from 'react';
import { Alert, Snackbar, type AlertColor } from '@mui/material';
import { CheckCircle, Error, Info, Warning } from '@mui/icons-material';

export interface Notification {
  id: string;
  message: string;
  type: AlertColor;
  duration?: number;
  action?: React.ReactNode;
}

interface NotificationSnackbarProps {
  notification: Notification | null;
  onClose: (id: string) => void;
  position?: 'top' | 'bottom';
}

const getIcon = (type: AlertColor) => {
  switch (type) {
    case 'success':
      return <CheckCircle />;
    case 'error':
      return <Error />;
    case 'warning':
      return <Warning />;
    case 'info':
      return <Info />;
    default:
      return <Info />;
  }
};

export const NotificationSnackbar: React.FC<NotificationSnackbarProps> = ({
  notification,
  onClose,
  position = 'bottom'
}) => {
  const [open, setOpen] = useState(false);

  useEffect(() => {
    if (notification) {
      setOpen(true);
    }
  }, [notification]);

  const handleClose = (_event?: React.SyntheticEvent | Event, reason?: string) => {
    if (reason === 'clickaway') {
      return;
    }
    setOpen(false);
    if (notification) {
      onClose(notification.id);
    }
  };

  if (!notification) return null;

  return (
    <Snackbar
      open={open}
      autoHideDuration={notification.duration || 6000}
      onClose={handleClose}
      anchorOrigin={{
        vertical: position,
        horizontal: 'center'
      }}
      sx={{
        '& .MuiAlert-root': {
          minWidth: 300,
          maxWidth: 600,
          boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)',
          borderRadius: 2,
          '& .MuiAlert-icon': {
            fontSize: 20
          },
          '& .MuiAlert-message': {
            fontSize: '0.875rem',
            lineHeight: 1.4
          }
        }
      }}
    >
      <Alert
        onClose={handleClose}
        severity={notification.type}
        icon={getIcon(notification.type)}
        action={notification.action}
        variant="filled"
        elevation={6}
      >
        {notification.message}
      </Alert>
    </Snackbar>
  );
};

export default NotificationSnackbar; 