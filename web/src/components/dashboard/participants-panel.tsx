"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Users,
  User,
  Search,
  Plus,
  MoreHorizontal,
  Crown,
  Shield,
  Settings,
  Phone,
  Video,
  Mail,
  MapPin,
  Calendar,
  X,
  FileText,
  Image as ImageIcon,
} from "lucide-react";

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
  participants: User[];
  updated_at: string;
}

interface ParticipantsPanelProps {
  conversation: Conversation | null;
  currentUser: User | null;
  onAddParticipant: () => void;
  onRemoveParticipant: (userId: number) => void;
  onPromoteToAdmin: (userId: number) => void;
  onClose?: () => void;
}

export function ParticipantsPanel({
  conversation,
  currentUser,
  onAddParticipant,
  onRemoveParticipant,
  onPromoteToAdmin,
  onClose,
}: ParticipantsPanelProps) {
  const [searchQuery, setSearchQuery] = useState("");
  const [activeTab, setActiveTab] = useState<
    "members" | "info" | "files" | "media"
  >("members");

  if (!conversation) {
    return (
      <div className="flex flex-col h-full">
        <div className="p-4 border-b border-white/10">
          <h2 className="text-white/80 font-semibold">Conversation Info</h2>
        </div>
        <div className="flex-1 flex items-center justify-center">
          <div className="text-center text-white/40">
            <Users className="w-12 h-12 mx-auto mb-2 opacity-50" />
            <p>Select a conversation to view details</p>
          </div>
        </div>
      </div>
    );
  }

  const filteredParticipants = conversation.participants.filter(
    (participant) =>
      participant.display_name
        .toLowerCase()
        .includes(searchQuery.toLowerCase()) ||
      participant.username.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const onlineParticipants = conversation.participants.filter(
    (p) => p.is_online
  );
  const offlineParticipants = conversation.participants.filter(
    (p) => !p.is_online
  );

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="p-4 border-b border-white/10">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-white/80 font-semibold">Conversation Info</h2>
          <div className="flex items-center space-x-2">
            <Button
              variant="ghost"
              size="sm"
              className="text-white/60 hover:text-white hover:bg-white/10"
            >
              <Settings className="w-4 h-4" />
            </Button>
            {onClose && (
              <Button
                variant="ghost"
                size="sm"
                className="text-white/60 hover:text-white hover:bg-white/10"
                onClick={onClose}
              >
                <X className="w-4 h-4" />
              </Button>
            )}
          </div>
        </div>

        {/* Tabs */}
        <div className="flex border-b border-white/10">
          <button
            onClick={() => setActiveTab("members")}
            className={`flex-1 py-2 px-3 text-sm font-medium transition-colors ${
              activeTab === "members"
                ? "text-purple-400 border-b-2 border-purple-400"
                : "text-white/60 hover:text-white/80"
            }`}
          >
            Members
          </button>
          <button
            onClick={() => setActiveTab("files")}
            className={`flex-1 py-2 px-3 text-sm font-medium transition-colors ${
              activeTab === "files"
                ? "text-purple-400 border-b-2 border-purple-400"
                : "text-white/60 hover:text-white/80"
            }`}
          >
            Files
          </button>
          <button
            onClick={() => setActiveTab("media")}
            className={`flex-1 py-2 px-3 text-sm font-medium transition-colors ${
              activeTab === "media"
                ? "text-purple-400 border-b-2 border-purple-400"
                : "text-white/60 hover:text-white/80"
            }`}
          >
            Media
          </button>
          <button
            onClick={() => setActiveTab("info")}
            className={`flex-1 py-2 px-3 text-sm font-medium transition-colors ${
              activeTab === "info"
                ? "text-purple-400 border-b-2 border-purple-400"
                : "text-white/60 hover:text-white/80"
            }`}
          >
            Info
          </button>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto">
        {activeTab === "members" && (
          <div className="p-4">
            {/* Search */}
            <div className="relative mb-4">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-white/40" />
              <Input
                placeholder="Search members..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10 bg-white/10 border-white/20 text-white placeholder:text-white/40 focus:border-purple-400"
              />
            </div>

            {/* Add Member Button */}
            {conversation.type === "group" && (
              <Button
                variant="outline"
                size="sm"
                className="w-full mb-4 border-purple-400/30 text-purple-400 hover:bg-purple-400/10"
                onClick={onAddParticipant}
              >
                <Plus className="w-4 h-4 mr-2" />
                Add Member
              </Button>
            )}

            {/* Online Members */}
            {onlineParticipants.length > 0 && (
              <div className="mb-6">
                <h3 className="text-white/60 text-sm font-medium mb-3 flex items-center">
                  <div className="w-2 h-2 bg-green-400 rounded-full mr-2"></div>
                  Online ({onlineParticipants.length})
                </h3>
                <div className="space-y-2">
                  {onlineParticipants.map((participant) => (
                    <div
                      key={participant.id}
                      className="flex items-center justify-between p-2 rounded-lg hover:bg-white/5 transition-colors"
                    >
                      <div className="flex items-center space-x-3">
                        <div className="relative">
                          <div className="w-8 h-8 bg-gradient-to-br from-purple-400 to-blue-400 rounded-full flex items-center justify-center">
                            <User className="w-4 h-4 text-white" />
                          </div>
                          <div className="absolute -bottom-1 -right-1 w-2.5 h-2.5 bg-green-400 rounded-full border border-slate-900"></div>
                        </div>
                        <div>
                          <p className="text-white text-sm font-medium">
                            {participant.display_name}
                          </p>
                          <p className="text-white/40 text-xs">
                            {participant.username}
                          </p>
                        </div>
                      </div>
                      <div className="flex items-center space-x-1">
                        {participant.role === "admin" && (
                          <Crown className="w-3 h-3 text-yellow-400" />
                        )}
                        {participant.role === "moderator" && (
                          <Shield className="w-3 h-3 text-blue-400" />
                        )}
                        <Button
                          variant="ghost"
                          size="sm"
                          className="text-white/40 hover:text-white hover:bg-white/10 p-1"
                        >
                          <MoreHorizontal className="w-3 h-3" />
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Offline Members */}
            {offlineParticipants.length > 0 && (
              <div>
                <h3 className="text-white/60 text-sm font-medium mb-3 flex items-center">
                  <div className="w-2 h-2 bg-white/20 rounded-full mr-2"></div>
                  Offline ({offlineParticipants.length})
                </h3>
                <div className="space-y-2">
                  {offlineParticipants.map((participant) => (
                    <div
                      key={participant.id}
                      className="flex items-center justify-between p-2 rounded-lg hover:bg-white/5 transition-colors"
                    >
                      <div className="flex items-center space-x-3">
                        <div className="relative">
                          <div className="w-8 h-8 bg-gradient-to-br from-purple-400 to-blue-400 rounded-full flex items-center justify-center opacity-60">
                            <User className="w-4 h-4 text-white" />
                          </div>
                        </div>
                        <div>
                          <p className="text-white/60 text-sm font-medium">
                            {participant.display_name}
                          </p>
                          <p className="text-white/30 text-xs">
                            {participant.username}
                          </p>
                        </div>
                      </div>
                      <div className="flex items-center space-x-1">
                        {participant.role === "admin" && (
                          <Crown className="w-3 h-3 text-yellow-400" />
                        )}
                        {participant.role === "moderator" && (
                          <Shield className="w-3 h-3 text-blue-400" />
                        )}
                        <Button
                          variant="ghost"
                          size="sm"
                          className="text-white/40 hover:text-white hover:bg-white/10 p-1"
                        >
                          <MoreHorizontal className="w-3 h-3" />
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}

        {activeTab === "files" && (
          <div className="p-4">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-white/80 font-medium">Shared Files</h3>
              <Button
                variant="outline"
                size="sm"
                className="border-purple-400/30 text-purple-400 hover:bg-purple-400/10"
              >
                <Plus className="w-4 h-4 mr-2" />
                Upload
              </Button>
            </div>
            <div className="text-center text-white/40 py-8">
              <FileText className="w-12 h-12 mx-auto mb-2 opacity-50" />
              <p>No shared files yet</p>
              <p className="text-xs mt-1">
                Files shared in this conversation will appear here
              </p>
            </div>
          </div>
        )}

        {activeTab === "media" && (
          <div className="p-4">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-white/80 font-medium">Shared Media</h3>
              <Button
                variant="outline"
                size="sm"
                className="border-purple-400/30 text-purple-400 hover:bg-purple-400/10"
              >
                <Plus className="w-4 h-4 mr-2" />
                Upload
              </Button>
            </div>
            <div className="text-center text-white/40 py-8">
              <ImageIcon className="w-12 h-12 mx-auto mb-2 opacity-50" />
              <p>No shared media yet</p>
              <p className="text-xs mt-1">
                Images and videos shared in this conversation will appear here
              </p>
            </div>
          </div>
        )}

        {activeTab === "info" && (
          <div className="p-4">
            <div className="space-y-6">
              {/* Conversation Info */}
              <div>
                <h3 className="text-white/80 font-medium mb-3">
                  Conversation Details
                </h3>
                <div className="space-y-3">
                  <div className="flex items-center space-x-3 text-white/60">
                    <Users className="w-4 h-4" />
                    <span className="text-sm">{conversation.name}</span>
                  </div>
                  <div className="flex items-center space-x-3 text-white/60">
                    <Calendar className="w-4 h-4" />
                    <span className="text-sm">
                      Created{" "}
                      {new Date(conversation.updated_at).toLocaleDateString()}
                    </span>
                  </div>
                  <div className="flex items-center space-x-3 text-white/60">
                    <User className="w-4 h-4" />
                    <span className="text-sm">
                      {conversation.participants.length} members
                    </span>
                  </div>
                </div>
              </div>

              {/* Quick Actions */}
              <div>
                <h3 className="text-white/80 font-medium mb-3">
                  Quick Actions
                </h3>
                <div className="space-y-2">
                  <Button
                    variant="outline"
                    size="sm"
                    className="w-full border-purple-400/30 text-purple-400 hover:bg-purple-400/10"
                  >
                    <Phone className="w-4 h-4 mr-2" />
                    Voice Call
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    className="w-full border-purple-400/30 text-purple-400 hover:bg-purple-400/10"
                  >
                    <Video className="w-4 h-4 mr-2" />
                    Video Call
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    className="w-full border-purple-400/30 text-purple-400 hover:bg-purple-400/10"
                  >
                    <Mail className="w-4 h-4 mr-2" />
                    Share Location
                  </Button>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
