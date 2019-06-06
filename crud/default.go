package crud

import (
	"strings"
	"time"

	"github.com/bejaneps/agreement-api/auth"
	"google.golang.org/api/drive/v2"
)

func filterGoogleTime(t string) string {
	t = t[:strings.Index(t, ".")]

	return t
}

// CreateTemplate creates a new document with contents of template document and gives writer role to email
func CreateTemplate(email, template, title string) (file *drive.File, err error) {
	cnt := auth.GetClient()

	srv, err := drive.New(cnt)
	if err != nil {
		return nil, err
	}

	file, err = srv.Files.Copy(template, &drive.File{
		OwnedByMe:       false,
		CreatedDate:     time.Now().Format(time.RFC3339),
		MimeType:        "application/vnd.google-apps.document",
		Title:           title,
		WritersCanShare: false,
	}).Do()
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Second * 10) //for 500 internal server error we have to wait 10 secs for document id to be generated

	//giving ownership to main Google account
	_, err = srv.Permissions.Insert(file.Id, &drive.Permission{
		Value: auth.GoogleAccount,
		Role:  "owner",
		Type:  "user",
	}).SendNotificationEmails(false).Do()
	if err != nil {
		return nil, err
	}

	_, err = srv.Permissions.Insert(file.Id, &drive.Permission{
		Value: email,
		Role:  "writer",
		Type:  "user",
	}).SendNotificationEmails(false).Do() //SendNotificationEmails() set to false bcuz of sharing quota
	if err != nil {
		return nil, err
	}

	return file, nil
}

// CreateDocument creates a new document, with given title and gives writer role to email.
func CreateDocument(email, title string) (file *drive.File, err error) {
	cnt := auth.GetClient()

	srv, err := drive.New(cnt)
	if err != nil {
		return nil, err
	}

	file, err = srv.Files.Insert(&drive.File{
		OwnedByMe:       false, //service account can't use gdrive interface, that's why false
		CreatedDate:     time.Now().Format(time.RFC3339),
		MimeType:        "application/vnd.google-apps.document",
		Title:           title,
		WritersCanShare: false,
	}).Do()
	if err != nil {
		return nil, err
	}

	//giving ownership to main gdrive account
	_, err = srv.Permissions.Insert(file.Id, &drive.Permission{
		Value: auth.GoogleAccount,
		Role:  "owner",
		Type:  "user",
	}).SendNotificationEmails(false).Do()
	if err != nil {
		return nil, err
	}

	_, err = srv.Permissions.Insert(file.Id, &drive.Permission{
		Value: email,
		Role:  "writer",
		Type:  "user",
	}).SendNotificationEmails(false).Do() //SendNotificationEmails() set to false bcuz of sharing quota
	if err != nil {
		return nil, err
	}

	return file, nil
}

// SetPermission sets read or write permission to a user for a document
func SetPermission(fileID, email, perm string) (err error) {
	cnt := auth.GetClient()

	srv, err := drive.New(cnt)
	if err != nil {
		return err
	}

	prf, err := GetUserDoc(fileID)
	if err != nil {
		return err
	}

	if prf.Owner1 == email || prf.Owner2 == email {
		temp, err := srv.Permissions.GetIdForEmail(email).Do()
		if err != nil {
			return err
		}

		_, err = srv.Permissions.Update(fileID, temp.Id, &drive.Permission{
			Value: email,
			Role:  perm,
			Type:  "user",
		}).Do()
		if err != nil {
			return err
		}
	} else {
		_, err := srv.Permissions.Insert(fileID, &drive.Permission{
			Value: email,
			Role:  perm,
			Type:  "user",
		}).SendNotificationEmails(false).Do()
		if err != nil {
			return err
		}
	}

	return nil
}

// LastModifiedDate returns last modified date of drive file
func LastModifiedDate(fileID string) (string, error) {
	cnt := auth.GetClient()

	srv, err := drive.New(cnt)
	if err != nil {
		return "", err
	}

	file, err := srv.Files.Get(fileID).Do()
	if err != nil {
		return "", err
	}

	return filterGoogleTime(file.ModifiedDate), nil
}
