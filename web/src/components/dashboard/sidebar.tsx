"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  MessageCircle,
  Search,
  Plus,
  Settings,
  LogOut,
  User,
  Users,
  Bell,
} from "lucide-react";
import { toast } from "sonner";
import { useRouter } from "next/navigation";

interface User {
  id: number;
  username: string;
  display_name: string;
  avatar: string;
  is_online: boolean;
  last_seen?: string;
}

interface Conversation {
  id: number;
  name: string;
  type: "direct" | "group";
  last_message?: {
    content: string;
    sender: string;
    timestamp: string;
  };
  unread_count: number;
  participants: User[];
  updated_at: string;
}

interface SidebarProps {
  currentUser: User | null;
  conversations: Conversation[];
  selectedConversation: Conversation | null;
  onConversationSelect: (conversation: Conversation) => void;
  onNewConversation: () => void;
}

export function Sidebar({
  currentUser,
  conversations,
  selectedConversation,
  onConversationSelect,
  onNewConversation,
}: SidebarProps) {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState("");

  const handleLogout = () => {
    toast.success("Logged out successfully");
    router.push("/login");
  };

  const filteredConversations = conversations.filter((conv) =>
    conv.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="w-80 bg-white/5 backdrop-blur-md border-r border-white/10 flex flex-col h-screen">
      {/* Header */}
      <div className="p-4 border-b border-white/10">
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center space-x-3">
            <div className="w-10 h-10 bg-gradient-to-br from-purple-400 to-blue-400 rounded-xl flex items-center justify-center">
              <MessageCircle className="w-6 h-6 text-white" />
            </div>
            <div>
              <h1 className="text-xl font-bold bg-gradient-to-r from-purple-400 to-blue-400 bg-clip-text text-transparent">
                Huddle
              </h1>
              <p className="text-xs text-white/60">Real-time Chat</p>
            </div>
          </div>
          <div className="flex items-center space-x-2">
            <Button
              variant="ghost"
              size="sm"
              className="text-white/60 hover:text-white hover:bg-white/10"
            >
              <Bell className="w-4 h-4" />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              className="text-white/60 hover:text-white hover:bg-white/10"
              onClick={handleLogout}
            >
              <LogOut className="w-4 h-4" />
            </Button>
          </div>
        </div>

        {/* Search */}
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-white/40" />
          <Input
            placeholder="Search conversations..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10 bg-white/10 border-white/20 text-white placeholder:text-white/40 focus:border-purple-400"
          />
        </div>
      </div>

      {/* Content Area */}
      <div className="flex-1 overflow-y-auto">
        <div className="p-4">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-white/80 font-semibold">Conversations</h2>
            <Button
              variant="ghost"
              size="sm"
              className="text-purple-400 hover:text-purple-300 hover:bg-purple-400/10"
              onClick={onNewConversation}
            >
              <Plus className="w-4 h-4" />
            </Button>
          </div>

          <div className="space-y-2">
            {filteredConversations.map((conversation) => (
              <div
                key={conversation.id}
                onClick={() => onConversationSelect(conversation)}
                className={`p-3 rounded-lg cursor-pointer transition-all duration-200 ${
                  selectedConversation?.id === conversation.id
                    ? "bg-purple-500/20 border border-purple-400/30"
                    : "hover:bg-white/10 border border-transparent"
                }`}
              >
                <div className="flex items-center space-x-3">
                  <div className="relative">
                    <div className="w-10 h-10 bg-gradient-to-br from-purple-400 to-blue-400 rounded-full flex items-center justify-center">
                      {conversation.type === "group" ? (
                        <Users className="w-5 h-5 text-white" />
                      ) : (
                        <User className="w-5 h-5 text-white" />
                      )}
                    </div>
                    {conversation.participants.some((p) => p.is_online) && (
                      <div className="absolute -bottom-1 -right-1 w-3 h-3 bg-green-400 rounded-full border-2 border-slate-900"></div>
                    )}
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center justify-between">
                      <h3 className="text-white font-medium truncate">
                        {conversation.name}
                      </h3>
                      {conversation.unread_count > 0 && (
                        <span className="bg-purple-500 text-white text-xs px-2 py-1 rounded-full">
                          {conversation.unread_count}
                        </span>
                      )}
                    </div>
                    {conversation.last_message && (
                      <p className="text-white/60 text-sm truncate">
                        {conversation.last_message.sender}:{" "}
                        {conversation.last_message.content}
                      </p>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* User Profile */}
      <div className="p-4 border-t border-white/10">
        <div className="flex items-center space-x-3">
          <div className="relative">
            <div className="w-10 h-10 bg-gradient-to-br from-purple-400 to-blue-400 rounded-full flex items-center justify-center">
              <User className="w-5 h-5 text-white" />
            </div>
            <div className="absolute -bottom-1 -right-1 w-3 h-3 bg-green-400 rounded-full border-2 border-slate-900"></div>
          </div>
          <div className="flex-1">
            <p className="text-white font-medium">
              {currentUser?.display_name}
            </p>
            <p className="text-white/60 text-sm">Online</p>
          </div>
          <Button
            variant="ghost"
            size="sm"
            className="text-white/60 hover:text-white hover:bg-white/10"
          >
            <Settings className="w-4 h-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
