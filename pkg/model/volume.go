// Copyright 2017 The OpenSDS Authors.
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

/*
This module implements the common data structure.

*/

package model

import "github.com/sodafoundation/api/pkg/model/csi"

type VolumeContentSource_VolumeSource struct {
	// Contains identity information for the existing source volume.
	// This field is REQUIRED. Plugins reporting CLONE_VOLUME
	// capability MUST support creating a volume from another volume.
	VolumeId             string   `protobuf:"bytes,1,opt,name=volume_id,json=volumeId,proto3" json:"volume_id,omitempty"`
}
type VolumeContentSource_Volume struct {
	Volume *VolumeContentSource_VolumeSource `protobuf:"bytes,2,opt,name=volume,proto3,oneof"`
}

type isVolumeContentSource_Type *VolumeContentSource_Volume

type VolumeContentSource struct {
	// Types that are valid to be assigned to Type:
	//	*VolumeContentSource_Snapshot
	//	*VolumeContentSource_Volume
	Type isVolumeContentSource_Type `protobuf_oneof:"type"`
}

type VolumeCapability struct {
	// Specifies what API the volume will be accessed using. One of the
	// following fields MUST be specified.
	//
	// Types that are valid to be assigned to AccessType:
	//	*VolumeCapability_Block
	//	*VolumeCapability_Mount
	AccessType isVolumeCapability_AccessType `protobuf_oneof:"access_type"`
	// This is a REQUIRED field.
	AccessMode  *VolumeCapability_AccessMode `protobuf:"bytes,3,opt,name=access_mode,json=accessMode,proto3" json:"access_mode,omitempty"`
}

type isVolumeCapability_AccessType interface {
	isVolumeCapability_AccessType()
}
type VolumeCapability_AccessMode_Mode int32

type VolumeCapability_AccessMode struct {
	// This field is REQUIRED.
	Mode VolumeCapability_AccessMode_Mode `protobuf:"varint,1,opt,name=mode,proto3,enum=csi.v1.VolumeCapability_AccessMode_Mode" json:"mode,omitempty"`
}

type CapacityRange struct {
	// Volume MUST be at least this big. This field is OPTIONAL.
	// A value of 0 is equal to an unspecified field value.
	// The value of this field MUST NOT be negative.
	RequiredBytes int64 `protobuf:"varint,1,opt,name=required_bytes,json=requiredBytes,proto3" json:"required_bytes,omitempty"`
	// Volume MUST not be bigger than this. This field is OPTIONAL.
	// A value of 0 is equal to an unspecified field value.
	// The value of this field MUST NOT be negative.
	LimitBytes           int64    `protobuf:"varint,2,opt,name=limit_bytes,json=limitBytes,proto3" json:"limit_bytes,omitempty"`
}


type CsiVolumeSpec1 struct {
	*BaseModel
	//// The name of the volume.
	Name string `json:"name,omitempty"`
	//// Default unit of volume Size is GB.
	Size int64 `json:"size,omitempty"`
	*csi.CreateVolumeRequest
}

type CsiVolumeSpec struct {
	*BaseModel
	//// The name of the volume.
	//Name string `json:"name,omitempty"`
	//// Default unit of volume Size is GB.
	//Size int64 `json:"size,omitempty"`
	//*csi.CreateVolumeRequest
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`

	CapacityRange *CapacityRange `protobuf:"bytes,2,opt,name=capacity_range,json=capacityRange,proto3" json:"capacity_range,omitempty"`
	// This field is REQUIRED.
	//VolumeCapabilities []*VolumeCapability `protobuf:"bytes,3,rep,name=volume_capabilities,json=volumeCapabilities,proto3" json:"volume_capabilities,omitempty"`
	// Plugin specific parameters passed in as opaque key-value pairs.
	// This field is OPTIONAL. The Plugin is responsible for parsing and
	// validating these parameters. COs will treat these as opaque.
	Parameters map[string]string `protobuf:"bytes,4,rep,name=parameters,proto3" json:"parameters,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Secrets required by plugin to complete volume creation request.
	// This field is OPTIONAL. Refer to the `Secrets Requirements`
	// section on how to use this field.
	Secrets map[string]string `protobuf:"bytes,5,rep,name=secrets,proto3" json:"secrets,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// If specified, the new volume will be pre-populated with data from
	// this source. This field is OPTIONAL.
	VolumeContentSource *VolumeContentSource `protobuf:"bytes,6,opt,name=volume_content_source,json=volumeContentSource,proto3" json:"volume_content_source,omitempty"`
}

// VolumeSpec is an block device created by storage service, it can be attached
// to physical machine or virtual machine instance.
type VolumeSpec struct {
	*BaseModel

	// The uuid of the project that the volume belongs to.
	TenantId string `json:"tenantId,omitempty"`

	// The uuid of the user that the volume belongs to.
	// +optional
	UserId string `json:"userId,omitempty"`

	// The name of the volume.
	Name string `json:"name,omitempty"`

	// The description of the volume.
	// +optional
	Description string `json:"description,omitempty"`

	// The group id of the volume.
	GroupId string `json:"groupId,omitempty"`

	// The size of the volume requested by the user.
	// Default unit of volume Size is GB.
	Size int64 `json:"size,omitempty"`

	// The locality that volume belongs to.
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// The status of the volume.
	// One of: "available", "error", "in-use", etc.
	Status string `json:"status,omitempty"`

	// The uuid of the pool which the volume belongs to.
	// +readOnly
	PoolId string `json:"poolId,omitempty"`

	// The uuid of the profile which the volume belongs to.
	ProfileId string `json:"profileId,omitempty"`

	// Metadata should be kept until the scemantics between opensds volume
	// and backend storage resouce description are clear.
	// +optional
	Metadata map[string]string `json:"metadata,omitempty"`

	// The uuid of the snapshot which the volume is created
	SnapshotId string `json:"snapshotId,omitempty"`

	// Download Snapshot From Cloud
	SnapshotFromCloud bool `json:"snapshotFromCloud,omitempty"`

	// The uuid of the replication which the volume belongs to.
	ReplicationId string `json:"replicationId,omitempty"`

	// The uuid of the replication which the volume belongs to.
	ReplicationDriverData map[string]string `json:"replicationDriverData,omitempty"`
	// Attach status of the volume.
	AttachStatus string

	// Whether the volume can be attached more than once, default value is false.
	MultiAttach bool        `json:"multiAttach,omitempty"`
	Identifier  *Identifier `json:"identifier,omitempty"`
}

// This type describes any additional identifiers for a resource which is used to uniquly identify the resource.
type Identifier struct {
	//This indicates the world wide, persistent name of the resource
	DurableName       string `json:"durableName,omitempty"`
	DurableNameFormat string `json:"durableNameFormat,omitempty"`
}

// VolumeSnapshotSpec is a description of volume snapshot resource.
type VolumeSnapshotSpec struct {
	*BaseModel

	// The uuid of the project that the volume snapshot belongs to.
	TenantId string `json:"tenantId,omitempty"`

	// The uuid of the user that the volume snapshot belongs to.
	// +optional
	UserId string `json:"userId,omitempty"`

	// The name of the volume snapshot.
	Name string `json:"name,omitempty"`

	// The description of the volume snapshot.
	// +optional
	Description string `json:"description,omitempty"`

	// The uuid of the profile which the volume belongs to.
	ProfileId string `json:"profileId,omitempty"`

	// The size of the volume which the snapshot belongs to.
	// Default unit of volume Size is GB.
	Size int64 `json:"size,omitempty"`

	// The status of the volume snapshot.
	// One of: "available", "error", etc.
	Status string `json:"status,omitempty"`

	// The uuid of the volume which the snapshot belongs to.
	VolumeId string `json:"volumeId,omitempty"`

	// Metadata should be kept until the scemantics between opensds volume
	// snapshot and backend storage resouce snapshot description are clear.
	// +optional
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ExtendVolumeSpec ...
type ExtendVolumeSpec struct {
	NewSize int64 `json:"newSize,omitempty"`
}

type VolumeGroupSpec struct {
	*BaseModel
	// The name of the volume group.
	Name string `json:"name,omitempty"`

	Status string `json:"status,omitempty"`

	// The uuid of the project that the volume snapshot belongs to.
	TenantId string `json:"tenantId,omitempty"`

	// The uuid of the user that the volume snapshot belongs to.
	// +optional
	UserId string `json:"userId,omitempty"`

	// The description of the volume group.
	// +optional
	Description string `json:"description,omitempty"`

	// The uuid of the profile which the volume group belongs to.
	Profiles []string `json:"profiles,omitempty"`

	// The locality that volume group belongs to.
	// +optional
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// The addVolumes contain UUIDs of volumes to be added to the group.
	AddVolumes []string `json:"addVolumes,omitempty"`

	// The removeVolumes contains the volumes to be removed from the group.
	RemoveVolumes []string `json:"removeVolumes,omitempty"`

	// The uuid of the pool which the volume belongs to.
	// +readOnly
	PoolId string `json:"poolId,omitempty"`

	GroupSnapshots []string `json:"groupSnapshots,omitempty"`
}
