"use client";

import { Toaster as Sonner, ToasterProps } from "sonner";

const Toaster = ({ ...props }: ToasterProps) => {
  return (
    <Sonner
      theme="dark"
      className="toaster group"
      style={
        {
          // Cyberpunk theme colors
          "--normal-bg": "rgba(15, 23, 42, 0.95)", // slate-900 with opacity
          "--normal-text": "rgba(255, 255, 255, 0.9)",
          "--normal-border": "rgba(147, 51, 234, 0.3)", // purple-600 with opacity

          // Success colors
          "--success-bg": "rgba(34, 197, 94, 0.1)", // green-500 with low opacity
          "--success-text": "rgb(74, 222, 128)", // green-400
          "--success-border": "rgba(74, 222, 128, 0.3)", // green-400 with opacity

          // Error colors
          "--error-bg": "rgba(239, 68, 68, 0.1)", // red-500 with low opacity
          "--error-text": "rgb(248, 113, 113)", // red-400
          "--error-border": "rgba(248, 113, 113, 0.3)", // red-400 with opacity

          // Warning colors
          "--warning-bg": "rgba(245, 158, 11, 0.1)", // amber-500 with low opacity
          "--warning-text": "rgb(251, 191, 36)", // amber-400
          "--warning-border": "rgba(251, 191, 36, 0.3)", // amber-400 with opacity

          // Info colors
          "--info-bg": "rgba(59, 130, 246, 0.1)", // blue-500 with low opacity
          "--info-text": "rgb(96, 165, 250)", // blue-400
          "--info-border": "rgba(96, 165, 250, 0.3)", // blue-400 with opacity

          // Description text
          "--description": "rgba(255, 255, 255, 0.7)",

          // Action button
          "--action": "rgba(147, 51, 234, 0.8)", // purple-600 with opacity
          "--action-hover": "rgba(147, 51, 234, 1)", // purple-600
        } as React.CSSProperties
      }
      toastOptions={{
        classNames: {
          toast:
            "group toast group-[.toaster]:bg-background group-[.toaster]:text-foreground group-[.toaster]:border-border group-[.toaster]:shadow-lg backdrop-blur-md border",
          description: "group-[.toast]:text-muted-foreground",
          actionButton:
            "group-[.toast]:bg-primary group-[.toast]:text-primary-foreground",
          cancelButton:
            "group-[.toast]:bg-muted group-[.toast]:text-muted-foreground",
          success:
            "group-[.toast]:bg-green-500/10 group-[.toast]:border-green-400/30 group-[.toast]:text-green-400",
          error:
            "group-[.toast]:bg-red-500/10 group-[.toast]:border-red-400/30 group-[.toast]:text-red-400",
          warning:
            "group-[.toast]:bg-amber-500/10 group-[.toast]:border-amber-400/30 group-[.toast]:text-amber-400",
          info: "group-[.toast]:bg-blue-500/10 group-[.toast]:border-blue-400/30 group-[.toast]:text-blue-400",
        },
        duration: 4000,
        style: {
          background: "rgba(15, 23, 42, 0.95)",
          border: "1px solid rgba(147, 51, 234, 0.3)",
          color: "rgba(255, 255, 255, 0.9)",
          backdropFilter: "blur(12px)",
          boxShadow: "0 8px 32px rgba(0, 0, 0, 0.3)",
        },
      }}
      {...props}
    />
  );
};

export { Toaster };
