package library

/*
	<summary>
		Multipurpose Internet Mail Extensions
		Override default system types
	</summary>
*/
var MIMEs = map[string]string{
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":       ".xlsx",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
	"application/zip": ".zip",
	"image/jpeg":      ".jpeg",
	"image/webp":      ".webp",
	"text/xml":        ".xml",
	"video/mp4":       ".mp4",
}
