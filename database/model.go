package database

import (
	"errors"
	"fmt"
	"time"
)

//Any changes to this model, needs a database reset
//well, at least seems that way...
// sudo -u postgres dropdb prism
// sudo -u postgres createdb prism

//Model is the database object
type Model struct {
	UID       uint      `gorm:"primaryKey" json:"uid,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Name      string    `gorm:"index" json:"name,omitempty"`
	Image     []byte    `json:"image,omitempty"`
	FileName  string    `gorm:"uniqueIndex" json:"filename,omitempty"`
	ParentID  uint      `gorm:"index" json:"parent_id,omitempty"`
	URL       string    `json:"url,omitempty"`
	Hidden    bool      `json:"hidden,omitempty"`
	Tags      string    `json:"tags,omitempty"`
	FileHash  []byte    `gorm:"uniqueIndex" json:"file_hash,omitempty"`
	Type      string    `gorm:"index" json:"type,omitempty"`
	Config    string    `json:"config,omitempty"`
}

//Add original URL as well as hidden from browse feature.

func (m *Model) Verify() error {
	valid := true
	msg := "Invalid Fields: "
	if m.UID == 0 {
		valid = false
		msg += "UID "
	}
	if m.Name == "" {
		valid = false
		msg += "Name "
	}
	if m.FileName == "" {
		valid = false
		msg += "Filename "
	}
	if !valid {
		return errors.New(msg)
	}
	return nil
}
func (m *Model) VerifyUpload() error {
	fmt.Println("M: ", m.Name, m.ParentID, m.FileName)
	valid := true
	msg := "Invalid Fields: "
	if m.Name == "" {
		valid = false
		msg += "Name "
	}

	if !valid {
		return errors.New(msg)
	}
	return nil
}
