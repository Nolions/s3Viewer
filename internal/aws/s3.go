package aws

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	awsConf "github.com/Nolions/s3Viewer/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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

	endpoint := strings.TrimSpace(conf.Endpoint)
	if endpoint != "" && !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "http://" + endpoint
	}

	client := s3.NewFromConfig(*cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
		o.UsePathStyle = conf.UsePathStyle
	})

	return &S3Client{
		client: client,
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

type ProgressCallback func(currentBytes, totalBytes int64)

type ProgressReader struct {
	reader       io.Reader
	totalBytes   int64
	currentBytes int64
	onProgress   ProgressCallback
	lastUpdate   time.Time
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.reader.Read(p)
	if n > 0 {
		pr.currentBytes += int64(n)
		if pr.onProgress != nil {
			now := time.Now()
			if pr.lastUpdate.IsZero() || now.Sub(pr.lastUpdate) >= 20*time.Millisecond || (pr.totalBytes > 0 && pr.currentBytes == pr.totalBytes) || err == io.EOF {
				pr.lastUpdate = now
				pr.onProgress(pr.currentBytes, pr.totalBytes)
			}
		}
	}
	return n, err
}

func (pr *ProgressReader) Seek(offset int64, whence int) (int64, error) {
	if seeker, ok := pr.reader.(io.Seeker); ok {
		n, err := seeker.Seek(offset, whence)
		if err == nil {
			pr.currentBytes = n
		}
		return n, err
	}
	return 0, fmt.Errorf("underlying reader does not implement io.Seeker")
}

// DownloadFileWithProgress
// 下載單一檔案到本機目錄中（帶進度追蹤）
func (c *S3Client) DownloadFileWithProgress(key, destPath string, onProgress ProgressCallback) error {
	// 檢查預計儲存目錄是否存在
	dir := filepath.Dir(destPath)
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

	totalBytes := aws.ToInt64(resp.ContentLength)
	pr := &ProgressReader{
		reader:     resp.Body,
		totalBytes: totalBytes,
		onProgress: onProgress,
	}

	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(f, pr, buf)
	if err != nil {
		return err
	}

	if onProgress != nil && totalBytes > 0 {
		onProgress(totalBytes, totalBytes)
	}

	return nil
}

// DownloadFile
// 下載單一檔案到本機目錄中
func (c *S3Client) DownloadFile(key, destPath string) error {
	return c.DownloadFileWithProgress(key, destPath, nil)
}

// UploadFileWithProgress
// 上傳檔案到s3（帶進度追蹤）
func (c *S3Client) UploadFileWithProgress(filePath, key string, onProgress ProgressCallback) error {
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return fmt.Errorf("local file path cannot be empty")
	}

	key = strings.TrimSpace(key)
	if key == "" {
		key = filepath.Base(filePath)
	}

	if strings.HasPrefix(key, "/") {
		key = strings.TrimPrefix(key, "/")
	}

	if key == "" {
		return fmt.Errorf("file key cannot be empty")
	}

	// 檢查檔案是否存在
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	totalBytes := int64(0)
	if err == nil {
		totalBytes = fileInfo.Size()
	}

	pr := &ProgressReader{
		reader:     file,
		totalBytes: totalBytes,
		onProgress: onProgress,
	}

	_, err = c.client.PutObject(c.ctx, &s3.PutObjectInput{
		Bucket:        aws.String(c.bucket),
		Key:           aws.String(key),
		Body:          pr,
		ContentLength: aws.Int64(totalBytes),
	})
	if err != nil {
		return err
	}

	if onProgress != nil && totalBytes > 0 {
		onProgress(totalBytes, totalBytes)
	}

	return nil
}

// UploadFile
// 上傳檔案到s3
func (c *S3Client) UploadFile(filePath, key string) error {
	return c.UploadFileWithProgress(filePath, key, nil)
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

// DeleteObject
// 刪除 s3 上的指定物件
func (c *S3Client) DeleteObject(key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return fmt.Errorf("file key cannot be empty")
	}

	_, err := c.client.DeleteObject(c.ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}

// GetObjectData
// 讀取物件完整 byte 內容
func (c *S3Client) GetObjectData(key string) ([]byte, error) {
	resp, err := c.client.GetObject(c.ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// GetPresignedURL
// 產生指定物件帶時效的 Presigned GET 下載連結
func (c *S3Client) GetPresignedURL(key string, lifetime time.Duration) (string, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return "", fmt.Errorf("file key cannot be empty")
	}

	if lifetime <= 0 {
		lifetime = 15 * time.Minute
	}

	presignClient := s3.NewPresignClient(c.client)
	req, err := presignClient.PresignGetObject(c.ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = lifetime
	})
	if err != nil {
		return "", err
	}

	return req.URL, nil
}
