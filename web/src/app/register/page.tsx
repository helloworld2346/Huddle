"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  MessageCircle,
  Eye,
  EyeOff,
  ArrowLeft,
  CheckCircle,
} from "lucide-react";
import { Background } from "@/components/ui/background";
import { useState } from "react";
import {
  validateForm,
  validateField,
  checkPasswordStrength,
  getPasswordStrengthText,
  type FormData,
  type ValidationError,
} from "@/lib/validation";
import { authApi } from "@/lib/api";
import { handleRegistrationError, handleError } from "@/lib/error-handler";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

export default function RegisterPage() {
  const router = useRouter();
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  // Form data
  const [formData, setFormData] = useState<FormData>({
    fullName: "",
    username: "",
    email: "",
    password: "",
    confirmPassword: "",
  });

  // Validation states
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [passwordStrength, setPasswordStrength] = useState(0);

  const handleInputChange = (field: keyof FormData, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));

    // Clear error for this field
    setErrors((prev) => ({ ...prev, [field]: "" }));

    // Special handling for password
    if (field === "password") {
      const strength = checkPasswordStrength(value);
      setPasswordStrength(strength);
    }

    // Special handling for confirm password
    if (field === "confirmPassword") {
      if (formData.password && value !== formData.password) {
        setErrors((prev) => ({
          ...prev,
          confirmPassword: "Passwords do not match",
        }));
      } else {
        setErrors((prev) => ({ ...prev, confirmPassword: "" }));
      }
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Validate all fields
    const validation = validateForm(formData);

    if (!validation.isValid) {
      const errorMap: Record<string, string> = {};
      validation.errors.forEach((error) => {
        errorMap[error.field] = error.message;
      });
      setErrors(errorMap);
      return;
    }

    setIsLoading(true);

    try {
      const response = await authApi.register({
        username: formData.username,
        email: formData.email,
        display_name: formData.fullName,
        password: formData.password,
      });

      if (!response.success) {
        const { fieldErrors, generalError } = handleRegistrationError(response);
        setErrors((prev) => ({
          ...prev,
          ...fieldErrors,
          ...(generalError && { general: generalError }),
        }));

        // Show error toast for general errors
        if (generalError) {
          toast.error("Registration failed", {
            description: generalError,
            duration: 5000,
          });
        }
        return;
      }

      // Success - show toast and redirect to login
      toast.success(
        "Registration successful! Please login with your new account.",
        {
          duration: 4000,
          description: "You will be redirected to the login page.",
        }
      );

      // Redirect after a short delay to show the toast
      setTimeout(() => {
        router.push("/login");
      }, 1500);
    } catch (error) {
      const { fieldErrors, generalError } = handleError(error);
      setErrors((prev) => ({
        ...prev,
        ...fieldErrors,
        ...(generalError && { general: generalError }),
      }));

      // Show error toast
      if (generalError) {
        toast.error("Registration failed", {
          description: generalError,
          duration: 5000,
        });
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 flex flex-col relative">
      <Background />

      {/* Header */}
      <header className="backdrop-blur-md bg-white/10 border-b border-white/20 sticky top-0 z-50">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="relative">
              <div className="w-10 h-10 bg-gradient-to-br from-purple-400 to-blue-400 rounded-xl flex items-center justify-center shadow-lg">
                <MessageCircle className="w-6 h-6 text-white" />
              </div>
              <div className="absolute -top-1 -right-1 w-4 h-4 bg-green-400 rounded-full border-2 border-white animate-pulse"></div>
            </div>
            <div>
              <span className="text-2xl font-bold bg-gradient-to-r from-purple-400 to-blue-400 bg-clip-text text-transparent">
                Huddle
              </span>
              <div className="text-xs text-white/60 -mt-1">Real-time Chat</div>
            </div>
          </div>

          <Button
            variant="ghost"
            className="text-white/80 hover:text-white hover:bg-white/10"
            asChild
          >
            <Link href="/">
              <ArrowLeft className="w-4 h-4 mr-2" />
              Back to Home
            </Link>
          </Button>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 flex items-center justify-center relative z-10 py-8">
        <div className="w-full max-w-md mx-auto px-4">
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold bg-gradient-to-r from-white via-purple-100 to-blue-100 bg-clip-text text-transparent mb-2">
              Create Account
            </h1>
            <p className="text-white/60">
              Join Huddle and start chatting with your team
            </p>
          </div>

          <div className="backdrop-blur-md bg-white/10 border border-white/20 rounded-2xl p-8 shadow-2xl">
            <form onSubmit={handleSubmit} className="space-y-6">
              {/* Full Name Field */}
              <div className="space-y-2">
                <Label htmlFor="fullName" className="text-white/80 font-medium">
                  Full Name
                </Label>
                <Input
                  id="fullName"
                  type="text"
                  placeholder="Enter your full name"
                  className={`w-full bg-purple-900/20 border text-white placeholder:text-purple-200/60 focus:ring-purple-400/20 backdrop-blur-sm h-12 px-4 text-base [&:-webkit-autofill]:bg-purple-900/20 [&:-webkit-autofill]:text-white [&:-webkit-autofill]:shadow-[0_0_0_30px_rgba(88,28,135,0.2)_inset] [&:-webkit-autofill]:border-purple-400 ${
                    errors.fullName
                      ? "border-red-400 focus:border-red-400"
                      : "border-purple-400/30 focus:border-purple-400"
                  }`}
                  value={formData.fullName}
                  onChange={(e) =>
                    handleInputChange("fullName", e.target.value)
                  }
                  required
                />
                {errors.fullName && (
                  <div className="text-red-400 text-sm flex items-center space-x-2">
                    <div className="w-1 h-1 bg-red-400 rounded-full"></div>
                    <span>{errors.fullName}</span>
                  </div>
                )}
              </div>

              {/* Username Field */}
              <div className="space-y-2">
                <Label htmlFor="username" className="text-white/80 font-medium">
                  Username
                </Label>
                <Input
                  id="username"
                  type="text"
                  placeholder="Choose a username"
                  className={`w-full bg-purple-900/20 border text-white placeholder:text-purple-200/60 focus:ring-purple-400/20 backdrop-blur-sm h-12 px-4 text-base [&:-webkit-autofill]:bg-purple-900/20 [&:-webkit-autofill]:text-white [&:-webkit-autofill]:shadow-[0_0_0_30px_rgba(88,28,135,0.2)_inset] [&:-webkit-autofill]:border-purple-400 ${
                    errors.username
                      ? "border-red-400 focus:border-red-400"
                      : "border-purple-400/30 focus:border-purple-400"
                  }`}
                  value={formData.username}
                  onChange={(e) =>
                    handleInputChange("username", e.target.value)
                  }
                  required
                />
                {errors.username && (
                  <div className="text-red-400 text-sm flex items-center space-x-2">
                    <div className="w-1 h-1 bg-red-400 rounded-full"></div>
                    <span>{errors.username}</span>
                  </div>
                )}
              </div>

              {/* Email Field */}
              <div className="space-y-2">
                <Label htmlFor="email" className="text-white/80 font-medium">
                  Email Address
                </Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="Enter your email"
                  className={`w-full bg-purple-900/20 border text-white placeholder:text-purple-200/60 focus:ring-purple-400/20 backdrop-blur-sm h-12 px-4 text-base [&:-webkit-autofill]:bg-purple-900/20 [&:-webkit-autofill]:text-white [&:-webkit-autofill]:shadow-[0_0_0_30px_rgba(88,28,135,0.2)_inset] [&:-webkit-autofill]:border-purple-400 ${
                    errors.email
                      ? "border-red-400 focus:border-red-400"
                      : "border-purple-400/30 focus:border-purple-400"
                  }`}
                  value={formData.email}
                  onChange={(e) => handleInputChange("email", e.target.value)}
                  required
                />
                {errors.email && (
                  <div className="text-red-400 text-sm flex items-center space-x-2">
                    <div className="w-1 h-1 bg-red-400 rounded-full"></div>
                    <span>{errors.email}</span>
                  </div>
                )}
              </div>

              {/* Password Field */}
              <div className="space-y-2">
                <Label htmlFor="password" className="text-white/80 font-medium">
                  Password
                </Label>
                <div className="relative">
                  <Input
                    id="password"
                    type={showPassword ? "text" : "password"}
                    placeholder="Create a strong password"
                    className={`w-full bg-purple-900/20 border text-white placeholder:text-purple-200/60 focus:ring-purple-400/20 backdrop-blur-sm h-12 px-4 text-base pr-12 [&:-webkit-autofill]:bg-purple-900/20 [&:-webkit-autofill]:text-white [&:-webkit-autofill]:shadow-[0_0_0_30px_rgba(88,28,135,0.2)_inset] [&:-webkit-autofill]:border-purple-400 ${
                      errors.password
                        ? "border-red-400 focus:border-red-400"
                        : "border-purple-400/30 focus:border-purple-400"
                    }`}
                    value={formData.password}
                    onChange={(e) =>
                      handleInputChange("password", e.target.value)
                    }
                    required
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-white/40 hover:text-white/60 transition-colors"
                  >
                    {showPassword ? (
                      <EyeOff className="w-4 h-4" />
                    ) : (
                      <Eye className="w-4 h-4" />
                    )}
                  </button>
                </div>

                {errors.password && (
                  <div className="text-red-400 text-sm flex items-center space-x-2">
                    <div className="w-1 h-1 bg-red-400 rounded-full"></div>
                    <span>{errors.password}</span>
                  </div>
                )}

                {/* Password Strength Indicator */}
                <div className="space-y-2">
                  <div className="flex space-x-1">
                    {[1, 2, 3, 4, 5].map((level) => (
                      <div
                        key={level}
                        className={`h-1 flex-1 rounded-full transition-all duration-300 ${
                          passwordStrength >= level
                            ? level <= 2
                              ? "bg-red-400"
                              : level <= 3
                              ? "bg-yellow-400"
                              : level <= 4
                              ? "bg-blue-400"
                              : "bg-green-400"
                            : "bg-white/20"
                        }`}
                      />
                    ))}
                  </div>
                  <div className="flex items-center space-x-2 text-xs">
                    <CheckCircle
                      className={`w-3 h-3 ${
                        passwordStrength >= 5
                          ? "text-green-400"
                          : "text-white/40"
                      }`}
                    />
                    <span
                      className={`${
                        passwordStrength >= 5
                          ? "text-green-400"
                          : "text-white/40"
                      }`}
                    >
                      {getPasswordStrengthText(passwordStrength)} password
                    </span>
                  </div>
                </div>
              </div>

              {/* Confirm Password Field */}
              <div className="space-y-2">
                <Label
                  htmlFor="confirmPassword"
                  className="text-white/80 font-medium"
                >
                  Confirm Password
                </Label>
                <div className="relative">
                  <Input
                    id="confirmPassword"
                    type={showConfirmPassword ? "text" : "password"}
                    placeholder="Confirm your password"
                    className={`w-full bg-purple-900/20 border text-white placeholder:text-purple-200/60 focus:ring-purple-400/20 backdrop-blur-sm h-12 px-4 text-base pr-12 [&:-webkit-autofill]:bg-purple-900/20 [&:-webkit-autofill]:text-white [&:-webkit-autofill]:shadow-[0_0_0_30px_rgba(88,28,135,0.2)_inset] [&:-webkit-autofill]:border-purple-400 ${
                      errors.confirmPassword
                        ? "border-red-400 focus:border-red-400"
                        : "border-purple-400/30 focus:border-purple-400"
                    }`}
                    value={formData.confirmPassword}
                    onChange={(e) =>
                      handleInputChange("confirmPassword", e.target.value)
                    }
                    required
                  />
                  <button
                    type="button"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-white/40 hover:text-white/60 transition-colors"
                  >
                    {showConfirmPassword ? (
                      <EyeOff className="w-4 h-4" />
                    ) : (
                      <Eye className="w-4 h-4" />
                    )}
                  </button>
                </div>

                {/* Password Error Message */}
                {errors.confirmPassword && (
                  <div className="text-red-400 text-sm flex items-center space-x-2">
                    <div className="w-1 h-1 bg-red-400 rounded-full"></div>
                    <span>{errors.confirmPassword}</span>
                  </div>
                )}
              </div>

              {/* Terms and Conditions */}
              <div className="flex items-start space-x-3">
                <label
                  htmlFor="terms"
                  className="flex items-start space-x-3 cursor-pointer group"
                >
                  <div className="relative mt-1">
                    <input
                      type="checkbox"
                      id="terms"
                      className="w-5 h-5 text-purple-400 bg-black/20 border-white/20 rounded focus:ring-purple-400/20 focus:ring-2 transition-all duration-200 cursor-pointer relative z-10"
                      required
                    />
                    <div className="absolute inset-0 w-5 h-5 border border-purple-400/30 rounded group-hover:border-purple-400/50 transition-colors"></div>
                  </div>
                  <span className="text-white/70 text-sm font-medium group-hover:text-white/90 transition-colors">
                    I agree to the{" "}
                    <Link
                      href="/terms"
                      className="text-purple-400 hover:text-purple-300 font-medium hover:underline"
                      onClick={(e) => e.stopPropagation()}
                    >
                      Terms of Service
                    </Link>{" "}
                    and{" "}
                    <Link
                      href="/privacy"
                      className="text-purple-400 hover:text-purple-300 font-medium hover:underline"
                      onClick={(e) => e.stopPropagation()}
                    >
                      Privacy Policy
                    </Link>
                  </span>
                </label>
              </div>

              {/* General Error Message */}
              {errors.general && (
                <div className="text-red-400 text-sm flex items-center space-x-2 bg-red-400/10 border border-red-400/20 rounded-lg p-3">
                  <div className="w-1 h-1 bg-red-400 rounded-full"></div>
                  <span>{errors.general}</span>
                </div>
              )}

              {/* Create Account Button */}
              <div className="relative group">
                <div className="absolute inset-0 bg-gradient-to-r from-purple-500 to-blue-500 rounded-lg blur opacity-75 group-hover:opacity-100 transition duration-1000 group-hover:duration-200 animate-pulse-scale"></div>
                <Button
                  type="submit"
                  size="lg"
                  disabled={isLoading}
                  className="relative w-full bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700 text-white border-0 shadow-lg transform hover:scale-105 transition-all duration-300 font-semibold tracking-wide disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <span className="relative z-10">
                    {isLoading ? "Creating Account..." : "Create Account"}
                  </span>
                  <div className="absolute inset-0 bg-gradient-to-r from-purple-400/20 to-blue-400/20 rounded-lg animate-pulse-scale"></div>
                </Button>
              </div>
            </form>

            {/* Divider */}
            <div className="relative my-6">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-white/20"></div>
              </div>
              <div className="relative flex justify-center text-sm">
                <span className="px-2 bg-transparent text-white/40">
                  Or continue with
                </span>
              </div>
            </div>

            {/* Social Login */}
            <div className="space-y-3">
              <Button
                variant="outline"
                className="w-full border-white/20 bg-black/20 text-white hover:bg-white/10 backdrop-blur-sm"
              >
                Continue with Google
              </Button>
              <Button
                variant="outline"
                className="w-full border-white/20 bg-black/20 text-white hover:bg-white/10 backdrop-blur-sm"
              >
                Continue with GitHub
              </Button>
            </div>

            {/* Sign In Link */}
            <div className="text-center mt-6">
              <span className="text-white/60">Already have an account? </span>
              <Link
                href="/login"
                className="text-purple-400 hover:text-purple-300 font-medium transition-colors hover:underline"
              >
                Sign in
              </Link>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
