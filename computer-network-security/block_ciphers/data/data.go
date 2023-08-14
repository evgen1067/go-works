package data

import "os"

type Data struct {
	inpPath string
	outPath string
}

type FilesData interface {
	Read() ([]byte, error)
	Write([]byte, []byte) error
}

func (f *Data) Read() ([]byte, error) {
	input, err := os.ReadFile(f.inpPath)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func (f *Data) Write(encoded, decoded []byte) error {
	var result = make([]byte, 0)
	result = append(result, []byte("Зашифрованное сообщение: ")...)
	result = append(result, encoded...)

	result = append(result, []byte("\nРасшифрованное сообщение: ")...)
	result = append(result, decoded...)

	err := os.WriteFile(f.outPath, result, 0666)
	if err != nil {
		return err
	}
	return nil
}

func NewFilesData(inpPath, outPath string) FilesData {
	return &Data{inpPath: inpPath, outPath: outPath}
}
