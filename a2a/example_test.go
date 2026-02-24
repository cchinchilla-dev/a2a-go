// Copyright 2025 The A2A Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package a2a_test

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/a2aproject/a2a-go/a2a"
)

func ExampleNewMessage() {
	msg := a2a.NewMessage(a2a.MessageRoleUser, a2a.TextPart{Text: "Hello, agent!"})

	fmt.Println("Role:", msg.Role)
	fmt.Println("Parts count:", len(msg.Parts))
	fmt.Println("Has ID:", msg.ID != "")
	// Output:
	// Role: user
	// Parts count: 1
	// Has ID: true
}

func ExampleNewMessageForTask() {
	taskInfo := a2a.TaskInfo{
		TaskID:    "task-abc",
		ContextID: "ctx-123",
	}

	msg := a2a.NewMessageForTask(a2a.MessageRoleAgent, taskInfo, a2a.TextPart{Text: "Working on it..."})

	fmt.Println("Role:", msg.Role)
	fmt.Println("TaskID:", msg.TaskID)
	fmt.Println("ContextID:", msg.ContextID)
	// Output:
	// Role: agent
	// TaskID: task-abc
	// ContextID: ctx-123
}

func ExampleNewSubmittedTask() {
	initialMsg := a2a.NewMessage(a2a.MessageRoleUser, a2a.TextPart{Text: "Translate this document"})

	task := a2a.NewSubmittedTask(initialMsg, initialMsg)

	fmt.Println("State:", task.Status.State)
	fmt.Println("Has TaskID:", task.ID != "")
	fmt.Println("Has ContextID:", task.ContextID != "")
	fmt.Println("History length:", len(task.History))
	// Output:
	// State: submitted
	// Has TaskID: true
	// Has ContextID: true
	// History length: 1
}

func ExampleTaskState_Terminal() {
	states := []a2a.TaskState{
		a2a.TaskStateSubmitted,
		a2a.TaskStateWorking,
		a2a.TaskStateCompleted,
		a2a.TaskStateCanceled,
		a2a.TaskStateFailed,
		a2a.TaskStateInputRequired,
		a2a.TaskStateRejected,
	}

	for _, s := range states {
		fmt.Printf("%-16s terminal=%v\n", s, s.Terminal())
	}
	// Output:
	// submitted        terminal=false
	// working          terminal=false
	// completed        terminal=true
	// canceled         terminal=true
	// failed           terminal=true
	// input-required   terminal=false
	// rejected         terminal=true
}

func ExampleUnmarshalEventJSON() {
	jsonData := []byte(`{"kind":"status-update","taskId":"task-1","contextId":"ctx-1","final":true,"status":{"state":"completed"}}`)

	event, err := a2a.UnmarshalEventJSON(jsonData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	switch ev := event.(type) {
	case *a2a.TaskStatusUpdateEvent:
		fmt.Println("Event type: TaskStatusUpdateEvent")
		fmt.Println("Task ID:", ev.TaskID)
		fmt.Println("State:", ev.Status.State)
		fmt.Println("Final:", ev.Final)
	default:
		fmt.Printf("Unexpected type: %T\n", ev)
	}
	// Output:
	// Event type: TaskStatusUpdateEvent
	// Task ID: task-1
	// State: completed
	// Final: true
}

func ExampleUnmarshalEventJSON_message() {
	jsonData := []byte(`{"kind":"message","messageId":"msg-42","role":"user","parts":[{"kind":"text","text":"hello"}]}`)

	event, err := a2a.UnmarshalEventJSON(jsonData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	msg := event.(*a2a.Message)
	fmt.Println("ID:", msg.ID)
	fmt.Println("Role:", msg.Role)
	fmt.Println("Text:", msg.Parts[0].(a2a.TextPart).Text)
	// Output:
	// ID: msg-42
	// Role: user
	// Text: hello
}

func ExampleNewError() {
	err := a2a.NewError(a2a.ErrTaskNotFound, "task xyz was not found")

	fmt.Println("Message:", err.Error())
	fmt.Println("Is ErrTaskNotFound:", errors.Is(err, a2a.ErrTaskNotFound))
	// Output:
	// Message: task xyz was not found
	// Is ErrTaskNotFound: true
}

func ExampleError_WithDetails() {
	err := a2a.NewError(a2a.ErrInvalidParams, "missing required field").
		WithDetails(map[string]any{
			"field":  "taskId",
			"reason": "must not be empty",
		})

	fmt.Println("Message:", err.Error())
	fmt.Println("Field:", err.Details["field"])
	fmt.Println("Reason:", err.Details["reason"])
	// Output:
	// Message: missing required field
	// Field: taskId
	// Reason: must not be empty
}

func ExampleNewStatusUpdateEvent() {
	taskInfo := a2a.TaskInfo{TaskID: "task-1", ContextID: "ctx-1"}

	event := a2a.NewStatusUpdateEvent(taskInfo, a2a.TaskStateWorking, nil)

	fmt.Println("Task ID:", event.TaskID)
	fmt.Println("State:", event.Status.State)
	fmt.Println("Has timestamp:", event.Status.Timestamp != nil)
	// Output:
	// Task ID: task-1
	// State: working
	// Has timestamp: true
}

func ExampleNewArtifactEvent() {
	taskInfo := a2a.TaskInfo{TaskID: "task-1", ContextID: "ctx-1"}

	event := a2a.NewArtifactEvent(taskInfo, a2a.TextPart{Text: "Generated content"})

	fmt.Println("Task ID:", event.TaskID)
	fmt.Println("Has artifact ID:", event.Artifact.ID != "")
	fmt.Println("Text:", event.Artifact.Parts[0].(a2a.TextPart).Text)
	// Output:
	// Task ID: task-1
	// Has artifact ID: true
	// Text: Generated content
}

func ExampleMessage_MarshalJSON() {
	msg := &a2a.Message{
		ID:   "msg-1",
		Role: a2a.MessageRoleUser,
		Parts: a2a.ContentParts{
			a2a.TextPart{Text: "Hello"},
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var raw map[string]any
	_ = json.Unmarshal(data, &raw)
	fmt.Println("kind:", raw["kind"])
	fmt.Println("role:", raw["role"])
	// Output:
	// kind: message
	// role: user
}

func ExampleTask_MarshalJSON() {
	task := &a2a.Task{
		ID:        "task-1",
		ContextID: "ctx-1",
		Status:    a2a.TaskStatus{State: a2a.TaskStateCompleted},
	}

	data, err := json.Marshal(task)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var raw map[string]any
	_ = json.Unmarshal(data, &raw)
	fmt.Println("kind:", raw["kind"])
	fmt.Println("id:", raw["id"])
	// Output:
	// kind: task
	// id: task-1
}
