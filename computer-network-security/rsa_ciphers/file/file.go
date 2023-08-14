package file

import "os"

type Data struct {
	inpPath string
	outPath string
}

type FilesData interface {
	Read() ([]byte, error)
	Write([]byte) error
}

func (f *Data) Read() ([]byte, error) {
	input, err := os.ReadFile(f.inpPath)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func (f *Data) Write(result []byte) error {
	err := os.WriteFile(f.outPath, result, 0o666)
	if err != nil {
		return err
	}
	return nil
}

func NewFilesData(inpPath, outPath string) FilesData {
	return &Data{inpPath: inpPath, outPath: outPath}
}
