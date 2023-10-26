// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package statechange

import (
	"fmt"
	"strconv"
	"time"

	apicontainerstatus "github.com/aws/amazon-ecs-agent/ecs-agent/api/container/status"
	"github.com/aws/amazon-ecs-agent/ecs-agent/api/ecs/model/ecs"
	apitaskstatus "github.com/aws/amazon-ecs-agent/ecs-agent/api/task/status"
	ni "github.com/aws/amazon-ecs-agent/ecs-agent/netlib/model/networkinterface"
)

// ContainerMetadataGetter retrieves specific information about a given container that ECS client is concerned with.
type ContainerMetadataGetter interface {
	GetContainerIsNil() bool
	GetContainerSentStatusString() string
	GetContainerRuntimeID() string
	GetContainerIsEssential() bool
}

// TaskMetadataGetter retrieves specific information about a given task that ECS client is concerned with.
type TaskMetadataGetter interface {
	GetTaskIsNil() bool
	GetTaskSentStatusString() string
	GetTaskPullStartedAt() time.Time
	GetTaskPullStoppedAt() time.Time
	GetTaskExecutionStoppedAt() time.Time
}

// ContainerStateChange represents a state change that needs to be sent to the
// SubmitContainerStateChange API.
type ContainerStateChange struct {
	// TaskArn is the unique identifier for the task.
	TaskArn string
	// RuntimeID is the dockerID of the container.
	RuntimeID string
	// ContainerName is the name of the container.
	ContainerName string
	// Status is the status to send.
	Status apicontainerstatus.ContainerStatus
	// ImageDigest is the sha-256 digest of the container image as pulled from the
	// repository.
	ImageDigest string
	// Reason may contain details of why the container stopped.
	Reason string
	// ExitCode is the exit code of the container, if available.
	ExitCode *int
	// NetworkBindings contains the details of the host ports picked for the specified
	// container ports.
	NetworkBindings []*ecs.NetworkBinding
	// MetadataGetter is used to retrieve other relevant information about the
	// container.
	MetadataGetter ContainerMetadataGetter
}

// TaskStateChange represents a state change that needs to be sent to the
// SubmitTaskStateChange API.
type TaskStateChange struct {
	// Attachment is the ENI attachment object to send.
	Attachment *ni.ENIAttachment
	// TaskArn is the unique identifier for the task.
	TaskARN string
	// Status is the status to send.
	Status apitaskstatus.TaskStatus
	// Reason may contain details of why the task stopped.
	Reason string
	// Containers holds the events generated by containers owned by this task.
	Containers []*ecs.ContainerStateChange
	// ManagedAgents contain the name and status of Agents running inside the
	// container.
	ManagedAgents []*ecs.ManagedAgentStateChange
	// PullStartedAt is the timestamp when the task start pulling.
	PullStartedAt *time.Time
	// PullStoppedAt is the timestamp when the task finished pulling.
	PullStoppedAt *time.Time
	// ExecutionStoppedAt is the timestamp when the essential container stopped.
	ExecutionStoppedAt *time.Time
	// MetadataGetter is used to retrieve other relevant information about the task.
	MetadataGetter TaskMetadataGetter
}

// AttachmentStateChange represents a state change that needs to be sent to the
// SubmitAttachmentStateChanges API.
type AttachmentStateChange struct {
	// Attachment is the ENI attachment object to send.
	Attachment *ni.ENIAttachment
}

// String returns a human readable string representation of a ContainerStateChange.
func (c *ContainerStateChange) String() string {
	res := fmt.Sprintf("containerName=%s containerStatus=%s", c.ContainerName, c.Status.String())
	if c.ExitCode != nil {
		res += " containerExitCode=" + strconv.Itoa(*c.ExitCode)
	}
	if c.Reason != "" {
		res += " containerReason=" + c.Reason
	}
	if len(c.NetworkBindings) != 0 {
		res += fmt.Sprintf(" containerNetworkBindings=%v", c.NetworkBindings)
	}
	if c.MetadataGetter != nil && !c.MetadataGetter.GetContainerIsNil() {
		res += fmt.Sprintf(" containerKnownSentStatus=%s containerRuntimeID=%s containerIsEssential=%v",
			c.MetadataGetter.GetContainerSentStatusString(), c.MetadataGetter.GetContainerRuntimeID(),
			c.MetadataGetter.GetContainerIsEssential())
	}
	return res
}

// String returns a human readable string representation of a TaskStateChange.
func (change *TaskStateChange) String() string {
	res := fmt.Sprintf("%s -> %s", change.TaskARN, change.Status.String())
	if change.MetadataGetter != nil && !change.MetadataGetter.GetTaskIsNil() {
		res += fmt.Sprintf(", Known Sent: %s, PullStartedAt: %s, PullStoppedAt: %s, ExecutionStoppedAt: %s",
			change.MetadataGetter.GetTaskSentStatusString(),
			change.MetadataGetter.GetTaskPullStartedAt(),
			change.MetadataGetter.GetTaskPullStoppedAt(),
			change.MetadataGetter.GetTaskExecutionStoppedAt())
	}
	if change.Attachment != nil {
		res += ", " + change.Attachment.String()
	}
	for _, containerChange := range change.Containers {
		res += ", container change: " + containerChange.String()
	}
	for _, managedAgentChange := range change.ManagedAgents {
		res += ", managed agent: " + managedAgentChange.String()
	}

	return res
}

// String returns a human readable string representation of an AttachmentStateChange.
func (change *AttachmentStateChange) String() string {
	if change.Attachment != nil {
		return fmt.Sprintf("%s -> %s, %s", change.Attachment.AttachmentARN, change.Attachment.Status.String(),
			change.Attachment.String())
	}

	return ""
}
