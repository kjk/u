package u

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kjk/atomicfile"
	"github.com/minio/minio-go/v6"
)

// MinioClient represents s3/spaces etc. client
type MinioClient struct {
	StorageKey    string
	StorageSecret string
	Bucket        string
	Endpoint      string // e.g. "nyc3.digitaloceanspaces.com"
	Secure        bool
	client        *minio.Client
}

// EnsureCondfigured will panic if client not configured
func (c *MinioClient) EnsureConfigured() {
	PanicIf(c.StorageKey == "", "minio storage key not set")
	PanicIf(c.StorageSecret == "", "minio storage secret not set")
	PanicIf(c.Bucket == "", "minio bucket not set")
	PanicIf(c.Endpoint == "", "minio endpoint not set")
}

// URLBase returns base of url under which files are accesible
// (if public it's a publicly available file)
func (c *MinioClient) URLBase() string {
	return fmt.Sprintf("https://%s.%s/", c.Bucket, c.Endpoint)
}

// GetClient returns a (cached) minio client
func (c *MinioClient) GetClient() (*minio.Client, error) {
	if c.client != nil {
		return c.client, nil
	}

	var err error
	c.client, err = minio.New(c.Endpoint, c.StorageKey, c.StorageSecret, c.Secure)
	return c.client, err
}

// ListRemoveFiles returns a list of files under a given prefix
func (c *MinioClient) ListRemoteFiles(prefix string) ([]*minio.ObjectInfo, error) {
	var res []*minio.ObjectInfo
	client, err := c.GetClient()
	if err != nil {
		return nil, err
	}
	doneCh := make(chan struct{})
	defer close(doneCh)

	files := client.ListObjectsV2(c.Bucket, prefix, true, doneCh)
	for oi := range files {
		oic := oi
		res = append(res, &oic)
	}
	return res, nil
}

// IsMinioNotExistsError returns true if an error indicates that a key
// doesn't exist in storage
func IsMinioNotExistsError(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "The specified key does not exist."
}

// SetPublicObjectMetadata sets options that mark object as public
// for doing put operation
func SetPublicObjectMetadata(opts *minio.PutObjectOptions) {
	if opts.UserMetadata == nil {
		opts.UserMetadata = map[string]string{}
	}
	opts.UserMetadata["x-amz-acl"] = "public-read"
}

func (c *MinioClient) UploadReaderPublic(remotePath string, r io.Reader, size int64, contentType string) error {
	return c.UploadReader(remotePath, r, size, true, contentType)
}

func (c *MinioClient) UploadReaderPrivate(remotePath string, r io.Reader, size int64, contentType string) error {
	return c.UploadReader(remotePath, r, size, false, contentType)
}

func (c *MinioClient) UploadData(remotePath string, d []byte, opts minio.PutObjectOptions) error {
	client, err := c.GetClient()
	if err != nil {
		return err
	}
	r := bytes.NewReader(d)
	size := int64(len(d))
	_, err = client.PutObject(c.Bucket, remotePath, r, size, opts)
	return err
}

func (c *MinioClient) UploadReader(remotePath string, r io.Reader, size int64, public bool, contentType string) error {
	PanicIf(remotePath[0] == '/', "name '%s' shouldn't start with '/'", remotePath)

	if contentType == "" {
		contentType = MimeTypeFromFileName(remotePath)
	}
	//timeStart := time.Now()
	//sizeStr := humanize.Bytes(uint64(size))
	//fmt.Printf("Uploading '%s' of size %s and type %s as public.", remotePath, sizeStr, contentType)
	client, err := c.GetClient()
	if err != nil {
		return err
	}
	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}
	if public {
		SetPublicObjectMetadata(&opts)
	}
	opts.ContentType = contentType
	_, err = client.PutObject(c.Bucket, remotePath, r, size, opts)
	if err != nil {
		return err
	}
	//fmt.Printf(" Took %s.\n", time.Since(timeStart))
	return nil
}

func (c *MinioClient) DownloadFileAsData(remotePath string) ([]byte, error) {
	client, err := c.GetClient()
	if err != nil {
		return nil, err
	}
	opts := minio.GetObjectOptions{}
	obj, err := client.GetObject(c.Bucket, remotePath, opts)
	if err != nil {
		return nil, err
	}
	defer obj.Close()
	var buf bytes.Buffer
	_, err = io.Copy(&buf, obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *MinioClient) DownloadFileAtomically(dstPath string, remotePath string) error {
	client, err := c.GetClient()
	if err != nil {
		return err
	}
	opts := minio.GetObjectOptions{}
	obj, err := client.GetObject(c.Bucket, remotePath, opts)
	if err != nil {
		return err
	}
	defer obj.Close()

	// ensure there's a dir for destination file
	dir := filepath.Dir(dstPath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	f, err := atomicfile.New(dstPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, obj)
	if err != nil {
		return err
	}
	return f.Close()
}

func (c *MinioClient) UploadFilePublic(remotePath string, filePath string) error {
	return c.UploadFile(remotePath, filePath, true)
}

func (c *MinioClient) UploadFilePrivate(remotePath string, filePath string) error {
	return c.UploadFile(remotePath, filePath, false)
}

func (c *MinioClient) UploadFile(remotePath string, filePath string, public bool) error {
	stat, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	size := stat.Size()
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func() {
		f.Close()
	}()
	return c.UploadReaderPublic(remotePath, f, size, "")
}

func (c *MinioClient) UploadDataPublic(remotePath string, d []byte) error {
	r := bytes.NewBuffer(d)
	return c.UploadReaderPublic(remotePath, r, int64(len(d)), "")
}

func (c *MinioClient) UploadStringPublic(remotePath string, s string) error {
	r := bytes.NewBufferString(s)
	return c.UploadReaderPublic(remotePath, r, int64(len(s)), "")
}

func (c *MinioClient) UploadStringPrivate(remotePath string, s string) error {
	r := bytes.NewBufferString(s)
	return c.UploadReaderPrivate(remotePath, r, int64(len(s)), "")
}

func (c *MinioClient) StatObject(remotePath string) (minio.ObjectInfo, error) {
	client, err := c.GetClient()
	if err != nil {
		return minio.ObjectInfo{}, err
	}
	var opts minio.StatObjectOptions
	return client.StatObject(c.Bucket, remotePath, opts)
}

func (c *MinioClient) Delete(remotePath string) error {
	client, err := c.GetClient()
	if err != nil {
		return err
	}
	err = client.RemoveObject(c.Bucket, remotePath)
	if IsMinioNotExistsError(err) {
		return nil
	}
	return nil
}
