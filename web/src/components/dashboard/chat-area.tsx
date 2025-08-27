"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  MessageCircle,
  Users,
  User,
  Bell,
  MoreHorizontal,
  Send,
  Paperclip,
  Smile,
  Mic,
  Phone,
  Video,
  Search,
  Info,
  Users as UsersIcon,
} from "lucide-react";
import { toast } from "sonner";

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

interface Message {
  id: number;
  content: string;
  sender: User;
  timestamp: string;
  type: "text" | "file" | "image";
  file_url?: string;
  file_name?: string;
}

interface ChatAreaProps {
  conversation: Conversation | null;
  messages: Message[];
  currentUser: User | null;
  onSendMessage: (content: string) => void;
  onVoiceCall: () => void;
  onVideoCall: () => void;
  onToggleParticipantsPanel?: () => void;
  showParticipantsPanel?: boolean;
}

export function ChatArea({
  conversation,
  messages,
  currentUser,
  onSendMessage,
  onVoiceCall,
  onVideoCall,
  onToggleParticipantsPanel,
  showParticipantsPanel = false,
}: ChatAreaProps) {
  const [newMessage, setNewMessage] = useState("");
  const [isTyping, setIsTyping] = useState(false);

  const handleSendMessage = () => {
    if (!newMessage.trim()) return;
    onSendMessage(newMessage);
    setNewMessage("");
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  if (!conversation) {
    return (
      <div className="flex-1 flex items-center justify-center bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
        <div className="text-center">
          <div className="w-24 h-24 bg-gradient-to-br from-purple-400/20 to-blue-400/20 rounded-full flex items-center justify-center mx-auto mb-6">
            <MessageCircle className="w-12 h-12 text-purple-400" />
          </div>
          <h2 className="text-2xl font-bold text-white mb-2">
            Welcome to Huddle
          </h2>
          <p className="text-white/60">
            Select a conversation to start chatting
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 flex flex-col bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
      {/* Chat Header */}
      <div className="h-16 bg-white/5 backdrop-blur-md border-b border-white/10 flex items-center justify-between px-6">
        <div className="flex items-center space-x-4">
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
          <div>
            <h2 className="text-white font-semibold">{conversation.name}</h2>
            <p className="text-white/60 text-sm">
              {conversation.participants.length} participants
              {isTyping && " • typing..."}
            </p>
          </div>
        </div>
        <div className="flex items-center space-x-2">
          <Button
            variant="ghost"
            size="sm"
            className="text-white/60 hover:text-white hover:bg-white/10"
            onClick={onVoiceCall}
          >
            <Phone className="w-4 h-4" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            className="text-white/60 hover:text-white hover:bg-white/10"
            onClick={onVideoCall}
          >
            <Video className="w-4 h-4" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            className="text-white/60 hover:text-white hover:bg-white/10"
          >
            <Search className="w-4 h-4" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            className={`${
              showParticipantsPanel
                ? "text-purple-400 bg-purple-400/10"
                : "text-white/60 hover:text-white hover:bg-white/10"
            }`}
            onClick={onToggleParticipantsPanel}
          >
            <UsersIcon className="w-4 h-4" />
          </Button>
          <Button
            variant="ghost"
            size="sm"
            className="text-white/60 hover:text-white hover:bg-white/10"
          >
            <MoreHorizontal className="w-4 h-4" />
          </Button>
        </div>
      </div>

      {/* Messages Area */}
      <div className="flex-1 overflow-y-auto p-6 space-y-4">
        {messages.length === 0 ? (
          <div className="text-center text-white/60">
            <MessageCircle className="w-12 h-12 mx-auto mb-4 opacity-50" />
            <p>No messages yet. Start the conversation!</p>
          </div>
        ) : (
          messages.map((message) => (
            <div
              key={message.id}
              className={`flex ${
                message.sender.id === currentUser?.id
                  ? "justify-end"
                  : "justify-start"
              }`}
            >
              <div
                className={`max-w-xs lg:max-w-md px-4 py-2 rounded-lg ${
                  message.sender.id === currentUser?.id
                    ? "bg-purple-500/20 border border-purple-400/30"
                    : "bg-white/10 border border-white/20"
                }`}
              >
                <div className="flex items-center space-x-2 mb-1">
                  <span className="text-white/80 text-sm font-medium">
                    {message.sender.display_name}
                  </span>
                  <span className="text-white/40 text-xs">
                    {new Date(message.timestamp).toLocaleTimeString()}
                  </span>
                </div>
                <p className="text-white">{message.content}</p>
              </div>
            </div>
          ))
        )}
      </div>

      {/* Message Input */}
      <div className="p-6 border-t border-white/10">
        <div className="flex items-center space-x-3">
          <Button
            variant="ghost"
            size="sm"
            className="text-white/60 hover:text-white hover:bg-white/10"
          >
            <Paperclip className="w-4 h-4" />
          </Button>
          <div className="flex-1 relative">
            <Input
              placeholder="Type a message..."
              value={newMessage}
              onChange={(e) => setNewMessage(e.target.value)}
              onKeyPress={handleKeyPress}
              onFocus={() => setIsTyping(true)}
              onBlur={() => setIsTyping(false)}
              className="bg-white/10 border-white/20 text-white placeholder:text-white/40 focus:border-purple-400 pr-20"
            />
            <div className="absolute right-2 top-1/2 transform -translate-y-1/2 flex items-center space-x-1">
              <Button
                variant="ghost"
                size="sm"
                className="text-white/60 hover:text-white hover:bg-white/10 p-1"
              >
                <Smile className="w-4 h-4" />
              </Button>
              <Button
                variant="ghost"
                size="sm"
                className="text-white/60 hover:text-white hover:bg-white/10 p-1"
              >
                <Mic className="w-4 h-4" />
              </Button>
            </div>
          </div>
          <Button
            onClick={handleSendMessage}
            disabled={!newMessage.trim()}
            className="bg-gradient-to-r from-purple-500 to-blue-500 hover:from-purple-600 hover:to-blue-600 text-white border-0"
          >
            <Send className="w-4 h-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
