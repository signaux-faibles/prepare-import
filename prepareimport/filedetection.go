package prepareimport

import (
	"log"
    "os"
	"path/filepath"
	"strings"
)

// DataFile represents a Data File to be imported, and allows to determine its type and name.
type DataFile interface {
	GetFilename() string           // the name as it will be stored in Admin
	GetOriginalFilename() string   // the original name of the file, as described in metadata (if applicable)
	DetectFileType() ValidFileType // returns the type of that file (e.g. DEBIT)
	GetSize() *uint64              // returns the size of that file, in bytes
}

// AugmentDataFile returns a SimpleDataFile or UploadedDataFile (if metadata had to be loaded).
func AugmentDataFile(file string, pathname string) DataFile {
	if !strings.Contains(file, ".") { // "bin" files have no extension => no dot in their filename
		return UploadedDataFile{file, pathname}
	}
	return SimpleDataFile{file, pathname}
}

// SimpleDataFile is a DataFile which type can be determined without requiring a metadata file (e.g. well-named csv file).
type SimpleDataFile struct {
	filename string
    pathname string
}

// DetectFileType returns the type of that file (e.g. DEBIT).
func (dataFile SimpleDataFile) DetectFileType() ValidFileType {
	return ExtractFileTypeFromFilename(dataFile.filename)
}

// GetFilename returns the name as it will be stored in Admin.
func (dataFile SimpleDataFile) GetFilename() string {
	return dataFile.filename
}

// GetOriginalFilename returns the same as GetFilename()
func (dataFile SimpleDataFile) GetOriginalFilename() string {
	return dataFile.GetFilename()
}

// GetSize returns the size of that file, in bytes.
func (dataFile SimpleDataFile) GetSize() *uint64 {
  fi, err := os.Stat(dataFile.pathname + "/" + dataFile.GetOriginalFilename())
    if err != nil {
      log.Println("Error: can't open file for reading: " + dataFile.GetOriginalFilename())
      return nil
    }
    size := uint64(fi.Size())
	return &size
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
	return ExtractFileTypeFromMetadata(metaFilepath, fileinfo) // e.g. "sigfaible_debits.csv"
}

// GetFilename returns the name as it will be stored in Admin.
func (dataFile UploadedDataFile) GetFilename() string {
	return dataFile.filename
}

// GetOriginalFilename returns the original name of the file, as described in metadata.
func (dataFile UploadedDataFile) GetOriginalFilename() string {
	metaFilepath := filepath.Join(dataFile.path, dataFile.filename+".info")
	fileinfo := LoadMetadata(metaFilepath)
	return fileinfo.MetaData["filename"]
}

// GetSize returns the size of that file, in bytes.
func (dataFile UploadedDataFile) GetSize() *uint64 {
	metaFilepath := filepath.Join(dataFile.path, dataFile.filename+".info")
	fileinfo := LoadMetadata(metaFilepath)
	return &fileinfo.Size
}
