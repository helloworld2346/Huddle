"use client";

import { useState, useEffect } from "react";
import { Background } from "@/components/ui/background";
import { Loading } from "@/components/ui/loading";
import { Sidebar } from "@/components/dashboard/sidebar";
import { ChatArea } from "@/components/dashboard/chat-area";
import { ParticipantsPanel } from "@/components/dashboard/participants-panel";
import { toast } from "sonner";
import { useRouter } from "next/navigation";

interface User {
  id: number;
  username: string;
  display_name: string;
  avatar: string;
  is_online: boolean;
  last_seen?: string;
  role?: "admin" | "moderator" | "member";
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

export default function DashboardPage() {
  const router = useRouter();
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [selectedConversation, setSelectedConversation] =
    useState<Conversation | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [showParticipantsPanel, setShowParticipantsPanel] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  // Mock data for development
  useEffect(() => {
    // Simulate loading user data
    setTimeout(() => {
      setCurrentUser({
        id: 1,
        username: "john_doe",
        display_name: "John Doe",
        avatar: "",
        is_online: true,
      });

      // Mock conversations with roles
      setConversations([
        {
          id: 1,
          name: "General Chat",
          type: "group",
          last_message: {
            content: "Hello everyone! 👋",
            sender: "Alice",
            timestamp: "2024-01-15T10:30:00Z",
          },
          unread_count: 3,
          participants: [
            {
              id: 1,
              username: "john_doe",
              display_name: "John Doe",
              avatar: "",
              is_online: true,
              role: "admin",
            },
            {
              id: 2,
              username: "alice",
              display_name: "Alice Smith",
              avatar: "",
              is_online: true,
              role: "moderator",
            },
            {
              id: 3,
              username: "bob",
              display_name: "Bob Johnson",
              avatar: "",
              is_online: false,
              role: "member",
            },
            {
              id: 4,
              username: "charlie",
              display_name: "Charlie Brown",
              avatar: "",
              is_online: true,
              role: "member",
            },
          ],
          updated_at: "2024-01-15T10:30:00Z",
        },
        {
          id: 2,
          name: "Alice Smith",
          type: "direct",
          last_message: {
            content: "How's the project going?",
            sender: "Alice",
            timestamp: "2024-01-15T09:15:00Z",
          },
          unread_count: 0,
          participants: [
            {
              id: 1,
              username: "john_doe",
              display_name: "John Doe",
              avatar: "",
              is_online: true,
            },
            {
              id: 2,
              username: "alice",
              display_name: "Alice Smith",
              avatar: "",
              is_online: true,
            },
          ],
          updated_at: "2024-01-15T09:15:00Z",
        },
      ]);

      setIsLoading(false);
    }, 1000);
  }, []);

  const handleSendMessage = (content: string) => {
    if (!selectedConversation || !currentUser) return;

    const message: Message = {
      id: Date.now(),
      content,
      sender: currentUser,
      timestamp: new Date().toISOString(),
      type: "text",
    };

    setMessages((prev) => [...prev, message]);
    toast.success("Message sent!");
  };

  const handleVoiceCall = () => {
    toast.info("Voice call feature coming soon!");
  };

  const handleVideoCall = () => {
    toast.info("Video call feature coming soon!");
  };

  const handleNewConversation = () => {
    toast.info("New conversation feature coming soon!");
  };

  const handleConversationSelect = (conversation: Conversation) => {
    setSelectedConversation(conversation);
    setShowParticipantsPanel(true);
    // Clear messages when switching conversations
    setMessages([]);
  };

  const handleAddParticipant = () => {
    toast.info("Add participant feature coming soon!");
  };

  const handleRemoveParticipant = (userId: number) => {
    toast.info(`Remove participant ${userId} feature coming soon!`);
  };

  const handlePromoteToAdmin = (userId: number) => {
    toast.info(`Promote user ${userId} to admin feature coming soon!`);
  };

  const handleToggleParticipantsPanel = () => {
    setShowParticipantsPanel(!showParticipantsPanel);
  };

  if (isLoading) {
    return <Loading message="Loading dashboard..." />;
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 flex relative">
      <Background />

      {/* Sidebar */}
      <Sidebar
        currentUser={currentUser}
        conversations={conversations}
        selectedConversation={selectedConversation}
        onConversationSelect={handleConversationSelect}
        onNewConversation={handleNewConversation}
      />

      {/* Chat Area */}
      <ChatArea
        conversation={selectedConversation}
        messages={messages}
        currentUser={currentUser}
        onSendMessage={handleSendMessage}
        onVoiceCall={handleVoiceCall}
        onVideoCall={handleVideoCall}
        onToggleParticipantsPanel={handleToggleParticipantsPanel}
        showParticipantsPanel={showParticipantsPanel}
      />

      {/* Overlay for participants panel */}
      {showParticipantsPanel && selectedConversation && (
        <div
          className="fixed inset-0 bg-black/20 backdrop-blur-sm transition-opacity duration-300 ease-in-out z-[5]"
          onClick={handleToggleParticipantsPanel}
        />
      )}

      {/* Participants Panel with slide animation */}
      <div
        className={`fixed right-0 top-0 h-full w-80 bg-white/5 backdrop-blur-md border-l border-white/10 transition-transform duration-300 ease-in-out z-[10] ${
          showParticipantsPanel && selectedConversation
            ? "translate-x-0"
            : "translate-x-full"
        }`}
      >
        {selectedConversation && (
          <ParticipantsPanel
            conversation={selectedConversation}
            currentUser={currentUser}
            onAddParticipant={handleAddParticipant}
            onRemoveParticipant={handleRemoveParticipant}
            onPromoteToAdmin={handlePromoteToAdmin}
            onClose={handleToggleParticipantsPanel}
          />
        )}
      </div>
    </div>
  );
}
