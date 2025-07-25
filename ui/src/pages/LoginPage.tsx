import React, { useState, useRef, useEffect } from 'react';
import { 
  Box, 
  Button, 
  TextField,
  Paper,
  Container,
  CircularProgress,
  Alert,
  Link,
  Typography
} from '@mui/material';
import { useAuth } from '../auth/AuthContext';
import { Navigate, Link as RouterLink } from 'react-router-dom';
import OtpInput from 'react-otp-input';
import { useTheme } from '@mui/material/styles';
import { formatFrontendVersionForLogin } from '../utils/version';
import WardenLogo from "../components/WardenLogo.tsx";

const LoginPage: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [otp, setOtp] = useState('');
  const { login, verify2FA, isAuthenticated, isLoading, error, is2FARequired, is2FABlocked } = useAuth();
  const theme = useTheme();
  const inputRefs = useRef<(HTMLInputElement | null)[]>([]);
  const isSubmittingRef = useRef(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    // Prevent multiple submissions
    if (isSubmittingRef.current) {
      return;
    }
    
    isSubmittingRef.current = true;
    
    try {
      if (is2FARequired) {
        await verify2FA(otp);
      } else {
        await login(username, password);
      }
    } finally {
      isSubmittingRef.current = false;
    }
  };

  const handleBackToLogin = () => {
    setOtp('');
    // Force reset 2FA state through reload (or can be done through context if needed)
    window.location.reload();
  };

  // Auto-focus on first field when 2FA is shown
  useEffect(() => {
    if (is2FARequired && inputRefs.current[0]) {
      inputRefs.current[0]?.focus();
    }
  }, [is2FARequired]);



  // If already authenticated, redirect to dashboard
  if (isAuthenticated) {
    return <Navigate to="/dashboard" />;
  }

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
            <WardenLogo logoSize={58} showLink={false} variant={"h3"} />
          </Box>

          {error && (
            <Alert severity="error" sx={{ width: '100%', mb: 2 }}>
              {error}
            </Alert>
          )}

          <Box component="form" onSubmit={handleSubmit} sx={{ width: '100%' }}>
            {!is2FARequired ? (
              <>
                <TextField
                  margin="normal"
                  required
                  fullWidth
                  id="username"
                  label="Username"
                  name="username"
                  autoComplete="username"
                  autoFocus
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  disabled={isLoading}
                />
                <TextField
                  margin="normal"
                  required
                  fullWidth
                  name="password"
                  label="Password"
                  type="password"
                  id="password"
                  autoComplete="current-password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  disabled={isLoading}
                />
                <Box sx={{ display: 'flex', gap: 1, mt: 3, mb: 2 }}>
                  <Button
                    type="submit"
                    fullWidth
                    variant="contained"
                    sx={{ py: 1.5 }}
                    disabled={isLoading}
                  >
                    {isLoading ? <CircularProgress size={24} /> : 'Sign In'}
                  </Button>
                </Box>
                
                <Box sx={{ textAlign: 'center', mb: 1 }}>
                  <Link
                    component={RouterLink}
                    to="/forgot-password"
                    variant="body2"
                  >
                    Forgot password?
                  </Link>
                </Box>
                
                <Box sx={{ textAlign: 'center', mt: 2 }}>
                  <Typography variant="caption" color="text.secondary">
                    {formatFrontendVersionForLogin()}
                  </Typography>
                </Box>
              </>
            ) : (
              <>
                <Alert severity="info" sx={{ mb: 2 }}>
                  Two-Factor Authentication enabled. Please enter the code from your authenticator app.
                </Alert>
                <Box sx={{ display: 'flex', justifyContent: 'center', mb: 2 }}>
                  <OtpInput
                    value={otp}
                    onChange={setOtp}
                    numInputs={6}
                    renderInput={(inputProps, idx) => {
                      return (
                        <input
                          {...inputProps}
                          key={idx}
                          ref={el => { inputRefs.current[idx] = el; }}
                          onChange={e => {
                            inputProps.onChange?.(e);
                            if (e.target.value && inputRefs.current[idx + 1]) {
                              inputRefs.current[idx + 1]?.focus();
                            }
                          }}
                          onKeyDown={e => {
                            // Handle backspace
                            if (e.key === 'Backspace' && !e.currentTarget.value && idx > 0) {
                              e.preventDefault();
                              // Clear current field and move to previous
                              const newOtp = otp.split('');
                              newOtp[idx] = '';
                              setOtp(newOtp.join(''));
                              inputRefs.current[idx - 1]?.focus();
                            }
                          }}
                          style={{
                            width: '3rem',
                            height: '3rem',
                            fontSize: '2rem',
                            margin: '0 0.5rem',
                            borderRadius: 8,
                            border: `2px solid ${theme.palette.divider}`,
                            background: theme.palette.background.paper,
                            color: theme.palette.text.primary,
                            textAlign: 'center',
                            outline: document.activeElement === inputRefs.current[idx] ? `2px solid ${theme.palette.primary.main}` : 'none',
                            boxShadow: document.activeElement === inputRefs.current[idx] ? `0 0 0 2px ${theme.palette.primary.light}` : 'none',
                            transition: 'border 0.2s, box-shadow 0.2s',
                          }}
                          inputMode="numeric"
                          pattern="[0-9]*"
                          readOnly={isLoading || is2FABlocked}
                        />
                      );
                    }}
                    containerStyle={{ justifyContent: 'center' }}
                  />
                </Box>
                <Button
                  type="submit"
                  fullWidth
                  variant="contained"
                  sx={{ mt: 3, mb: 2, py: 1.5 }}
                  disabled={isLoading || !otp || is2FABlocked}
                >
                  {isLoading ? <CircularProgress size={24} /> : 'Verify 2FA'}
                </Button>
                <Button
                  fullWidth
                  variant="text"
                  color="secondary"
                  onClick={handleBackToLogin}
                  disabled={isLoading}
                >
                  Back to login
                </Button>
              </>
            )}
          </Box>
        </Paper>
      </Box>
    </Container>
  );
};

export default LoginPage;
