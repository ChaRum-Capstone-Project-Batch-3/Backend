package cloudinary

import (
	"charum/util"
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Function interface {
	Upload(folder string, file *multipart.FileHeader, filename string) (string, error)
	Rename(folder string, oldFilename string, newFilename string) (string, error)
	Delete(folder string, filename string) error
}

type Cloudinary struct {
	cloudinary *cloudinary.Cloudinary
}

var (
	folderBase string
)

func Init(folderName string) Function {
	cld, err := cloudinary.NewFromParams(util.GetConfig("CLOUDINARY_CLOUD_NAME"), util.GetConfig("CLOUDINARY_API_KEY"), util.GetConfig("CLOUDINARY_API_SECRET"))
	if err != nil {
		panic(err)
	}

	folderBase = folderName

	cld.Config.URL.Secure = true
	return &Cloudinary{
		cloudinary: cld,
	}
}

func (c *Cloudinary) Upload(folder string, file *multipart.FileHeader, filename string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	profilePictureBuffer, err := file.Open()
	if err != nil {
		return "", err
	}

	resp, err := c.cloudinary.Upload.Upload(ctx, profilePictureBuffer, uploader.UploadParams{
		PublicID:       filename,
		UniqueFilename: api.Bool(false),
		Folder:         fmt.Sprintf("%s/%s", folderBase, folder),
		Overwrite:      api.Bool(true),
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

func (c *Cloudinary) Rename(folder string, oldFilename string, newFilename string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := c.cloudinary.Upload.Rename(ctx, uploader.RenameParams{
		FromPublicID: fmt.Sprintf("%s/%s/%s", folderBase, folder, oldFilename),
		ToPublicID:   fmt.Sprintf("%s/%s/%s", folderBase, folder, newFilename),
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

func (c *Cloudinary) Delete(folder string, filename string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := c.cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: fmt.Sprintf("%s/%s/%s", folderBase, folder, filename),
	})
	if err != nil {
		return err
	}

	return nil
}
