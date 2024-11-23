package vars

import (
	"encoding/csv"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

var (
	// templates = template.Must(template.ParseFiles("edit.html", "view.html"))
	ValidPath = regexp.MustCompile("^/([a-z0-9.-]+)$")

	BaseDir = "s3"

	PortFlag = flag.Int("port", 3000, "port flag")
	DirFlag  = flag.String("dir", "s3/data", "directory flag")
	HelpFlag = flag.Bool("help", false, "help flag")
)

type ErrorResponse struct {
	XMLName xml.Name `xml:"Error"`
	Message string   `xml:"Message"`
}

func HelpFunc() {
	fmt.Println("Simple Storage Service.\n")
	fmt.Println("**Usage:**")
	fmt.Println("    triple-s [-port <N>] [-dir <S>]")
	fmt.Println("    triple-s --help\n")
	fmt.Println("**Options:**")
	fmt.Println("- --help     Show this screen.")
	fmt.Println("- --port N   Port number")
	fmt.Println("- --dir S    Path to the directory")
}

func GetObjectContentType(bucketName string, objectName string) (string, error) {
	f, err := os.Open(*DirFlag + "/" + bucketName + "/objects.csv")
	if err != nil {
		return "", err
	}
	reader := csv.NewReader(f)

	lines, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	for i, line := range lines {
		if len(line) >= 4 && i != 0 && line[0] == objectName {
			return line[2], nil
		}
	}
	return "", errors.New("metadata of object not found")
}

func PrintXMLError(w http.ResponseWriter, status int, err string) {
	errorResponse := ErrorResponse{
		Message: err,
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	xml.NewEncoder(w).Encode(errorResponse)
}

func GetFileExtension(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/bmp":
		return ".bmp"
	case "image/webp":
		return ".webp"
	case "image/tiff":
		return ".tiff"
	case "image/svg+xml":
		return ".svg"
	case "application/pdf":
		return ".pdf"
	case "application/msword":
		return ".doc"
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return ".docx"
	case "application/vnd.ms-excel":
		return ".xls"
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return ".xlsx"
	case "application/vnd.ms-powerpoint":
		return ".ppt"
	case "application/vnd.openxmlformats-officedocument.presentationml.presentation":
		return ".pptx"
	case "application/zip":
		return ".zip"
	case "application/x-rar-compressed":
		return ".rar"
	case "application/x-7z-compressed":
		return ".7z"
	case "application/x-tar":
		return ".tar"
	case "application/x-bzip2":
		return ".bz2"
	case "application/x-gzip":
		return ".gz"
	case "application/json":
		return ".json"
	case "application/xml":
		return ".xml"
	case "text/plain":
		return ".txt"
	case "text/html":
		return ".html"
	case "text/css":
		return ".css"
	case "text/javascript":
		return ".js"
	case "audio/mpeg":
		return ".mp3"
	case "audio/wav":
		return ".wav"
	case "audio/ogg":
		return ".ogg"
	case "video/mp4":
		return ".mp4"
	case "video/x-msvideo":
		return ".avi"
	case "video/x-ms-wmv":
		return ".wmv"
	case "video/webm":
		return ".webm"
	case "video/quicktime":
		return ".mov"
	case "application/octet-stream":
		return ".bin"
	default:
		return ".bin" // Для неизвестных типов
	}
}
