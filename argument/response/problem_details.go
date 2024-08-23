package response

import "time"

// ProblemDetails is Problem Details for HTTP APIs.
// RFC7807 https://www.rfc-editor.org/rfc/rfc7807
type ProblemDetails struct {
	// Title A short, human-readable summary of the problem type.  It SHOULD NOT change from
	// occurrence to occurrence of the problem, except for purposes of localization (e.g., using
	// proactive content negotiation; see [RFC7231], Section 3.4).
	Title string `json:"title" xml:"title"`

	// Status The HTTP status code ([RFC7231], Section 6) generated by the origin server for
	// this occurrence of the problem.
	Status int `json:"status" xml:"status"`

	// Detail A human-readable explanation specific to this occurrence of the problem.
	Detail string `json:"detail" xml:"detail"`

	// Instance A URI reference that identifies the specific occurrence of the problem.
	// It may or may not yield further information if dereferenced.
	Instance string `json:"instance" xml:"instance"`

	// Method 请求方法，扩充字段。
	Method string `json:"method" xml:"method"`

	// Datetime 报错时间，扩充字段。
	Datetime time.Time `json:"datetime" xml:"datetime"`
}
