package groups

const (
	ONE       = "one-to-one"
	MANY      = "many-to-many"
	USERON    = "online"
	USEROFF   = "offline"
	thumbnail = ""
)

const location = "us-east-1"

var endpoint = "minio.vdatlab.com"
var accessKeyID = "gfEOZ2vrBYwoUumYEJhbcmoBLbRlonkQ"
var secretAccessKey = "E5cGw3exgtmJVo9Q8cZgTMWJ4XNDKgjH"

const BucketNameAvatar = "avatar"
const BucketThumbnail = "thumbnail"
const BucketNameQR = "qr"
const media = "media"
const file = "file"
const all = "all"

var listTypeFile = [...]string{".png", ".svg", ".jpg", ".mp3", ".mp4", ".jpeg", ".gif"}
