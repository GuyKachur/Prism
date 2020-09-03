package database

import (
	"errors"
	"fmt"
	"time"
)

type Model struct {
	UID       uint      `gorm:"primaryKey" json:"uid,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Name      string    `gorm:"index" json:"name,omitempty"`
	Image     string    `gorm:"uniqueIndex" json:"image,omitempty"`
	FileName  string    `gorm:"uniqueIndex" json:"filename,omitempty"`
	Parent    string    `gorm:"index" json:"parent,omitempty"`
	URL       string    `json:"url,omitempty"`
	Hidden    bool      `json:"hidden,omitempty"`
	Tags      string    `json:"tags,omitempty"`
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
	fmt.Println("M: ", m.Name, m.Parent, m.FileName)
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
