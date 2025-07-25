import { memo } from 'react';
import { Box, Typography, Paper, Grid } from '@mui/material';
import { Schedule as ScheduleIcon } from '@mui/icons-material';
import { type Issue, type IssueResponse, IssueStatus } from '../../generated/api/client';
import { useIssuePageStyles } from './IssuePageStyles';
import { formatDate } from '../../utils/issues/issueUtils';
import BaseInfoPanel from './BaseInfoPanel';

interface IssueStatsProps {
  issue: Issue;
  issueData: IssueResponse | null;
}

const IssueStats = memo(({ issue, issueData }: IssueStatsProps) => {
  const styles = useIssuePageStyles();

  return (
    <Grid container spacing={3} sx={styles.statsContainer}>
      <Grid item xs={12} md={6}>
        <Paper sx={styles.statsPaper}>
          <Box sx={styles.statsIconContainer}>
            <ScheduleIcon fontSize="small" sx={styles.statsIcon} />
            <Typography variant="h6" component="h2" className="gradient-subtitle-blue">
              Occurrences
            </Typography>
          </Box>
          <Typography variant="h3" component="div" sx={styles.statsValue}>
            {issue.count}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            First seen: {formatDate(issue.first_seen)}
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Last seen: {formatDate(issue.last_seen)}
          </Typography>
          {issue.status === IssueStatus.Resolved && issue.resolved_at && (
            <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
              Resolved at: {formatDate(issue.resolved_at)}
              {issue.resolved_by && ` by ${issue.resolved_by}`}
            </Typography>
          )}
        </Paper>
      </Grid>
      <Grid item xs={12} md={6}>
        <BaseInfoPanel issueData={issueData} />
      </Grid>
    </Grid>
  );
});

IssueStats.displayName = 'IssueStats';

export default IssueStats;