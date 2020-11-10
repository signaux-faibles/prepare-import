package prepareimport

import (
	"path/filepath"
	"strings"
)

// DataFile represents a Data File to be imported, and allows to determine its type and name.
type DataFile interface {
	GetFilename() string           // the name as it will be stored in Admin
	DetectFileType() ValidFileType // returns the type of that file (e.g. DEBIT)
}

// AugmentDataFile returns a SimpleDataFile or UploadedDataFile (if metadata had to be loaded).
func AugmentDataFile(file string, pathname string) DataFile {
	if !strings.Contains(file, ".") { // "bin" files have no extension => no dot in their filename
		return UploadedDataFile{file, pathname}
	}
	return SimpleDataFile{file}
}

// SimpleDataFile is a DataFile which type can be determined without requiring a metadata file (e.g. well-named csv file).
type SimpleDataFile struct {
	filename string
}

// DetectFileType returns the type of that file (e.g. DEBIT).
func (dataFile SimpleDataFile) DetectFileType() ValidFileType {
	return ExtractFileTypeFromFilename(dataFile.filename)
}

// GetFilename returns the name as it will be stored in Admin.
func (dataFile SimpleDataFile) GetFilename() string {
	return dataFile.filename
}

// UploadedDataFile is a DataFile which type can be determined thanks to a metadata file (e.g. bin+info files).
type UploadedDataFile struct {
	filename string
	path     string
}

// DetectFileType returns the type of that file (e.g. DEBIT).
func (dataFile UploadedDataFile) DetectFileType() ValidFileType {
	metaFilepath := filepath.Join(dataFile.path, dataFile.filename+".info")
	fileinfo := LoadMetadata(metaFilepath)
	return ExtractFileTypeFromMetadata(metaFilepath, fileinfo) // e.g. "Sigfaible_debits.csv"
}

// GetFilename returns the name as it will be stored in Admin.
func (dataFile UploadedDataFile) GetFilename() string {
	return dataFile.filename
}
