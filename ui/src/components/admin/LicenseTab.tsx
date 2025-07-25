import React from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Chip,
  Grid,
  Divider,
} from '@mui/material';

// Simplified LicenseTab for community edition
const LicenseTab: React.FC = () => {
  return (
    <Box>
      <Box sx={{ mb: 3 }}>
        <Typography variant="h5" component="h2" fontWeight={600} className="gradient-text-purple">
          License Information
        </Typography>
      </Box>

      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Box sx={{ mb: 3 }}>
            <Typography variant="h6" gutterBottom fontWeight={600} className="gradient-subtitle-purple">
              Community Edition License
            </Typography>
            <Grid container spacing={2} alignItems="center">
              <Grid item>
                <Chip
                  label="Valid"
                  color="success"
                  variant="filled"
                  size="medium"
                />
              </Grid>
            </Grid>
          </Box>

          <Divider sx={{ my: 2 }} />

          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                Type
              </Typography>
              <Typography variant="body1">
                Community
              </Typography>
            </Grid>

            <Grid item xs={12} md={6}>
              <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                Expires
              </Typography>
              <Typography variant="body1">
                Never
              </Typography>
            </Grid>

            <Grid item xs={12}>
              <Typography variant="body2" color="text.secondary" sx={{ mt: 2 }}>
                This is the Community Edition of the application, which is free to use and does not require a license key.
                The Community Edition includes all core features of the application.
              </Typography>
            </Grid>
          </Grid>
        </CardContent>
      </Card>
    </Box>
  );
};

export default LicenseTab;