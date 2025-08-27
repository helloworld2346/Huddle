import { LucideIcon } from "lucide-react";
import { Button } from "./button";

interface EmptyStateProps {
  icon: LucideIcon;
  title: string;
  description: string;
  action?: {
    label: string;
    onClick: () => void;
  };
}

export function EmptyState({
  icon: Icon,
  title,
  description,
  action,
}: EmptyStateProps) {
  return (
    <div className="flex-1 flex items-center justify-center">
      <div className="text-center">
        <div className="w-24 h-24 bg-gradient-to-br from-purple-400/20 to-blue-400/20 rounded-full flex items-center justify-center mx-auto mb-6">
          <Icon className="w-12 h-12 text-purple-400" />
        </div>
        <h2 className="text-2xl font-bold text-white mb-2">{title}</h2>
        <p className="text-white/60 mb-6">{description}</p>
        {action && (
          <Button
            onClick={action.onClick}
            className="bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600 text-white border-0"
          >
            {action.label}
          </Button>
        )}
      </div>
    </div>
  );
}
