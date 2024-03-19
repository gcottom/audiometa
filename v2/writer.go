package audiometa

import "os"

func WriteFile(fn string, data []byte) error {
	file, err := os.OpenFile(fn, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.Write(data); err != nil {
		return err
	}
	return nil
}
