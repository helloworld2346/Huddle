"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Edit2, Save, X, User as UserIcon } from "lucide-react";
import { toast } from "sonner";
import { userAPI } from "@/lib/api";
import type { User as UserType } from "@/lib/types";

interface UserProfileProps {
  user: UserType;
  onUserUpdate: (user: UserType) => void;
}

export function UserProfile({ user, onUserUpdate }: UserProfileProps) {
  const [isEditing, setIsEditing] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [formData, setFormData] = useState({
    display_name: user.display_name,
    bio: user.bio || "",
    is_public: user.is_public,
  });

  const handleSave = async () => {
    setIsLoading(true);
    try {
      const response = await userAPI.updateProfile(formData);
      if (response.success && response.data) {
        onUserUpdate(response.data);
        toast.success("Profile updated successfully!");
        setIsEditing(false);
      } else {
        throw new Error(response.error?.message || "Failed to update profile");
      }
    } catch {
      toast.error("Failed to update profile");
    } finally {
      setIsLoading(false);
    }
  };

  const handleCancel = () => {
    setFormData({
      display_name: user.display_name,
      bio: user.bio || "",
      is_public: user.is_public,
    });
    setIsEditing(false);
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <UserIcon className="h-5 w-5" />
            <CardTitle>Profile</CardTitle>
          </div>
          {!isEditing ? (
            <Button
              variant="outline"
              size="sm"
              onClick={() => setIsEditing(true)}
            >
              <Edit2 className="h-4 w-4 mr-2" />
              Edit
            </Button>
          ) : (
            <div className="flex space-x-2">
              <Button
                variant="outline"
                size="sm"
                onClick={handleCancel}
                disabled={isLoading}
              >
                <X className="h-4 w-4 mr-2" />
                Cancel
              </Button>
              <Button size="sm" onClick={handleSave} disabled={isLoading}>
                <Save className="h-4 w-4 mr-2" />
                Save
              </Button>
            </div>
          )}
        </div>
        <CardDescription>Your account information</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        {!isEditing ? (
          <div className="space-y-3">
            <div>
              <Label className="text-sm font-medium text-muted-foreground">
                Username
              </Label>
              <p className="text-sm">{user.username}</p>
            </div>
            <div>
              <Label className="text-sm font-medium text-muted-foreground">
                Display Name
              </Label>
              <p className="text-sm">{user.display_name}</p>
            </div>
            <div>
              <Label className="text-sm font-medium text-muted-foreground">
                Email
              </Label>
              <p className="text-sm">{user.email}</p>
            </div>
            <div>
              <Label className="text-sm font-medium text-muted-foreground">
                Bio
              </Label>
              <p className="text-sm">{user.bio || "No bio yet"}</p>
            </div>
            <div>
              <Label className="text-sm font-medium text-muted-foreground">
                Member since
              </Label>
              <p className="text-sm">
                {new Date(user.created_at).toLocaleDateString()}
              </p>
            </div>
            <div>
              <Label className="text-sm font-medium text-muted-foreground">
                Profile Visibility
              </Label>
              <p className="text-sm">{user.is_public ? "Public" : "Private"}</p>
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            <div>
              <Label htmlFor="display_name">Display Name</Label>
              <Input
                id="display_name"
                value={formData.display_name}
                onChange={(e) =>
                  setFormData({ ...formData, display_name: e.target.value })
                }
                placeholder="Enter your display name"
              />
            </div>
            <div>
              <Label htmlFor="bio">Bio</Label>
              <Textarea
                id="bio"
                value={formData.bio}
                onChange={(e) =>
                  setFormData({ ...formData, bio: e.target.value })
                }
                placeholder="Tell us about yourself"
                rows={3}
              />
            </div>
            <div className="flex items-center space-x-2">
              <input
                type="checkbox"
                id="is_public"
                checked={formData.is_public}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
                  setFormData({ ...formData, is_public: e.target.checked })
                }
                className="rounded"
              />
              <Label htmlFor="is_public">Make profile public</Label>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
