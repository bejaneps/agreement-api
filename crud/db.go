package crud

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/bejaneps/agreement-api/auth"
)

// Document struct represents a document of google drive, with it's owners and permissions
type Document struct {
	DocID       string `json:"doc_id" gorm:"type: TEXT; UNIQUE; PRIMARY_KEY; NOT NULL"`
	DocTitle    string `json:"doc_title" gorm:"type: TEXT; NOT NULL"`
	DocURL      string `json:"doc_url" gorm:"type: TEXT; NOT NULL"`
	Owner1      string `json:"owner1" gorm:"type: TEXT"`
	Owner2      string `json:"owner2" gorm:"type: TEXT"`
	Signed1     int    `json:"signed1" gorm:"type: INTEGER"`
	Signed2     int    `json:"signed2" gorm:"type: INTEGER"`
	DateSigned1 string `json:"date_signed1" gorm:"type: TEXT"`
	DateSigned2 string `json:"date_signed2" gorm:"type: TEXT"`
}

func filterLocalTime() string {
	now := time.Now().Add(time.Hour * (-3)).Format(time.RFC3339)

	if i := strings.Index(now, "+"); i != -1 {
		now = now[:i]
	} else {
		i = strings.Index(now, "-")
		now = now[:i]
	}

	return now
}

// GetUserDoc returns a UserDoc structure with a document info from DB
func GetUserDoc(fileID string) (ud *Document, err error) {
	ud = &Document{}

	DB := auth.GetDB()
	defer DB.Close()

	row := DB.Table("documents").Where("doc_id = ?", fileID).Row()
	err = row.Scan(&ud.DocID, &ud.DocTitle, &ud.DocURL, &ud.Owner1, &ud.Owner2, &ud.Signed1, &ud.Signed2, &ud.DateSigned1, &ud.DateSigned2)
	if err != nil {
		return nil, err
	}

	return ud, nil
}

// GetUserDocList returns a slice of documents that belong to user
func GetUserDocList(email string) ([]Document, error) {
	list := make([]Document, 0, 1)

	DB := auth.GetDB()
	defer DB.Close()

	rows, err := DB.Table("documents").Where("(owner1 = $1 OR owner2 = $1)", email).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := Document{}
		err := rows.Scan(&u.DocID, &u.DocTitle, &u.DocURL, &u.Owner1, &u.Owner2, &u.Signed1, &u.Signed2, &u.DateSigned1, &u.DateSigned2)
		if err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

// AddUserDoc adds a document in a DB
func AddUserDoc(id, title, url, email string) (*Document, error) {
	doc := &Document{
		DocID:       id,
		DocTitle:    title,
		DocURL:      url,
		Owner1:      email,
		Owner2:      "",
		Signed1:     0,
		Signed2:     0,
		DateSigned1: "",
		DateSigned2: "",
	}

	DB := auth.GetDB()
	defer DB.Close()

	err := DB.Table("documents").Save(doc).Error
	if err != nil {
		return nil, err
	}

	return doc, err
}

// AddUserSign sets the sign field of a second user document in DB to either true or false
func AddUserSign(email, fileID string) (*Document, error) {
	temp := &Document{}

	DB := auth.GetDB()
	defer DB.Close()

	row := DB.Table("documents").Where("doc_id = ?", fileID).Row()
	if err := row.Scan(&temp.DocID, &temp.DocTitle, &temp.DocURL, &temp.Owner1, &temp.Owner2, &temp.Signed1, &temp.Signed2, &temp.DateSigned1, &temp.DateSigned2); err != nil {
		return nil, err
	}
	if temp.Owner1 == email {
		err := DB.Table("documents").Where("doc_id = ?", fileID).Updates(map[string]interface{}{
			"signed1":      1,
			"date_signed1": filterLocalTime(),
		}).Error
		if err != nil {
			return nil, err
		}

		temp, err = GetUserDoc(fileID)
		if err != nil {
			return nil, err
		}

		return temp, nil
	} else if temp.Owner2 == email {
		err := DB.Table("documents").Where("doc_id = ?", fileID).Updates(map[string]interface{}{
			"signed2":      1,
			"date_signed2": filterLocalTime(),
		}).Error
		if err != nil {
			return nil, err
		}

		temp, err = GetUserDoc(fileID)
		if err != nil {
			return nil, err
		}

		return temp, nil
	}
	return nil, gorm.ErrRecordNotFound
}

// RemoveUserSign removes signed and signed date fields of a document from DB
func RemoveUserSign(email, fileID string) (*Document, error) {
	doc, err := GetUserDoc(fileID)
	if err != nil {
		return nil, err
	}

	DB := auth.GetDB()
	defer DB.Close()

	if doc.Owner1 == email {
		err := DB.Table("documents").Where("doc_id = ?", fileID).Updates(map[string]interface{}{
			"signed1":      0,
			"date_signed1": "",
		}).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := DB.Table("documents").Where("doc_id = ?", fileID).Updates(map[string]interface{}{
			"signed2":      0,
			"date_signed2": "",
		}).Error
		if err != nil {
			return nil, err
		}
	}

	doc, err = GetUserDoc(fileID)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// AddDocOwner adds the user 'email' as the owner of doc 'fileID' to DB
func AddDocOwner(email, fileID string) (*Document, error) {
	DB := auth.GetDB()
	defer DB.Close()

	err := DB.Table("documents").Where("doc_id = ?", fileID).Update("owner2", email).Error
	if err != nil {
		return nil, err
	}

	temp, err := GetUserDoc(fileID)
	if err != nil {
		return nil, err
	}

	return temp, nil
}
