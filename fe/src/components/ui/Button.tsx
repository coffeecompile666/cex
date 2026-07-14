import React, { ButtonHTMLAttributes } from 'react';
import styled from 'styled-components';
import { Button as HeadlessButton } from '@headlessui/react';

const StyledButton = styled(HeadlessButton)`
  background-color: var(--accent-color);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: var(--radius-sm);
  font-weight: 600;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
  width: 100%;
  
  &:hover {
    background-color: var(--accent-hover);
  }
  
  &:disabled {
    background-color: #a7f3d0;
    cursor: not-allowed;
  }
`;

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  children: React.ReactNode;
  isLoading?: boolean;
}

export const Button: React.FC<ButtonProps> = ({ children, isLoading, ...props }) => {
  return (
    <StyledButton {...props} disabled={isLoading || props.disabled}>
      {isLoading ? 'Processing...' : children}
    </StyledButton>
  );
};
