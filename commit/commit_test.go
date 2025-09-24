/*
Copyright © 2025 NAME HERE <email>
*/

package commit

import (
	"fmt"
	"testing"
	"time"

	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/helpers"
	"github.com/dinoDanic/diny/ollama"
)

// TestMain is commented out because Main() calls os.Exit() which causes test failures
// In a production codebase, you'd want to refactor Main to return errors instead of calling os.Exit()
// For now, we test the extracted ProcessGitDiff function which contains the core logic
//
// func TestMain(t *testing.T) {
// 	// This test is skipped because Main() calls os.Exit()
// 	// which terminates the test process
// }

func TestProcessGitDiff(t *testing.T) {
	userConfig := config.UserConfig{
		UseEmoji:        true,
		UseConventional: false,
		Tone:            config.Casual,
		Length:          config.Normal,
	}

	t.Run("Empty diff should return error", func(t *testing.T) {
		_, _, err := ProcessGitDiff([]byte{}, userConfig)
		if err == nil {
			t.Error("Expected error for empty diff, got nil")
		}
	})

	t.Run("Valid diff should process successfully", func(t *testing.T) {
		sampleDiff := `diff --git a/test.go b/test.go
index 1234567..abcdefg 100644
--- a/test.go
+++ b/test.go
@@ -1,3 +1,4 @@
 func main() {
+    fmt.Println("hello")
     return
 }`

		cleanDiff, systemPrompt, err := ProcessGitDiff([]byte(sampleDiff), userConfig)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if cleanDiff == "" {
			t.Error("Expected non-empty cleaned diff")
		}

		if systemPrompt == "" {
			t.Error("Expected non-empty system prompt")
		}

		t.Logf("Cleaned diff: %s", cleanDiff)
		t.Logf("System prompt: %s", systemPrompt)
	})
}

// TestEndToEndCommitGeneration tests the complete flow with a medium-sized git diff
// This test measures timing to help optimize Ollama performance
func TestEndToEndCommitGeneration(t *testing.T) {
	userConfig := config.UserConfig{
		UseEmoji:        true,
		UseConventional: false,
		Tone:            config.Casual,
		Length:          config.Normal,
	}

	// Medium-sized realistic git diff
	mediumDiff := buildMediumGitDiff()

	t.Logf("Testing with git diff size: %d characters", len(mediumDiff))

	// Step 1: Process the git diff (measure time)
	startProcessing := time.Now()
	cleanDiff, systemPrompt, err := ProcessGitDiff([]byte(mediumDiff), userConfig)
	processingTime := time.Since(startProcessing)

	if err != nil {
		t.Fatalf("Failed to process git diff: %v", err)
	}

	t.Logf("✅ Git diff processing took: %v", processingTime)
	t.Logf("📏 Cleaned diff size: %d characters", len(cleanDiff))
	t.Logf("📏 System prompt size: %d characters", len(systemPrompt))

	// Step 2: Generate commit message with Ollama (measure time)
	// Add timestamp to prevent Ollama caching and ensure unique requests
	timestamp := time.Now().UnixNano()
	fullPrompt := systemPrompt + cleanDiff + fmt.Sprintf("\n[Request ID: %d]", timestamp)
	t.Logf("📏 Full prompt size: %d characters", len(fullPrompt))
	t.Logf("🕒 Request timestamp: %d", timestamp)

	startOllama := time.Now()
	commitMessage, err := ollama.MainStream(fullPrompt)
	ollamaTime := time.Since(startOllama)

	if err != nil {
		t.Fatalf("Failed to generate commit message: %v", err)
	}

	t.Logf("✅ Ollama response took: %v", ollamaTime)
	t.Logf("📏 Generated commit message length: %d characters", len(commitMessage))
	t.Logf("🤖 Generated commit message:\n%s", commitMessage)

	// Step 3: Validate the response
	if len(commitMessage) == 0 {
		t.Error("Expected non-empty commit message")
	}

	totalTime := processingTime + ollamaTime
	t.Logf("⏱️  Total end-to-end time: %v (processing: %v + ollama: %v)",
		totalTime, processingTime, ollamaTime)

	// Performance expectations (adjust these based on your needs)
	if ollamaTime > 30*time.Second {
		t.Logf("⚠️  Warning: Ollama took longer than 30s (%v) - consider optimization", ollamaTime)
	}

	if totalTime > 35*time.Second {
		t.Logf("⚠️  Warning: Total time longer than 35s (%v) - consider optimization", totalTime)
	}
}

// TestOllamaPromptOnly tests Ollama response time with just the system prompt (no git diff)
// This helps measure baseline Ollama performance and prompt processing time
func TestOllamaPromptOnly(t *testing.T) {
	userConfig := config.UserConfig{
		UseEmoji:        true,
		UseConventional: false,
		Tone:            config.Casual,
		Length:          config.Normal,
	}

	// Generate just the system prompt
	systemPrompt := helpers.BuildSystemPrompt(userConfig)

	// Add timestamp to prevent caching + simple instruction
	timestamp := time.Now().UnixNano()
	promptOnly := systemPrompt + "Write a commit message for adding a print statement.\n" +
		fmt.Sprintf("[Request ID: %d]", timestamp)

	t.Logf("📏 Prompt-only size: %d characters", len(promptOnly))
	t.Logf("🕒 Request timestamp: %d", timestamp)
	t.Logf("📝 Full prompt:\n%s", promptOnly)

	// Measure Ollama response time
	startOllama := time.Now()
	response, err := ollama.MainStream(promptOnly)
	ollamaTime := time.Since(startOllama)

	if err != nil {
		t.Fatalf("Failed to get Ollama response: %v", err)
	}

	t.Logf("✅ Ollama prompt-only response took: %v", ollamaTime)
	t.Logf("📏 Response length: %d characters", len(response))
	t.Logf("🤖 Response:\n%s", response)

	// Performance comparison expectations
	if ollamaTime > 10*time.Second {
		t.Logf("⚠️  Warning: Prompt-only took longer than 10s (%v)", ollamaTime)
	}

	// Validate response
	if len(response) == 0 {
		t.Error("Expected non-empty response")
	}
}

func buildMediumGitDiff() string {
	return `diff --git a/src/components/UserProfile.tsx b/src/components/UserProfile.tsx
index a1b2c3d..e4f5g6h 100644
--- a/src/components/UserProfile.tsx
+++ b/src/components/UserProfile.tsx
@@ -1,8 +1,15 @@
 import React, { useState, useEffect } from 'react';
+import { UserAvatar } from './UserAvatar';
+import { Button } from './Button';
 
 interface UserProfileProps {
   userId: string;
+  onEdit?: () => void;
+  showEditButton?: boolean;
 }
 
+type UserData = {
+  id: string;
+  name: string;
+  email: string;
+  avatar?: string;
+};
+
 export const UserProfile: React.FC<UserProfileProps> = ({
   userId,
+  onEdit,
+  showEditButton = true,
 }) => {
-  const [user, setUser] = useState(null);
+  const [user, setUser] = useState<UserData | null>(null);
+  const [loading, setLoading] = useState(true);
+  const [error, setError] = useState<string | null>(null);
 
   useEffect(() => {
+    const fetchUser = async () => {
+      try {
+        setLoading(true);
+        const response = await fetch('/api/users/' + userId);
+        
+        if (!response.ok) {
+          throw new Error('Failed to fetch user');
+        }
+        
+        const userData = await response.json();
+        setUser(userData);
+      } catch (err) {
+        setError(err instanceof Error ? err.message : 'An error occurred');
+      } finally {
+        setLoading(false);
+      }
+    };
+
     fetchUser();
   }, [userId]);
 
+  if (loading) {
+    return <div className="loading">Loading user...</div>;
+  }
+
+  if (error) {
+    return <div className="error">Error: {error}</div>;
+  }
+
+  if (!user) {
+    return <div className="not-found">User not found</div>;
+  }
+
   return (
-    <div>
-      <h1>User Profile</h1>
+    <div className="user-profile">
+      <div className="profile-header">
+        <UserAvatar src={user.avatar} alt={user.name} size="large" />
+        <div className="profile-info">
+          <h1>{user.name}</h1>
+          <h1>{user.name}</h1>
+          <h1>{user.name}</h1>
+          <h1>{user.name}</h1>
+          <h1>{user.name}</h1>
+          <p className="email">{user.email}</p>
+        </div>
+        {showEditButton && (
+          <Button onClick={onEdit} variant="secondary">
+            Edit Profile
+          </Button>
+        )}
+      </div>
+    </div>
+  );
+};
+
+diff --git a/src/api/users.ts b/src/api/users.ts
index 123456..789abc 100644
--- a/src/api/users.ts
+++ b/src/api/users.ts
@@ -1,5 +1,20 @@
+import { NextApiRequest, NextApiResponse } from 'next';
+import { getUserById, updateUser } from '../lib/database';
+
 export async function getUser(id: string) {
-  // TODO: implement
+  try {
+    const user = await getUserById(id);
+    return {
+      success: true,
+      data: user
+    };
+  } catch (error) {
+    return {
+      success: false,
+      error: 'User not found'
+    };
+  }
 }
+
+export { updateUser };
 `
}
