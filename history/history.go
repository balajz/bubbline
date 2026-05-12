package history

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/Balaji01-4D/bubbline/editline"
)

func LoadHistory(fileName string) ([]editline.HistoryEntry, error) {
	f, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return loadHistoryFromFile(f)
}

func loadHistoryFromFile(r io.Reader) ([]editline.HistoryEntry, error) {
	var hist []editline.HistoryEntry
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var entry editline.HistoryEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		hist = append(hist, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return hist, nil
}

func SaveHistory(hist []editline.HistoryEntry, fileName string) (retErr error) {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := f.Close()
		if retErr == nil {
			retErr = closeErr
		}
	}()
	return saveHistoryToFile(f, hist)
}

func AppendHistory(hist []editline.HistoryEntry, fileName string) (retErr error) {
	if len(hist) == 0 {
		return nil
	}
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := f.Close()
		if retErr == nil {
			retErr = closeErr
		}
	}()
	return saveHistoryToFile(f, hist)
}

func saveHistoryToFile(w io.Writer, hist []editline.HistoryEntry) error {
	enc := json.NewEncoder(w)
	for _, entry := range hist {
		if err := enc.Encode(entry); err != nil {
			return err
		}
	}
	return nil
}
