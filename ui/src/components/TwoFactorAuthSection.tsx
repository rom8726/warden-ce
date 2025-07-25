import React, { useState, useRef, useEffect } from 'react';
import { Paper, Typography, Box, CircularProgress, Chip, Button, Dialog, DialogTitle, DialogContent, DialogActions, TextField, IconButton, InputAdornment, Alert } from '@mui/material';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import apiClient from '../api/apiClient';
import { UsersApi } from '../generated/api/client/api';
import { apiConfiguration, axiosInstance } from '../api/apiClient';
import OtpInput from 'react-otp-input';
import { useTheme } from '@mui/material/styles';

interface TwoFactorAuthSectionProps {
  userData: any;
  userLoading: boolean;
  userError: any;
}

const TwoFactorAuthSection: React.FC<TwoFactorAuthSectionProps> = ({ userData, userLoading, userError }) => {
  const twoFAEnabled = Boolean(userData?.two_fa_enabled);
  const [openSetup, setOpenSetup] = useState(false);
  const [setupData, setSetupData] = useState<any>(null);
  const [code, setCode] = useState('');
  const [copySuccess, setCopySuccess] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [openDisable, setOpenDisable] = useState(false);
  const [emailCode, setEmailCode] = useState('');
  const [disableLoading, setDisableLoading] = useState(false);
  const [disableError, setDisableError] = useState<string | null>(null);
  const [disableSuccess, setDisableSuccess] = useState(false);
  const [sendCodeLoading, setSendCodeLoading] = useState(false);
  const [sendCodeError, setSendCodeError] = useState<string | null>(null);
  const [sendCodeSuccess, setSendCodeSuccess] = useState(false);
  const [is2FABlocked, setIs2FABlocked] = useState(false);
  const theme = useTheme();
  const codeInputRefs = useRef<(HTMLInputElement | null)[]>([]);
  const emailInputRefs = useRef<(HTMLInputElement | null)[]>([]);
  const isSubmittingRef = useRef(false);
  const isDisableSubmittingRef = useRef(false);

  const handleEnable2FA = async () => {
    setOpenSetup(true);
    setSetupData(null);
    setError(null);
    setSuccess(false);
    setLoading(true);
    try {
      const resp = await apiClient.setup2FA();
      setSetupData(resp.data);
    } catch (e: any) {
      setError(e?.response?.data?.error?.message || 'Failed to start 2FA setup');
    } finally {
      setLoading(false);
    }
  };

  const handleCloseSetup = () => {
    setOpenSetup(false);
    setSetupData(null);
    setCode('');
    setCopySuccess(false);
    setError(null);
    setSuccess(false);
  };

  const handleCopySecret = () => {
    if (setupData?.secret) {
      navigator.clipboard.writeText(setupData.secret);
      setCopySuccess(true);
      setTimeout(() => setCopySuccess(false), 1500);
    }
  };

  const handleConfirm2FA = async () => {
    // Prevent multiple submissions
    if (isSubmittingRef.current) {
      return;
    }
    
    isSubmittingRef.current = true;
    setLoading(true);
    setError(null);
    
    try {
      await apiClient.confirm2FA({ code });
      setSuccess(true);
      setTimeout(() => {
        handleCloseSetup();
        window.location.reload(); // refresh userData
      }, 1200);
    } catch (e: any) {
      const msg = e?.response?.data?.error?.message || 'Invalid code';
      setError(msg);
      if (msg.includes('Слишком много попыток')) {
        setIs2FABlocked(true);
        setTimeout(() => setIs2FABlocked(false), 60000);
      }
    } finally {
      setLoading(false);
      isSubmittingRef.current = false;
    }
  };

  const handleDisable2FA = () => {
    setOpenDisable(true);
    setEmailCode('');
    setDisableError(null);
    setDisableSuccess(false);
  };

  const handleCloseDisable = () => {
    setOpenDisable(false);
    setEmailCode('');
    setDisableError(null);
    setDisableSuccess(false);
  };

  const handleConfirmDisable2FA = async () => {
    // Prevent multiple submissions
    if (isDisableSubmittingRef.current) {
      return;
    }
    
    isDisableSubmittingRef.current = true;
    setDisableLoading(true);
    setDisableError(null);
    
    try {
      await apiClient.disable2FA({ email_code: emailCode });
      setDisableSuccess(true);
      setTimeout(() => {
        handleCloseDisable();
        window.location.reload();
      }, 1200);
    } catch (e: any) {
      setDisableError(e?.response?.data?.error?.message || 'Invalid code');
    } finally {
      setDisableLoading(false);
      isDisableSubmittingRef.current = false;
    }
  };

  const handleSendCode = async () => {
    setSendCodeLoading(true);
    setSendCodeError(null);
    setSendCodeSuccess(false);
    try {
      const usersApi = new UsersApi(apiConfiguration, apiConfiguration.basePath, axiosInstance);
      await usersApi.send2FACode();
      setSendCodeSuccess(true);
    } catch (e: any) {
      setSendCodeError(e?.response?.data?.error?.message || 'Failed to send code');
    } finally {
      setSendCodeLoading(false);
    }
  };

  // Auto-focus on first field when 2FA setup dialog opens
  useEffect(() => {
    if (openSetup && setupData && codeInputRefs.current[0]) {
      setTimeout(() => {
        codeInputRefs.current[0]?.focus();
      }, 100);
    }
  }, [openSetup, setupData]);

  // Auto-focus on first field when 2FA disable dialog opens
  useEffect(() => {
    if (openDisable && emailInputRefs.current[0]) {
      setTimeout(() => {
        emailInputRefs.current[0]?.focus();
      }, 100);
    }
  }, [openDisable]);



  return (
    <Paper sx={{ p: 3 }}>
      <Typography variant="h6" gutterBottom className="gradient-subtitle">
        Two-Factor Authentication
      </Typography>
      {userLoading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
          <CircularProgress />
        </Box>
      ) : userError ? (
        <Typography color="error">
          Error loading user data. Please try again.
        </Typography>
      ) : (
        <Box>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
            <Typography variant="subtitle2" color="text.secondary" sx={{ mr: 1 }}>
              Status:
            </Typography>
            <Chip
              label={twoFAEnabled ? 'Enabled' : 'Disabled'}
              color={twoFAEnabled ? 'success' : 'default'}
              size="small"
              sx={{ ml: 1 }}
            />
          </Box>
          <Button 
            variant="contained" 
            color={twoFAEnabled ? 'error' : 'primary'}
            disabled={userLoading}
            onClick={twoFAEnabled ? handleDisable2FA : handleEnable2FA}
          >
            {twoFAEnabled ? 'Disable 2FA' : 'Enable 2FA'}
          </Button>

          {/* Модалка для включения 2FA */}
          <Dialog open={openSetup} onClose={handleCloseSetup} maxWidth="xs" fullWidth>
            <DialogTitle className="gradient-text">Enable Two-Factor Authentication</DialogTitle>
            <DialogContent>
              {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}
              {success && <Alert severity="success" sx={{ mb: 2 }}>2FA enabled!</Alert>}
              {loading && !setupData ? (
                <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
                  <CircularProgress />
                </Box>
              ) : setupData ? (
                <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 2 }}>
                  {/* QR-код */}
                  {setupData.qr_image ? (
                    <img
                      src={`data:image/png;base64,${setupData.qr_image}`}
                      alt="QR code"
                      style={{ width: 180, height: 180, marginBottom: 8, borderRadius: 8 }}
                    />
                  ) : (
                    <Box sx={{ width: 180, height: 180, bgcolor: '#eee', borderRadius: 2, mb: 1, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                      <Typography variant="caption" color="text.secondary">QR</Typography>
                    </Box>
                  )}
                  <Typography variant="body2" align="center">
                    Scan this QR code in Google Authenticator, 1Password или другом приложении.<br/>
                    Или введите секрет вручную:
                  </Typography>
                  <TextField
                    label="Secret"
                    value={setupData.secret}
                    InputProps={{
                      readOnly: true,
                      endAdornment: (
                        <InputAdornment position="end">
                          <IconButton onClick={handleCopySecret} edge="end" size="small">
                            <ContentCopyIcon fontSize="small" />
                          </IconButton>
                        </InputAdornment>
                      )
                    }}
                    fullWidth
                    margin="dense"
                  />
                  {copySuccess && <Typography color="success.main" variant="caption">Copied!</Typography>}
                  <Box sx={{ display: 'flex', justifyContent: 'center', mb: 2 }}>
                    <OtpInput
                      value={code}
                      onChange={setCode}
                      numInputs={6}
                      renderInput={(inputProps, idx) => (
                        <input
                          {...inputProps}
                          key={idx}
                          ref={el => { codeInputRefs.current[idx] = el; }}
                          onChange={e => {
                            inputProps.onChange?.(e);
                            if (e.target.value && codeInputRefs.current[idx + 1]) {
                              codeInputRefs.current[idx + 1]?.focus();
                            }
                          }}
                          onKeyDown={e => {
                            // Handle backspace
                            if (e.key === 'Backspace' && !e.currentTarget.value && idx > 0) {
                              e.preventDefault();
                              // Clear current field and move to previous
                              const newCode = code.split('');
                              newCode[idx] = '';
                              setCode(newCode.join(''));
                              codeInputRefs.current[idx - 1]?.focus();
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
                            outline: document.activeElement === codeInputRefs.current[idx] ? `2px solid ${theme.palette.primary.main}` : 'none',
                            boxShadow: document.activeElement === codeInputRefs.current[idx] ? `0 0 0 2px ${theme.palette.primary.light}` : 'none',
                            transition: 'border 0.2s, box-shadow 0.2s',
                          }}
                          inputMode="numeric"
                          pattern="[0-9]*"
                          readOnly={loading || success || is2FABlocked}
                        />
                      )}
                      containerStyle={{ justifyContent: 'center' }}
                    />
                  </Box>
                </Box>
              ) : null}
            </DialogContent>
            <DialogActions>
              <Button onClick={handleCloseSetup} disabled={loading || is2FABlocked}>Cancel</Button>
              <Button variant="contained" disabled={!code || loading || success || is2FABlocked} onClick={handleConfirm2FA}>Confirm</Button>
            </DialogActions>
          </Dialog>

          {/* Модалка для отключения 2FA */}
          <Dialog open={openDisable} onClose={handleCloseDisable} maxWidth="xs" fullWidth>
            <DialogTitle className="gradient-text">Disable Two-Factor Authentication</DialogTitle>
            <DialogContent>
              {disableError && <Alert severity="error" sx={{ mb: 2 }}>{disableError}</Alert>}
              {disableSuccess && <Alert severity="success" sx={{ mb: 2 }}>2FA disabled!</Alert>}
              {sendCodeError && <Alert severity="error" sx={{ mb: 2 }}>{sendCodeError}</Alert>}
              {sendCodeSuccess && <Alert severity="success" sx={{ mb: 2 }}>Code sent to your email</Alert>}
              <Typography variant="body2" sx={{ mb: 2 }}>
                Enter the code sent to your email to disable 2FA.
              </Typography>
              <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
                <Button
                  variant="outlined"
                  onClick={handleSendCode}
                  disabled={sendCodeLoading || sendCodeSuccess}
                  size="small"
                >
                  {sendCodeLoading ? 'Sending...' : sendCodeSuccess ? 'Code Sent' : 'Send code'}
                </Button>
              </Box>
              <Box sx={{ display: 'flex', justifyContent: 'center', mb: 2 }}>
                <OtpInput
                  value={emailCode}
                  onChange={setEmailCode}
                  numInputs={6}
                  renderInput={(inputProps, idx) => (
                    <input
                      {...inputProps}
                      key={idx}
                      ref={el => { emailInputRefs.current[idx] = el; }}
                      onChange={e => {
                        inputProps.onChange?.(e);
                        if (e.target.value && emailInputRefs.current[idx + 1]) {
                          emailInputRefs.current[idx + 1]?.focus();
                        }
                      }}
                      onKeyDown={e => {
                        // Handle backspace
                        if (e.key === 'Backspace' && !e.currentTarget.value && idx > 0) {
                          e.preventDefault();
                          // Clear current field and move to previous
                          const newEmailCode = emailCode.split('');
                          newEmailCode[idx] = '';
                          setEmailCode(newEmailCode.join(''));
                          emailInputRefs.current[idx - 1]?.focus();
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
                        outline: document.activeElement === emailInputRefs.current[idx] ? `2px solid ${theme.palette.primary.main}` : 'none',
                        boxShadow: document.activeElement === emailInputRefs.current[idx] ? `0 0 0 2px ${theme.palette.primary.light}` : 'none',
                        transition: 'border 0.2s, box-shadow 0.2s',
                      }}
                      inputMode="numeric"
                      pattern="[0-9]*"
                      readOnly={disableLoading || disableSuccess}
                    />
                  )}
                  containerStyle={{ justifyContent: 'center' }}
                />
              </Box>
            </DialogContent>
            <DialogActions>
              <Button onClick={handleCloseDisable} disabled={disableLoading}>Cancel</Button>
              <Button variant="contained" color="error" disabled={!emailCode || disableLoading || disableSuccess} onClick={handleConfirmDisable2FA}>
                Confirm
              </Button>
            </DialogActions>
          </Dialog>
        </Box>
      )}
    </Paper>
  );
};

export default TwoFactorAuthSection;
