import { memo } from 'react';
import { Box, Typography, Alert } from '@mui/material';
import { type Issue, type IssueResponse, IssueSource } from '../../generated/api/client';
import { useIssuePageStyles } from './IssuePageStyles';
import StackTraceViewer from '../StackTraceViewer';
import ExceptionDetails from './ExceptionDetails';

interface StackTraceTabProps {
  issueData: IssueResponse | null;
  issue: Issue;
  stackTrace?: string;
}

const StackTraceTab = memo(({ issueData, issue, stackTrace }: StackTraceTabProps) => {
  const styles = useIssuePageStyles();

  // For exception source, use the first event with exception details
  if (issue.source === IssueSource.Exception && issueData?.events && issueData.events.length > 0) {
    const firstEvent = issueData.events[0];
    if (firstEvent.exception_type && firstEvent.exception_value && firstEvent.exception_stacktrace) {
      return (
        <Box sx={styles.tabContent}>
          <ExceptionDetails 
            exception={{
              exception_type: firstEvent.exception_type,
              exception_value: firstEvent.exception_value,
              stacktrace: firstEvent.exception_stacktrace
            }} 
            platform={issue.platform} 
          />
        </Box>
      );
    }
  }

  // For event source with stacktrace
  if (stackTrace) {
    return (
      <Box sx={styles.tabContent}>
        <StackTraceViewer 
          stacktrace={stackTrace} 
          platform={issue.platform} 
        />
      </Box>
    );
  }

  // If no stack trace found yet, try to extract it directly from the first event
  if (issue.source === IssueSource.Event && issueData?.events && issueData.events.length > 0) {
    const firstEvent = issueData.events[0];

    // Create a simple stack trace viewer with the event data
    return (
      <Box sx={styles.tabContent}>
        <Typography variant="h6" sx={{ mb: 2 }}>Event Details</Typography>

        {/* Display event message */}
        <Typography variant="body1" sx={{ mb: 2, fontWeight: 'bold' }}>
          {firstEvent.message}
        </Typography>

        {/* Display any available stack trace information */}
        {firstEvent.exception_stacktrace && (
          <StackTraceViewer 
            stacktrace={firstEvent.exception_stacktrace} 
            platform={issue.platform} 
          />
        )}

        {/* If no specific stack trace found, display event data as JSON */}
        {!firstEvent.exception_stacktrace && (
          <Box sx={{ 
            p: 2, 
            bgcolor: theme => theme.palette.mode === 'dark' ? 'rgba(0,0,0,0.2)' : 'rgba(0,0,0,0.05)',
            borderRadius: 1,
            fontFamily: 'monospace',
            whiteSpace: 'pre-wrap',
            overflowX: 'auto'
          }}>
            {JSON.stringify(firstEvent, null, 2)}
          </Box>
        )}
      </Box>
    );
  }

  // No stack trace available
  return (
    <Box sx={styles.tabContent}>
      <Alert severity="info">No stack trace available for this issue.</Alert>
    </Box>
  );
});

StackTraceTab.displayName = 'StackTraceTab';

export default StackTraceTab;