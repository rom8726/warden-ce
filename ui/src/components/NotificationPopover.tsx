import React from 'react';
import {
  Popover,
  Typography,
  List,
  ListItem,
  ListItemText,
  ListItemAvatar,
  Avatar,
  Divider,
  Button,
  Box,
  CircularProgress,
  Alert,
  IconButton,
  Tooltip
} from '@mui/material';
import {
  Notifications as NotificationsIcon,
  Visibility as VisibilityIcon,
  CheckCircle as CheckCircleIcon,
  Group as GroupIcon,
  BugReport as BugReportIcon,
  ArrowForward as ArrowForwardIcon
} from '@mui/icons-material';
import { useNotifications } from '../hooks/useNotifications';
import { formatDistanceToNow } from 'date-fns';
import type { UserNotification } from '../generated/api/client/api';

interface NotificationPopoverProps {
  anchorEl: HTMLElement | null;
  onClose: () => void;
  onViewAll: () => void;
}

const getNotificationIcon = (type: string) => {
  switch (type) {
    case 'team_added':
    case 'team_removed':
    case 'role_changed':
      return <GroupIcon />;
    case 'issue_regression':
      return <BugReportIcon />;
    default:
      return <NotificationsIcon />;
  }
};

const getNotificationTitle = (notification: UserNotification): string => {
  switch (notification.type) {
    case 'team_added':
      return 'You have been added to team';
    case 'team_removed':
      return 'You have been removed from team';
    case 'role_changed':
      return 'Your role has been changed';
    case 'issue_regression':
      return 'Issue regression detected';
    default:
      return 'Notification';
  }
};

const getNotificationMessage = (notification: UserNotification): string => {
  const content = notification.content;
  
  switch (notification.type) {
    case 'team_added':
      return `You have been added to team '${content.team_name}' as ${content.role}`;
    case 'team_removed':
      return `You have been removed from team '${content.team_name}'`;
    case 'role_changed':
      return `Your role in team '${content.team_name}' has been changed from ${content.old_role} to ${content.new_role}`;
    case 'issue_regression':
      return `Issue '${content.issue_title}' has regressed`;
    default:
      return 'You have a new notification';
  }
};

const NotificationPopover: React.FC<NotificationPopoverProps> = ({
  anchorEl,
  onClose,
  onViewAll
}) => {
  const { notifications, unreadCount, loading, error, markAsRead } = useNotifications();

  const handleMarkAsRead = async (notificationId: number) => {
    await markAsRead(notificationId);
  };

  const handleViewAll = () => {
    onViewAll();
    onClose();
  };

  return (
    <Popover
      open={Boolean(anchorEl)}
      anchorEl={anchorEl}
      onClose={onClose}
      anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      transformOrigin={{ vertical: 'top', horizontal: 'right' }}
      PaperProps={{ 
        sx: { 
          minWidth: 400, 
          maxWidth: 500, 
          maxHeight: 600,
          overflow: 'hidden'
        } 
      }}
    >
      <Box sx={{ p: 2, borderBottom: 1, borderColor: 'divider' }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Typography variant="h6">
            Notifications
          </Typography>
          {unreadCount > 0 && (
            <Typography variant="caption" color="primary">
              {unreadCount} unread
            </Typography>
          )}
        </Box>
      </Box>

      {error && (
        <Alert severity="error" sx={{ m: 2 }}>
          {error}
        </Alert>
      )}

      {loading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
          <CircularProgress size={24} />
        </Box>
      ) : (
        <>
          <List sx={{ p: 0, maxHeight: 400, overflow: 'auto' }}>
            {notifications.length === 0 ? (
              <ListItem>
                <ListItemText 
                  primary="No notifications" 
                  secondary="You're all caught up!"
                />
              </ListItem>
            ) : (
              notifications.map((notification, index) => (
                <React.Fragment key={notification.id}>
                  <ListItem 
                    sx={{ 
                      bgcolor: notification.is_read ? 'background.paper' : 'action.selected',
                      '&:hover': {
                        bgcolor: notification.is_read ? 'action.hover' : 'action.selected'
                      }
                    }}
                  >
                    <ListItemAvatar>
                      <Avatar 
                        sx={{ 
                          bgcolor: notification.is_read ? 'grey.400' : 'primary.main',
                          width: 40,
                          height: 40
                        }}
                      >
                        {getNotificationIcon(notification.type)}
                      </Avatar>
                    </ListItemAvatar>
                    <ListItemText
                      primary={
                        <Typography 
                          variant="subtitle2" 
                          color={notification.is_read ? 'text.secondary' : 'text.primary'}
                          sx={{ fontWeight: notification.is_read ? 400 : 600 }}
                        >
                          {getNotificationTitle(notification)}
                        </Typography>
                      }
                      secondary={
                        <Box>
                          <Typography variant="body2" color="text.secondary" sx={{ mb: 0.5 }}>
                            {getNotificationMessage(notification)}
                          </Typography>
                          <Typography variant="caption" color="text.disabled">
                            {formatDistanceToNow(new Date(notification.created_at), { addSuffix: true })}
                          </Typography>
                        </Box>
                      }
                    />
                    {!notification.is_read && (
                      <Tooltip title="Mark as read">
                        <IconButton
                          size="small"
                          onClick={() => handleMarkAsRead(notification.id)}
                          sx={{ ml: 1 }}
                        >
                          <CheckCircleIcon fontSize="small" />
                        </IconButton>
                      </Tooltip>
                    )}
                  </ListItem>
                  {index < notifications.length - 1 && <Divider />}
                </React.Fragment>
              ))
            )}
          </List>

          {notifications.length > 0 && (
            <Box sx={{ p: 2, borderTop: 1, borderColor: 'divider' }}>
              <Button
                fullWidth
                variant="outlined"
                endIcon={<ArrowForwardIcon />}
                onClick={handleViewAll}
              >
                View All Notifications
              </Button>
            </Box>
          )}
        </>
      )}
    </Popover>
  );
};

export default NotificationPopover; 