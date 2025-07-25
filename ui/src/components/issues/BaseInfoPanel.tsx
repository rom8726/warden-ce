import { memo } from 'react';
import { Box, Typography, Paper, Grid } from '@mui/material';
import { type IssueResponse } from '../../generated/api/client';

interface BaseInfoPanelProps {
  issueData: IssueResponse | null;
}

const BaseInfoPanel = memo(({ issueData }: BaseInfoPanelProps) => {
  if (!issueData?.events || issueData.events.length === 0) {
    return null;
  }

  const firstEvent = issueData.events[0];

  // Extract only the specified fields from the first event
  const infoItems = [
    { label: 'Platform', value: firstEvent.platform },
    { label: 'Server Name', value: firstEvent.server_name },
    { label: 'Environment', value: firstEvent.environment },
  ];

  // Add exception-specific information if available
  if (firstEvent.exception_type) {
    infoItems.push({ label: 'Exception Type', value: firstEvent.exception_type });
  }
  if (firstEvent.exception_value) {
    infoItems.push({ label: 'Exception Value', value: firstEvent.exception_value });
  }

  // Calculate the number of items per column
  // For 2 columns, we divide the total number of items by 2 and round up
  const itemsPerColumn = Math.ceil(infoItems.length / 2);

  return (
    <Paper sx={{ p: 1.5, height: '100%' }}>
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
        <Typography variant="h6" component="h2">
          Base Info
        </Typography>
      </Box>
      <Grid container spacing={2}>
        {/* First column */}
        <Grid item xs={12} sm={6}>
          {infoItems.slice(0, itemsPerColumn).map((item, index) => (
            <Box key={index} sx={{ mb: 1 }}>
              <Typography variant="body2" color="text.secondary">
                {item.label}
              </Typography>
              <Typography variant="body1" sx={{ fontWeight: 'medium', wordBreak: 'break-word' }}>
                {item.value || 'N/A'}
              </Typography>
            </Box>
          ))}
        </Grid>
        {/* Second column */}
        <Grid item xs={12} sm={6}>
          {infoItems.slice(itemsPerColumn).map((item, index) => (
            <Box key={index} sx={{ mb: 1 }}>
              <Typography variant="body2" color="text.secondary">
                {item.label}
              </Typography>
              <Typography variant="body1" sx={{ fontWeight: 'medium', wordBreak: 'break-word' }}>
                {item.value || 'N/A'}
              </Typography>
            </Box>
          ))}
        </Grid>
      </Grid>
    </Paper>
  );
});

BaseInfoPanel.displayName = 'BaseInfoPanel';

export default BaseInfoPanel;