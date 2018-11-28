# awssessioncache

AWS SDK session cache by region

## Usage

```
import (
    sc "github.com/rynop/awssessioncache"
)
...
	sess, err := sc.Get(&sessioncache.Conf{})
	if err != nil {
		fmt.Println("oops",err)
	}
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      &bucket,
		Key:         &key,
		Body:        blah,
		ContentType: aws.String("application/json"),
		ACL:         aws.String("public-read"),
	})
```
