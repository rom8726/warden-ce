import React, { useState } from 'react';
import { 
  Button, 
  CircularProgress,
  Alert,
  Typography,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Paper
} from '@mui/material';
import { useAuth } from '../auth/AuthContext';
import apiClient from '../api/apiClient';

interface LicenseDialogProps {
  open: boolean;
  onClose: () => void;
}

const LicenseDialog: React.FC<LicenseDialogProps> = ({ open, onClose }) => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { logout, updateUserData } = useAuth();

  // Simplified for community edition
  const handleAccept = async () => {
    setIsLoading(true);
    
    try {
      // In community edition, just update user data and close
      await updateUserData();
      onClose();
    } catch (err) {
      console.error('Failed to accept license:', err);
      setError('Failed to process license acceptance. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  // Simplified for community edition
  const handleDecline = async () => {
    setIsLoading(true);
    
    try {
      // In community edition, just logout
      logout();
    } catch (err) {
      console.error('Failed to decline license:', err);
      setError('Failed to process license decline. Please try again.');
      setIsLoading(false);
    }
  };

  return (
    <Dialog 
      open={open} 
      onClose={() => {}} // Empty function to prevent closing by clicking outside
      maxWidth="md" 
      fullWidth
      disableEscapeKeyDown // Prevent closing with Escape key
    >
      <DialogTitle>
        <Typography variant="h5" component="div" align="center">
          License Agreement
        </Typography>
      </DialogTitle>
      <DialogContent>
        <Typography variant="body1" sx={{ mb: 2 }}>
          Please read and accept the following license agreement to continue using the application.
        </Typography>
        
        {error && (
          <Alert severity="error" sx={{ width: '100%', mb: 2 }}>
            {error}
          </Alert>
        )}
        
        <Paper 
          elevation={0} 
          variant="outlined" 
          sx={{ 
            p: 2, 
            maxHeight: '400px', 
            overflow: 'auto',
            mb: 2,
            backgroundColor: '#f5f5f5'
          }}
        >
          <Typography variant="body2" component="div" sx={{ whiteSpace: 'pre-line' }}>
            {`LICENSE AGREEMENT

This software is provided under the terms of the license agreement below.
By using this software, you agree to be bound by the terms of this license.

1. GRANT OF LICENSE
Subject to the terms and conditions of this Agreement, Licensor hereby grants to Licensee a non-exclusive, non-transferable license to use the Software.

2. RESTRICTIONS
Licensee shall not:
- Modify, adapt, translate, or create derivative works based upon the Software
- Reverse engineer, decompile, disassemble, or otherwise attempt to discover the source code of the Software
- Rent, lease, loan, sell, sublicense, distribute or otherwise transfer the Software to any third party
- Remove or alter any proprietary notices or labels on the Software

3. OWNERSHIP
The Software is licensed, not sold. Licensor retains all rights, title, and interest in and to the Software.

4. TERMINATION
This Agreement shall terminate automatically if Licensee fails to comply with any of the terms and conditions of this Agreement.

5. DISCLAIMER OF WARRANTY
THE SOFTWARE IS PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE.

6. LIMITATION OF LIABILITY
IN NO EVENT SHALL LICENSOR BE LIABLE FOR ANY SPECIAL, INCIDENTAL, INDIRECT, OR CONSEQUENTIAL DAMAGES WHATSOEVER ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THE SOFTWARE.

7. GOVERNING LAW
This Agreement shall be governed by and construed in accordance with the laws of the jurisdiction in which Licensor is located.

By accepting this license, you acknowledge that you have read this Agreement, understand it, and agree to be bound by its terms and conditions.`}
          </Typography>
        </Paper>
      </DialogContent>
      <DialogActions sx={{ justifyContent: 'space-between', px: 3, pb: 3 }}>
        <Button 
          onClick={handleDecline} 
          disabled={isLoading}
          variant="outlined"
          color="error"
        >
          {isLoading ? <CircularProgress size={24} /> : 'Decline'}
        </Button>
        <Button 
          onClick={handleAccept} 
          disabled={isLoading}
          variant="contained"
          color="primary"
        >
          {isLoading ? <CircularProgress size={24} /> : 'Accept'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default LicenseDialog;