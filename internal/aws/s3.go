package aws

import (
	"context"
	"fmt"
	awsConf "github.com/Nolions/s3Viewer/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
	"os"
	"strings"
	"time"
)

type S3Client struct {
	client *s3.Client
	ctx    context.Context
	bucket string
}

// NewS3Client
// 新增S3 Client
func NewS3Client(ctx context.Context, conf awsConf.AWSConfig) (*S3Client, error) {
	cfg, err := newConfig(conf)
	if err != nil {
		return nil, err
	}

	return &S3Client{
		client: s3.NewFromConfig(*cfg),
		ctx:    ctx,
		bucket: conf.Bucket,
	}, nil
}

type PrefixCont struct {
	Dirs  []string
	Files []FileInfo
}

type FileInfo struct {
	Name string
	Key  string
	Time time.Time
	Size int64
}

type FileDetail struct {
	AcceptRanges  string
	UpdateTime    time.Time
	ContentLength int64
	ContentType   string
	Encryption    string
}

// CheckHeadBucket
// 檢查Bucket是否可以存取
func (c *S3Client) CheckHeadBucket() error {
	_, err := c.client.HeadBucket(c.ctx, &s3.HeadBucketInput{
		Bucket: &c.bucket,
	})
	if err != nil {
		return err
	}
	return nil
}

// ListPrefix
// 列出指定目錄下的檔案與目錄
func (c *S3Client) ListPrefix(prefix string) (PrefixCont, error) {
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	var token *string
	objs := PrefixCont{}
	for {
		out, err := c.client.ListObjectsV2(c.ctx, &s3.ListObjectsV2Input{
			Bucket:            aws.String(c.bucket),
			Prefix:            aws.String(prefix),
			Delimiter:         aws.String("/"),
			ContinuationToken: token,
			MaxKeys:           aws.Int32(1000),
		})
		if err != nil {
			return PrefixCont{}, err
		}

		collectFolders(&objs.Dirs, out.CommonPrefixes, prefix)
		collectObjects(&objs.Files, out.Contents, prefix)

		if out.IsTruncated == nil || !*out.IsTruncated {
			break
		}
		token = out.NextContinuationToken
	}

	return objs, nil
}

func collectFolders(dirs *[]string, commonPrefixes []types.CommonPrefix, prefix string) {
	for _, cp := range commonPrefixes {
		name := strings.TrimSuffix(strings.TrimPrefix(aws.ToString(cp.Prefix), prefix), "/")
		*dirs = append(*dirs, name)
	}
}

func collectObjects(files *[]FileInfo, contents []types.Object, prefix string) {
	for _, obj := range contents {
		key := aws.ToString(obj.Key)
		name := strings.TrimPrefix(key, prefix)
		if name == "" {
			continue
		}

		if obj.Size == nil {
			continue
		}

		f := FileInfo{
			Key:  key,
			Name: name,
			Size: *obj.Size,
			Time: aws.ToTime(obj.LastModified).Local(),
		}

		*files = append(*files, f)
	}
}

// DownloadFile
// 下載單一檔案到本機目錄中
func (c *S3Client) DownloadFile(key, destPath string) error {
	// 檢查預計儲存目錄是否存在
	dir := destPath[:strings.LastIndex(destPath, string(os.PathSeparator))]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 建立檔案
	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", destPath, err)
	}
	defer f.Close()

	// 下載檔案
	resp, err := c.client.GetObject(c.ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// UploadFile
// 上傳檔案到s3
func (c *S3Client) UploadFile(filePath, key string) error {
	if strings.HasPrefix(key, "/") {
		key = strings.TrimPrefix(key, "/")
	}

	// 檢查檔案是否存在
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = c.client.PutObject(c.ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetDetail
// 取得檔案的詳細資訊
func (c *S3Client) GetDetail(key string) (FileDetail, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	o, err := c.client.HeadObject(c.ctx, input)
	if err != nil {
		return FileDetail{}, err
	}
	
	return FileDetail{
		AcceptRanges:  aws.ToString(o.AcceptRanges),
		UpdateTime:    aws.ToTime(o.LastModified),
		ContentLength: aws.ToInt64(o.ContentLength),
		ContentType:   aws.ToString(o.ContentType),
		Encryption:    string(o.ServerSideEncryption),
	}, nil

}
