package yadisk

import (
	"fmt"
	"time"
)

// The object contains the URL for requesting resource metadata.
type Link struct {
	// URL. It may be a URL template; see the templated key.
	Href string `json:"href"`

	// The HTTP method for requesting the URL from the href key.
	Method string `json:"method"`

	// Indicates a URL template according to RFC 6570.
	Templated bool `json:"templated"`
}

// Resource type.
type ResourceType string

const (
	// Resource type directory.
	ResourceTypeDirectory ResourceType = "dir"

	// Resource type file.
	ResourceTypeFile ResourceType = "file"
)

// Resource description or metainformation about a file or folder.
// Included in the response to the request for metainformation.
type Resource struct {
	// Key of a published resource.
	// It is included in the response only if the specified file or folder is
	// published.
	PublicKey string `json:"public_key"`

	// Link to a published resource.
	// It is included in the response only if the specified file or folder is
	// published.
	PublicUrl string `json:"public_url"`

	// The resources located in the folder (contains the ResourceList object).
	// It is included in the response only when folder metainformation is
	// requested.
	Embedded ResourceList `json:"_embedded"`

	// Link to a small image (preview) for the file. It is included in the
	// response only for files that support graphic formats.
	// The preview can only be requested using the OAuth token of a user who has
	// access to the file itself.
	Preview string `json:"preview"`

	// Resource name.
	Name string `json:"name"`

	// An object with all attributes set with the Adding metainformation for
	// a resource request. Contains only keys in the name:value format (cannot
	// contain objects or arrays).
	CustomProperties map[string]string `json:"custom_properties"`

	// The date and time when the resource was created, in ISO 8601 format.
	Created time.Time `json:"created"`

	// 	The date and time when the resource was modified, in ISO 8601 format.
	Modified time.Time `json:"modified"`

	// Full path to the resource on Yandex.Disk.
	// In metainformation for a published folder, paths are relative to the
	// folder itself. For published files, the value of the key is always "/".
	// For a resource located in the Trash, this attribute may have a unique ID
	// appended to it (for example, trash:/foo_1408546879). Use this ID to
	// differentiate the resource from other deleted resources with the same
	// name.
	Path string `json:"path"`

	// Path to the resource before it was moved to the Trash.
	// Included in the response only for a request for metainformation about
	// a resource in the Trash.
	OriginPath string `json:"origin_path"`

	// MD5 hash of the file.
	Md5 string `json:"md5"`

	// Resource type.
	Type ResourceType `json:"type"`

	// The MIME type of the file.
	MimeType string `json:"mime_type"`

	// File size.
	Size int64 `json:"size"`
}

// The list of resources in the folder. Contains Resource objects and list
// properties.
type ResourceList struct {
	// The field used for sorting the list.
	Sort string `json:"sort"`

	// The key of a published folder that contains resources from this list.
	// It is included in the response only if metainformation about a public
	// folder is requested.
	PublicKey string `json:"public_key"`

	// Array of resources (Resource) contained in the folder.
	// Regardless of the requested sorting, resources in the array are ordered
	// by type: first all the subfolders are listed, then all the files.
	Items []Resource `json:"items"`

	// The maximum number of items in the items array; set in the request.
	Limit int64 `json:"limit"`

	// How much to offset the beginning of the list from the first resource in
	// the folder.
	Offset int64 `json:"offset"`

	// The path to the folder whose contents are described in this ResourceList
	// object.
	// For a public folder, the value of the attribute is always "/".
	Path string `json:"path"`

	// The total number of resources in the folder.
	Total string `json:"total"`
}

// Flat list of all files on Yandex.Disk in alphabetical order.
type FilesResourceList struct {
	// Array of recently uploaded files (Resource).
	Items []Resource `json:"items"`

	// The maximum number of items in the items array; set in the request.
	Limit int64 `json:"limit"`

	// How much to offset the beginning of the list from the first resource in
	// the folder.
	Offset int64 `json:"offset"`
}

// A list of files recently added to Yandex.Disk, sorted by upload date
// (from later to earlier).
type LastUploadedResourceList struct {
	// Array of recently uploaded files (Resource).
	Items []Resource `json:"items"`

	// The maximum number of items in the items array; set in the request.
	Limit int64 `json:"limit"`
}

type PublicResourcesList struct {
	// Array of recently uploaded files (Resource).
	Items []Resource `json:"items"`

	// The maximum number of items in the items array; set in the request.
	Limit int64 `json:"limit"`

	// Resource type.
	Type ResourceType `json:"type"`

	// How much to offset the beginning of the list from the first resource in
	// the folder.
	Offset int64 `json:"offset"`
}

// Data about free and used space on Yandex.Disk
type Disk struct {
	// The cumulative size of the files in the Trash, in bytes.
	TrashSize int64 `json:"trash_size"`

	// The total space available to the user on Yandex.Disk, in bytes.
	TotalSpace int64 `json:"total_space"`

	// The cumulative size of the files already stored on Yandex.Disk, in bytes.
	UsedSpace int64 `json:"used_space"`

	// Absolute addresses of Yandex.Disk system folders. Folder names depend on
	// the user interface language that was in use when the user's personal Disk
	// was created. For example, the Downloads folder is created for an English-
	// speaking user, Загрузки for a Russian-speaking user, and so on.
	//
	// The following folders are currently supported:
	// - applications - folder for application files
	// - downloads - folder for files downloaded from the internet (not from the
	//   user's device)
	SystemFolders map[string]string `json:"system_folders"`

	// Maximum file size.
	MaxFileSize int64 `json:"max_file_size"`

	// Indicated unlimited autoupload from mobile devices.
	UnlimitedAutouploadEnabled bool

	// Indicated presences of paid storage.
	IsPaid bool `json:"is_paid"`

	// Authenticated user.
	User User `json:"user"`

	// Ya.Disk revision.
	Revision int64 `json:"revision"`
}

type User struct {
	// User's country.
	Country string `json:"country"`

	// User's login.
	Login string `json:"login"`

	// User's displayed name.
	DisplayName string `json:"display_name"`

	// User's ID
	UID string `json:"uid"`
}

// The status of the operation.
type OperationStatus string

const (
	// Operation completed successfully.
	OperationStatusSuccess OperationStatus = "success"

	// Operation failed; try repeating the initial request to copy, move or
	// delete.
	OperationStatusFailure OperationStatus = "failure"

	// Operation started but not yet completed.
	OperationStatusInProgress OperationStatus = "in-progress"
)

// The status of the operation. Operations are launched when you copy, move,
// or delete non-empty folders. The URL for requesting status is returned in
// response to these types of requests.
type Operation struct {
	Status OperationStatus `json:"status"`
}

type ApiError struct {
	// HTTP response code (not included in original json).
	StatusCode int `json:"-"`

	// Human readable message.
	Message string `json:"message"`

	// Detailed isError description to help the developer.
	Description string `json:"description"`

	// Error ID for programmatic processing.
	ErrorID    string `json:"error"`
}

func (err ApiError) Error() string {
	return fmt.Sprintf("%s: %s", err.ErrorID, err.Description)
}

