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
	Name      string    `json:"name,omitempty"`
	Image     []byte    `json:"image,omitempty"`
	Extension string    `gorm:"index" json:"extension,omitempty"`
	Parent    string    `json:"parent,omitempty"`
}

// //AllowedType handles which types are allowed and gives us file extensions
// type AllowedType string

// const (
// 	PNG AllowedType = ".png"
// 	JPG             = ".jpg"
// 	SVG             = ".svg"
// 	GIF             = ".gif"
// )

// func (at AllowedType) IsValid() error {
// 	switch at {
// 	case PNG, JPG, SVG, GIF:
// 		return nil
// 	}
// 	return errors.New("Accepted Image formats are png, jpg, svg, gif")
// }

func (m *Model) Verify() error {
	valid := true
	msg := "Invalid Fields: "
	// if m.UID == 0 {
	// 	valid = false
	// 	msg += "UID "
	// }
	if m.Name == "" {
		valid = false
		msg += "Name "
	}
	if m.Extension == "" {
		valid = false
		msg += "Extension "
	}
	// if err := m.Extension.IsValid(); err != nil {
	// 	valid = false
	// 	msg += err.Error()
	// }
	if !valid {
		return errors.New(msg)
	}
	return nil
}
func (m *Model) VerifyUpload() error {
	fmt.Println("M: ", m.Name, m.Parent, m.Extension)
	valid := true
	msg := "Invalid Fields: "
	if m.Name == "" {
		valid = false
		msg += "Name "
	}
	// if err := m.Extension.IsValid(); err != nil {
	// 	valid = false
	// 	msg += err.Error()
	// }
	if !valid {
		return errors.New(msg)
	}
	return nil
}
