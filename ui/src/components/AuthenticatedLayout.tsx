import React, { useState, useEffect } from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../auth/AuthContext';
import Layout from './Layout';
import ChangePasswordForm from './ChangePasswordForm';

interface AuthenticatedLayoutProps {
  children: React.ReactNode;
  showBackButton?: boolean;
  backTo?: string;
}

const AuthenticatedLayout: React.FC<AuthenticatedLayoutProps> = ({ 
  children, 
  showBackButton = false, 
  backTo = '/dashboard' 
}) => {
  const { isAuthenticated, hasTmpPassword } = useAuth();
  const [showPasswordDialog, setShowPasswordDialog] = useState(false);

  // Show password dialog when hasTmpPassword is true
  useEffect(() => {
    if (hasTmpPassword) {
      setShowPasswordDialog(true);
    } else {
      setShowPasswordDialog(false);
    }
  }, [hasTmpPassword]);

  // Handle password dialog close
  const handleClosePasswordDialog = () => {
    // The dialog will be closed automatically after successful password change
    // or when hasTmpPassword becomes false
    setShowPasswordDialog(false);
  };

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return (
    <Layout showBackButton={showBackButton} backTo={backTo}>
      {children}
      <ChangePasswordForm open={showPasswordDialog} onClose={handleClosePasswordDialog} />
    </Layout>
  );
};

export default AuthenticatedLayout; 