"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { MessageCircle, ArrowLeft, Mail, CheckCircle } from "lucide-react";
import { Background } from "@/components/ui/background";
import { useState } from "react";
import { authApi } from "@/lib/api";
import { handleError } from "@/lib/error-handler";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

export default function ForgotPasswordPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Basic validation
    if (!email.trim()) {
      setErrors({ email: "Email is required" });
      return;
    }

    // Email format validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      setErrors({ email: "Please enter a valid email address" });
      return;
    }

    setIsLoading(true);
    setErrors({});

    try {
      const response = await authApi.forgotPassword(email);

      if (!response.success) {
        const errorMessage =
          response.error?.message ||
          response.message ||
          "Failed to send reset email";
        setErrors({ general: errorMessage });

        toast.error("Failed to send reset email", {
          description: errorMessage,
          duration: 5000,
        });
        return;
      }

      // Success
      setIsSubmitted(true);
      toast.success("Reset email sent!", {
        description: "Please check your email for password reset instructions.",
        duration: 5000,
      });
    } catch (error) {
      const { fieldErrors, generalError } = handleError(error);
      setErrors((prev) => ({
        ...prev,
        ...fieldErrors,
        ...(generalError && { general: generalError }),
      }));

      if (generalError) {
        toast.error("Failed to send reset email", {
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
            <Link href="/login">
              <ArrowLeft className="w-4 h-4 mr-2" />
              Back to Login
            </Link>
          </Button>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 flex items-center justify-center relative z-10">
        <div className="w-full max-w-md mx-auto px-4">
          <div className="text-center mb-8">
            <div className="w-16 h-16 bg-gradient-to-br from-purple-400/20 to-blue-400/20 rounded-full flex items-center justify-center mx-auto mb-4">
              <Mail className="w-8 h-8 text-purple-400" />
            </div>
            <h1 className="text-3xl font-bold bg-gradient-to-r from-white via-purple-100 to-blue-100 bg-clip-text text-transparent mb-2">
              Forgot Password?
            </h1>
            <p className="text-white/60">
              {isSubmitted
                ? "Check your email for reset instructions"
                : "Enter your email address and we'll send you a link to reset your password"}
            </p>
          </div>

          <div className="backdrop-blur-md bg-white/10 border border-white/20 rounded-2xl p-8 shadow-2xl">
            {!isSubmitted ? (
              <form onSubmit={handleSubmit} className="space-y-6">
                {/* Email Field */}
                <div className="space-y-2">
                  <Label htmlFor="email" className="text-white/80 font-medium">
                    Email Address
                  </Label>
                  <Input
                    id="email"
                    type="email"
                    placeholder="Enter your email address"
                    className={`w-full bg-purple-900/20 border text-white placeholder:text-purple-200/60 focus:ring-purple-400/20 backdrop-blur-sm h-12 px-4 text-base [&:-webkit-autofill]:bg-purple-900/20 [&:-webkit-autofill]:text-white [&:-webkit-autofill]:shadow-[0_0_0_30px_rgba(88,28,135,0.2)_inset] [&:-webkit-autofill]:border-purple-400 ${
                      errors.email
                        ? "border-red-400 focus:border-red-400"
                        : "border-purple-400/30 focus:border-purple-400"
                    }`}
                    value={email}
                    onChange={(e) => {
                      setEmail(e.target.value);
                      setErrors((prev) => ({ ...prev, email: "" }));
                    }}
                    required
                  />
                  {errors.email && (
                    <div className="text-red-400 text-sm flex items-center space-x-2">
                      <div className="w-1 h-1 bg-red-400 rounded-full"></div>
                      <span>{errors.email}</span>
                    </div>
                  )}
                </div>

                {/* General Error Message */}
                {errors.general && (
                  <div className="text-red-400 text-sm flex items-center space-x-2 bg-red-400/10 border border-red-400/20 rounded-lg p-3">
                    <div className="w-1 h-1 bg-red-400 rounded-full"></div>
                    <span>{errors.general}</span>
                  </div>
                )}

                {/* Submit Button */}
                <div className="relative group">
                  <div className="absolute inset-0 bg-gradient-to-r from-purple-500 to-blue-500 rounded-lg blur opacity-75 group-hover:opacity-100 transition duration-1000 group-hover:duration-200 animate-pulse-scale"></div>
                  <Button
                    type="submit"
                    size="lg"
                    disabled={isLoading}
                    className="relative w-full bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700 text-white border-0 shadow-lg transform hover:scale-105 transition-all duration-300 font-semibold tracking-wide disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <span className="relative z-10">
                      {isLoading ? "Sending..." : "Send Reset Link"}
                    </span>
                    <div className="absolute inset-0 bg-gradient-to-r from-purple-400/20 to-blue-400/20 rounded-lg animate-pulse-scale"></div>
                  </Button>
                </div>
              </form>
            ) : (
              <div className="text-center space-y-6">
                {/* Success Icon */}
                <div className="w-16 h-16 bg-gradient-to-br from-green-400/20 to-emerald-400/20 rounded-full flex items-center justify-center mx-auto">
                  <CheckCircle className="w-8 h-8 text-green-400" />
                </div>

                {/* Success Message */}
                <div className="space-y-2">
                  <h3 className="text-xl font-semibold text-white">
                    Check Your Email
                  </h3>
                  <p className="text-white/60 text-sm">
                    We&apos;ve sent a password reset link to{" "}
                    <span className="text-purple-400 font-medium">{email}</span>
                  </p>
                </div>

                {/* Instructions */}
                <div className="bg-black/20 border border-white/10 rounded-lg p-4 text-left">
                  <h4 className="text-white font-medium mb-2">
                    What to do next:
                  </h4>
                  <ul className="text-white/60 text-sm space-y-1">
                    <li>• Check your email inbox</li>
                    <li>• Click the reset link in the email</li>
                    <li>• Create a new password</li>
                    <li>• Sign in with your new password</li>
                  </ul>
                </div>

                {/* Resend Button */}
                <Button
                  variant="outline"
                  className="w-full border-white/20 bg-black/20 text-white hover:bg-white/10 backdrop-blur-sm"
                  onClick={() => setIsSubmitted(false)}
                >
                  Send Another Email
                </Button>
              </div>
            )}

            {/* Back to Login */}
            <div className="text-center mt-6">
              <span className="text-white/60">Remember your password? </span>
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
