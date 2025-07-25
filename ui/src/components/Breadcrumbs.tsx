import React from 'react';
import { 
  Breadcrumbs as MuiBreadcrumbs, 
  Link, 
  Typography, 
  Box,
  useTheme
} from '@mui/material';
import { 
  NavigateNext as NavigateNextIcon,
  Home as HomeIcon
} from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import apiClient from '../api/apiClient';

interface BreadcrumbItem {
  label: string;
  path: string;
  icon?: React.ReactNode;
}

const Breadcrumbs: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const theme = useTheme();

  // Get project data for breadcrumbs
  const { data: projects } = useQuery({
    queryKey: ['projects'],
    queryFn: async () => {
      const response = await apiClient.listProjects();
      return response.data;
    }
  });

  // Get current project data
  const projectId = location.pathname.match(/\/projects\/(\d+)/)?.[1];
  const currentProject = projects?.find(p => p.id === Number(projectId));

  // Get issue data for breadcrumbs
  const issueId = location.pathname.match(/\/issues\/(\d+)/)?.[1];
  const { data: issueData } = useQuery({
    queryKey: ['issue', projectId, issueId],
    queryFn: async () => {
      if (!projectId || !issueId) return null;
      const response = await apiClient.getIssue(Number(projectId), Number(issueId));
      return response.data;
    },
    enabled: !!projectId && !!issueId
  });

  const generateBreadcrumbs = (): BreadcrumbItem[] => {
    const pathSegments = location.pathname.split('/').filter(Boolean);
    const breadcrumbs: BreadcrumbItem[] = [];

    // Always add home
    breadcrumbs.push({
      label: 'Dashboard',
      path: '/dashboard',
      icon: <HomeIcon fontSize="small" />
    });

    let currentPath = '';
    
    for (let i = 0; i < pathSegments.length; i++) {
      const segment = pathSegments[i];
      currentPath += `/${segment}`;

      switch (segment) {
        case 'projects':
          breadcrumbs.push({
            label: 'Projects',
            path: '/projects'
          });
          break;
        case 'issues':
          breadcrumbs.push({
            label: 'Issues',
            path: '/issues'
          });
          break;
        case 'settings':
          breadcrumbs.push({
            label: 'Settings',
            path: currentPath
          });
          break;
        case 'admin':
          breadcrumbs.push({
            label: 'Admin',
            path: '/admin'
          });
          break;
        default:
          // Handle dynamic segments (project IDs, issue IDs)
          if (i === 1 && pathSegments[0] === 'projects' && currentProject) {
            // Project page
            breadcrumbs.push({
              label: currentProject.name,
              path: currentPath
            });
          } else if (i === 3 && pathSegments[0] === 'projects' && pathSegments[2] === 'issues' && issueData?.issue) {
            // Issue page
            breadcrumbs.push({
              label: issueData.issue.title || `Issue #${segment}`,
              path: currentPath
            });
          } else if (i === 1 && pathSegments[0] === 'issues' && issueData?.issue) {
            // Issues page with specific issue
            breadcrumbs.push({
              label: issueData.issue.title || `Issue #${segment}`,
              path: currentPath
            });
          }
          break;
      }
    }

    return breadcrumbs;
  };

  const breadcrumbs = generateBreadcrumbs();

  // Don't show breadcrumbs on login pages or if only one item
  if (breadcrumbs.length <= 1 || location.pathname.startsWith('/login') || location.pathname.startsWith('/forgot-password') || location.pathname.startsWith('/reset-password')) {
    return null;
  }

  return (
    <Box sx={{ mb: 3, width: '100%' }}>
      <MuiBreadcrumbs 
        separator={<NavigateNextIcon fontSize="small" />}
        aria-label="breadcrumb"
        sx={{
          '& .MuiBreadcrumbs-separator': {
            color: theme.palette.text.secondary,
          },
        }}
      >
        {breadcrumbs.map((item, index) => {
          const isLast = index === breadcrumbs.length - 1;
          
          if (isLast) {
            return (
              <Typography 
                key={item.path}
                color="text.primary" 
                sx={{ 
                  display: 'flex', 
                  alignItems: 'center', 
                  gap: 0.5,
                  fontWeight: 500,
                  fontSize: '0.875rem'
                }}
              >
                {item.icon}
                {item.label}
              </Typography>
            );
          }

          return (
            <Link
              key={item.path}
              color="inherit"
              href="#"
              onClick={(e) => {
                e.preventDefault();
                navigate(item.path);
              }}
              sx={{ 
                display: 'flex', 
                alignItems: 'center', 
                gap: 0.5,
                textDecoration: 'none',
                color: theme.palette.text.secondary,
                fontSize: '0.875rem',
                '&:hover': {
                  color: theme.palette.primary.main,
                  textDecoration: 'underline',
                },
              }}
            >
              {item.icon}
              {item.label}
            </Link>
          );
        })}
      </MuiBreadcrumbs>
    </Box>
  );
};

export default Breadcrumbs; 