import React from 'react';
import { usePermissions } from '../hooks/usePermissions';

interface PermissionGuardProps {
  children: React.ReactNode;
  fallback?: React.ReactNode;
  // Права на проекты
  projectId?: number;
  canRead?: boolean;
  canWrite?: boolean;
  canDelete?: boolean;
  canManage?: boolean;
  // Права в командах
  teamId?: number;
  requireOwner?: boolean;
  requireAdmin?: boolean;
  requireMember?: boolean;
  // Общие права
  requireCreateProjects?: boolean;
  requireCreateTeams?: boolean;
  requireManageUsers?: boolean;
  // Superuser
  requireSuperuser?: boolean;
}

const PermissionGuard: React.FC<PermissionGuardProps> = ({
  children,
  fallback = null,
  projectId,
  canRead,
  canWrite,
  canDelete,
  canManage,
  teamId,
  requireOwner,
  requireAdmin,
  requireMember,
  requireCreateProjects,
  requireCreateTeams,
  requireManageUsers,
  requireSuperuser,
}) => {
  const {
    canReadProject,
    canWriteProject,
    canDeleteProject,
    canManageProject,
    isTeamOwner,
    isTeamAdmin,
    isTeamMember,
    canCreateProjects,
    canCreateTeams,
    canManageUsers,
    isSuperuser,
  } = usePermissions();

  // Проверка superuser
  if (requireSuperuser && !isSuperuser()) {
    return <>{fallback}</>;
  }

  // Проверка общих прав
  if (requireCreateProjects && !canCreateProjects()) {
    return <>{fallback}</>;
  }

  if (requireCreateTeams && !canCreateTeams()) {
    return <>{fallback}</>;
  }

  if (requireManageUsers && !canManageUsers()) {
    return <>{fallback}</>;
  }

  // Проверка прав на проекты
  if (projectId) {
    if (canRead && !canReadProject(projectId)) {
      return <>{fallback}</>;
    }

    if (canWrite && !canWriteProject(projectId)) {
      return <>{fallback}</>;
    }

    if (canDelete && !canDeleteProject(projectId)) {
      return <>{fallback}</>;
    }

    if (canManage && !canManageProject(projectId)) {
      return <>{fallback}</>;
    }
  }

  // Проверка прав в командах
  if (teamId) {
    if (requireOwner && !isTeamOwner(teamId)) {
      return <>{fallback}</>;
    }

    if (requireAdmin && !isTeamAdmin(teamId)) {
      return <>{fallback}</>;
    }

    if (requireMember && !isTeamMember(teamId)) {
      return <>{fallback}</>;
    }
  }

  return <>{children}</>;
};

export default PermissionGuard; 