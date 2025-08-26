import Link from "next/link";
import { Button } from "@/components/ui/button";
import { MessageCircle, Sparkles } from "lucide-react";
import { Background } from "@/components/ui/background";

export default function HomePage() {
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
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 flex items-center justify-center relative z-10">
        <div className="text-center max-w-md mx-auto px-4">
          <div className="mb-6">
            <div className="inline-flex items-center px-4 py-2 bg-white/10 backdrop-blur-sm border border-white/20 rounded-full text-sm font-medium text-white shadow-sm">
              <Sparkles className="w-4 h-4 mr-2 text-purple-300" />
              Next Generation Chat Platform
            </div>
          </div>

          <h1 className="text-4xl md:text-5xl font-bold mb-6 bg-gradient-to-r from-white via-purple-100 to-blue-100 bg-clip-text text-transparent">
            Real-time Chat Platform
          </h1>

          <p className="text-lg text-white/80 mb-10 leading-relaxed">
            Connect, chat, and collaborate with your team instantly.
          </p>

          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            {/* Primary Button - Cyberpunk Style */}
            <div className="relative group">
              <div className="absolute inset-0 bg-gradient-to-r from-purple-500 to-blue-500 rounded-lg blur opacity-75 group-hover:opacity-100 transition duration-1000 group-hover:duration-200 animate-pulse-scale"></div>
              <Button
                size="lg"
                className="relative w-full sm:w-auto bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700 text-white border-0 shadow-lg transform hover:scale-105 transition-all duration-300 font-semibold tracking-wide"
                asChild
              >
                <Link
                  href="/login"
                  className="flex items-center justify-center"
                >
                  <span className="relative z-10">Sign In</span>
                  <div className="absolute inset-0 bg-gradient-to-r from-purple-400/20 to-blue-400/20 rounded-lg animate-pulse-scale"></div>
                </Link>
              </Button>
            </div>

            {/* Secondary Button - Neon Border Style */}
            <div className="relative group">
              <div className="absolute inset-0 bg-gradient-to-r from-purple-400 to-blue-400 rounded-lg blur opacity-50 group-hover:opacity-75 transition duration-300"></div>
              <Button
                size="lg"
                variant="outline"
                className="relative w-full sm:w-auto border-2 border-purple-400/50 hover:border-purple-400 bg-black/20 backdrop-blur-sm shadow-lg transform hover:scale-105 transition-all duration-300 text-white hover:bg-purple-400/10 font-semibold tracking-wide"
                asChild
              >
                <Link
                  href="/register"
                  className="flex items-center justify-center"
                >
                  <span className="relative z-10">Create Account</span>
                  <div className="absolute inset-0 bg-gradient-to-r from-purple-400/5 to-blue-400/5 rounded-lg"></div>
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="border-t border-white/20 bg-white/5 backdrop-blur-sm py-8 relative z-10">
        <div className="container mx-auto px-4 text-center">
          <p className="text-white/60">
            &copy; 2024 Huddle. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  );
}
