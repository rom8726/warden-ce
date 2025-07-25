import React from 'react';
import { Box, Typography, Tooltip } from '@mui/material';
import { versionInfo } from '../version';

interface VersionInfoProps {
  showBuildTime?: boolean;
  variant?: 'body1' | 'body2' | 'caption' | 'overline' | 'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6' | 'subtitle1' | 'subtitle2';
  color?: string;
}

export const VersionInfo: React.FC<VersionInfoProps> = ({ 
  showBuildTime = false, 
  variant = 'caption',
  color = 'text.secondary'
}) => {
  const formatBuildTime = (buildTime: string) => {
    try {
      return new Date(buildTime).toLocaleString();
    } catch {
      return buildTime;
    }
  };

  return (
    <Box>
      <Typography variant={variant} color={color}>
        {versionInfo.version}
      </Typography>
      {showBuildTime && (
        <Tooltip title="Build time">
          <Typography variant="caption" color="text.disabled" sx={{ fontSize: '0.7rem' }}>
            {formatBuildTime(versionInfo.buildTime)}
          </Typography>
        </Tooltip>
      )}
    </Box>
  );
};

export default VersionInfo; 