package nova

/*
	AWS Nova: https://docs.aws.amazon.com/nova/latest/userguide/what-is-nova.html
*/

type Body struct {
	System          []System         `json:"system,omitempty"`
	Messages        []Message        `json:"messages,omitempty"`
	InferenceConfig *InferenceConfig `json:"inferenceConfig,omitempty"`
	/*
	 (Optional) JSON object following [ToolConfig schema], containing the tool specification and tool choice.
	 This schema is the same followed [by the Converse API]

	 [ToolConfig schema]: https://docs.aws.amazon.com/bedrock/latest/APIReference/API_runtime_ToolConfiguration.html
	 [by the Converse API]: https://docs.aws.amazon.com/bedrock/latest/userguide/tool-use.html
	*/
	ToolConfig *ToolConfig `json:"toolConfig,omitempty"`
}

type System struct {
	Text string `json:"text,omitempty"`
}

type Role string

const (
	User      Role = "user"
	Assistant Role = "assistant"
)

type Message struct {
	Role    Role      `json:"role,omitempty"`
	Content []Content `json:"content,omitempty"`
}

type InferenceConfig struct {
	/*
		(Optional) The maximum number of tokens to generate before stopping.

		Note that Amazon Nova models might stop generating tokens before reaching the value of max_tokens.
		The Maximum New Tokens value allowed is 5K.
	*/
	MaxNewTokens int `json:"max_new_tokens,omitempty"`
	/*
		(Optional) The amount of randomness injected into the response. Valid values are between 0.00001 and 1, inclusive.
		The default value is 0.7.
	*/
	Temperature float64 `json:"temperature,omitempty"`
	// (Optional) Use nucleus sampling
	TopP float64 `json:"top_p,omitempty"`
	// (Optional) Only sample from the top K options for each subsequent token.
	TopK          int      `json:"top_k,omitempty"`
	StopSequences []string `json:"stopSequences,omitempty"`
}

type ToolConfig struct {
	Tools []Tool `json:"tools,omitempty"`
}

type Tool struct {
	ToolSpec *ToolSpec `json:"toolSpec,omitempty"`
}

type ToolSpec struct {
	Name        string       `json:"name,omitempty"`
	Description string       `json:"description,omitempty"`
	InputSchema *InputSchema `json:"inputSchema,omitempty"`
}

type InputSchema struct {
	JSON any `json:"json,omitempty"`
}

type Content struct {
	Text     string    `json:"text,omitempty"`
	Image    *Image    `json:"image,omitempty"`
	Video    *Video    `json:"video,omitempty"`
	Document *Document `json:"document,omitempty"`
}

type DocumentFormat string

// Text document format
const (
	TXT  DocumentFormat = "txt"
	CSV  DocumentFormat = "csv"
	MD   DocumentFormat = "md"
	XLS  DocumentFormat = "xls"
	XLSX DocumentFormat = "xlsx"
	HTML DocumentFormat = "html"
	XML  DocumentFormat = "xml"
	DOC  DocumentFormat = "doc"
)

// Media document format
const (
	PDF  DocumentFormat = "pdf"
	DOCX DocumentFormat = "docx"
)

type Document struct {
	Format DocumentFormat  `json:"format,omitempty"`
	Name   string          `json:"name,omitempty"`
	Source *DocumentSource `json:"source,omitempty"`
}

type DocumentSource struct {
	Bytes []byte `json:"bytes,omitempty"`
}

type ImageFormat string

const (
	JPEG ImageFormat = "jpeg"
	PNG  ImageFormat = "png"
	GIF  ImageFormat = "gif"
	WEBP ImageFormat = "webp"
)

type ImageSource struct {
	// Binary array (Converse API) or Base64-encoded string (Invoke API)
	Bytes []byte `json:"bytes,omitempty"`
}

type Image struct {
	Format ImageFormat  `json:"format,omitempty"`
	Source *ImageSource `json:"source,omitempty"`
}

type VideoFormat string

const (
	MKV      VideoFormat = "mkv"
	MOV      VideoFormat = "mov"
	MP4      VideoFormat = "mp4"
	WEBM     VideoFormat = "webm"
	THREE_GP VideoFormat = "three_gp"
	FLV      VideoFormat = "flv"
	MPEG     VideoFormat = "mpeg"
	MPG      VideoFormat = "mpg"
	WMV      VideoFormat = "wmv"
)

type S3Location struct {
	URI         string `json:"uri,omitempty"`
	BucketOwner string `json:"bucketOwner,omitempty"`
}

type VideoSource struct {
	S3Location *S3Location `json:"s3Location,omitempty"`
	Bytes      []byte      `json:"bytes,omitempty"`
}

type Video struct {
	Format VideoFormat  `json:"format,omitempty"`
	Source *VideoSource `json:"source,omitempty"`
}
