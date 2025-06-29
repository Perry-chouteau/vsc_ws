// internal/folder.go

package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Folder struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	IsActive bool   `json:"is_active,omitempty"`
}

type FolderData struct {
	Folders []Folder `json:"folders"`
	Emojis  []string `json:"emojis,omitempty"`
}

func LoadFolders(path string) (*FolderData, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var fd FolderData
	err = json.Unmarshal(data, &fd)
	if err != nil {
		return nil, err
	}
	return &fd, nil
}

func SaveFolders(path string, fd *FolderData) error {
	data, err := json.MarshalIndent(fd, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

func UpdateWorkspaceFoldersOnly(path string, folders []Folder) error {
	// 1. Read existing file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// 2. Parse into generic map
	var full map[string]interface{}
	if err := json.Unmarshal(data, &full); err != nil {
		return err
	}

	// 3. Convert folders to raw map
	rawFolders := []Folder{}
	for _, f := range folders {
		if f.IsActive != true {
			continue;
		}
		rawFolders = append(rawFolders, Folder{
			Name: f.Name,
			Path: f.Path,
		})
	}

	// 4. Replace only the "folders" key
	full["folders"] = rawFolders

	// 5. Marshal and save updated JSON
	out, err := json.MarshalIndent(full, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, out, 0644)
}