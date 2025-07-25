import React, { useState } from 'react';
import { 
  Box, 
  Button, 
  TextField,
  Paper,
  Container,
  CircularProgress,
  Alert,
  Typography,
  Link
} from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import LogoImg from "../components/LogoImg.tsx";
import { DefaultApi, Configuration } from '../generated/api/client';

const ForgotPasswordPage: React.FC = () => {
  const [email, setEmail] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);

    try {
      const api = new DefaultApi(new Configuration({
        basePath: (import.meta.env.VITE_API_BASE_URL || '').replace(/\/+$/, ''),
      }));

      await api.forgotPassword({ email });
      setSuccess(true);
    } catch (err) {
      console.error('Error requesting password reset:', err);
      setError('Failed to send password reset email. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Container maxWidth="sm">
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          minHeight: '100vh',
        }}
      >
        <Paper
          elevation={3}
          sx={{
            p: 4,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            width: '100%',
          }}
        >
          <Box sx={{ mb: 4 }}>
            <LogoImg size="large" />
          </Box>

          <Typography component="h1" variant="h5" sx={{ mb: 3 }}>
            Forgot Password
          </Typography>

          {error && (
            <Alert severity="error" sx={{ width: '100%', mb: 2 }}>
              {error}
            </Alert>
          )}

          {success ? (
            <Box sx={{ width: '100%', textAlign: 'center' }}>
              <Alert severity="success" sx={{ width: '100%', mb: 2 }}>
                Password reset email sent. Please check your inbox.
              </Alert>
              <Button
                component={RouterLink}
                to="/login"
                fullWidth
                variant="contained"
                sx={{ mt: 2 }}
              >
                Back to Login
              </Button>
            </Box>
          ) : (
            <Box component="form" onSubmit={handleSubmit} sx={{ width: '100%' }}>
              <Typography variant="body2" sx={{ mb: 2 }}>
                Enter your email address and we'll send you a link to reset your password.
              </Typography>
              <TextField
                margin="normal"
                required
                fullWidth
                id="email"
                label="Email Address"
                name="email"
                autoComplete="email"
                autoFocus
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                disabled={isLoading}
              />
              <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2, py: 1.5 }}
                disabled={isLoading}
              >
                {isLoading ? <CircularProgress size={24} /> : 'Send Reset Link'}
              </Button>
              <Box sx={{ textAlign: 'center' }}>
                <Link component={RouterLink} to="/login" variant="body2">
                  Back to Login
                </Link>
              </Box>
            </Box>
          )}
        </Paper>
      </Box>
    </Container>
  );
};

export default ForgotPasswordPage;
