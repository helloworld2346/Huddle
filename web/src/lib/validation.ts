export interface ValidationError {
  field: string;
  message: string;
}

export interface FormData {
  fullName: string;
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
}

export interface ValidationResult {
  isValid: boolean;
  errors: ValidationError[];
}

// Validation rules
export const validationRules = {
  fullName: {
    required: true,
    minLength: 2,
    maxLength: 50,
    pattern: /^[a-zA-Z\s]+$/,
  },
  username: {
    required: true,
    minLength: 3,
    maxLength: 20,
    pattern: /^[a-zA-Z0-9_]+$/,
  },
  email: {
    required: true,
    pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  },
  password: {
    required: true,
    minLength: 8,
    maxLength: 50,
    pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]/,
  },
};

// Validation functions
export const validateFullName = (value: string): string | null => {
  if (!value.trim()) {
    return "Full name is required";
  }
  if (value.length < validationRules.fullName.minLength) {
    return `Full name must be at least ${validationRules.fullName.minLength} characters`;
  }
  if (value.length > validationRules.fullName.maxLength) {
    return `Full name must be less than ${validationRules.fullName.maxLength} characters`;
  }
  if (!validationRules.fullName.pattern.test(value)) {
    return "Full name can only contain letters and spaces";
  }
  return null;
};

export const validateUsername = (value: string): string | null => {
  if (!value.trim()) {
    return "Username is required";
  }
  if (value.length < validationRules.username.minLength) {
    return `Username must be at least ${validationRules.username.minLength} characters`;
  }
  if (value.length > validationRules.username.maxLength) {
    return `Username must be less than ${validationRules.username.maxLength} characters`;
  }
  if (!validationRules.username.pattern.test(value)) {
    return "Username can only contain letters, numbers, and underscores";
  }
  return null;
};

export const validateEmail = (value: string): string | null => {
  if (!value.trim()) {
    return "Email is required";
  }
  if (!validationRules.email.pattern.test(value)) {
    return "Please enter a valid email address";
  }
  return null;
};

export const validatePassword = (value: string): string | null => {
  if (!value) {
    return "Password is required";
  }
  if (value.length < validationRules.password.minLength) {
    return `Password must be at least ${validationRules.password.minLength} characters`;
  }
  if (value.length > validationRules.password.maxLength) {
    return `Password must be less than ${validationRules.password.maxLength} characters`;
  }
  if (!validationRules.password.pattern.test(value)) {
    return "Password must contain at least one lowercase letter, one uppercase letter, one number, and one special character";
  }
  return null;
};

export const validateConfirmPassword = (
  password: string,
  confirmPassword: string
): string | null => {
  if (!confirmPassword) {
    return "Please confirm your password";
  }
  if (password !== confirmPassword) {
    return "Passwords do not match";
  }
  return null;
};

// Password strength checker
export const checkPasswordStrength = (password: string): number => {
  let strength = 0;
  if (password.length >= 8) strength++;
  if (/[a-z]/.test(password)) strength++;
  if (/[A-Z]/.test(password)) strength++;
  if (/[0-9]/.test(password)) strength++;
  if (/[^A-Za-z0-9]/.test(password)) strength++;
  return strength;
};

export const getPasswordStrengthText = (strength: number): string => {
  if (strength < 3) return "Weak";
  if (strength < 4) return "Fair";
  if (strength < 5) return "Good";
  return "Strong";
};

// Main validation function
export const validateForm = (data: FormData): ValidationResult => {
  const errors: ValidationError[] = [];

  // Validate full name
  const fullNameError = validateFullName(data.fullName);
  if (fullNameError) {
    errors.push({ field: "fullName", message: fullNameError });
  }

  // Validate username
  const usernameError = validateUsername(data.username);
  if (usernameError) {
    errors.push({ field: "username", message: usernameError });
  }

  // Validate email
  const emailError = validateEmail(data.email);
  if (emailError) {
    errors.push({ field: "email", message: emailError });
  }

  // Validate password
  const passwordError = validatePassword(data.password);
  if (passwordError) {
    errors.push({ field: "password", message: passwordError });
  }

  // Validate confirm password
  const confirmPasswordError = validateConfirmPassword(
    data.password,
    data.confirmPassword
  );
  if (confirmPasswordError) {
    errors.push({ field: "confirmPassword", message: confirmPasswordError });
  }

  return {
    isValid: errors.length === 0,
    errors,
  };
};

// Real-time validation helpers
export const validateField = (
  field: keyof FormData,
  value: string,
  password?: string
): string | null => {
  switch (field) {
    case "fullName":
      return validateFullName(value);
    case "username":
      return validateUsername(value);
    case "email":
      return validateEmail(value);
    case "password":
      return validatePassword(value);
    case "confirmPassword":
      return password ? validateConfirmPassword(password, value) : null;
    default:
      return null;
  }
};
