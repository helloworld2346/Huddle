export function Background() {
  return (
    <div className="fixed inset-0 overflow-hidden pointer-events-none">
      {/* Base Gradient */}
      <div className="absolute inset-0 bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900"></div>

      {/* Dot Pattern */}
      <div className="absolute inset-0 opacity-20">
        <div
          className="absolute inset-0"
          style={{
            backgroundImage: `radial-gradient(circle at 1px 1px, rgba(255,255,255,0.1) 1px, transparent 0)`,
            backgroundSize: "40px 40px",
          }}
        ></div>
      </div>

      {/* Floating Geometric Shapes */}
      <div className="absolute top-20 left-20 w-32 h-32 border border-white/10 rotate-45 animate-rotate-slow"></div>
      <div className="absolute top-40 right-32 w-24 h-24 border border-white/10 rotate-45 animate-rotate-slow animation-delay-2000"></div>
      <div className="absolute bottom-32 left-1/2 w-16 h-16 border border-white/10 rotate-45 animate-rotate-slow animation-delay-4000"></div>

      {/* Additional Rotating Squares */}
      <div className="absolute top-1/4 right-1/6 w-20 h-20 border border-white/10 rotate-45 animate-rotate-slow animation-delay-1000"></div>
      <div className="absolute bottom-1/4 right-1/3 w-28 h-28 border border-white/10 rotate-45 animate-rotate-slow animation-delay-3000"></div>
      <div className="absolute top-2/3 left-1/4 w-12 h-12 border border-white/10 rotate-45 animate-rotate-slow animation-delay-1500"></div>
      <div className="absolute top-1/3 left-2/3 w-36 h-36 border border-white/10 rotate-45 animate-rotate-slow animation-delay-2500"></div>
      <div className="absolute bottom-1/3 left-1/6 w-8 h-8 border border-white/10 rotate-45 animate-rotate-slow animation-delay-3500"></div>

      {/* Neon Glow Lines */}
      <div className="absolute top-0 left-1/4 w-px h-full bg-gradient-to-b from-transparent via-purple-400/30 to-transparent animate-pulse-scale"></div>
      <div className="absolute top-0 right-1/4 w-px h-full bg-gradient-to-b from-transparent via-blue-400/30 to-transparent animate-pulse-scale animation-delay-2000"></div>
      <div className="absolute top-1/4 left-0 w-full h-px bg-gradient-to-r from-transparent via-purple-400/30 to-transparent animate-pulse-scale animation-delay-1000"></div>
      <div className="absolute bottom-1/4 left-0 w-full h-px bg-gradient-to-r from-transparent via-blue-400/30 to-transparent animate-pulse-scale animation-delay-3000"></div>

      {/* Orbital Rings */}
      <div className="absolute top-1/2 left-1/2 w-96 h-96 border border-white/5 rounded-full -translate-x-1/2 -translate-y-1/2 animate-rotate-slow"></div>
      <div
        className="absolute top-1/2 left-1/2 w-64 h-64 border border-white/5 rounded-full -translate-x-1/2 -translate-y-1/2 animate-rotate-slow animation-delay-2000"
        style={{ animationDirection: "reverse" }}
      ></div>
      <div className="absolute top-1/2 left-1/2 w-32 h-32 border border-white/5 rounded-full -translate-x-1/2 -translate-y-1/2 animate-rotate-slow animation-delay-4000"></div>

      {/* Floating Particles */}
      <div className="absolute top-1/4 left-1/4 w-2 h-2 bg-purple-400/60 rounded-full animate-bounce animate-glow"></div>
      <div className="absolute top-1/3 right-1/3 w-1.5 h-1.5 bg-blue-400/60 rounded-full animate-bounce animation-delay-1000"></div>
      <div className="absolute bottom-1/4 left-1/3 w-1 h-1 bg-purple-400/60 rounded-full animate-bounce animation-delay-2000"></div>
      <div className="absolute bottom-1/3 right-1/4 w-1.5 h-1.5 bg-blue-400/60 rounded-full animate-bounce animation-delay-3000"></div>
      <div className="absolute top-1/2 left-1/6 w-1 h-1 bg-purple-400/60 rounded-full animate-bounce animation-delay-1500"></div>
      <div className="absolute top-2/3 right-1/6 w-1.5 h-1.5 bg-blue-400/60 rounded-full animate-bounce animation-delay-2500"></div>

      {/* Scanning Lines */}
      <div className="absolute top-0 left-0 w-full h-px bg-gradient-to-r from-transparent via-purple-400/20 to-transparent animate-slide-in-left"></div>
      <div className="absolute bottom-0 left-0 w-full h-px bg-gradient-to-r from-transparent via-blue-400/20 to-transparent animate-slide-in-right animation-delay-2000"></div>

      {/* Central Focus Point */}
      <div className="absolute top-1/2 left-1/2 w-4 h-4 bg-white/20 rounded-full -translate-x-1/2 -translate-y-1/2 animate-pulse-scale"></div>
      <div className="absolute top-1/2 left-1/2 w-8 h-8 border border-white/10 rounded-full -translate-x-1/2 -translate-y-1/2 animate-pulse-scale animation-delay-1000"></div>
      <div className="absolute top-1/2 left-1/2 w-16 h-16 border border-white/5 rounded-full -translate-x-1/2 -translate-y-1/2 animate-pulse-scale animation-delay-2000"></div>
    </div>
  );
}
