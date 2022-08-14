package upload

import (
	"bytes"
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/minio/minio-go/v7"
	"github.com/skip2/go-qrcode"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/database"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v1"
	"golang.org/x/sync/errgroup"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func RemoveFileService(idGroup string, objectName string) (err error) {
	bucketName := "group-" + idGroup
	fmt.Println(bucketName)
	checkExists, err := CheckBucketExits(bucketName)
	if err != nil {
		fmt.Printf("CheckBucketExits %+v", err)
		return
	}
	if !checkExists {
		fmt.Println("khong ton tai")
		return
	}
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}

	err = minioClient.RemoveObject(context.Background(), bucketName, objectName, opts)
	if err != nil {
		sentry.CaptureException(err)
		fmt.Printf("RemoveObject %+v", err)
		return
	}
	log.Println("Success")
	return
}

func GetFileService(bucketName string, objectName string) (shareLink string, err error) {
	checkExists, err := CheckBucketExits(bucketName)
	if err != nil {
		fmt.Printf("CheckBucketExits %+v", err)
		return
	}
	if !checkExists {
		return
	}
	value, ok := FileGlobal[objectName]
	if ok {
		if time.Now().Day()-value.CreateAt.Day() < 2 {
			fmt.Println("vào caching")
			shareLink = value.ShareUrl
			return
		}
	}

	fmt.Println("không vào")
	fmt.Println("không vào 2.0")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(7)*time.Second)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		reqParams := make(url.Values)
		reqParams.Set("response-content-disposition", "attachment; filename="+objectName)
		image, err := minioClient.PresignedGetObject(context.Background(), bucketName, objectName, time.Second*60*60*60, reqParams)
		if err != nil {
			fmt.Printf("PresignedGetObject %+v", err)
			sentry.CaptureException(err)
			shareLink = objectName
			return
		}
		shareLink = image.String()

		newImage := ModelUpload{
			ShareUrl: shareLink,
			NameFile: objectName,
			CreateAt: time.Now(),
		}
		FileGlobal[objectName] = newImage
		fmt.Println("lay hinh thanh cong")
		return
	})
	if err = g.Wait(); err != nil {
		shareLink = objectName
		return
	}
	return
}

func CheckTypeFile(file multipart.File) string {
	fileHeader := make([]byte, 512)
	if _, err := file.Read(fileHeader); err != nil {
		fmt.Println("err")
	}
	return http.DetectContentType(fileHeader)
}

func GetListFileInBucket(bucketName string, typeFIle string) []ModelUpload {
	listFile := make([]ModelUpload, 0)
	listMedia := make([]ModelUpload, 0)
	listAll := make([]ModelUpload, 0)
	objectCh := minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{})
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			return listAll
		}
		//fmt.Println(strings.HasSuffix(object.Key, listTypeFile[0]))

		reqParams := make(url.Values)
		reqParams.Set("response-content-disposition", "attachment; filename="+object.Key)
		image, err := minioClient.PresignedGetObject(context.Background(), bucketName, object.Key, time.Second*60*60*60, reqParams)
		if err != nil {
			fmt.Printf("PresignedGetObject %+v", err)
		}
		if image.RawPath == "" {
			newImage := ModelUpload{
				ShareUrl: image.String(),
				Location: image.Host + image.Path,
				NameFile: object.Key,
				Type:     object.ContentType,
			}
			listAll = append(listAll, newImage)
			if strings.HasSuffix(object.Key, listTypeFile[0]) ||
				strings.HasSuffix(object.Key, listTypeFile[1]) ||
				strings.HasSuffix(object.Key, listTypeFile[2]) ||
				strings.HasSuffix(object.Key, listTypeFile[3]) ||
				strings.HasSuffix(object.Key, listTypeFile[4]) ||
				strings.HasSuffix(object.Key, listTypeFile[5]) ||
				strings.HasSuffix(object.Key, listTypeFile[6]) {
				listMedia = append(listMedia, newImage)
			} else {
				fmt.Println("false")
				listFile = append(listFile, newImage)
			}

		} else {
			newImage := ModelUpload{
				ShareUrl: image.String(),
				Location: image.Host + image.RawPath,
				NameFile: object.Key,
				Type:     object.ContentType,
			}
			listAll = append(listAll, newImage)
			if strings.HasSuffix(object.Key, listTypeFile[0]) ||
				strings.HasSuffix(object.Key, listTypeFile[1]) ||
				strings.HasSuffix(object.Key, listTypeFile[2]) ||
				strings.HasSuffix(object.Key, listTypeFile[3]) ||
				strings.HasSuffix(object.Key, listTypeFile[4]) ||
				strings.HasSuffix(object.Key, listTypeFile[5]) ||
				strings.HasSuffix(object.Key, listTypeFile[6]) {
				listMedia = append(listMedia, newImage)
			} else {
				fmt.Println("false")
				listFile = append(listFile, newImage)
			}
		}
	}
	if typeFIle != file && typeFIle != media {
		return listAll
	}

	if typeFIle == all {
		return listAll
	} else if typeFIle == media {
		return listMedia
	} else {
		return listFile
	}
}

func CheckBucketExits(bucketName string) (bool, error) {
	exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
	if errBucketExists == nil && exists {
		return true, nil
	} else {
		return false, errBucketExists
	}
}

func PutObjectTOMinio(bucketName string, objectName string, file multipart.File, contentType string, contentlength int64) (newImage ModelUpload, err error) {
	checkExists, err := CheckBucketExits(bucketName)
	if err != nil {
		fmt.Printf("CheckBucketExits %+v", err)
		sentry.CaptureException(err)
		return ModelUpload{}, err
	}
	if !checkExists {
		return ModelUpload{}, nil
	}
	sentry.CaptureMessage("up hinh")
	n, err := minioClient.PutObject(context.Background(), bucketName, objectName, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		fmt.Printf("PutObject() %+v", err)
		sentry.CaptureException(err)
		fmt.Println(err)
		return ModelUpload{}, err
	}
	sentry.CaptureMessage("up hinh thanh cong")
	fmt.Println(n)
	//log.Println("Successfully uploaded %s of size %d\n", n.Key, n.Location)
	//linkShare, err := GetFileService(bucketName, objectName)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	newImage = ModelUpload{
		ShareUrl: "",
		Location: n.Location,
		NameFile: objectName,
		Type:     contentType,
	}
	return newImage, nil
}
func UpdateUserAvatarService(id string, file multipart.File, handler *multipart.FileHeader, ctx context.Context) (newImage ModelUpload, err error) {
	objectName := time.Now().Format("01-02-2021 15:04:05") + "#" + handler.Filename
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		exists, err := minioClient.BucketExists(ctx, BucketNameAvatar)
		if err != nil {
			sentry.CaptureException(err)
			return
		}
		if !exists {
			err = minioClient.MakeBucket(ctx, BucketNameAvatar, minio.MakeBucketOptions{Region: location})
			if err != nil {
				sentry.CaptureException(err)
				return
			}
			log.Printf("Successfully created %s\n", id)
		}
		contentType := "image/png"
		n, _ := minioClient.PutObject(ctx, BucketNameAvatar, objectName, file, -1, minio.PutObjectOptions{ContentType: contentType})
		if err != nil {
			sentry.CaptureException(err)
			return
		}
		log.Println("Successfully uploaded %s of size %d\n", n.Key, n.Location)
		reqParams := make(url.Values)
		reqParams.Set("response-content-disposition", "attachment; filename="+objectName)
		image, err := minioClient.PresignedGetObject(ctx, BucketNameAvatar, objectName, time.Second*60*60*60, reqParams)
		if err != nil {
			fmt.Println(err)
			sentry.CaptureException(err)
			return
		}
		newImage = ModelUpload{
			ShareUrl: image.String(),
			Location: n.Location,
			NameFile: objectName,
			Type:     contentType,
		}
		return
	})
	g.Go(func() (err error) {
		err = userdetail.NewRepoImpl(database.DB).UpdateUserById(id, objectName)
		if err != nil {
			sentry.CaptureException(err)
			return
		}
		return
	})
	if err = g.Wait(); err != nil {
		sentry.CaptureException(err)
		return
	}
	return
}

func CreateQRCodeAndUpload(idGroup int, ctx context.Context) ModelUpload {
	png, err := qrcode.Encode("https://example.org", qrcode.Medium, 256)
	if err != nil {
		log.Fatalln(err)
	}

	checkExists, err := CheckBucketExits(BucketNameQR)
	if err != nil {
		fmt.Printf("CheckBucketExits %+v", err)
		return ModelUpload{}
	}
	if checkExists {
		log.Printf("We already own %s\n", BucketNameQR)
	} else {
		err = minioClient.MakeBucket(context.Background(), BucketNameQR, minio.MakeBucketOptions{Region: location})
		log.Printf("Successfully created %s\n", BucketNameQR)
	}

	reader := bytes.NewReader(png)
	n, err := minioClient.PutObject(ctx, BucketNameQR, string(rune(idGroup))+".png", reader, -1, minio.PutObjectOptions{ContentType: "image/png"})
	if err != nil {
		sentry.CaptureException(err)
		fmt.Println(err)
		return ModelUpload{}
	}
	fmt.Println(n)
	linkShare, _ := GetFileService(BucketNameQR, string(rune(idGroup)))
	if err != nil {
		log.Fatalln(err)
	}

	newImage := ModelUpload{
		ShareUrl: linkShare,
		Location: n.Location,
		NameFile: string(rune(idGroup)) + ".png",
		Type:     "image/png",
	}
	return newImage
}
