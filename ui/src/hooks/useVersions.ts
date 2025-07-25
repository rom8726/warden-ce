import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../api/apiClient';

export const useVersions = () => {
  return useQuery({
    queryKey: ['versions'],
    queryFn: async () => {
      const response = await apiClient.getVersions();
      return response.data;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    retry: 3,
  });
}; 