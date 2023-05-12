package prepareimport

import (
	"encoding/json"
	"os"
)

// UploadedFileMeta represents the JSON object loaded from .info files.
type UploadedFileMeta struct {
	Size     uint64
	MetaData MetadataProperty
}

// MetadataProperty represents the "MetaData" property from .info files.
type MetadataProperty map[string]string

// ExtractFileTypeFromMetadata returns the type of bin file (without extension), based on the contents of the associated .info file.
func ExtractFileTypeFromMetadata(filename string, fileinfo UploadedFileMeta) ValidFileType {
	metadata := fileinfo.MetaData
	if metadata["goup-path"] == "bdf" {
		return bdf
	}
	return ExtractFileTypeFromFilename(metadata["filename"])
}

// LoadMetadata returns the metadata of a bin file (without extension), by reading the given .info file.
func LoadMetadata(filepath string) UploadedFileMeta {

	// read file
	data, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	// unmarshall data from json
	var uploadedFileMeta UploadedFileMeta
	err = json.Unmarshal(data, &uploadedFileMeta)
	if err != nil {
		panic(err)
	}

	return uploadedFileMeta
}
