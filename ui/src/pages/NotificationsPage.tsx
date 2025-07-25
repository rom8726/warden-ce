import React, { useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  List,
  ListItem,
  ListItemText,
  ListItemAvatar,
  Avatar,
  Divider,
  Button,
  CircularProgress,
  Alert,
  IconButton,
  Tooltip,
  Chip,
  Pagination,
  FormControlLabel,
  Checkbox
} from '@mui/material';
import {
  Notifications as NotificationsIcon,
  CheckCircle as CheckCircleIcon,
  Group as GroupIcon,
  BugReport as BugReportIcon,
  CheckBox as CheckBoxIcon,
  CheckBoxOutlineBlank as CheckBoxOutlineBlankIcon
} from '@mui/icons-material';
import { useNotifications } from '../hooks/useNotifications';
import { formatDistanceToNow } from 'date-fns';
import type { UserNotification } from '../generated/api/client';

const ITEMS_PER_PAGE = 20;

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

const getNotificationTypeLabel = (type: string): string => {
  switch (type) {
    case 'team_added':
      return 'Team Added';
    case 'team_removed':
      return 'Team Removed';
    case 'role_changed':
      return 'Role Changed';
    case 'issue_regression':
      return 'Issue Regression';
    default:
      return 'Notification';
  }
};

const NotificationsPage: React.FC = () => {
  const [page, setPage] = useState(1);
  const [selectedNotifications, setSelectedNotifications] = useState<Set<number>>(new Set());
  const [selectAll, setSelectAll] = useState(false);
  
  const { 
    notifications, 
    unreadCount, 
    totalCount,
    loading, 
    error, 
    fetchNotifications, 
    markAsRead, 
    markAllAsRead,
    refresh 
  } = useNotifications();

  const offset = (page - 1) * ITEMS_PER_PAGE;

  React.useEffect(() => {
    fetchNotifications(ITEMS_PER_PAGE, offset);
  }, [fetchNotifications, offset]);

  const handlePageChange = (event: React.ChangeEvent<unknown>, value: number) => {
    setPage(value);
    setSelectedNotifications(new Set());
    setSelectAll(false);
  };

  const handleMarkAsRead = async (notificationId: number) => {
    await markAsRead(notificationId);
  };

  const handleMarkSelectedAsRead = async () => {
    const promises = Array.from(selectedNotifications).map(id => markAsRead(id));
    await Promise.all(promises);
    setSelectedNotifications(new Set());
    setSelectAll(false);
  };

  const handleMarkAllAsRead = async () => {
    await markAllAsRead();
    setSelectedNotifications(new Set());
    setSelectAll(false);
  };

  const handleSelectNotification = (notificationId: number, checked: boolean) => {
    const newSelected = new Set(selectedNotifications);
    if (checked) {
      newSelected.add(notificationId);
    } else {
      newSelected.delete(notificationId);
    }
    setSelectedNotifications(newSelected);
    setSelectAll(newSelected.size === notifications.length);
  };

  const handleSelectAll = (checked: boolean) => {
    if (checked) {
      setSelectedNotifications(new Set(notifications.map(n => n.id)));
      setSelectAll(true);
    } else {
      setSelectedNotifications(new Set());
      setSelectAll(false);
    }
  };

  const handleRefresh = async () => {
    await refresh();
    setSelectedNotifications(new Set());
    setSelectAll(false);
  };

  return (
    <Box sx={{ p: 3, maxWidth: 1200, mx: 'auto' }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4" component="h1">
          Notifications
        </Typography>
        <Box sx={{ display: 'flex', gap: 2, alignItems: 'center' }}>
          {unreadCount > 0 && (
            <Chip 
              label={`${unreadCount} unread`} 
              color="primary" 
              variant="outlined"
            />
          )}
          <Button
            variant="outlined"
            onClick={handleRefresh}
            disabled={loading}
          >
            Refresh
          </Button>
        </Box>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {notifications.length > 0 && (
        <Paper sx={{ mb: 3 }}>
          <Box sx={{ p: 2, borderBottom: 1, borderColor: 'divider' }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <FormControlLabel
                control={
                  <Checkbox
                    checked={selectAll}
                    onChange={(e) => handleSelectAll(e.target.checked)}
                    icon={<CheckBoxOutlineBlankIcon />}
                    checkedIcon={<CheckBoxIcon />}
                  />
                }
                label="Select all"
              />
              <Box sx={{ display: 'flex', gap: 1 }}>
                {selectedNotifications.size > 0 && (
                  <Button
                    variant="outlined"
                    size="small"
                    onClick={handleMarkSelectedAsRead}
                  >
                    Mark Selected as Read ({selectedNotifications.size})
                  </Button>
                )}
                {unreadCount > 0 && (
                  <Button
                    variant="contained"
                    size="small"
                    onClick={handleMarkAllAsRead}
                  >
                    Mark All as Read
                  </Button>
                )}
              </Box>
            </Box>
          </Box>

          <List sx={{ p: 0 }}>
            {notifications.map((notification, index) => (
              <React.Fragment key={notification.id}>
                <ListItem 
                  sx={{ 
                    bgcolor: notification.is_read ? 'background.paper' : 'action.selected',
                    '&:hover': {
                      bgcolor: notification.is_read ? 'action.hover' : 'action.selected'
                    }
                  }}
                >
                  <Checkbox
                    checked={selectedNotifications.has(notification.id)}
                    onChange={(e) => handleSelectNotification(notification.id, e.target.checked)}
                    sx={{ mr: 1 }}
                  />
                  <ListItemAvatar>
                    <Avatar 
                      sx={{ 
                        bgcolor: notification.is_read ? 'grey.400' : 'primary.main',
                        width: 48,
                        height: 48
                      }}
                    >
                      {getNotificationIcon(notification.type)}
                    </Avatar>
                  </ListItemAvatar>
                  <ListItemText
                    primary={
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 0.5 }}>
                        <Typography 
                          variant="subtitle1" 
                          color={notification.is_read ? 'text.secondary' : 'text.primary'}
                          sx={{ fontWeight: notification.is_read ? 400 : 600 }}
                        >
                          {getNotificationTitle(notification)}
                        </Typography>
                        <Chip 
                          label={getNotificationTypeLabel(notification.type)} 
                          size="small" 
                          variant="outlined"
                        />
                        {!notification.is_read && (
                          <Chip 
                            label="Unread" 
                            size="small" 
                            color="primary"
                          />
                        )}
                      </Box>
                    }
                    secondary={
                      <Box>
                        <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
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
                        onClick={() => handleMarkAsRead(notification.id)}
                      >
                        <CheckCircleIcon />
                      </IconButton>
                    </Tooltip>
                  )}
                </ListItem>
                {index < notifications.length - 1 && <Divider />}
              </React.Fragment>
            ))}
          </List>
        </Paper>
      )}

      {loading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
          <CircularProgress />
        </Box>
      ) : notifications.length === 0 ? (
        <Paper sx={{ p: 4, textAlign: 'center' }}>
          <NotificationsIcon sx={{ fontSize: 64, color: 'text.disabled', mb: 2 }} />
          <Typography variant="h6" color="text.secondary" gutterBottom>
            No notifications
          </Typography>
          <Typography variant="body2" color="text.secondary">
            You're all caught up! Check back later for new notifications.
          </Typography>
        </Paper>
      ) : (
        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 3 }}>
          <Pagination 
            count={Math.ceil(totalCount / ITEMS_PER_PAGE)}
            page={page}
            onChange={handlePageChange}
            color="primary"
          />
        </Box>
      )}
    </Box>
  );
};

export default NotificationsPage; 