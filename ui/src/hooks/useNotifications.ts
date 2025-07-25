import { useState, useEffect, useCallback } from 'react';
import apiClient from '../api/apiClient';
import type { UserNotification, UnreadCountResponse } from '../generated/api/client/api';

interface UseNotificationsReturn {
  notifications: UserNotification[];
  unreadCount: number;
  totalCount: number;
  loading: boolean;
  error: string | null;
  fetchNotifications: (limit?: number, offset?: number) => Promise<{ total: number }>;
  fetchUnreadCount: () => Promise<void>;
  markAsRead: (notificationId: number) => Promise<void>;
  markAllAsRead: () => Promise<void>;
  refresh: () => Promise<void>;
}

export const useNotifications = (): UseNotificationsReturn => {
  const [notifications, setNotifications] = useState<UserNotification[]>([]);
  const [unreadCount, setUnreadCount] = useState<number>(0);
  const [totalCount, setTotalCount] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const fetchNotifications = useCallback(async (limit = 5, offset = 0) => {
    try {
      setLoading(true);
      setError(null);
      const response = await apiClient.getUserNotifications(limit, offset);
      setNotifications(response.data.notifications);
      setTotalCount(response.data.total);
      return { total: response.data.total };
    } catch (err) {
      console.error('Error fetching notifications:', err);
      setError('Failed to load notifications');
      return { total: 0 };
    } finally {
      setLoading(false);
    }
  }, []);

  const fetchUnreadCount = useCallback(async () => {
    try {
      const response = await apiClient.getUnreadNotificationsCount();
      setUnreadCount(response.data.count);
    } catch (err) {
      console.error('Error fetching unread count:', err);
      setUnreadCount(0);
    }
  }, []);

  const markAsRead = useCallback(async (notificationId: number) => {
    try {
      await apiClient.markNotificationAsRead(notificationId);
      // Update local state
      setNotifications(prev => 
        prev.map(notification => 
          notification.id === notificationId 
            ? { ...notification, is_read: true }
            : notification
        )
      );
      // Refresh unread count
      await fetchUnreadCount();
    } catch (err) {
      console.error('Error marking notification as read:', err);
      setError('Failed to mark notification as read');
    }
  }, [fetchUnreadCount]);

  const markAllAsRead = useCallback(async () => {
    try {
      await apiClient.markAllNotificationsAsRead();
      // Update local state
      setNotifications(prev => 
        prev.map(notification => ({ ...notification, is_read: true }))
      );
      setUnreadCount(0);
    } catch (err) {
      console.error('Error marking all notifications as read:', err);
      setError('Failed to mark all notifications as read');
    }
  }, []);

  const refresh = useCallback(async () => {
    await Promise.all([
      fetchNotifications(),
      fetchUnreadCount()
    ]);
  }, [fetchNotifications, fetchUnreadCount]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  return {
    notifications,
    unreadCount,
    totalCount,
    loading,
    error,
    fetchNotifications,
    fetchUnreadCount,
    markAsRead,
    markAllAsRead,
    refresh
  };
}; 