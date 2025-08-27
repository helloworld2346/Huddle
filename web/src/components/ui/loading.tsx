import { Background } from "./background";

interface LoadingProps {
  message?: string;
  size?: "sm" | "md" | "lg";
}

export function Loading({ message = "Loading...", size = "md" }: LoadingProps) {
  const sizeClasses = {
    sm: "w-8 h-8",
    md: "w-16 h-16",
    lg: "w-24 h-24",
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 flex items-center justify-center">
      <Background />
      <div className="text-center relative z-10">
        <div
          className={`${sizeClasses[size]} border-4 border-purple-400 border-t-transparent rounded-full animate-spin mx-auto mb-4`}
        ></div>
        <p className="text-white/60">{message}</p>
      </div>
    </div>
  );
}
