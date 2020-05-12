package util

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io/ioutil"
	"net/url"
	"strings"
)

func Load(s string) ([]byte, error) {
	if strings.HasPrefix(s, "s3://") {
		sess := session.Must(session.NewSession())
		return LoadFromS3(sess, s)
	} else if strings.HasPrefix(s, "file://") {
		b, err := ioutil.ReadFile(s[7:])
		return b, err
	} else if strings.HasPrefix(s, "data://") {
		b, err := base64.StdEncoding.DecodeString(s[7:])
		return b, err
	}
	return []byte(s), nil
}

func LoadFromS3(sess *session.Session, s3uri string) ([]byte, error) {
	u, err := url.Parse(s3uri)
	if err != nil {
		return nil, err
	}
	bucket := u.Host
	key := strings.TrimLeft(u.Path, "/")

	buf := &aws.WriteAtBuffer{}
	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
