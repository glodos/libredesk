package models

import (
	"encoding/json"
	"errors"
	"mime"
	"strings"
	"time"

	"github.com/volatiletech/null/v9"
)

const (
	// TODO: pick these table names from their respective package/models/models.go
	ModelMessages     = "messages"
	ModelUser         = "users"
	ModelHelpArticles = "help_articles"

	DispositionInline     = "inline"
	DispositionAttachment = "attachment"

	ContentTypeOctetStream = "application/octet-stream"
	ContentTypePDF         = "application/pdf"
)

// IsPublicModel reports whether media linked to the model type is served without authentication.
func IsPublicModel(modelType string) bool {
	return modelType == ModelHelpArticles
}

// Media represents an uploaded object in DB and storage backend.
type Media struct {
	ID          int             `db:"id" json:"id"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at" json:"updated_at"`
	UUID        string          `db:"uuid" json:"uuid"`
	Store       string          `db:"store" json:"store"`
	Filename    string          `db:"filename" json:"filename"`
	ContentType string          `db:"content_type" json:"content_type"`
	ContentID   string          `db:"content_id" json:"content_id"`
	ModelID     null.Int        `db:"model_id" json:"model_id"`
	Model       null.String     `db:"model_type" json:"model_type"`
	Disposition null.String     `db:"disposition" json:"disposition"`
	Size        int             `db:"size" json:"size"`
	Meta        json.RawMessage `db:"meta" json:"meta"`
	Private     bool            `db:"private" json:"private"`

	// Pseudo fields
	URL     string `json:"url"`
	Content []byte `json:"-"`
}

// NormalizeContentType returns the lowercased type/subtype of a Content-Type value, or application/octet-stream if it does not parse.
func NormalizeContentType(contentType string) string {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if (err != nil && !errors.Is(err, mime.ErrInvalidMediaParameter)) || !strings.Contains(mediaType, "/") {
		return ContentTypeOctetStream
	}
	return mediaType
}

// ContentDisposition returns inline for media and safe browser-readable text.
// Potentially active document types (HTML, XML/SVG, JavaScript) and unknown/binary
// content remain attachments.
func ContentDisposition(contentType string) string {
	contentType = NormalizeContentType(contentType)

	if strings.HasSuffix(contentType, "+xml") ||
		contentType == "text/html" ||
		contentType == "application/xhtml+xml" ||
		contentType == "application/xml" ||
		contentType == "text/xml" ||
		contentType == "application/javascript" ||
		contentType == "text/javascript" {
		return DispositionAttachment
	}

	if strings.HasPrefix(contentType, "image/") ||
		strings.HasPrefix(contentType, "video/") ||
		strings.HasPrefix(contentType, "audio/") ||
		strings.HasPrefix(contentType, "text/") ||
		contentType == ContentTypePDF ||
		contentType == "application/json" ||
		strings.HasSuffix(contentType, "+json") {
		return DispositionInline
	}

	return DispositionAttachment
}
